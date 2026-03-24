package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/service"
)

func (h *CompatibilityHandler) GetStudentCoursewareList(c *gin.Context) {
	var courses []model.Course
	query := h.db.Order("created_at desc").Where("is_published = ?", true)
	if err := query.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取课件列表失败"})
		return
	}
	if len(courses) == 0 {
		_ = h.db.Order("created_at desc").Find(&courses).Error
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": courses})
}

func (h *CompatibilityHandler) StartStudentSession(c *gin.Context) {
	var req struct {
		UserID   string `json:"userId" binding:"required"`
		CourseID string `json:"courseId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	state := sessionState{
		SessionID:     "sess_" + uuid.NewString(),
		UserID:        req.UserID,
		CourseID:      req.CourseID,
		CurrentPage:   1,
		CurrentNodeID: "p1_n1",
		UpdatedAt:     time.Now(),
	}
	h.persistSession(state)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"sessionId": state.SessionID, "courseId": req.CourseID}})
}

func (h *CompatibilityHandler) UpdateStudentProgress(c *gin.Context) {
	var req struct {
		SessionID     string `json:"sessionId"`
		UserID        string `json:"userId"`
		CourseID      string `json:"courseId" binding:"required"`
		Page          int    `json:"page"`
		CurrentPage   int    `json:"currentPage"`
		NodeID        string `json:"nodeId"`
		CurrentNodeID string `json:"currentNodeId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	page := req.Page
	if page <= 0 {
		page = req.CurrentPage
	}
	if page <= 0 {
		page = 1
	}
	nodeID := strings.TrimSpace(req.NodeID)
	if nodeID == "" {
		nodeID = strings.TrimSpace(req.CurrentNodeID)
	}
	if nodeID == "" {
		nodeID = fmt.Sprintf("p%d_n1", page)
	}
	state := sessionState{
		SessionID:     req.SessionID,
		UserID:        req.UserID,
		CourseID:      req.CourseID,
		CurrentPage:   page,
		CurrentNodeID: nodeID,
		UpdatedAt:     time.Now(),
	}
	if state.SessionID == "" {
		state.SessionID = "sess_" + uuid.NewString()
	}
	h.persistSession(state)

	if strings.TrimSpace(req.UserID) != "" {
		var progress model.UserProgress
		err := h.db.Where("user_id = ? AND course_id = ?", req.UserID, req.CourseID).First(&progress).Error
		if err == nil {
			_ = h.db.Model(&progress).Updates(map[string]any{"last_page": page}).Error
		} else {
			_ = h.db.Create(&model.UserProgress{UserID: req.UserID, CourseID: req.CourseID, LastPage: page}).Error
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok", "data": gin.H{"sessionId": state.SessionID, "page": page, "nodeId": nodeID}})
}

func (h *CompatibilityHandler) GetStudentScript(c *gin.Context) {
	courseID := c.Param("courseId")
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "页码错误"})
		return
	}
	var coursePage model.CoursePage
	_ = h.db.Where("course_id = ? AND page_index = ?", courseID, page).First(&coursePage).Error
	nodes := buildScriptNodes(page, coursePage.ScriptText)
	pageSummary := strings.TrimSpace(coursePage.ScriptText)
	if len(pageSummary) > 80 {
		pageSummary = pageSummary[:80]
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"courseId": courseID, "page": page, "nodes": nodes, "page_summary": pageSummary}})
}

