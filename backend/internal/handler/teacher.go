package handler

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/service"
	"smart-teaching-backend/pkg/logger"
)

type TeacherHandler struct {
	db       *gorm.DB
	aiClient service.AIEngine
}

func NewTeacherHandler(db *gorm.DB, aiClient service.AIEngine) *TeacherHandler {
	return &TeacherHandler{db: db, aiClient: aiClient}
}

func (h *TeacherHandler) GetCoursewareList(c *gin.Context) {
	var courses []model.Course
	if err := h.db.Order("created_at desc").Find(&courses).Error; err != nil {
		logger.Errorf("获取课件列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取课件列表失败"})
		return
	}

	data := make([]gin.H, 0, len(courses))
	for _, course := range courses {
		data = append(data, gin.H{
			"id":           course.ID,
			"courseId":     course.ID,
			"title":        course.Title,
			"file_url":     course.FileURL,
			"fileType":     course.FileType,
			"file_type":    course.FileType,
			"total_page":   course.TotalPage,
			"is_published": course.IsPublished,
			"status":       publishStatus(course.IsPublished),
			"createdAt":    course.CreatedAt,
			"created_at":   course.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": data})
}

func (h *TeacherHandler) GetScript(c *gin.Context) {
	courseID := c.Param("courseId")
	pageNum, err := parsePageParam(c.Param("pageNum"), c.Param("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "页码必须是数字"})
		return
	}

	var coursePage model.CoursePage
	if err := h.db.Where("course_id = ? AND page_index = ?", courseID, pageNum).First(&coursePage).Error; err != nil {
		nodes := loadTeachingNodesByPage(h.db, courseID, pageNum)
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"courseId": courseID, "pageNum": pageNum, "page": pageNum, "content": buildPageContextFromTeachingNodes(nodes)}})
		return
	}

	content := strings.TrimSpace(coursePage.ScriptText)
	nodes := loadTeachingNodesByPage(h.db, courseID, pageNum)
	if content == "" {
		content = buildPageContextFromTeachingNodes(nodes)
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"courseId": courseID, "pageNum": pageNum, "page": pageNum, "content": content, "nodes": buildTeacherNodePayload(nodes)}})
}

