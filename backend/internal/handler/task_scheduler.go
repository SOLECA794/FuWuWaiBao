package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/service"
)

// TaskSchedulerHandler 任务调度处理器
type TaskSchedulerHandler struct {
	taskSchedulerService *service.TaskSchedulerService
}

// NewTaskSchedulerHandler 创建任务调度处理器
func NewTaskSchedulerHandler(taskSchedulerService *service.TaskSchedulerService) *TaskSchedulerHandler {
	return &TaskSchedulerHandler{
		taskSchedulerService: taskSchedulerService,
	}
}

// CreateScheduledTask 创建定时任务
func (tsh *TaskSchedulerHandler) CreateScheduledTask(c *gin.Context) {
	var req struct {
		UserID      uint   `json:"userId" binding:"required"`
		TaskType    string `json:"taskType" binding:"required"`
		TaskData    string `json:"taskData" binding:"required"`
		CronExpr    string `json:"cronExpr" binding:"required"`
		Description string `json:"description"`
		MaxRetries  int    `json:"maxRetries"`
		Priority    int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	task := &model.ScheduledTask{
		UserID:      req.UserID,
		TaskType:    req.TaskType,
		TaskData:    req.TaskData,
		CronExpr:    req.CronExpr,
		Description: req.Description,
		Status:      "active",
		MaxRetries:  req.MaxRetries,
		Priority:    req.Priority,
	}

	err := tsh.taskSchedulerService.CreateScheduledTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建定时任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": task})
}

// GetScheduledTasks 获取定时任务列表
func (tsh *TaskSchedulerHandler) GetScheduledTasks(c *gin.Context) {
	userIDStr := c.Query("userId")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")
	status := c.Query("status")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 userId"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "userId 参数无效"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "page 参数无效"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pageSize 参数无效"})
		return
	}

	tasks, total, err := tsh.taskSchedulerService.GetScheduledTasks(uint(userID), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取定时任务列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  tasks,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}

// GetScheduledTask 获取单个定时任务详情
func (tsh *TaskSchedulerHandler) GetScheduledTask(c *gin.Context) {
	taskIDStr := c.Param("id")

	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "任务ID无效"})
		return
	}

	task, err := tsh.taskSchedulerService.GetScheduledTask(uint(taskID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "定时任务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": task})
}

// UpdateScheduledTask 更新定时任务
func (tsh *TaskSchedulerHandler) UpdateScheduledTask(c *gin.Context) {
	taskIDStr := c.Param("id")

	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "任务ID无效"})
		return
	}

	var req struct {
		CronExpr    string `json:"cronExpr"`
		Description string `json:"description"`
		Status      string `json:"status"`
		MaxRetries  int    `json:"maxRetries"`
		Priority    int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := make(map[string]interface{})
	if req.CronExpr != "" {
		updates["cron_expr"] = req.CronExpr
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.MaxRetries > 0 {
		updates["max_retries"] = req.MaxRetries
	}
	if req.Priority > 0 {
		updates["priority"] = req.Priority
	}

	err = tsh.taskSchedulerService.UpdateScheduledTask(uint(taskID), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新定时任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

// DeleteScheduledTask 删除定时任务
func (tsh *TaskSchedulerHandler) DeleteScheduledTask(c *gin.Context) {
	taskIDStr := c.Param("id")

	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "任务ID无效"})
		return
	}

	err = tsh.taskSchedulerService.DeleteScheduledTask(uint(taskID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除定时任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

// ExecuteTaskNow 立即执行任务
func (tsh *TaskSchedulerHandler) ExecuteTaskNow(c *gin.Context) {
	taskIDStr := c.Param("id")

	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "任务ID无效"})
		return
	}

	err = tsh.taskSchedulerService.ExecuteTaskNow(uint(taskID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "执行任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "任务已提交执行"})
}

// GetTaskStatuses 获取任务执行状态列表
func (tsh *TaskSchedulerHandler) GetTaskStatuses(c *gin.Context) {
	userIDStr := c.Query("userId")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")
	status := c.Query("status")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 userId"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "userId 参数无效"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "page 参数无效"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pageSize 参数无效"})
		return
	}

	statuses, total, err := tsh.taskSchedulerService.GetTaskStatuses(uint(userID), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取任务状态列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  statuses,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}

// GetTaskStatus 获取单个任务执行状态
func (tsh *TaskSchedulerHandler) GetTaskStatus(c *gin.Context) {
	statusIDStr := c.Param("id")

	statusID, err := strconv.ParseUint(statusIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "状态ID无效"})
		return
	}

	status, err := tsh.taskSchedulerService.GetTaskStatus(uint(statusID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "任务状态不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": status})
}