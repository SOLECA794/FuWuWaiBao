package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/service"
)

func (h *CompatibilityHandler) GetStudentNotesV1(c *gin.Context) {
	studentID := strings.TrimSpace(c.Query("studentId"))
	if studentID == "" {
		studentID = strings.TrimSpace(c.GetHeader("X-Student-Id"))
	}
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 studentId"})
		return
	}
	courseID := strings.TrimSpace(c.Query("courseId"))
	pageNum, _ := strconv.Atoi(strings.TrimSpace(c.Query("pageNum")))
	page, _ := strconv.Atoi(strings.TrimSpace(c.Query("page")))
	pageSize, _ := strconv.Atoi(strings.TrimSpace(c.Query("pageSize")))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := h.db.Model(&model.StudentNote{}).Where("user_id = ?", studentID)
	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}
	if pageNum > 0 {
		query = query.Where("page_num = ?", pageNum)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询笔记失败"})
		return
	}
	var notes []model.StudentNote
	if err := query.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询笔记失败"})
		return
	}

	items := make([]gin.H, 0, len(notes))
	for _, note := range notes {
		items = append(items, gin.H{
			"noteId":    note.ID,
			"studentId": note.UserID,
			"courseId":  note.CourseID,
			"nodeId":    note.NodeID,
			"pageNum":   note.PageNum,
			"content":   note.Note,
			"createdAt": note.CreatedAt,
			"updatedAt": note.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data": gin.H{
			"items":      items,
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": int((total + int64(pageSize) - 1) / int64(pageSize)),
		},
	})
}

func (h *CompatibilityHandler) CreateFavoriteV1(c *gin.Context) {
	var req struct {
		StudentID string `json:"studentId" binding:"required"`
		CourseID  string `json:"courseId" binding:"required"`
		NodeID    string `json:"nodeId"`
		PageNum   int    `json:"pageNum"`
		Title     string `json:"title"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if strings.TrimSpace(req.NodeID) == "" {
		req.NodeID = fmt.Sprintf("p%d_n1", req.PageNum)
	}

	var existing model.NodeFavorite
	if err := h.db.Where("user_id = ? AND course_id = ? AND node_id = ? AND page_index = ?", req.StudentID, req.CourseID, req.NodeID, req.PageNum).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已收藏", "data": gin.H{"favoriteId": existing.ID}})
		return
	}

	record := model.NodeFavorite{
		UserID:    req.StudentID,
		CourseID:  req.CourseID,
		NodeID:    req.NodeID,
		PageIndex: req.PageNum,
		Title:     strings.TrimSpace(req.Title),
	}
	if err := h.db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "收藏失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "收藏成功", "data": gin.H{"favoriteId": record.ID}})
}

func (h *CompatibilityHandler) ListFavoritesV1(c *gin.Context) {
	studentID := strings.TrimSpace(c.Query("studentId"))
	if studentID == "" {
		studentID = strings.TrimSpace(c.GetHeader("X-Student-Id"))
	}
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 studentId"})
		return
	}
	courseID := strings.TrimSpace(c.Query("courseId"))
	page, _ := strconv.Atoi(strings.TrimSpace(c.Query("page")))
	pageSize, _ := strconv.Atoi(strings.TrimSpace(c.Query("pageSize")))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	query := h.db.Model(&model.NodeFavorite{}).Where("user_id = ?", studentID)
	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询收藏失败"})
		return
	}
	var records []model.NodeFavorite
	if err := query.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询收藏失败"})
		return
	}
	items := make([]gin.H, 0, len(records))
	for _, item := range records {
		items = append(items, gin.H{
			"favoriteId": item.ID,
			"studentId":  item.UserID,
			"courseId":   item.CourseID,
			"nodeId":     item.NodeID,
			"pageNum":    item.PageIndex,
			"title":      item.Title,
			"createdAt":  item.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data": gin.H{
			"items":      items,
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": int((total + int64(pageSize) - 1) / int64(pageSize)),
		},
	})
}

func (h *CompatibilityHandler) DeleteFavoriteV1(c *gin.Context) {
	favoriteID := c.Param("favoriteId")
	if strings.TrimSpace(favoriteID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 favoriteId"})
		return
	}
	studentID := strings.TrimSpace(c.Query("studentId"))
	if studentID == "" {
		studentID = strings.TrimSpace(c.GetHeader("X-Student-Id"))
	}
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 studentId"})
		return
	}

	result := h.db.Delete(&model.NodeFavorite{}, "id = ? AND user_id = ?", favoriteID, studentID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除收藏失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "收藏不存在或无权限"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *CompatibilityHandler) ExplainNodeV1(c *gin.Context) {
	nodeID := strings.TrimSpace(c.Param("nodeId"))
	var req struct {
		CourseID string `json:"courseId" binding:"required"`
		PageNum  int    `json:"pageNum"`
		Question string `json:"question"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if req.PageNum <= 0 {
		req.PageNum = parsePageFromNodeID(nodeID)
	}
	if req.PageNum <= 0 {
		req.PageNum = 1
	}

	question := strings.TrimSpace(req.Question)
	if question == "" {
		question = fmt.Sprintf("请对知识点 %s 进行专项讲解，先用一句话总结，再给出一个引导问题。", nodeID)
	}
	contextText := buildPageContextFromTeachingNodes(loadTeachingNodesByPage(h.db, req.CourseID, req.PageNum))
	if strings.TrimSpace(contextText) == "" {
		contextText = fmt.Sprintf("课程 %s 第 %d 页，节点 %s", req.CourseID, req.PageNum, nodeID)
	}

	resp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{
		Question:    question,
		CurrentPage: req.PageNum,
		Context:     contextText,
		Mode:        "llm",
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"code": 503, "message": "专项讲解暂不可用"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{
		"nodeId":    nodeID,
		"courseId":  req.CourseID,
		"pageNum":   req.PageNum,
		"content":   resp.Answer,
		"followUp":  resp.FollowUpSuggestion,
		"sourcePage": max(req.PageNum, resp.SourcePage),
	}})
}