func (h *TeacherHandler) UpdateScript(c *gin.Context) {
	courseID := c.Param("courseId")
	pageNum, err := parsePageParam(c.Param("pageNum"), c.Param("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "页码必须是数字"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := h.upsertScript(courseID, pageNum, req.Content); err != nil {
		logger.Errorf("保存讲稿失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "保存成功"})
}

func (h *TeacherHandler) GetTeachingNodes(c *gin.Context) {
	courseID := c.Param("courseId")
	pageNum, err := parsePageParam(c.Param("pageNum"), c.Param("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "页码必须是数字"})
		return
	}

	nodes := loadTeachingNodesByPage(h.db, courseID, pageNum)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"courseId": courseID, "pageNum": pageNum, "nodes": buildTeacherNodePayload(nodes)}})
}

func (h *TeacherHandler) UpdateTeachingNodes(c *gin.Context) {
	courseID := c.Param("courseId")
	pageNum, err := parsePageParam(c.Param("pageNum"), c.Param("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "页码必须是数字"})
		return
	}

	var req struct {
		Nodes []teacherNodeUpsertRequest `json:"nodes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	content, savedNodes, err := h.replaceTeachingNodes(courseID, pageNum, req.Nodes)
	if err != nil {
		logger.Errorf("保存节点讲稿失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "节点保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "保存成功", "data": gin.H{"courseId": courseID, "pageNum": pageNum, "content": content, "nodes": buildTeacherNodePayload(savedNodes)}})
}

func (h *TeacherHandler) SaveScript(c *gin.Context) {
	var req struct {
		CourseID string `json:"courseId" binding:"required"`
		Page     int    `json:"page" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := h.upsertScript(req.CourseID, req.Page, req.Content); err != nil {
		logger.Errorf("保存讲稿失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "保存成功"})
}

func (h *TeacherHandler) AIGenerateScript(c *gin.Context) {
	courseID := c.Param("courseId")
	if courseID == "" {
		courseID = strings.TrimSpace(c.PostForm("courseId"))
	}

	var req struct {
		PageNum    int    `json:"pageNum"`
		Page       int    `json:"page"`
		Mode       string `json:"mode"`
		CourseID   string `json:"courseId"`
		CourseName string `json:"courseName"`
	}
	_ = c.ShouldBindJSON(&req)
	if courseID == "" {
		courseID = req.CourseID
	}
	pageNum := req.PageNum
	if pageNum <= 0 {
		pageNum = req.Page
	}
	if pageNum <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var course model.Course
	if err := h.db.First(&course, "id = ?", courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "课件不存在"})
		return
	}

	courseName := course.Title
	if strings.TrimSpace(req.CourseName) != "" {
		courseName = req.CourseName
	}

	var page model.CoursePage
	contextText := ""
	if err := h.db.Where("course_id = ? AND page_index = ?", courseID, pageNum).First(&page).Error; err == nil {
		contextText = pageContextText(page)
	}
	if strings.TrimSpace(contextText) == "" {
		contextText = fmt.Sprintf("课程：%s，第 %d 页", courseName, pageNum)
	}

	script := ""
	mindmapMarkdown := ""
	teachingNodes := loadTeachingNodesByPage(h.db, courseID, pageNum)
	if h.aiClient != nil {
		if generatedScript, generatedMindmap, usedNodes, err := generateAndStoreTeachingNodeScripts(c.Request.Context(), h.db, h.aiClient, courseName, defaultTeacherMode(req.Mode), teachingNodes); usedNodes && err == nil {
			script = generatedScript
			mindmapMarkdown = generatedMindmap
		} else {
			resp, genErr := h.aiClient.GenerateScript(c.Request.Context(), service.GenerateScriptRequest{Page: pageNum, Content: contextText, CourseName: courseName, Mode: defaultTeacherMode(req.Mode)})
			if genErr == nil && strings.TrimSpace(resp.Script) != "" {
				script = resp.Script
				mindmapMarkdown = resp.MindmapMarkdown
			}
		}
	}
	if strings.TrimSpace(script) == "" {
		script = generateMockScript(courseName, pageNum)
	}

	if err := h.upsertScript(courseID, pageNum, script); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"courseId": courseID, "pageNum": pageNum, "page": pageNum, "content": script, "mindmapMarkdown": mindmapMarkdown}})
}

func (h *TeacherHandler) GeneratePageAudio(c *gin.Context) {
	courseID := c.Param("courseId")
	pageNum, err := parsePageParam(c.Param("pageNum"), c.Param("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "页码必须是数字"})
		return
	}

	var req struct {
		VoiceType string `json:"voiceType"`
		Format    string `json:"audioFormat"`
		Provider  string `json:"provider"`
	}
	_ = c.ShouldBindJSON(&req)

	payload, err := ensurePlaybackAudioAssets(c.Request.Context(), h.db, h.aiClient, courseID, pageNum, req.VoiceType, req.Format, req.Provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "音频元数据已生成", "data": payload})
}

func (h *TeacherHandler) PublishCourseware(c *gin.Context) {
	courseID := c.Param("courseId")
	var req struct {
		CourseID string `json:"courseId"`
		Scope    string `json:"scope"`
	}
	_ = c.ShouldBindJSON(&req)
	if courseID == "" {
		courseID = req.CourseID
	}
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 courseId"})
		return
	}

	var course model.Course
	if err := h.db.First(&course, "id = ?", courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "课件不存在"})
		return
	}

	scope := strings.TrimSpace(req.Scope)
	if scope == "" {
		scope = "all"
	}
	now := time.Now()
	if err := h.db.Model(&course).Updates(map[string]any{"is_published": true, "publish_scope": scope, "published_at": now}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "发布状态写入失败"})
		return
	}

	logger.Infof("课件发布成功: courseId=%s, scope=%s, title=%s", courseID, scope, course.Title)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "发布成功", "data": gin.H{"courseId": courseID, "scope": scope, "publishedAt": now.Format(time.RFC3339)}})
}

func (h *TeacherHandler) GetStudentStats(c *gin.Context) {
	courseID := c.Param("courseId")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": h.buildClassStats(courseID)})
}

func (h *TeacherHandler) GetClassStats(c *gin.Context) {
	courseID := c.Param("courseId")
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": h.buildClassStats(courseID)})
}

