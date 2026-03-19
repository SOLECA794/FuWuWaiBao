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
		SessionID:     uuid.NewString(),
		UserID:        req.UserID,
		CourseID:      req.CourseID,
		CurrentPage:   1,
		CurrentNodeID: "p1_n1",
		UpdatedAt:     time.Now(),
	}
	h.persistSession(state)
	syncDialogueSessionState(h.db, state.SessionID, req.UserID, req.CourseID, 1, state.CurrentNodeID, 0)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"sessionId": state.SessionID, "courseId": req.CourseID}})
}

func (h *CompatibilityHandler) UpdateStudentProgress(c *gin.Context) {
	var req struct {
		SessionID      string `json:"sessionId"`
		UserID         string `json:"userId"`
		CourseID       string `json:"courseId" binding:"required"`
		Page           int    `json:"page"`
		CurrentPage    int    `json:"currentPage"`
		NodeID         string `json:"nodeId"`
		CurrentNodeID  string `json:"currentNodeId"`
		CurrentTimeSec int    `json:"currentTimeSec"`
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
		state.SessionID = uuid.NewString()
	}
	h.persistSession(state)
	syncDialogueSessionState(h.db, state.SessionID, req.UserID, req.CourseID, page, nodeID, req.CurrentTimeSec)

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
	displayText := pageDisplayText(coursePage)
	teachingNodes := loadTeachingNodesByPage(h.db, courseID, page)
	nodes := buildPlaybackNodesFromTeachingNodes(teachingNodes)
	if len(nodes) == 0 {
		nodes = buildScriptNodes(page, displayText)
	}
	audioMeta := buildPlaybackAudioMeta(courseID, page, teachingNodes)
	pageSummary := strings.TrimSpace(displayText)
	if len(teachingNodes) > 0 {
		pageSummary = buildTeachingNodePageSummary(teachingNodes)
	}
	if len(pageSummary) > 80 {
		pageSummary = pageSummary[:80]
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"courseId": courseID, "page": page, "nodes": nodes, "page_summary": pageSummary, "audio_meta": audioMeta, "playback_mode": audioMeta["playback_mode"]}})
}

func (h *CompatibilityHandler) GetStudentPlaybackAudio(c *gin.Context) {
	courseID := c.Param("courseId")
	page, err := strconv.Atoi(c.Param("pageNum"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "页码错误"})
		return
	}

	nodes := loadTeachingNodesByPage(h.db, courseID, page)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": buildPlaybackAudioMeta(courseID, page, nodes)})
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
	session := h.loadSession(req.SessionID)
	userID := defaultStudentID(func() string {
		if session != nil {
			return session.UserID
		}
		return ""
	}())
	sessionID := strings.TrimSpace(req.SessionID)
	if sessionID == "" && session != nil {
		sessionID = session.SessionID
	}
	if sessionID == "" {
		sessionID = uuid.NewString()
	}
	nodeID := defaultStringValue(req.NodeID, fmt.Sprintf("p%d_n1", req.Page))
	historySummary, recentTurns := buildDialogueContext(h.db, sessionID, 4)
	var coursePage model.CoursePage
	contextText := ""
	if err := h.db.Where("course_id = ? AND page_index = ?", req.CourseID, req.Page).First(&coursePage).Error; err == nil {
		contextText = pageContextText(coursePage)
	}
	nodeScopedContext := buildNodeScopedContext(h.db, req.CourseID, req.Page, nodeID)
	if strings.TrimSpace(nodeScopedContext) != "" {
		contextText = nodeScopedContext
	}
	if strings.TrimSpace(contextText) == "" {
		contextText = buildPageContextFromTeachingNodes(loadTeachingNodesByPage(h.db, req.CourseID, req.Page))
	}
	aiResp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{Question: req.Question, CurrentPage: req.Page, Context: contextText, Mode: "llm", SessionID: sessionID, HistorySummary: historySummary, RecentTurns: recentTurns})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"code": 503, "message": "AI服务暂不可用"})
		return
	}
	resumePage := maxInt(aiResp.ResumePage, req.Page)
	resumeNodeID := resolveResumeNodeIDByCourse(h.db, req.CourseID, nodeID, resumePage, aiResp.Intent.NeedReteach)
	resumeSec := resolvePlaybackResumeSec(h.db, req.CourseID, resumePage, resumeNodeID)
	appendDialogueTurn(h.db, sessionID, userID, req.CourseID, req.Page, nodeID, req.Question, aiResp.Answer, maxInt(aiResp.SourcePage, req.Page), aiResp.Intent.NeedReteach, aiResp.FollowUpSuggestion)
	if userID != "" {
		_ = h.db.Create(&model.QuestionLog{UserID: userID, CourseID: req.CourseID, PageIndex: req.Page, NodeID: nodeID, Question: req.Question, Answer: aiResp.Answer}).Error
	}
	state := sessionState{SessionID: sessionID, UserID: userID, CourseID: req.CourseID, CurrentPage: resumePage, CurrentNodeID: resumeNodeID, UpdatedAt: time.Now()}
	h.persistSession(state)
	syncDialogueSessionState(h.db, sessionID, userID, req.CourseID, resumePage, resumeNodeID, resumeSec)

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
		"session_id":           sessionID,
		"need_reteach":         aiResp.Intent.NeedReteach,
		"source_page":          maxInt(aiResp.SourcePage, req.Page),
		"source_node_id":       nodeID,
		"resume_page":          resumePage,
		"resume_node_id":       resumeNodeID,
		"resume_sec":           resumeSec,
		"follow_up_suggestion": aiResp.FollowUpSuggestion,
	})
}

func (h *CompatibilityHandler) GetStudentPlaybackAudioV1(c *gin.Context) {
	h.GetStudentPlaybackAudio(c)
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
		NodeID    string `json:"nodeId"`
		PageNum   int    `json:"pageNum" binding:"required"`
		Content   string `json:"content" binding:"required"`
		X         int    `json:"x"`
		Y         int    `json:"y"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	nodeID := strings.TrimSpace(req.NodeID)
	if nodeID == "" {
		nodeID = fmt.Sprintf("p%d_n1", req.PageNum)
	}
	note := model.StudentNote{UserID: req.StudentID, CourseID: courseID, NodeID: nodeID, PageNum: req.PageNum, Note: req.Content}
	if err := h.db.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "保存成功", "data": gin.H{"noteId": note.ID, "nodeId": nodeID}})
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