func (h *CompatibilityHandler) StreamStudentQA(c *gin.Context) {
	var req struct {
		SessionID string `json:"sessionId"`
		CourseID  string `json:"courseId" binding:"required"`
		Page      int    `json:"page" binding:"required"`
		NodeID    string `json:"nodeId"`
		Question  string `json:"question" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	var coursePage model.CoursePage
	contextText := ""
	if err := h.db.Where("course_id = ? AND page_index = ?", req.CourseID, req.Page).First(&coursePage).Error; err == nil {
		contextText = coursePage.ScriptText
	}
	aiResp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{Question: req.Question, CurrentPage: req.Page, Context: contextText, Mode: "llm"})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"code": 503, "message": "AI服务暂不可用"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "当前连接不支持流式输出"})
		return
	}

	writeSSE := func(event string, payload any) {
		body, _ := json.Marshal(payload)
		_, _ = fmt.Fprintf(c.Writer, "event: %s\n", event)
		_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", string(body))
		flusher.Flush()
	}

	parts := strings.Fields(aiResp.Answer)
	if len(parts) == 0 {
		for _, r := range []rune(aiResp.Answer) {
			writeSSE("token", gin.H{"text": string(r)})
		}
	} else {
		for _, part := range parts {
			writeSSE("token", gin.H{"text": part + " "})
		}
	}
	writeSSE("sentence", gin.H{"text": aiResp.Answer})
	writeSSE("final", gin.H{
		"need_reteach":   aiResp.Intent.NeedReteach,
		"source_page":    aiResp.SourcePage,
		"resume_page":    aiResp.ResumePage,
		"resume_node_id": nextNodeID(req.NodeID, req.Page),
	})
}

func weakPointVirtualID(page int) string {
	return fmt.Sprintf("wp-page-%d", page)
}

func parseWeakPointVirtualID(id string) int {
	if strings.HasPrefix(id, "wp-page-") {
		page, _ := strconv.Atoi(strings.TrimPrefix(id, "wp-page-"))
		return page
	}
	return 0
}

func (h *CompatibilityHandler) listWeakPoints(studentID, courseID string) []gin.H {
	type pageStat struct {
		PageIndex int
		Count     int
	}
	var stats []pageStat
	h.db.Table("question_logs").Select("page_index, count(*) as count").Where("user_id = ? AND course_id = ?", studentID, courseID).Group("page_index").Order("count desc").Limit(5).Scan(&stats)
	result := make([]gin.H, 0, len(stats))
	for _, stat := range stats {
		mastery := 100 - stat.Count*12
		if mastery < 0 {
			mastery = 0
		}
		result = append(result, gin.H{"weakPointId": weakPointVirtualID(stat.PageIndex), "name": fmt.Sprintf("第%d页知识点", stat.PageIndex), "count": stat.Count, "mastery": mastery, "pageIndex": stat.PageIndex})
	}
	return result
}

func (h *CompatibilityHandler) GetWeakPointsV1(c *gin.Context) {
	courseID := c.Param("courseId")
	studentID := c.Query("studentId")
	if studentID == "" {
		studentID = c.GetHeader("X-Student-Id")
	}
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 studentId"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": h.listWeakPoints(studentID, courseID)})
}

func (h *CompatibilityHandler) ExplainWeakPointV1(c *gin.Context) {
	weakPointID := c.Param("weakPointId")
	page := parseWeakPointVirtualID(weakPointID)
	name := c.Query("name")
	if name == "" {
		if page > 0 {
			name = fmt.Sprintf("第%d页知识点", page)
		} else {
			name = weakPointID
		}
	}
	resp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{Question: "请用教学口吻讲解知识点：" + name, CurrentPage: max(page, 1), Context: name, Mode: "llm"})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"code": 503, "message": "AI讲解暂不可用"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"weakPointId": weakPointID, "title": name, "content": resp.Answer}})
}

func (h *CompatibilityHandler) GenerateWeakPointTestV1(c *gin.Context) {
	weakPointID := c.Param("weakPointId")
	page := parseWeakPointVirtualID(weakPointID)
	var req struct {
		StudentID     string `json:"studentId"`
		QuestionType  string `json:"questionType"`
		WeakPointName string `json:"weakPointName"`
	}
	_ = c.ShouldBindJSON(&req)
	name := strings.TrimSpace(req.WeakPointName)
	if name == "" {
		if page > 0 {
			name = fmt.Sprintf("第%d页知识点", page)
		} else {
			name = weakPointID
		}
	}
	prompt := fmt.Sprintf("请为薄弱点“%s”生成一道单选题，返回严格JSON：{\"content\":\"题干\",\"options\":[\"A.xxx\",\"B.xxx\",\"C.xxx\",\"D.xxx\"],\"answer\":\"A\",\"explanation\":\"解析\"}", name)
	resp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{Question: prompt, CurrentPage: max(page, 1), Context: name, Mode: "llm"})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"code": 503, "message": "AI习题生成失败"})
		return
	}
	testData, err := parseWeakPointTest(resp.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "AI习题返回格式异常"})
		return
	}
	optionsJSON, _ := json.Marshal(testData.Options)
	questionRecord := model.Question{WeakPointID: weakPointID, QuestionType: defaultString(req.QuestionType, "single"), Content: testData.Content, Options: string(optionsJSON), Answer: testData.Answer, Explanation: testData.Explanation, Difficulty: 2}
	if err := h.db.Create(&questionRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存题目失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"questionId": questionRecord.ID, "content": testData.Content, "options": testData.Options, "type": questionRecord.QuestionType}})
}

func (h *CompatibilityHandler) CheckAnswerV1(c *gin.Context) {
	questionID := c.Param("questionId")
	var req struct {
		StudentID  string `json:"studentId" binding:"required"`
		UserAnswer string `json:"userAnswer" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	var question model.Question
	if err := h.db.First(&question, "id = ?", questionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "题目不存在"})
		return
	}
	correctAnswer := strings.TrimSpace(question.Answer)
	userAnswer := strings.TrimSpace(req.UserAnswer)
	isCorrect := strings.EqualFold(userAnswer, correctAnswer)
	if !isCorrect && len(correctAnswer) == 1 {
		isCorrect = strings.HasPrefix(strings.ToUpper(userAnswer), strings.ToUpper(correctAnswer))
	}
	masteryDelta := -5
	if isCorrect {
		masteryDelta = 10
	}
	_ = h.db.Create(&model.AnswerRecord{StudentID: req.StudentID, QuestionID: questionID, UserAnswer: req.UserAnswer, IsCorrect: isCorrect, MasteryDelta: masteryDelta}).Error
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"isCorrect": isCorrect, "correctAnswer": correctAnswer, "explanation": question.Explanation, "masteryDelta": masteryDelta}})
}