func (h *TeacherHandler) GetQuestionRecords(c *gin.Context) {
	courseID := c.Param("courseId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	type questionRecord struct {
		ID          string    `json:"id"`
		UserID      string    `json:"user_id"`
		PageIndex   int       `json:"page_index"`
		Question    string    `json:"question"`
		Answer      string    `json:"answer"`
		NeedReteach bool      `json:"need_reteach"`
		CreatedAt   time.Time `json:"created_at"`
	}

	var total int64
	h.db.Model(&model.DialogueTurn{}).Where("course_id = ?", courseID).Count(&total)

	list := make([]questionRecord, 0)
	if total > 0 {
		_ = h.db.Model(&model.DialogueTurn{}).
			Select("id, user_id, page_index, question, answer, need_reteach, created_at").
			Where("course_id = ?", courseID).
			Order("created_at desc").
			Offset(offset).
			Limit(pageSize).
			Scan(&list).Error
	} else {
		h.db.Model(&model.QuestionLog{}).Where("course_id = ?", courseID).Count(&total)
		_ = h.db.Model(&model.QuestionLog{}).
			Select("id, user_id, page_index, question, answer, created_at").
			Where("course_id = ?", courseID).
			Order("created_at desc").
			Offset(offset).
			Limit(pageSize).
			Scan(&list).Error
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"list": list, "total": total, "page": page, "pageSize": pageSize}})
}

func (h *TeacherHandler) GetCardData(c *gin.Context) {
	courseID := c.Param("courseId")
	var course model.Course
	if err := h.db.First(&course, "id = ?", courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "课件不存在"})
		return
	}

	stats := h.buildTeacherPageStats(courseID, course.TotalPage)

	topCandidates := append([]teacherPageStat(nil), stats...)
	sort.Slice(topCandidates, func(i, j int) bool {
		return topCandidates[i].CardIndex > topCandidates[j].CardIndex
	})
	topPages := make([]gin.H, 0, minTeacherInt(5, len(topCandidates)))
	for i := 0; i < minTeacherInt(5, len(topCandidates)); i++ {
		stat := topCandidates[i]
		ratio := stat.CardIndex * 10
		if ratio > 100 {
			ratio = 100
		}
		topPages = append(topPages, gin.H{"page": stat.PageIndex, "value": stat.QuestionCount, "ratio": ratio, "needReteachCount": stat.NeedReteachCount, "sessionCount": stat.SessionCount})
	}

	classStats := h.buildClassStats(courseID)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"pageStats": stats, "topPages": topPages, "totalQuestions": classStats["totalQuestions"], "reteachCount": classStats["reteachCount"], "activeSessions": classStats["activeSessions"]}})
}

func (h *TeacherHandler) upsertScript(courseID string, pageNum int, content string) error {
	return upsertTeacherScriptWithDB(h.db, courseID, pageNum, content)
}

func upsertTeacherScriptWithDB(db *gorm.DB, courseID string, pageNum int, content string) error {
	var coursePage model.CoursePage
	err := db.Where("course_id = ? AND page_index = ?", courseID, pageNum).First(&coursePage).Error
	if err != nil {
		return db.Create(&model.CoursePage{CourseID: courseID, PageIndex: pageNum, ScriptText: content}).Error
	}
	return db.Model(&coursePage).Update("script_text", content).Error
}