func (h *CompatibilityHandler) GeneratePracticeV1(c *gin.Context) {
	var req struct {
		StudentID  string `json:"studentId" binding:"required"`
		CourseID   string `json:"courseId" binding:"required"`
		NodeID     string `json:"nodeId"`
		PageNum    int    `json:"pageNum"`
		Difficulty int    `json:"difficulty"`
		Count      int    `json:"count"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if req.Count <= 0 || req.Count > 10 {
		req.Count = 3
	}
	if req.Difficulty <= 0 {
		req.Difficulty = 2
	}
	if req.PageNum <= 0 {
		req.PageNum = parsePageFromNodeID(req.NodeID)
	}
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if strings.TrimSpace(req.NodeID) == "" {
		req.NodeID = fmt.Sprintf("p%d_n1", req.PageNum)
	}

	questionIDs := make([]string, 0, req.Count)
	items := make([]gin.H, 0, req.Count)
	questionPayloads := h.generatePracticePayloadsWithFallback(c, req.CourseID, req.NodeID, req.PageNum, req.Count)
	for idx, payload := range questionPayloads {
		optionsJSON, _ := json.Marshal(payload.Options)
		question := model.Question{
			WeakPointID:  req.NodeID,
			QuestionType: "single",
			Content:      payload.Content,
			Options:      string(optionsJSON),
			Answer:       payload.Answer,
			Explanation:  payload.Explanation,
			Difficulty:   req.Difficulty,
		}
		if err := h.db.Create(&question).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成练习题失败"})
			return
		}
		questionIDs = append(questionIDs, question.ID)
		items = append(items, gin.H{
			"questionId": question.ID,
			"nodeId":     req.NodeID,
			"type":       question.QuestionType,
			"content":    question.Content,
			"options":    payload.Options,
			"fromAI":     idx < len(questionPayloads),
		})
	}

	questionIDsJSON, _ := json.Marshal(questionIDs)
	task := model.PracticeTask{
		UserID:      req.StudentID,
		CourseID:    req.CourseID,
		NodeID:      req.NodeID,
		PageIndex:   req.PageNum,
		Difficulty:  req.Difficulty,
		QuestionIDs: string(questionIDsJSON),
		Status:      "pending",
	}
	if err := h.db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建练习任务失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "生成成功", "data": gin.H{
		"taskId":      task.ID,
		"courseId":    req.CourseID,
		"nodeId":      req.NodeID,
		"pageNum":     req.PageNum,
		"difficulty":  req.Difficulty,
		"questionCount": len(items),
		"questions":   items,
	}})
}

func (h *CompatibilityHandler) SubmitPracticeV1(c *gin.Context) {
	var req struct {
		TaskID    string `json:"taskId" binding:"required"`
		StudentID string `json:"studentId" binding:"required"`
		Answers   []struct {
			QuestionID string `json:"questionId"`
			UserAnswer string `json:"userAnswer"`
		} `json:"answers" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var task model.PracticeTask
	if err := h.db.First(&task, "id = ?", req.TaskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "练习任务不存在"})
		return
	}

	// 幂等：同一 student + task 已提交过则直接返回历史结果
	var existing model.PracticeAttempt
	if err := h.db.Where("task_id = ? AND user_id = ?", task.ID, req.StudentID).Order("created_at desc").First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已提交（幂等返回）", "data": gin.H{
			"taskId":       task.ID,
			"attemptId":    existing.ID,
			"courseId":     existing.CourseID,
			"nodeId":       existing.NodeID,
			"totalCount":   existing.TotalCount,
			"correctCount": existing.CorrectCount,
			"score":        existing.Score,
			"idempotent":   true,
		}})
		return
	}

	total := len(req.Answers)
	correct := 0
	details := make([]gin.H, 0, total)
	for _, answer := range req.Answers {
		var question model.Question
		if err := h.db.First(&question, "id = ?", answer.QuestionID).Error; err != nil {
			details = append(details, gin.H{"questionId": answer.QuestionID, "error": "题目不存在"})
			continue
		}
		ok := strings.EqualFold(strings.TrimSpace(answer.UserAnswer), strings.TrimSpace(question.Answer))
		if !ok && len(strings.TrimSpace(question.Answer)) == 1 {
			ok = strings.HasPrefix(strings.ToUpper(strings.TrimSpace(answer.UserAnswer)), strings.ToUpper(strings.TrimSpace(question.Answer)))
		}
		if ok {
			correct++
		}
		masteryDelta := -5
		if ok {
			masteryDelta = 10
		}
		_ = h.db.Create(&model.AnswerRecord{
			StudentID:    req.StudentID,
			QuestionID:   question.ID,
			UserAnswer:   answer.UserAnswer,
			IsCorrect:    ok,
			MasteryDelta: masteryDelta,
		}).Error
		details = append(details, gin.H{
			"questionId":    question.ID,
			"isCorrect":     ok,
			"correctAnswer": question.Answer,
			"explanation":   question.Explanation,
		})
	}

	score := 0
	if total > 0 {
		score = (correct * 100) / total
	}
	answersJSON, _ := json.Marshal(req.Answers)
	attempt := model.PracticeAttempt{
		TaskID:       task.ID,
		UserID:       req.StudentID,
		CourseID:     task.CourseID,
		NodeID:       task.NodeID,
		TotalCount:   total,
		CorrectCount: correct,
		Score:        score,
		AnswersJSON:  string(answersJSON),
	}
	_ = h.db.Create(&attempt).Error
	_ = h.db.Model(&task).Update("status", "completed").Error

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "提交成功", "data": gin.H{
		"taskId":       task.ID,
		"attemptId":    attempt.ID,
		"courseId":     task.CourseID,
		"nodeId":       task.NodeID,
		"totalCount":   total,
		"correctCount": correct,
		"score":        score,
		"details":      details,
	}})
}

