package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/pkg/logger"
)

type StudentHandler struct {
	db *gorm.DB
}

func NewStudentHandler(db *gorm.DB) *StudentHandler {
	return &StudentHandler{db: db}
}

// ==================== 请求参数结构体 ====================

// GetCoursewarePageRequest 获取课件页面请求
type GetCoursewarePageRequest struct {
	CourseID    string `json:"courseId" binding:"required"`
	CurrentPage int    `json:"currentPage" binding:"required"`
}

// AIQuestionRequest AI提问请求
type AIQuestionRequest struct {
	CourseID string `json:"courseId" binding:"required"`
	PageNum  int    `json:"pageNum" binding:"required"`
	Question string `json:"question" binding:"required"`
	Type     string `json:"type"` // text/voice/image
}

// TraceQuestionRequest 溯源提问请求
type TraceQuestionRequest struct {
	CourseID string  `json:"courseId" binding:"required"`
	PageNum  int     `json:"pageNum" binding:"required"`
	X        float64 `json:"x" binding:"required"`
	Y        float64 `json:"y" binding:"required"`
	Question string  `json:"question" binding:"required"`
}

// StudyDataRequest 学习数据请求
type StudyDataRequest struct {
	StudentID string `form:"studentId" binding:"required"`
	CourseID  string `form:"courseId" binding:"required"`
}

// BreakpointRequest 断点请求
type BreakpointRequest struct {
	StudentID string `form:"studentId" binding:"required"`
	CourseID  string `form:"courseId" binding:"required"`
}

// SaveNoteRequest 保存笔记请求
type SaveNoteRequest struct {
	StudentID string `json:"studentId" binding:"required"`
	PageNum   int    `json:"pageNum" binding:"required"`
	Note      string `json:"note" binding:"required"`
}

// ==================== 1. 获取课件页面 ====================

// GetCoursewarePage 获取指定页码的课件信息
// POST /api/student/courseware/page
func (h *StudentHandler) GetCoursewarePage(c *gin.Context) {
	var req GetCoursewarePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必填参数: courseId 或 currentPage",
		})
		return
	}

	// 查询课件信息
	var course model.Course
	if err := h.db.First(&course, "id = ?", req.CourseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课件不存在",
		})
		return
	}

	// 查询指定页面的讲稿
	var coursePage model.CoursePage
	err := h.db.Where("course_id = ? AND page_index = ?", req.CourseID, req.CurrentPage).First(&coursePage).Error

	content := ""
	if err == nil {
		content = coursePage.ScriptText
	}

	// TODO: 这里可以调用 AI 生成更丰富的内容
	data := gin.H{
		"courseId":    req.CourseID,
		"currentPage": req.CurrentPage,
		"content":     content,
		"title":       course.Title,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data":    data,
	})
}

// ==================== 2. 多模态提问 ====================

// AskAIQuestion 多模态提问
// POST /api/student/ai/question
func (h *StudentHandler) AskAIQuestion(c *gin.Context) {
	var req AIQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必填参数: courseId, pageNum 或 question",
		})
		return
	}

	// 获取当前页面的讲稿作为上下文
	var coursePage model.CoursePage
	context := ""
	if err := h.db.Where("course_id = ? AND page_index = ?", req.CourseID, req.PageNum).First(&coursePage).Error; err == nil {
		context = coursePage.ScriptText
	}

	// TODO: 调用 AI 大模型 API 获取答案
	// 这里先返回模拟数据
	answer := generateAIAnswer(req.Question, context)

	// 记录提问日志
	log := model.QuestionLog{
		UserID:    c.GetString("userId"), // 需要从 JWT 中获取
		CourseID:  req.CourseID,
		PageIndex: req.PageNum,
		Question:  req.Question,
		Answer:    answer,
	}
	h.db.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"answer": answer,
		},
	})
}

// ==================== 3. 溯源定位提问 ====================

