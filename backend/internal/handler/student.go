package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/service"
	"smart-teaching-backend/pkg/logger"
)

type StudentHandler struct {
	db       *gorm.DB
	aiClient service.AIEngine
}

func NewStudentHandler(db *gorm.DB, aiClient service.AIEngine) *StudentHandler {
	return &StudentHandler{db: db, aiClient: aiClient}
}

type GetCoursewarePageRequest struct {
	CourseID    string `json:"courseId" binding:"required"`
	CurrentPage int    `json:"currentPage" binding:"required"`
}

type AIQuestionRequest struct {
	StudentID string `json:"studentId"`
	CourseID  string `json:"courseId" binding:"required"`
	PageNum   int    `json:"pageNum" binding:"required"`
	Question  string `json:"question" binding:"required"`
	Type      string `json:"type"`
}

type TraceQuestionRequest struct {
	StudentID string  `json:"studentId"`
	CourseID  string  `json:"courseId" binding:"required"`
	PageNum   int     `json:"pageNum" binding:"required"`
	X         float64 `json:"x" binding:"required"`
	Y         float64 `json:"y" binding:"required"`
	Question  string  `json:"question" binding:"required"`
}

type UpdateBreakpointRequest struct {
	StudentID string `json:"studentId"`
	UserID    string `json:"userId"`
	CourseID  string `json:"courseId"`
	LastPage  int    `json:"lastPage"`
	PageNum   int    `json:"pageNum"`
}

type SaveNoteRequest struct {
	StudentID string `json:"studentId"`
	UserID    string `json:"userId"`
	CourseID  string `json:"courseId"`
	PageNum   int    `json:"pageNum" binding:"required"`
	Note      string `json:"note"`
	Content   string `json:"content"`
}

