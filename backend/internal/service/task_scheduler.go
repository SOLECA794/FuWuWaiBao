package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/repository"
)

// TaskSchedulerService 定时任务调度服务
type TaskSchedulerService struct {
	db         *gorm.DB
	cron       *cron.Cron
	redis      *repository.RedisClient
	taskQueue  chan *model.ScheduledTask
	workerPool chan struct{} // 限制并发数
}

// NewTaskSchedulerService 创建定时任务调度服务
func NewTaskSchedulerService(db *gorm.DB, redisClient *repository.RedisClient) *TaskSchedulerService {
	return &TaskSchedulerService{
		db:         db,
		cron:       cron.New(),
		redis:      redisClient,
		taskQueue:  make(chan *model.ScheduledTask, 1000), // 任务队列缓冲
		workerPool: make(chan struct{}, 10),               // 最多10个并发worker
	}
}

// Start 启动调度器
func (ts *TaskSchedulerService) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动cron调度器
	ts.cron.Start()

	// 启动worker池
	for i := 0; i < 10; i++ {
		go ts.worker(ctx)
	}

	// 每分钟检查一次待执行任务
	ts.cron.AddFunc("@every 1m", func() {
		ts.checkPendingTasks()
	})

	// 启动时加载现有任务
	ts.loadExistingTasks()

	log.Println("TaskSchedulerService started")
}

// Stop 停止调度器
func (ts *TaskSchedulerService) Stop() {
	ts.cron.Stop()
	close(ts.taskQueue)
	close(ts.workerPool)
}

// ScheduleTask 调度一个任务
func (ts *TaskSchedulerService) ScheduleTask(taskType string, taskData interface{}, scheduledAt time.Time, studentID string) error {
	data, err := json.Marshal(taskData)
	if err != nil {
		return err
	}

	task := &model.ScheduledTask{
		TaskType:    taskType,
		TaskData:    string(data),
		ScheduledAt: scheduledAt,
		StudentID:   studentID,
		Priority:    1,
	}

	if err := ts.db.Create(task).Error; err != nil {
		return err
	}

	// 如果是立即执行的任务，直接加入队列
	if scheduledAt.Before(time.Now().Add(time.Minute)) {
		select {
		case ts.taskQueue <- task:
		default:
			log.Printf("Task queue full, task %s will be processed later", task.ID)
		}
	}

	return nil
}

// ScheduleReviewReminders 调度复习提醒
func (ts *TaskScheduler) ScheduleReviewReminders() {
	var plans []model.ReviewPlan
	ts.db.Where("status = ? AND next_review_date IS NOT NULL", "active").Find(&plans)

	for _, plan := range plans {
		if plan.NextReviewDate != nil && plan.NextReviewDate.After(time.Now()) {
			// 提前1小时提醒
			remindAt := plan.NextReviewDate.Add(-time.Hour)

			taskData := map[string]interface{}{
				"plan_id":   plan.ID,
				"plan_name": plan.Name,
				"student_id": plan.StudentID,
			}

			ts.ScheduleTask("review_reminder", taskData, remindAt, plan.StudentID)
		}
	}
}

// SchedulePracticeGeneration 调度练习题生成
func (ts *TaskScheduler) SchedulePracticeGeneration() {
	var weakPoints []model.WeakPoint
	ts.db.Where("mastery_level < ?", 70).Find(&weakPoints)

	for _, wp := range weakPoints {
		// 为每个薄弱点生成练习题
		taskData := map[string]interface{}{
			"weak_point_id": wp.ID,
			"student_id":    wp.StudentID,
			"course_id":     wp.CourseID,
			"difficulty":    3, // 中等难度
			"count":         5, // 5道题
		}

		// 每天早上9点生成
		now := time.Now()
		scheduledAt := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
		if scheduledAt.Before(now) {
			scheduledAt = scheduledAt.Add(24 * time.Hour)
		}

		ts.ScheduleTask("practice_generation", taskData, scheduledAt, wp.StudentID)
	}
}

// worker 工作协程
func (ts *TaskSchedulerService) worker(ctx context.Context) {
	for {
		select {
		case task := <-ts.taskQueue:
			ts.workerPool <- struct{}{} // 获取worker slot
			go func(t *model.ScheduledTask) {
				defer func() { <-ts.workerPool }() // 释放worker slot
				ts.processTask(t)
			}(task)
		case <-ctx.Done():
			return
		}
	}
}