func (h *CompatibilityHandler) GetBreakpointV1(c *gin.Context) {
	courseID := c.Param("courseId")
	studentID := c.Query("studentId")
	if studentID == "" {
		studentID = c.GetHeader("X-Student-Id")
	}
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 studentId"})
		return
	}
	var progress model.UserProgress
	lastPage := 1
	if err := h.db.Where("user_id = ? AND course_id = ?", studentID, courseID).First(&progress).Error; err == nil && progress.LastPage > 0 {
		lastPage = progress.LastPage
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"courseId": courseID, "pageNum": lastPage}})
}

func (h *CompatibilityHandler) UpdateBreakpointV1(c *gin.Context) {
	courseID := c.Param("courseId")
	var req struct {
		StudentID string `json:"studentId" binding:"required"`
		PageNum   int    `json:"pageNum"`
		LastPage  int    `json:"lastPage"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	page := req.PageNum
	if page <= 0 {
		page = req.LastPage
	}
	if page <= 0 {
		page = 1
	}
	var progress model.UserProgress
	if err := h.db.Where("user_id = ? AND course_id = ?", req.StudentID, courseID).First(&progress).Error; err == nil {
		_ = h.db.Model(&progress).Update("last_page", page).Error
	} else {
		_ = h.db.Create(&model.UserProgress{UserID: req.StudentID, CourseID: courseID, LastPage: page}).Error
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok", "data": gin.H{"pageNum": page}})
}

func (h *CompatibilityHandler) SaveNoteV1(c *gin.Context) {
	courseID := c.Param("courseId")
	var req struct {
		StudentID string `json:"studentId" binding:"required"`
		PageNum   int    `json:"pageNum" binding:"required"`
		Content   string `json:"content" binding:"required"`
		X         int    `json:"x"`
		Y         int    `json:"y"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	note := model.StudentNote{UserID: req.StudentID, CourseID: courseID, PageNum: req.PageNum, Note: req.Content}
	if err := h.db.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "保存成功", "data": gin.H{"noteId": note.ID}})
}

