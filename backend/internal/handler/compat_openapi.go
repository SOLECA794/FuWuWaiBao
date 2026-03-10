package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/service"
)

func (h *CompatibilityHandler) OpenLessonParse(c *gin.Context) {
	var req struct {
		SchoolID          string `json:"schoolId"`
		UserID            string `json:"userId"`
		CourseID          string `json:"courseId"`
		FileType          string `json:"fileType" binding:"required"`
		FileURL           string `json:"fileUrl" binding:"required"`
		IsExtractKeyPoint bool   `json:"isExtractKeyPoint"`
		Enc               string `json:"enc"`
		Time              string `json:"time"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		openAPIError(c, http.StatusBadRequest, "参数错误")
		return
	}
	title := filepath.Base(req.FileURL)
	if decoded, err := url.QueryUnescape(title); err == nil {
		title = decoded
	}
	course := model.Course{BaseModel: model.BaseModel{ID: strings.TrimSpace(req.CourseID)}, Title: title, FileURL: req.FileURL, FileType: req.FileType}
	if course.ID == "" {
		course.ID = uuid.NewString()
	}
	var existing model.Course
	if err := h.db.First(&existing, "id = ?", course.ID).Error; err == nil {
		existing.Title = course.Title
		existing.FileURL = course.FileURL
		existing.FileType = course.FileType
		_ = h.db.Save(&existing).Error
		course = existing
	} else {
		_ = h.db.Create(&course).Error
	}
	openAPISuccess(c, "课件解析成功", gin.H{"parseId": course.ID, "fileInfo": gin.H{"fileName": title, "fileSize": 0, "pageCount": course.TotalPage}, "structurePreview": gin.H{"chapters": []gin.H{}}, "taskStatus": "processing"})
}

func (h *CompatibilityHandler) OpenGenerateScript(c *gin.Context) {
	var req struct {
		ParseID       string `json:"parseId" binding:"required"`
		TeachingStyle string `json:"teachingStyle"`
		SpeechSpeed   string `json:"speechSpeed"`
		CustomOpening string `json:"customOpening"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		openAPIError(c, http.StatusBadRequest, "参数错误")
		return
	}
	var course model.Course
	if err := h.db.First(&course, "id = ?", req.ParseID).Error; err != nil {
		openAPIError(c, http.StatusNotFound, "解析任务不存在")
		return
	}
	content := req.CustomOpening
	if strings.TrimSpace(content) == "" {
		content = "课程标题：" + course.Title
	}
	script := ""
	teachingNodes := loadTeachingNodesByPage(h.db, req.ParseID, 1)
	if generatedScript, _, usedNodes, err := generateAndStoreTeachingNodeScripts(c.Request.Context(), h.db, h.aiClient, course.Title, "llm", teachingNodes); usedNodes && err == nil {
		script = generatedScript
	} else {
		resp, err := h.aiClient.GenerateScript(c.Request.Context(), service.GenerateScriptRequest{Page: 1, Content: content, CourseName: course.Title, Mode: "llm"})
		if err != nil {
			openAPIError(c, http.StatusInternalServerError, "脚本生成失败")
			return
		}
		script = resp.Script
	}
	sections := buildScriptNodes(1, script)
	openAPISuccess(c, "脚本生成成功", gin.H{"scriptId": "script_" + req.ParseID, "scriptStructure": sections, "editUrl": "/teacher/script/" + req.ParseID + "/1", "audioGenerateUrl": "/api/v1/lesson/generateAudio"})
}

func (h *CompatibilityHandler) OpenGenerateAudio(c *gin.Context) {
	var req struct {
		ScriptID    string   `json:"scriptId" binding:"required"`
		VoiceType   string   `json:"voiceType"`
		AudioFormat string   `json:"audioFormat"`
		SectionIDs  []string `json:"sectionIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		openAPIError(c, http.StatusBadRequest, "参数错误")
		return
	}
	openAPISuccess(c, "语音合成成功", gin.H{"audioId": "audio_" + strings.TrimPrefix(req.ScriptID, "script_"), "audioUrl": "", "audioInfo": gin.H{"totalDuration": 0, "fileSize": 0, "format": defaultString(req.AudioFormat, "mp3"), "bitRate": 0}, "sectionAudios": []gin.H{}})
}

func (h *CompatibilityHandler) OpenQAInteract(c *gin.Context) {
	var req struct {
		SchoolID         string `json:"schoolId"`
		UserID           string `json:"userId"`
		CourseID         string `json:"courseId" binding:"required"`
		LessonID         string `json:"lessonId"`
		SessionID        string `json:"sessionId"`
		QuestionType     string `json:"questionType"`
		Question         string `json:"question" binding:"required"`
		CurrentSectionID string `json:"currentSectionId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		openAPIError(c, http.StatusBadRequest, "参数错误")
		return
	}
	page := parsePageFromNodeID(req.CurrentSectionID)
	if page <= 0 {
		page = 1
	}
	var coursePage model.CoursePage
	contextText := ""
	if err := h.db.Where("course_id = ? AND page_index = ?", req.CourseID, page).First(&coursePage).Error; err == nil {
		contextText = pageContextText(coursePage)
	}
	if strings.TrimSpace(contextText) == "" {
		contextText = buildPageContextFromTeachingNodes(loadTeachingNodesByPage(h.db, req.CourseID, page))
	}
	resp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{Question: req.Question, CurrentPage: page, Context: contextText, Mode: "llm"})
	if err != nil {
		openAPIError(c, http.StatusServiceUnavailable, "问答交互失败")
		return
	}
	openAPISuccess(c, "问答交互成功", gin.H{"answerId": "ans_" + uuid.NewString(), "answerContent": resp.Answer, "answerType": defaultString(req.QuestionType, "text"), "suggestions": []string{resp.FollowUpSuggestion}, "understandingLevel": understandingLevel(resp.Intent.NeedReteach)})
}

