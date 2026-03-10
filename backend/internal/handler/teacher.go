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
	if content == "" {
		content = buildPageContextFromTeachingNodes(loadTeachingNodesByPage(h.db, courseID, pageNum))
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"courseId": courseID, "pageNum": pageNum, "page": pageNum, "content": content}})
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
	var total int64
	h.db.Model(&model.QuestionLog{}).Where("course_id = ?", courseID).Count(&total)

	var logs []model.QuestionLog
	h.db.Where("course_id = ?", courseID).Order("created_at desc").Offset(offset).Limit(pageSize).Find(&logs)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"list": logs, "total": total, "page": page, "pageSize": pageSize}})
}

func (h *TeacherHandler) GetCardData(c *gin.Context) {
	courseID := c.Param("courseId")
	var course model.Course
	if err := h.db.First(&course, "id = ?", courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "课件不存在"})
		return
	}

	type pageStat struct {
		PageIndex     int     `json:"page"`
		QuestionCount int     `json:"questionCount"`
		StayTime      float64 `json:"stayTime"`
		CardIndex     float64 `json:"cardIndex"`
	}

	stats := make([]pageStat, 0)
	rows, err := h.db.Table("question_logs").Select("page_index, count(*) as question_count").Where("course_id = ?", courseID).Group("page_index").Order("page_index").Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var stat pageStat
			_ = rows.Scan(&stat.PageIndex, &stat.QuestionCount)
			stat.StayTime = 0
			stat.CardIndex = float64(stat.QuestionCount)
			stats = append(stats, stat)
		}
	}
	if len(stats) == 0 {
		for page := 1; page <= maxCoursePage(course.TotalPage); page++ {
			stats = append(stats, pageStat{PageIndex: page, QuestionCount: 0, StayTime: 0, CardIndex: 0})
		}
	}

	topCandidates := append([]pageStat(nil), stats...)
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
		topPages = append(topPages, gin.H{"page": stat.PageIndex, "value": stat.QuestionCount, "ratio": ratio})
	}

	var totalQuestions int64
	h.db.Model(&model.QuestionLog{}).Where("course_id = ?", courseID).Count(&totalQuestions)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "请求成功", "data": gin.H{"pageStats": stats, "topPages": topPages, "totalQuestions": totalQuestions}})
}

func (h *TeacherHandler) upsertScript(courseID string, pageNum int, content string) error {
	var coursePage model.CoursePage
	err := h.db.Where("course_id = ? AND page_index = ?", courseID, pageNum).First(&coursePage).Error
	if err != nil {
		return h.db.Create(&model.CoursePage{CourseID: courseID, PageIndex: pageNum, ScriptText: content}).Error
	}
	return h.db.Model(&coursePage).Update("script_text", content).Error
}

func (h *TeacherHandler) buildClassStats(courseID string) gin.H {
	type questionFreq struct {
		Page  int `json:"page"`
		Count int `json:"count"`
	}

	var questionFreqs []questionFreq
	h.db.Table("question_logs").Select("page_index as page, count(*) as count").Where("course_id = ?", courseID).Group("page_index").Order("page").Scan(&questionFreqs)

	var totalQuestions int64
	h.db.Model(&model.QuestionLog{}).Where("course_id = ?", courseID).Count(&totalQuestions)
	var activeUsers int64
	h.db.Table("question_logs").Where("course_id = ?", courseID).Distinct("user_id").Count(&activeUsers)

	questions := make([]string, 0)
	h.db.Table("question_logs").Where("course_id = ?", courseID).Pluck("question", &questions)
	keywords := generateKeywordStats(questions)

	pageStats := make([]gin.H, 0, len(questionFreqs))
	for _, item := range questionFreqs {
		pageStats = append(pageStats, gin.H{"page": item.Page, "count": item.Count})
	}

	return gin.H{"pageStayTime": []gin.H{}, "questionFreq": questionFreqs, "wordCloud": keywords, "pageStats": pageStats, "keywords": keywords, "totalQuestions": totalQuestions, "activeUsers": activeUsers}
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