func (h *CompatibilityHandler) GetStudentNotesV1(c *gin.Context) {
	studentID := strings.TrimSpace(c.Query("studentId"))
	courseID := strings.TrimSpace(c.Query("courseId"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 studentId"})
		return
	}
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	query := h.db.Model(&model.StudentNote{}).Where("user_id = ?", studentID)
	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}
	var total int64
	query.Count(&total)
	notes := []model.StudentNote{}
	query.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&notes)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"items": notes, "total": total, "page": pageNum, "pageSize": pageSize, "totalPages": (total+int64(pageSize)-1)/int64(pageSize)}})
}

func (h *CompatibilityHandler) AddFavoriteV1(c *gin.Context) {
	var req struct {
		StudentID string   `json:"studentId" binding:"required"`
		CourseID  string   `json:"courseId"`
		NodeID    string   `json:"nodeId"`
		PageNum   int      `json:"pageNum"`
		Title     string   `json:"title"`
		Tags      []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if req.StudentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 studentId"})
		return
	}
	tagsJSON, _ := json.Marshal(req.Tags)
	favorite := model.StudentFavorite{UserID: req.StudentID, CourseID: req.CourseID, NodeID: req.NodeID, PageNum: req.PageNum, Title: req.Title, Tags: string(tagsJSON)}
	existing := model.StudentFavorite{}
	err := h.db.Where("user_id = ? AND course_id = ? AND node_id = ?", req.StudentID, req.CourseID, req.NodeID).First(&existing).Error
	if err == nil {
		// Update tags if provided
		if len(req.Tags) > 0 {
			existing.Tags = string(tagsJSON)
			h.db.Save(&existing)
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"favoriteId": existing.ID, "message": "已收藏"}})
		return
	}
	if err := h.db.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "收藏失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"favoriteId": favorite.ID}})
}

func (h *CompatibilityHandler) GetFavoritesV1(c *gin.Context) {
	studentID := strings.TrimSpace(c.Query("studentId"))
	courseID := strings.TrimSpace(c.Query("courseId"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 studentId"})
		return
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	query := h.db.Model(&model.StudentFavorite{}).Where("user_id = ?", studentID)
	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}
	var total int64
	query.Count(&total)
	items := []model.StudentFavorite{}
	query.Order("created_at desc").Offset((page-1)*pageSize).Limit(pageSize).Find(&items)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"items": items, "total": total, "page": page, "pageSize": pageSize, "totalPages": (total+int64(pageSize)-1)/int64(pageSize)}})
}

func (h *CompatibilityHandler) DeleteFavoriteV1(c *gin.Context) {
	favoriteID := c.Param("favoriteId")
	studentID := strings.TrimSpace(c.Query("studentId"))
	if favoriteID == "" || studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少参数"})
		return
	}
	fav := model.StudentFavorite{}
	if err := h.db.First(&fav, "id = ?", favoriteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "收藏不存在"})
		return
	}
	if fav.UserID != studentID {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "收藏不存在"})
		return
	}
	if err := h.db.Delete(&fav).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
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
	if req.Count <= 0 {
		req.Count = 3
	}
	if req.Difficulty <= 0 {
		req.Difficulty = 2
	}
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	questions := make([]gin.H, 0, req.Count)
	for i := 1; i <= req.Count; i++ {
		q := gin.H{"questionId": fmt.Sprintf("q-%d-%d", time.Now().UnixNano(), i), "content": fmt.Sprintf("请解释第%d页第%d个知识点的含义。", req.PageNum, i), "options": []string{"A. 示例A", "B. 示例B", "C. 示例C", "D. 示例D"}, "answer": "A", "difficulty": req.Difficulty}
		questions = append(questions, q)
	}
	questionJSON, _ := json.Marshal(questions)
	task := model.PracticeTask{TaskID: "task_" + uuid.NewString(), UserID: req.StudentID, CourseID: req.CourseID, NodeID: req.NodeID, PageNum: req.PageNum, Difficulty: req.Difficulty, Count: req.Count, Questions: string(questionJSON)}
	if err := h.db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成练习失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"taskId": task.TaskID, "questions": questions, "questionCount": len(questions)}})
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
	if err := h.db.Where("task_id = ? AND user_id = ?", req.TaskID, req.StudentID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "练习任务不存在"})
		return
	}
	existing := model.PracticeAttempt{}
	if err := h.db.Where("task_id = ? AND user_id = ?", req.TaskID, req.StudentID).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"attemptId": existing.ID, "score": existing.Score, "correctCount": existing.Correct, "totalCount": existing.Total, "details": existing.Details}})
		return
	}
	details := make([]gin.H, 0, len(req.Answers))
	correctCount := 0
	for _, answer := range req.Answers {
		isCorrect := strings.EqualFold(strings.TrimSpace(answer.UserAnswer), "A")
		if isCorrect {
			correctCount++
		}
		details = append(details, gin.H{"questionId": answer.QuestionID, "userAnswer": answer.UserAnswer, "correct": isCorrect})
	}
	totalCount := len(req.Answers)
	score := 0
	if totalCount > 0 {
		score = int(float64(correctCount) / float64(totalCount) * 100)
	}
	detailsJSON, _ := json.Marshal(details)
	attempt := model.PracticeAttempt{TaskID: req.TaskID, UserID: req.StudentID, Score: score, Correct: correctCount, Total: totalCount, Details: string(detailsJSON)}
	if err := h.db.Create(&attempt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "提交失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"attemptId": attempt.ID, "score": score, "correctCount": correctCount, "totalCount": totalCount, "details": details}})
}