func (h *CompatibilityHandler) OpenVoiceToText(c *gin.Context) {
	var req struct {
		AudioURL string `json:"audioUrl"`
		Text     string `json:"text"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		openAPIError(c, http.StatusBadRequest, "参数错误")
		return
	}
	text := strings.TrimSpace(req.Text)
	openAPISuccess(c, "语音识别成功", gin.H{"text": text, "audioUrl": req.AudioURL})
}

func (h *CompatibilityHandler) OpenTrackProgress(c *gin.Context) {
	var req struct {
		UserID           string `json:"userId" binding:"required"`
		CourseID         string `json:"courseId" binding:"required"`
		CurrentPage      int    `json:"currentPage"`
		CurrentSectionID string `json:"currentSectionId"`
		SessionID        string `json:"sessionId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		openAPIError(c, http.StatusBadRequest, "参数错误")
		return
	}
	page := req.CurrentPage
	if page <= 0 {
		page = parsePageFromNodeID(req.CurrentSectionID)
	}
	if page <= 0 {
		page = 1
	}
	state := sessionState{SessionID: defaultString(req.SessionID, "sess_"+uuid.NewString()), UserID: req.UserID, CourseID: req.CourseID, CurrentPage: page, CurrentNodeID: defaultString(req.CurrentSectionID, fmt.Sprintf("p%d_n1", page)), UpdatedAt: time.Now()}
	h.persistSession(state)
	var progress model.UserProgress
	if err := h.db.Where("user_id = ? AND course_id = ?", req.UserID, req.CourseID).First(&progress).Error; err == nil {
		_ = h.db.Model(&progress).Update("last_page", page).Error
	} else {
		_ = h.db.Create(&model.UserProgress{UserID: req.UserID, CourseID: req.CourseID, LastPage: page}).Error
	}
	openAPISuccess(c, "进度记录成功", gin.H{"sessionId": state.SessionID, "currentPage": page, "currentSectionId": state.CurrentNodeID})
}

func (h *CompatibilityHandler) OpenAdjustProgress(c *gin.Context) {
	var req struct {
		CourseID           string `json:"courseId" binding:"required"`
		CurrentSectionID   string `json:"currentSectionId"`
		UnderstandingLevel string `json:"understandingLevel"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		openAPIError(c, http.StatusBadRequest, "参数错误")
		return
	}
	adjustType := "continue"
	supplementContent := ""
	if strings.EqualFold(req.UnderstandingLevel, "partial") || strings.EqualFold(req.UnderstandingLevel, "low") {
		adjustType = "reteach"
		supplementContent = "建议补充讲解当前节点的基础概念与示例。"
	}
	openAPISuccess(c, "节奏调整成功", gin.H{"adjustPlan": gin.H{"continueSectionId": nextNodeID(req.CurrentSectionID, parsePageFromNodeID(req.CurrentSectionID)), "adjustType": adjustType, "supplementContent": supplementContent}})
}

func (h *CompatibilityHandler) OpenSyncCourse(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		openAPIError(c, http.StatusBadRequest, "参数错误")
		return
	}
	openAPISuccess(c, "课程同步成功", gin.H{"received": req})
}

func (h *CompatibilityHandler) OpenSyncUser(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		openAPIError(c, http.StatusBadRequest, "参数错误")
		return
	}
	openAPISuccess(c, "用户同步成功", gin.H{"received": req})
}

func fileHeaderFromMultipart(file *multipart.FileHeader) *multipart.FileHeader {
	return file
}