// TraceAIQuestion 溯源定位提问
// POST /api/student/ai/traceQuestion
func (h *StudentHandler) TraceAIQuestion(c *gin.Context) {
	var req TraceQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必填参数: courseId, pageNum, x, y 或 question",
		})
		return
	}

	// TODO: 根据坐标定位到具体内容，调用 AI 生成针对性回答
	answer := generateTraceAnswer(req.Question, req.X, req.Y)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"answer": answer,
		},
	})
}

// ==================== 4. 获取学习数据 ====================

// GetStudentStudyData 获取学生学习数据
// GET /api/student/studyData
func (h *StudentHandler) GetStudentStudyData(c *gin.Context) {
	var req StudyDataRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必填参数: studentId 或 courseId",
		})
		return
	}

	// 查询该学生的提问记录
	var totalQuestions int64
	h.db.Model(&model.QuestionLog{}).
		Where("user_id = ? AND course_id = ?", req.StudentID, req.CourseID).
		Count(&totalQuestions)

	// 查询提问最多的页面（薄弱点）
	type WeakPoint struct {
		PageIndex int `json:"pageIndex"`
		Count     int `json:"count"`
	}
	var weakPoints []WeakPoint
	h.db.Table("question_logs").
		Select("page_index, count(*) as count").
		Where("user_id = ? AND course_id = ?", req.StudentID, req.CourseID).
		Group("page_index").
		Order("count desc").
		Limit(3).
		Scan(&weakPoints)

	// 计算专注度（简单算法：提问越少专注度越高）
	focusScore := 85
	if totalQuestions > 10 {
		focusScore = 70
	} else if totalQuestions > 5 {
		focusScore = 80
	}

	// 格式化薄弱点
	weakPointList := make([]string, 0)
	for _, wp := range weakPoints {
		weakPointList = append(weakPointList, "第"+strconv.Itoa(wp.PageIndex)+"页")
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"studentId":      req.StudentID,
			"courseId":       req.CourseID,
			"focusScore":     focusScore,
			"weakPoints":     weakPointList,
			"totalQuestions": totalQuestions,
		},
	})
}

// ==================== 5. 获取断点 ====================

// GetStudentBreakpoint 获取学习断点
// GET /api/student/breakpoint
func (h *StudentHandler) GetStudentBreakpoint(c *gin.Context) {
	var req BreakpointRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必填参数: studentId 或 courseId",
		})
		return
	}

	// 查询学生进度
	var progress model.UserProgress
	err := h.db.Where("user_id = ? AND course_id = ?", req.StudentID, req.CourseID).First(&progress).Error

	lastPageNum := 1
	if err == nil {
		lastPageNum = progress.LastPage
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"lastPageNum": lastPageNum,
		},
	})
}

// ==================== 6. 保存笔记 ====================

// SaveStudentNote 保存学生笔记
// POST /api/student/saveNote
func (h *StudentHandler) SaveStudentNote(c *gin.Context) {
	var req SaveNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必填参数: studentId, pageNum 或 note",
		})
		return
	}

	// TODO: 创建笔记表 model.StudentNote 并保存
	// 这里先用日志记录
	logger.Infof("保存笔记: studentId=%s, pageNum=%d, note=%s",
		req.StudentID, req.PageNum, req.Note)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "保存成功",
		"data": gin.H{
			"status": "saved",
		},
	})
}

// ==================== 辅助函数 ====================

// 生成 AI 回答（模拟）
func generateAIAnswer(question, context string) string {
	// TODO: 接入真实 AI 服务
	return "这是AI基于上下文生成的回答：" + question
}

// 生成溯源回答（模拟）
func generateTraceAnswer(question string, x, y float64) string {
	// TODO: 接入真实 AI 服务
	return "针对坐标(" + strconv.FormatFloat(x, 'f', 2, 64) + "," +
		strconv.FormatFloat(y, 'f', 2, 64) + ")的解答：" + question
}