func (h *CompatibilityHandler) generatePracticePayloadsWithFallback(c *gin.Context, courseID, nodeID string, pageNum, count int) []weakPointTestPayload {
	result := make([]weakPointTestPayload, 0, count)
	contextText := buildPageContextFromTeachingNodes(loadTeachingNodesByPage(h.db, courseID, pageNum))
	if strings.TrimSpace(contextText) == "" {
		contextText = fmt.Sprintf("课程 %s 第 %d 页 节点 %s", courseID, pageNum, nodeID)
	}

	if h.aiClient != nil {
		for i := 0; i < count; i++ {
			prompt := fmt.Sprintf("请围绕知识点 %s 生成 1 道单选题。严格返回 JSON：{\"content\":\"题干\",\"options\":[\"A.xxx\",\"B.xxx\",\"C.xxx\",\"D.xxx\"],\"answer\":\"A\",\"explanation\":\"解析\"}", nodeID)
			resp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{
				Question:    prompt,
				CurrentPage: pageNum,
				Context:     contextText,
				Mode:        "llm",
			})
			if err != nil {
				break
			}
			parsed, err := parseWeakPointTest(resp.Answer)
			if err != nil || parsed == nil || strings.TrimSpace(parsed.Content) == "" || len(parsed.Options) < 2 {
				break
			}
			if strings.TrimSpace(parsed.Answer) == "" {
				parsed.Answer = "A"
			}
			result = append(result, *parsed)
		}
	}

	// 兜底：不足数量用模板题补齐，保证接口稳定可用
	for i := len(result); i < count; i++ {
		result = append(result, weakPointTestPayload{
			Content:     fmt.Sprintf("【%s】练习题 %d：以下哪项最符合本节点核心内容？", nodeID, i+1),
			Options:     []string{"A. 核心概念的定义与应用", "B. 与主题无关的随机信息", "C. 完全错误的结论", "D. 无法判断"},
			Answer:      "A",
			Explanation: "本题用于巩固当前节点的核心概念，正确项为 A。",
		})
	}
	return result
}