// processTask 处理单个任务
func (ts *TaskSchedulerService) processTask(task *model.ScheduledTask) {
	// 更新任务状态
	task.Status = "processing"
	task.LastAttempt = &time.Time{}
	*task.LastAttempt = time.Now()
	ts.db.Save(task)

	// 创建任务状态记录
	taskStatus := &model.TaskStatus{
		TaskID:    task.ID,
		TaskType:  task.TaskType,
		StudentID: task.StudentID,
		Status:    "running",
		StartTime: task.LastAttempt,
	}
	ts.db.Create(taskStatus)

	var err error
	switch task.TaskType {
	case "review_reminder":
		err = ts.processReviewReminder(task)
	case "practice_generation":
		err = ts.processPracticeGeneration(task)
	case "notification":
		err = ts.processNotification(task)
	default:
		err = fmt.Errorf("unknown task type: %s", task.TaskType)
	}

	// 更新任务状态
	if err != nil {
		task.Status = "failed"
		task.ErrorMessage = err.Error()
		task.RetryCount++
		if task.RetryCount < task.MaxRetries {
			// 计算下次重试时间（指数退避）
			nextAttempt := time.Now().Add(time.Duration(1<<task.RetryCount) * time.Minute)
			task.NextAttempt = &nextAttempt
			task.Status = "pending"
		}
	} else {
		task.Status = "completed"
	}

	taskStatus.Status = task.Status
	taskStatus.EndTime = &time.Time{}
	*taskStatus.EndTime = time.Now()
	taskStatus.Message = task.ErrorMessage

	ts.db.Save(task)
	ts.db.Save(taskStatus)
}

// checkPendingTasks 检查待执行任务
func (ts *TaskSchedulerService) checkPendingTasks() {
	var tasks []model.ScheduledTask
	now := time.Now()

	ts.db.Where("status = ? AND scheduled_at <= ?", "pending", now).Find(&tasks)

	for _, task := range tasks {
		select {
		case ts.taskQueue <- &task:
		default:
			log.Printf("Task queue full, skipping task %s", task.ID)
		}
	}
}

// loadExistingTasks 加载现有任务
func (ts *TaskSchedulerService) loadExistingTasks() {
	var tasks []model.ScheduledTask
	ts.db.Where("status IN (?)", []string{"pending", "processing"}).Find(&tasks)

	for _, task := range tasks {
		if task.ScheduledAt.Before(time.Now().Add(time.Minute)) {
			select {
			case ts.taskQueue <- &task:
			default:
				log.Printf("Task queue full during startup, task %s will be processed later", task.ID)
			}
		}
	}
}

// processReviewReminder 处理复习提醒任务
func (ts *TaskSchedulerService) processReviewReminder(task *model.ScheduledTask) error {
	var taskData map[string]interface{}
	if err := json.Unmarshal([]byte(task.TaskData), &taskData); err != nil {
		return err
	}

	planID := taskData["plan_id"].(string)
	planName := taskData["plan_name"].(string)

	// 创建通知
	notification := &model.Notification{
		StudentID:   task.StudentID,
		Title:       "复习提醒",
		Content:     fmt.Sprintf("您的复习计划「%s」该复习了！", planName),
		Type:        "review_reminder",
		Priority:    "normal",
		RelatedID:   planID,
		RelatedType: "review_plan",
		Channels:    `["app"]`,
	}

	return ts.db.Create(notification).Error
}

// processPracticeGeneration 处理练习题生成任务
func (ts *TaskSchedulerService) processPracticeGeneration(task *model.ScheduledTask) error {
	var taskData map[string]interface{}
	if err := json.Unmarshal([]byte(task.TaskData), &taskData); err != nil {
		return err
	}

	weakPointID := taskData["weak_point_id"].(string)
	studentID := taskData["student_id"].(string)
	courseID := taskData["course_id"].(string)
	difficulty := int(taskData["difficulty"].(float64))
	count := int(taskData["count"].(float64))

	// 这里应该调用AI引擎生成练习题
	// 暂时创建空的练习任务作为占位符
	practiceTask := &model.PracticeTask{
		TaskID:    fmt.Sprintf("task_%d", time.Now().Unix()),
		UserID:    studentID,
		CourseID:  courseID,
		NodeID:    weakPointID, // 使用weakPointID作为nodeID
		Difficulty: difficulty,
		Count:     count,
		Questions: "[]", // 暂时为空
	}

	if err := ts.db.Create(practiceTask).Error; err != nil {
		return err
	}

	// 创建通知
	notification := &model.Notification{
		StudentID:   studentID,
		Title:       "新练习题",
		Content:     fmt.Sprintf("为您生成了%d道练习题，快来练习吧！", count),
		Type:        "practice_due",
		Priority:    "normal",
		RelatedID:   practiceTask.ID,
		RelatedType: "practice_task",
		Channels:    `["app"]`,
	}

	return ts.db.Create(notification).Error
}

// processNotification 处理通知任务
func (ts *TaskSchedulerService) processNotification(task *model.ScheduledTask) error {
	var taskData map[string]interface{}
	if err := json.Unmarshal([]byte(task.TaskData), &taskData); err != nil {
		return err
	}

	notificationID := taskData["notification_id"].(string)
	var notification model.Notification
	if err := ts.db.First(&notification, "id = ?", notificationID).Error; err != nil {
		return err
	}

	// 标记为已发送
	now := time.Now()
	notification.SentAt = &now
	notification.Status = "read" // 假设发送成功就标记为已读

	return ts.db.Save(&notification).Error
}