func (h *CompatibilityHandler) ExplainNodeV1(c *gin.Context) {
	nodeID := c.Param("nodeId")
	var req struct {
		CourseID string `json:"courseId" binding:"required"`
		PageNum  int    `json:"pageNum"`
		Question string `json:"question"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	question := req.Question
	if question == "" {
		question = fmt.Sprintf("请讲解节点 %s 的内容", nodeID)
	}
	resp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{Question: question, CurrentPage: req.PageNum, Context: fmt.Sprintf("节点ID:%s 课程ID:%s", nodeID, req.CourseID), Mode: "llm"})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"code": 503, "message": "AI讲解暂不可用"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"nodeId": nodeID, "explanation": resp.Answer, "sourcePage": resp.SourcePage, "sourceExcerpt": resp.SourceExcerpt}})
}

func (h *CompatibilityHandler) GetStudyStatsV1(c *gin.Context) {
	courseID := c.Param("courseId")
	studentID := c.Query("studentId")
	if studentID == "" {
		studentID = c.GetHeader("X-Student-Id")
	}
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 studentId"})
		return
	}
	var totalQuestions int64
	h.db.Model(&model.QuestionLog{}).Where("user_id = ? AND course_id = ?", studentID, courseID).Count(&totalQuestions)
	weakPoints := h.listWeakPoints(studentID, courseID)
	focusScore := 85
	if totalQuestions > 10 {
		focusScore = 70
	} else if totalQuestions > 5 {
		focusScore = 80
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"studentId": studentID, "courseId": courseID, "focusScore": focusScore, "totalQuestions": totalQuestions, "weakPoints": weakPoints, "noteCount": countStudentNotes(h.db, studentID, courseID)}})
}

func (h *CompatibilityHandler) GetStudentCoursewareListV1(c *gin.Context) {
	h.GetStudentCoursewareList(c)
}

func (h *CompatibilityHandler) StartStudentSessionV1(c *gin.Context) {
	h.StartStudentSession(c)
}

func (h *CompatibilityHandler) UpdateStudentProgressV1(c *gin.Context) {
	h.UpdateStudentProgress(c)
}

func (h *CompatibilityHandler) GetStudentScriptV1(c *gin.Context) {
	courseID := c.Param("courseId")
	pageStr := c.Param("pageNum")
	c.Params = append(filterParams(c.Params, "courseId"), gin.Param{Key: "courseId", Value: courseID})
	c.Params = append(filterParams(c.Params, "page"), gin.Param{Key: "page", Value: pageStr})
	h.GetStudentScript(c)
}

func (h *CompatibilityHandler) StreamStudentQAV1(c *gin.Context) {
	h.StreamStudentQA(c)
}

func (h *CompatibilityHandler) ParseKnowledgeV1(c *gin.Context) {
	NewWeakPointHandler(h.db, h.aiClient).ParseKnowledge(c)
}

func countStudentNotes(db *gorm.DB, studentID, courseID string) int64 {
	var total int64
	db.Model(&model.StudentNote{}).Where("user_id = ? AND course_id = ?", studentID, courseID).Count(&total)
	return total
}

// Review Plan Handlers

func (h *CompatibilityHandler) CreateReviewPlanV1(c *gin.Context) {
	var req struct {
		StudentID   string `json:"studentId" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Frequency   string `json:"frequency" binding:"required"` // daily, weekly, monthly
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	plan := model.ReviewPlan{
		StudentID:   req.StudentID,
		Name:        req.Name,
		Description: req.Description,
		Frequency:   req.Frequency,
		Status:      "active",
	}
	if err := h.db.Create(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建复习计划失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": plan})
}

func (h *CompatibilityHandler) GetReviewPlansV1(c *gin.Context) {
	studentID := c.Query("studentId")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 studentId"})
		return
	}
	var plans []model.ReviewPlan
	if err := h.db.Where("student_id = ?", studentID).Order("created_at desc").Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取复习计划失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": plans})
}

func (h *CompatibilityHandler) UpdateReviewPlanV1(c *gin.Context) {
	planID := c.Param("planId")
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Frequency   string `json:"frequency"`
		Status      string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	var plan model.ReviewPlan
	if err := h.db.First(&plan, "id = ?", planID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "复习计划不存在"})
		return
	}
	if req.Name != "" {
		plan.Name = req.Name
	}
	if req.Description != "" {
		plan.Description = req.Description
	}
	if req.Frequency != "" {
		plan.Frequency = req.Frequency
	}
	if req.Status != "" {
		plan.Status = req.Status
	}
	if err := h.db.Save(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新复习计划失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": plan})
}