func (h *CompatibilityHandler) GetNodeInsightsV1(c *gin.Context) {
	courseID := c.Param("courseId")
	if strings.TrimSpace(courseID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 courseId"})
		return
	}

	type pageStat struct {
		PageIndex int
		Count     int64
	}
	type nodePracticeStat struct {
		NodeID       string
		AttemptCount int64
		TotalCount   int64
		CorrectCount int64
	}
	type nodeReteachStat struct {
		NodeID       string
		ReteachCount int64
		TotalCount   int64
	}
	type trendRow struct {
		Day       string
		PageIndex int
		Count     int64
	}
	type nodeTrendRow struct {
		Day   string
		NodeID string
		Count int64
	}

	var qStats []pageStat
	var nStats []pageStat
	h.db.Model(&model.QuestionLog{}).
		Select("page_index, count(*) as count").
		Where("course_id = ?", courseID).
		Group("page_index").
		Scan(&qStats)
	h.db.Model(&model.StudentNote{}).
		Select("page_num as page_index, count(*) as count").
		Where("course_id = ?", courseID).
		Group("page_num").
		Scan(&nStats)

	var pages []model.CoursePage
	_ = h.db.Where("course_id = ?", courseID).Order("page_index asc").Find(&pages).Error

	// page -> primary node
	nodeByPage := map[int]string{}
	titleByNode := map[string]string{}
	for _, page := range pages {
		nodeID := fmt.Sprintf("p%d_n1", page.PageIndex)
		nodeTitle := fmt.Sprintf("第%d页", page.PageIndex)
		nodes := loadTeachingNodesByPage(h.db, courseID, page.PageIndex)
		if len(nodes) > 0 && strings.TrimSpace(nodes[0].NodeID) != "" {
			nodeID = nodes[0].NodeID
			if strings.TrimSpace(nodes[0].Title) != "" {
				nodeTitle = nodes[0].Title
			}
		}
		nodeByPage[page.PageIndex] = nodeID
		titleByNode[nodeID] = nodeTitle
	}

	questionByNode := map[string]int64{}
	noteByNode := map[string]int64{}
	for _, item := range qStats {
		nodeID := nodeByPage[item.PageIndex]
		if nodeID == "" {
			nodeID = fmt.Sprintf("p%d_n1", item.PageIndex)
		}
		questionByNode[nodeID] += item.Count
	}
	for _, item := range nStats {
		nodeID := nodeByPage[item.PageIndex]
		if nodeID == "" {
			nodeID = fmt.Sprintf("p%d_n1", item.PageIndex)
		}
		noteByNode[nodeID] += item.Count
	}

	// 练习正确率
	var pStats []nodePracticeStat
	h.db.Model(&model.PracticeAttempt{}).
		Select("node_id, count(*) as attempt_count, sum(total_count) as total_count, sum(correct_count) as correct_count").
		Where("course_id = ?", courseID).
		Group("node_id").
		Scan(&pStats)
	practiceByNode := map[string]nodePracticeStat{}
	for _, stat := range pStats {
		practiceByNode[strings.TrimSpace(stat.NodeID)] = stat
	}

	// 重讲率（need_reteach）
	var rStats []nodeReteachStat
	h.db.Model(&model.DialogueTurn{}).
		Select("node_id, sum(case when need_reteach = true then 1 else 0 end) as reteach_count, count(*) as total_count").
		Where("course_id = ?", courseID).
		Group("node_id").
		Scan(&rStats)
	reteachByNode := map[string]nodeReteachStat{}
	for _, stat := range rStats {
		reteachByNode[strings.TrimSpace(stat.NodeID)] = stat
	}

	// 近7天趋势（问题数 + 练习尝试数）
	questionTrendRaw := make([]trendRow, 0)
	h.db.Model(&model.QuestionLog{}).
		Select("to_char(created_at, 'YYYY-MM-DD') as day, page_index, count(*) as count").
		Where("course_id = ? AND created_at >= ?", courseID, time.Now().AddDate(0, 0, -6)).
		Group("day, page_index").
		Scan(&questionTrendRaw)
	questionTrendByNode := map[string]map[string]int64{}
	for _, row := range questionTrendRaw {
		nodeID := nodeByPage[row.PageIndex]
		if nodeID == "" {
			nodeID = fmt.Sprintf("p%d_n1", row.PageIndex)
		}
		if _, ok := questionTrendByNode[nodeID]; !ok {
			questionTrendByNode[nodeID] = map[string]int64{}
		}
		questionTrendByNode[nodeID][row.Day] += row.Count
	}

	practiceTrendRaw := make([]nodeTrendRow, 0)
	h.db.Model(&model.PracticeAttempt{}).
		Select("to_char(created_at, 'YYYY-MM-DD') as day, node_id, count(*) as count").
		Where("course_id = ? AND created_at >= ?", courseID, time.Now().AddDate(0, 0, -6)).
		Group("day, node_id").
		Scan(&practiceTrendRaw)
	practiceTrendByNode := map[string]map[string]int64{}
	for _, row := range practiceTrendRaw {
		nodeID := strings.TrimSpace(row.NodeID)
		if nodeID == "" {
			continue
		}
		if _, ok := practiceTrendByNode[nodeID]; !ok {
			practiceTrendByNode[nodeID] = map[string]int64{}
		}
		practiceTrendByNode[nodeID][row.Day] += row.Count
	}

	// 汇总并按 insightScore 降序
	items := make([]gin.H, 0, len(nodeByPage))
	for _, page := range pages {
		nodeID := nodeByPage[page.PageIndex]
		if nodeID == "" {
			nodeID = fmt.Sprintf("p%d_n1", page.PageIndex)
		}
		nodeTitle := titleByNode[nodeID]
		questionCount := questionByNode[nodeID]
		noteCount := noteByNode[nodeID]

		pStat := practiceByNode[nodeID]
		practiceAccuracy := 0.0
		if pStat.TotalCount > 0 {
			practiceAccuracy = float64(pStat.CorrectCount) * 100 / float64(pStat.TotalCount)
		}
		rStat := reteachByNode[nodeID]
		reteachRate := 0.0
		if rStat.TotalCount > 0 {
			reteachRate = float64(rStat.ReteachCount) * 100 / float64(rStat.TotalCount)
		}

		trend := make([]gin.H, 0, 7)
		for i := 6; i >= 0; i-- {
			day := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
			trend = append(trend, gin.H{
				"date":          day,
				"questionCount": questionTrendByNode[nodeID][day],
				"practiceCount": practiceTrendByNode[nodeID][day],
			})
		}

		insightScore := int(questionCount*2 + noteCount + int64(rStat.ReteachCount)*2)
		if practiceAccuracy > 0 {
			insightScore += int((100 - practiceAccuracy) / 10)
		}

		items = append(items, gin.H{
			"pageNum":            page.PageIndex,
			"nodeId":             nodeID,
			"nodeTitle":          nodeTitle,
			"questionCount":      questionCount,
			"noteCount":          noteCount,
			"practiceAttemptCount": pStat.AttemptCount,
			"practiceAccuracy":   fmt.Sprintf("%.2f", practiceAccuracy),
			"reteachCount":       rStat.ReteachCount,
			"reteachRate":        fmt.Sprintf("%.2f", reteachRate),
			"trend7d":            trend,
			"insightScore":       insightScore,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		li, _ := items[i]["insightScore"].(int)
		lj, _ := items[j]["insightScore"].(int)
		return li > lj
	})

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{
		"courseId": courseID,
		"items":    items,
		"total":    len(items),
	}})
}