func (h *StudentHandler) StartSession(c *gin.Context) {
	var req struct {
		UserID   string `json:"userId" binding:"required"`
		CourseID string `json:"courseId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	sessionID := "sess_" + uuid.NewString()
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"sessionId": sessionID, "courseId": req.CourseID}})
}

func (h *StudentHandler) UpdateProgress(c *gin.Context) {
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
	if strings.TrimSpace(req.UserID) != "" {
		h.saveProgress(req.UserID, req.CourseID, page)
	}

	nodeID := strings.TrimSpace(req.NodeID)
	if nodeID == "" {
		nodeID = strings.TrimSpace(req.CurrentNodeID)
	}
	if nodeID == "" {
		nodeID = fmt.Sprintf("p%d_n1", page)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok", "data": gin.H{"sessionId": req.SessionID, "page": page, "nodeId": nodeID}})
}

func (h *StudentHandler) GetBreakpoint(c *gin.Context) {
	courseID := c.Param("courseId")
	userID := strings.TrimSpace(c.Query("userId"))
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少参数: userId"})
		return
	}

	lastPage := h.getLastPage(userID, courseID)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"lastPageNum": lastPage, "pageNum": lastPage}})
}

func (h *StudentHandler) GetStudentBreakpoint(c *gin.Context) {
	courseID := strings.TrimSpace(c.Query("courseId"))
	studentID := strings.TrimSpace(c.Query("studentId"))
	if studentID == "" || courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少必填参数: studentId 或 courseId"})
		return
	}

	lastPage := h.getLastPage(studentID, courseID)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"lastPageNum": lastPage}})
}

func (h *StudentHandler) UpdateBreakpoint(c *gin.Context) {
	h.updateBreakpointCore(c, c.Param("courseId"))
}

func (h *StudentHandler) UpdateStudentBreakpoint(c *gin.Context) {
	var req UpdateBreakpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少必填参数: studentId, courseId 或 lastPage"})
		return
	}

	studentID := defaultStudentID(req.StudentID, req.UserID)
	page := req.LastPage
	if page <= 0 {
		page = req.PageNum
	}
	if page <= 0 {
		page = 1
	}
	if studentID == "" || strings.TrimSpace(req.CourseID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少必填参数: studentId, courseId 或 lastPage"})
		return
	}

	h.saveProgress(studentID, req.CourseID, page)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "断点更新成功", "data": gin.H{"lastPageNum": page}})
}

func (h *StudentHandler) SaveNote(c *gin.Context) {
	h.saveNoteCore(c, c.Param("courseId"))
}

func (h *StudentHandler) SaveStudentNote(c *gin.Context) {
	var req SaveNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少必填参数"})
		return
	}

	studentID := defaultStudentID(req.StudentID, req.UserID)
	content := strings.TrimSpace(req.Content)
	if content == "" {
		content = strings.TrimSpace(req.Note)
	}
	if studentID == "" || strings.TrimSpace(req.CourseID) == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少必填参数"})
		return
	}

	note := model.StudentNote{UserID: studentID, CourseID: req.CourseID, PageNum: req.PageNum, Note: content}
	if err := h.db.Create(&note).Error; err != nil {
		logger.Errorf("保存笔记失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "保存成功", "data": gin.H{"status": "saved", "noteId": note.ID}})
}

func (h *StudentHandler) GetCoursewarePage(c *gin.Context) {
	if courseID := c.Param("courseId"); courseID != "" {
		h.getCoursewarePageV1(c, courseID)
		return
	}

	var req GetCoursewarePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	content, title := h.loadPageScript(req.CourseID, req.CurrentPage)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": gin.H{"courseId": req.CourseID, "currentPage": req.CurrentPage, "content": content, "title": title}})
}

func (h *StudentHandler) AskAIQuestion(c *gin.Context) {
	var req AIQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少必填参数: courseId, pageNum 或 question"})
		return
	}

	resp := h.askQuestionWithFallback(c, req.CourseID, req.PageNum, req.Question)
	h.recordQuestion(defaultStudentID(req.StudentID, c.GetString("userId")), req.CourseID, req.PageNum, req.Question, resp.Answer)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": gin.H{"answer": resp.Answer, "sourcePage": resp.SourcePage, "sourceExcerpt": resp.SourceExcerpt, "needReteach": resp.Intent.NeedReteach, "followUpSuggestion": resp.FollowUpSuggestion, "aiUnavailable": resp.Fallback}})
}

func (h *StudentHandler) TraceAIQuestion(c *gin.Context) {
	var req TraceQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	question := fmt.Sprintf("%s（圈选坐标: %.2f, %.2f）", req.Question, req.X, req.Y)
	resp := h.askQuestionWithFallback(c, req.CourseID, req.PageNum, question)
	h.recordQuestion(defaultStudentID(req.StudentID, c.GetString("userId")), req.CourseID, req.PageNum, req.Question, resp.Answer)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "成功", "data": gin.H{"answer": resp.Answer, "sourcePage": resp.SourcePage, "sourceExcerpt": resp.SourceExcerpt}})
}

func (h *StudentHandler) QAStream(c *gin.Context) {
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

	resp := h.askQuestionWithFallback(c, req.CourseID, req.Page, req.Question)
	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "当前连接不支持流式输出"})
		return
	}

	writeEvent := func(event string, payload any) {
		body, _ := json.Marshal(payload)
		_, _ = fmt.Fprintf(c.Writer, "event: %s\n", event)
		_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", string(body))
		flusher.Flush()
	}

	parts := strings.Fields(resp.Answer)
	if len(parts) == 0 {
		for _, r := range []rune(resp.Answer) {
			writeEvent("token", gin.H{"text": string(r)})
		}
	} else {
		for _, part := range parts {
			writeEvent("token", gin.H{"text": part + " "})
		}
	}
	writeEvent("sentence", gin.H{"text": resp.Answer})
	writeEvent("final", gin.H{"need_reteach": resp.Intent.NeedReteach, "source_page": resp.SourcePage, "resume_page": maxInt(resp.ResumePage, req.Page), "resume_node_id": defaultStringValue(req.NodeID, fmt.Sprintf("p%d_n1", req.Page))})
	h.recordQuestion(c.GetString("userId"), req.CourseID, req.Page, req.Question, resp.Answer)
}

func (h *StudentHandler) GetStudentStudyData(c *gin.Context) {
	studentID := strings.TrimSpace(c.Query("studentId"))
	courseID := strings.TrimSpace(c.Query("courseId"))
	if studentID == "" || courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少必填参数: studentId 或 courseId"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": h.buildStudyStats(studentID, courseID)})
}

func (h *StudentHandler) GetPersonalStats(c *gin.Context) {
	courseID := c.Param("courseId")
	studentID := strings.TrimSpace(c.Query("userId"))
	if studentID == "" {
		studentID = strings.TrimSpace(c.Query("studentId"))
	}
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少参数: userId"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": h.buildStudyStats(studentID, courseID)})
}

func (h *StudentHandler) getCoursewarePageV1(c *gin.Context, courseID string) {
	pageNum, err := strconv.Atoi(c.Param("pageNum"))
	if err != nil || pageNum < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "页码必须是数字"})
		return
	}

	content, _ := h.loadPageScript(courseID, pageNum)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"courseId": courseID, "page": pageNum, "nodes": buildPlaybackNodes(pageNum, content), "page_summary": content}})
}

func (h *StudentHandler) updateBreakpointCore(c *gin.Context, courseID string) {
	var req UpdateBreakpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	studentID := defaultStudentID(req.StudentID, req.UserID)
	page := req.PageNum
	if page <= 0 {
		page = req.LastPage
	}
	if page <= 0 {
		page = 1
	}
	if studentID == "" || courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	h.saveProgress(studentID, courseID, page)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok", "data": gin.H{"lastPageNum": page, "pageNum": page}})
}

func (h *StudentHandler) saveNoteCore(c *gin.Context, courseID string) {
	var req SaveNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	studentID := defaultStudentID(req.StudentID, req.UserID)
	content := strings.TrimSpace(req.Content)
	if content == "" {
		content = strings.TrimSpace(req.Note)
	}
	if studentID == "" || courseID == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	note := model.StudentNote{UserID: studentID, CourseID: courseID, PageNum: req.PageNum, Note: content}
	if err := h.db.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "保存成功", "data": gin.H{"noteId": note.ID, "status": "saved"}})
}

type questionReply struct {
	Answer             string
	SourcePage         int
	SourceExcerpt      string
	ResumePage         int
	FollowUpSuggestion string
	Fallback           bool
	Intent             struct {
		NeedReteach bool
	}
}

func (h *StudentHandler) askQuestionWithFallback(c *gin.Context, courseID string, pageNum int, question string) questionReply {
	content, _ := h.loadPageScript(courseID, pageNum)
	result := questionReply{SourcePage: pageNum, ResumePage: pageNum}
	if h.aiClient != nil {
		resp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{Question: question, CurrentPage: pageNum, Context: content, Mode: "llm"})
		if err == nil && strings.TrimSpace(resp.Answer) != "" {
			result.Answer = resp.Answer
			result.SourcePage = maxInt(resp.SourcePage, pageNum)
			result.SourceExcerpt = resp.SourceExcerpt
			result.ResumePage = maxInt(resp.ResumePage, pageNum)
			result.FollowUpSuggestion = resp.FollowUpSuggestion
			result.Intent.NeedReteach = resp.Intent.NeedReteach
			return result
		}
	}

	preview := content
	if len([]rune(preview)) > 80 {
		runes := []rune(preview)
		preview = string(runes[:80]) + "..."
	}
	if strings.TrimSpace(preview) == "" {
		preview = "当前页暂无讲稿内容，可先让教师补充讲稿后再进行精确答疑。"
	}
	result.Answer = fmt.Sprintf("根据当前课件内容，关于“%s”的说明如下：%s", question, preview)
	result.SourceExcerpt = preview
	result.FollowUpSuggestion = "可以继续追问关键概念、公式来源或应用场景。"
	result.Fallback = true
	return result
}

func (h *StudentHandler) loadPageScript(courseID string, pageNum int) (string, string) {
	var course model.Course
	_ = h.db.First(&course, "id = ?", courseID).Error
	var coursePage model.CoursePage
	if err := h.db.Where("course_id = ? AND page_index = ?", courseID, pageNum).First(&coursePage).Error; err == nil {
		return strings.TrimSpace(coursePage.ScriptText), course.Title
	}
	return "", course.Title
}

func (h *StudentHandler) recordQuestion(studentID, courseID string, pageNum int, question, answer string) {
	studentID = strings.TrimSpace(studentID)
	if studentID == "" || strings.TrimSpace(courseID) == "" || strings.TrimSpace(question) == "" {
		return
	}
	_ = h.db.Create(&model.QuestionLog{UserID: studentID, CourseID: courseID, PageIndex: pageNum, Question: question, Answer: answer}).Error
}

func (h *StudentHandler) saveProgress(studentID, courseID string, page int) {
	var progress model.UserProgress
	err := h.db.Where("user_id = ? AND course_id = ?", studentID, courseID).First(&progress).Error
	if err == nil {
		_ = h.db.Model(&progress).Update("last_page", page).Error
		return
	}
	_ = h.db.Create(&model.UserProgress{UserID: studentID, CourseID: courseID, LastPage: page}).Error
}

func (h *StudentHandler) getLastPage(studentID, courseID string) int {
	var progress model.UserProgress
	if err := h.db.Where("user_id = ? AND course_id = ?", studentID, courseID).First(&progress).Error; err == nil && progress.LastPage > 0 {
		return progress.LastPage
	}
	return 1
}

func (h *StudentHandler) buildStudyStats(studentID, courseID string) gin.H {
	var totalQuestions int64
	h.db.Model(&model.QuestionLog{}).Where("user_id = ? AND course_id = ?", studentID, courseID).Count(&totalQuestions)
	var noteCount int64
	h.db.Model(&model.StudentNote{}).Where("user_id = ? AND course_id = ?", studentID, courseID).Count(&noteCount)
	weakPoints := buildStudentWeakPoints(h.db, studentID, courseID)
	focusScore := 90
	if totalQuestions > 10 {
		focusScore = 72
	} else if totalQuestions > 5 {
		focusScore = 82
	}

	mastery := gin.H{}
	for _, item := range weakPoints {
		name, _ := item["name"].(string)
		masteryLevel, _ := item["mastery"].(int)
		mastery[name] = masteryLevel
	}

	return gin.H{"studentId": studentID, "courseId": courseID, "focusScore": focusScore, "totalQuestions": totalQuestions, "weakPoints": weakPoints, "noteCount": noteCount, "studyHours": float64(maxInt(int(noteCount), 1)) * 0.5, "mastery": mastery}
}

func buildStudentWeakPoints(db *gorm.DB, studentID, courseID string) []gin.H {
	type pageStat struct {
		PageIndex int
		Count     int
	}
	var stats []pageStat
	db.Table("question_logs").Select("page_index, count(*) as count").Where("user_id = ? AND course_id = ?", studentID, courseID).Group("page_index").Order("count desc").Limit(5).Scan(&stats)
	result := make([]gin.H, 0, len(stats))
	for _, stat := range stats {
		mastery := 100 - stat.Count*12
		if mastery < 0 {
			mastery = 0
		}
		result = append(result, gin.H{"weakPointId": fmt.Sprintf("wp-page-%d", stat.PageIndex), "name": fmt.Sprintf("第%d页知识点", stat.PageIndex), "count": stat.Count, "mastery": mastery, "pageIndex": stat.PageIndex})
	}
	return result
}

func buildPlaybackNodes(page int, content string) []gin.H {
	content = strings.TrimSpace(content)
	if content == "" {
		return []gin.H{}
	}
	parts := strings.FieldsFunc(content, func(r rune) bool {
		return r == '\n' || r == '。' || r == '！' || r == '？'
	})
	nodes := make([]gin.H, 0, len(parts))
	for index, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		nodeType := "explain"
		if index == 0 {
			nodeType = "opening"
		} else if index == len(parts)-1 {
			nodeType = "transition"
		}
		nodes = append(nodes, gin.H{"node_id": fmt.Sprintf("p%d_n%d", page, index+1), "type": nodeType, "text": part})
	}
	return nodes
}

func defaultStudentID(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func defaultStringValue(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