func (h *TeacherHandler) buildClassStats(courseID string) gin.H {
	var course model.Course
	_ = h.db.Select("id", "total_page").First(&course, "id = ?", courseID).Error
	stats := h.buildTeacherPageStats(courseID, course.TotalPage)

	type questionFreq struct {
		Page  int `json:"page"`
		Count int `json:"count"`
	}

	questionFreqs := make([]questionFreq, 0, len(stats))
	pageStats := make([]gin.H, 0, len(stats))
	pageStayTime := make([]gin.H, 0, len(stats))
	hotPages := make([]gin.H, 0)
	totalQuestions := 0
	totalReteach := 0
	activeSessions := 0
	for _, item := range stats {
		questionFreqs = append(questionFreqs, questionFreq{Page: item.PageIndex, Count: item.QuestionCount})
		pageStats = append(pageStats, gin.H{
			"page":              item.PageIndex,
			"count":             item.QuestionCount,
			"questionCount":     item.QuestionCount,
			"dialogueCount":     item.DialogueCount,
			"sessionCount":      item.SessionCount,
			"needReteachCount":  item.NeedReteachCount,
			"stayTime":          item.StayTime,
			"cardIndex":         item.CardIndex,
			"avgTurns":          item.AvgTurns,
			"estimatedDuration": item.BaseDuration,
		})
		pageStayTime = append(pageStayTime, gin.H{"page": item.PageIndex, "seconds": item.StayTime})
		totalQuestions += item.QuestionCount
		totalReteach += item.NeedReteachCount
		activeSessions += item.SessionCount
	}

	topCandidates := append([]teacherPageStat(nil), stats...)
	sort.Slice(topCandidates, func(i, j int) bool {
		if topCandidates[i].QuestionCount == topCandidates[j].QuestionCount {
			return topCandidates[i].CardIndex > topCandidates[j].CardIndex
		}
		return topCandidates[i].QuestionCount > topCandidates[j].QuestionCount
	})
	for i := 0; i < minTeacherInt(3, len(topCandidates)); i++ {
		if topCandidates[i].QuestionCount == 0 && topCandidates[i].SessionCount == 0 {
			continue
		}
		hotPages = append(hotPages, gin.H{"page": topCandidates[i].PageIndex, "count": topCandidates[i].QuestionCount, "cardIndex": topCandidates[i].CardIndex})
	}

	questions := h.collectTeacherQuestions(courseID)
	keywords := generateKeywordStats(questions)
	activeUsers := h.countDistinctUsers(courseID)
	avgTurnsPerSession := 0.0
	if activeSessions > 0 {
		avgTurnsPerSession = roundTeacherFloat(float64(totalQuestions) / float64(activeSessions))
	}

	return gin.H{
		"pageStayTime":       pageStayTime,
		"questionFreq":       questionFreqs,
		"wordCloud":          keywords,
		"pageStats":          pageStats,
		"keywords":           keywords,
		"hotPages":           hotPages,
		"totalQuestions":     totalQuestions,
		"totalDialogueTurns": totalQuestions,
		"reteachCount":       totalReteach,
		"activeSessions":     activeSessions,
		"activeUsers":        activeUsers,
		"avgTurnsPerSession": avgTurnsPerSession,
	}
}

type teacherNodeUpsertRequest struct {
	ID                string `json:"id"`
	NodeID            string `json:"nodeId"`
	Title             string `json:"title"`
	Summary           string `json:"summary"`
	ScriptText        string `json:"scriptText"`
	ReteachScript     string `json:"reteachScript"`
	TransitionText    string `json:"transitionText"`
	EstimatedDuration int    `json:"estimatedDuration"`
	SortOrder         int    `json:"sortOrder"`
}

type teacherPageStat struct {
	PageIndex        int     `json:"page"`
	QuestionCount    int     `json:"questionCount"`
	DialogueCount    int     `json:"dialogueCount"`
	NeedReteachCount int     `json:"needReteachCount"`
	SessionCount     int     `json:"sessionCount"`
	BaseDuration     int     `json:"estimatedDuration"`
	StayTime         float64 `json:"stayTime"`
	CardIndex        float64 `json:"cardIndex"`
	AvgTurns         float64 `json:"avgTurns"`
}

type teacherQuestionStat struct {
	PageIndex     int `gorm:"column:page_index"`
	QuestionCount int `gorm:"column:question_count"`
}

type teacherDialogueStat struct {
	PageIndex        int `gorm:"column:page_index"`
	DialogueCount    int `gorm:"column:dialogue_count"`
	NeedReteachCount int `gorm:"column:need_reteach_count"`
	SessionCount     int `gorm:"column:session_count"`
}

type teacherDurationStat struct {
	PageIndex    int `gorm:"column:page_index"`
	BaseDuration int `gorm:"column:base_duration"`
}