func (h *CompatibilityHandler) DeleteReviewPlanV1(c *gin.Context) {
	planID := c.Param("planId")
	if err := h.db.Delete(&model.ReviewPlan{}, "id = ?", planID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除复习计划失败"})
		return
	}
	// Also delete associated items
	h.db.Delete(&model.ReviewPlanItem{}, "review_plan_id = ?", planID)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *CompatibilityHandler) AddReviewPlanItemV1(c *gin.Context) {
	var req struct {
		ReviewPlanID string `json:"reviewPlanId" binding:"required"`
		ItemType     string `json:"itemType" binding:"required"` // note, favorite
		ItemID       string `json:"itemId" binding:"required"`
		Priority     int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	item := model.ReviewPlanItem{
		ReviewPlanID: req.ReviewPlanID,
		ItemType:     req.ItemType,
		ItemID:       req.ItemID,
		Priority:     req.Priority,
	}
	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "添加复习项失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": item})
}

func (h *CompatibilityHandler) GetReviewPlanItemsV1(c *gin.Context) {
	planID := c.Param("planId")
	var items []model.ReviewPlanItem
	if err := h.db.Where("review_plan_id = ?", planID).Order("priority desc, created_at desc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取复习项失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": items})
}

func (h *CompatibilityHandler) UpdateReviewPlanItemV1(c *gin.Context) {
	itemID := c.Param("itemId")
	var req struct {
		Priority       int  `json:"priority"`
		LastReviewedAt *time.Time `json:"lastReviewedAt"`
		ReviewCount    int `json:"reviewCount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	var item model.ReviewPlanItem
	if err := h.db.First(&item, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "复习项不存在"})
		return
	}
	if req.Priority > 0 {
		item.Priority = req.Priority
	}
	if req.LastReviewedAt != nil {
		item.LastReviewedAt = req.LastReviewedAt
		item.ReviewCount = req.ReviewCount
	}
	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新复习项失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": item})
}

func (h *CompatibilityHandler) DeleteReviewPlanItemV1(c *gin.Context) {
	itemID := c.Param("itemId")
	if err := h.db.Delete(&model.ReviewPlanItem{}, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除复习项失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

// Delete Student Note
func (h *CompatibilityHandler) DeleteStudentNoteV1(c *gin.Context) {
	noteID := c.Param("noteId")
	studentID := c.Query("studentId")
	if noteID == "" || studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少参数"})
		return
	}
	note := model.StudentNote{}
	if err := h.db.First(&note, "id = ?", noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "笔记不存在"})
		return
	}
	if note.UserID != studentID {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "笔记不存在"})
		return
	}
	if err := h.db.Delete(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}