func (h *TeacherHandler) replaceTeachingNodes(courseID string, pageNum int, rawNodes []teacherNodeUpsertRequest) (string, []model.TeachingNode, error) {
	trimmed := make([]teacherNodeUpsertRequest, 0, len(rawNodes))
	for index, node := range rawNodes {
		text := strings.TrimSpace(node.ScriptText)
		title := strings.TrimSpace(node.Title)
		summary := strings.TrimSpace(node.Summary)
		if text == "" && title == "" && summary == "" {
			continue
		}
		node.SortOrder = index + 1
		node.ScriptText = text
		node.Title = title
		node.Summary = summary
		trimmed = append(trimmed, node)
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var existing []model.TeachingNode
		if err := tx.Where("course_id = ? AND page_index = ?", courseID, pageNum).Find(&existing).Error; err != nil {
			return err
		}
		existingByID := make(map[string]model.TeachingNode, len(existing))
		existingByNodeID := make(map[string]model.TeachingNode, len(existing))
		keepIDs := make(map[string]struct{}, len(trimmed))
		for _, item := range existing {
			existingByID[item.ID] = item
			existingByNodeID[item.NodeID] = item
		}

		for index, item := range trimmed {
			persisted := model.TeachingNode{}
			if item.ID != "" {
				persisted = existingByID[item.ID]
			}
			if persisted.ID == "" && item.NodeID != "" {
				persisted = existingByNodeID[item.NodeID]
			}
			if persisted.ID == "" {
				persisted = model.TeachingNode{
					CourseID:     courseID,
					PageIndex:    pageNum,
					NodeID:       buildTeacherNodeID(pageNum, index+1, item.NodeID),
					Title:        defaultTeacherNodeTitle(pageNum, index+1, item.Title),
					ChapterTitle: "第" + strconv.Itoa(pageNum) + "页",
				}
			}
			persisted.CourseID = courseID
			persisted.PageIndex = pageNum
			persisted.NodeID = buildTeacherNodeID(pageNum, index+1, firstTeacherNonEmpty(item.NodeID, persisted.NodeID))
			persisted.Title = defaultTeacherNodeTitle(pageNum, index+1, firstTeacherNonEmpty(item.Title, persisted.Title))
			persisted.Summary = item.Summary
			persisted.ScriptText = item.ScriptText
			persisted.ReteachScript = strings.TrimSpace(item.ReteachScript)
			persisted.TransitionText = strings.TrimSpace(item.TransitionText)
			persisted.EstimatedDuration = normalizeTeacherDuration(item.EstimatedDuration, persisted.ScriptText, persisted.Summary)
			persisted.SortOrder = index + 1

			if persisted.ID == "" {
				if err := tx.Create(&persisted).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(&model.TeachingNode{}).Where("id = ?", persisted.ID).Updates(map[string]any{
					"course_id":          persisted.CourseID,
					"page_index":         persisted.PageIndex,
					"node_id":            persisted.NodeID,
					"title":              persisted.Title,
					"summary":            persisted.Summary,
					"script_text":        persisted.ScriptText,
					"reteach_script":     persisted.ReteachScript,
					"transition_text":    persisted.TransitionText,
					"estimated_duration": persisted.EstimatedDuration,
					"sort_order":         persisted.SortOrder,
					"chapter_title":      persisted.ChapterTitle,
				}).Error; err != nil {
					return err
				}
			}
			keepIDs[persisted.ID] = struct{}{}
		}

		deleteIDs := make([]string, 0)
		for _, item := range existing {
			if _, ok := keepIDs[item.ID]; !ok {
				deleteIDs = append(deleteIDs, item.ID)
			}
		}
		if len(deleteIDs) > 0 {
			if err := tx.Where("id IN ?", deleteIDs).Delete(&model.TeachingNode{}).Error; err != nil {
				return err
			}
		}

		content := buildTeacherPageScriptFromNodes(trimmed)
		if err := upsertTeacherScriptWithDB(tx, courseID, pageNum, content); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", nil, err
	}

	savedNodes := loadTeachingNodesByPage(h.db, courseID, pageNum)
	return buildPageContextFromTeachingNodes(savedNodes), savedNodes, nil
}

func buildTeacherNodePayload(nodes []model.TeachingNode) []gin.H {
	result := make([]gin.H, 0, len(nodes))
	for index, node := range nodes {
		result = append(result, gin.H{
			"id":                node.ID,
			"nodeId":            node.NodeID,
			"page":              node.PageIndex,
			"type":              teacherNodeType(index, len(nodes)),
			"title":             node.Title,
			"summary":           node.Summary,
			"scriptText":        node.ScriptText,
			"reteachScript":     node.ReteachScript,
			"transitionText":    node.TransitionText,
			"estimatedDuration": playbackDurationSec(node),
			"sortOrder":         node.SortOrder,
		})
	}
	return result
}

func (h *TeacherHandler) buildTeacherPageStats(courseID string, totalPages int) []teacherPageStat {
	pageStats := make(map[int]*teacherPageStat)
	for page := 1; page <= maxCoursePage(totalPages); page++ {
		pageStats[page] = &teacherPageStat{PageIndex: page}
	}

	var questionStats []teacherQuestionStat
	_ = h.db.Table("question_logs").Select("page_index, count(*) as question_count").Where("course_id = ?", courseID).Group("page_index").Scan(&questionStats).Error
	for _, item := range questionStats {
		stat := ensureTeacherPageStat(pageStats, item.PageIndex)
		stat.QuestionCount = item.QuestionCount
	}

	var dialogueStats []teacherDialogueStat
	_ = h.db.Table("dialogue_turns").Select("page_index, count(*) as dialogue_count, sum(case when need_reteach then 1 else 0 end) as need_reteach_count, count(distinct session_id) as session_count").Where("course_id = ?", courseID).Group("page_index").Scan(&dialogueStats).Error
	for _, item := range dialogueStats {
		stat := ensureTeacherPageStat(pageStats, item.PageIndex)
		stat.DialogueCount = item.DialogueCount
		stat.NeedReteachCount = item.NeedReteachCount
		stat.SessionCount = item.SessionCount
		if stat.QuestionCount < item.DialogueCount {
			stat.QuestionCount = item.DialogueCount
		}
	}

	var durationStats []teacherDurationStat
	_ = h.db.Table("teaching_nodes").Select("page_index, coalesce(sum(case when estimated_duration > 0 then estimated_duration else 0 end), 0) as base_duration").Where("course_id = ?", courseID).Group("page_index").Scan(&durationStats).Error
	for _, item := range durationStats {
		stat := ensureTeacherPageStat(pageStats, item.PageIndex)
		stat.BaseDuration = item.BaseDuration
	}

	result := make([]teacherPageStat, 0, len(pageStats))
	for _, item := range pageStats {
		base := item.BaseDuration
		if base <= 0 {
			base = 30
		}
		sessionCount := item.SessionCount
		if sessionCount <= 0 {
			sessionCount = 1
		}
		engagementBoost := float64(item.DialogueCount*18+item.NeedReteachCount*35+item.QuestionCount*12) / float64(sessionCount)
		item.StayTime = roundTeacherFloat(float64(base) + engagementBoost)
		item.AvgTurns = roundTeacherFloat(float64(item.DialogueCount) / float64(sessionCount))
		item.CardIndex = roundTeacherCardIndex(item)
		result = append(result, *item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].PageIndex < result[j].PageIndex
	})
	return result
}

func ensureTeacherPageStat(pageStats map[int]*teacherPageStat, page int) *teacherPageStat {
	if page <= 0 {
		page = 1
	}
	if stat, ok := pageStats[page]; ok {
		return stat
	}
	stat := &teacherPageStat{PageIndex: page}
	pageStats[page] = stat
	return stat
}

func (h *TeacherHandler) collectTeacherQuestions(courseID string) []string {
	questions := make([]string, 0)
	logged := make([]string, 0)
	fromDialogue := make([]string, 0)
	_ = h.db.Table("question_logs").Where("course_id = ?", courseID).Pluck("question", &logged).Error
	_ = h.db.Table("dialogue_turns").Where("course_id = ?", courseID).Pluck("question", &fromDialogue).Error
	questions = append(questions, logged...)
	questions = append(questions, fromDialogue...)
	return questions
}

func (h *TeacherHandler) countDistinctUsers(courseID string) int64 {
	userSet := make(map[string]struct{})
	ids := make([]string, 0)
	_ = h.db.Table("question_logs").Where("course_id = ?", courseID).Distinct().Pluck("user_id", &ids).Error
	for _, item := range ids {
		if strings.TrimSpace(item) != "" {
			userSet[item] = struct{}{}
		}
	}
	ids = ids[:0]
	_ = h.db.Table("dialogue_turns").Where("course_id = ?", courseID).Distinct().Pluck("user_id", &ids).Error
	for _, item := range ids {
		if strings.TrimSpace(item) != "" {
			userSet[item] = struct{}{}
		}
	}
	return int64(len(userSet))
}

func buildTeacherPageScriptFromNodes(nodes []teacherNodeUpsertRequest) string {
	parts := make([]string, 0, len(nodes))
	for _, item := range nodes {
		text := strings.TrimSpace(item.ScriptText)
		if text == "" {
			text = strings.TrimSpace(item.Summary)
		}
		if text != "" {
			parts = append(parts, text)
		}
	}
	return strings.Join(parts, "\n\n")
}

func teacherNodeType(index, total int) string {
	if total <= 1 || index == 0 {
		return "opening"
	}
	if index == total-1 {
		return "transition"
	}
	return "explain"
}

func buildTeacherNodeID(pageNum, index int, raw string) string {
	raw = strings.TrimSpace(raw)
	if raw != "" {
		return raw
	}
	return fmt.Sprintf("p%d_n%d", pageNum, index)
}

func defaultTeacherNodeTitle(pageNum, index int, title string) string {
	title = strings.TrimSpace(title)
	if title != "" {
		return title
	}
	return fmt.Sprintf("第%d页节点%d", pageNum, index)
}

func normalizeTeacherDuration(duration int, content ...string) int {
	if duration > 0 {
		return duration
	}
	text := strings.Join(content, " ")
	base := len([]rune(strings.TrimSpace(text))) / 14
	if base < 20 {
		return 20
	}
	if base > 90 {
		return 90
	}
	return base
}

func firstTeacherNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func roundTeacherFloat(value float64) float64 {
	return mathRound(value*10) / 10
}

func roundTeacherCardIndex(stat *teacherPageStat) float64 {
	baseDuration := stat.BaseDuration
	if baseDuration <= 0 {
		baseDuration = 30
	}
	score := float64(stat.QuestionCount)*1.4 + float64(stat.NeedReteachCount)*2.3 + float64(stat.SessionCount)*0.8 + (stat.StayTime/float64(baseDuration))*1.5
	if score > 10 {
		score = 10
	}
	return roundTeacherFloat(score)
}

func mathRound(value float64) float64 {
	if value < 0 {
		return float64(int64(value - 0.5))
	}
	return float64(int64(value + 0.5))
}

func generateMockScript(courseName string, page int) string {
	templates := []string{
		"## %s 第%d页：课程导入\n\n### 教学目标\n- 了解本章节的核心概念\n- 建立知识框架\n- 明确学习重点\n\n### 讲解内容\n本页主要介绍本章节的背景、核心概念及学习路径。",
		"## %s 第%d页：深入讲解\n\n### 知识要点\n- 概念定义\n- 关键原理\n- 使用场景\n\n### 教学提示\n可以结合图示、案例和公式帮助学生理解。",
		"## %s 第%d页：案例分析\n\n### 案例拆解\n1. 场景说明\n2. 关键步骤\n3. 常见误区\n4. 课堂提问点\n\n### 小结\n帮助学生把抽象概念与真实场景建立联系。",
		"## %s 第%d页：总结回顾\n\n### 重点回顾\n1. 核心概念\n2. 方法步骤\n3. 易错点\n\n### 课堂互动\n建议通过一道练习题检验学生掌握情况。",
	}
	return fmt.Sprintf(templates[(page-1)%len(templates)], courseName, page)
}

func generateKeywordStats(questions []string) []gin.H {
	commonKeywords := []string{"依赖注入", "IoC", "AOP", "Spring", "微服务", "分布式", "事务", "缓存", "数据库", "接口", "fillna", "interpolate", "dropna", "缺失值", "异常值"}
	stats := make([]gin.H, 0)
	for _, keyword := range commonKeywords {
		count := 0
		for _, question := range questions {
			if strings.Contains(question, keyword) {
				count++
			}
		}
		if count > 0 {
			stats = append(stats, gin.H{"word": keyword, "count": count})
		}
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i]["count"].(int) > stats[j]["count"].(int)
	})
	if len(stats) < 5 {
		for _, keyword := range []string{"概念", "原理", "实现", "应用", "区别"} {
			stats = append(stats, gin.H{"word": keyword, "count": 1})
			if len(stats) >= 5 {
				break
			}
		}
	}
	return stats
}

func parsePageParam(values ...string) (int, error) {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}
		return strconv.Atoi(value)
	}
	return 0, fmt.Errorf("missing page")
}

func publishStatus(published bool) string {
	if published {
		return "published"
	}
	return "draft"
}

func defaultTeacherMode(mode string) string {
	if strings.TrimSpace(mode) == "" {
		return "llm"
	}
	return mode
}

func maxCoursePage(page int) int {
	if page > 0 {
		return page
	}
	return 1
}

func minTeacherInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
