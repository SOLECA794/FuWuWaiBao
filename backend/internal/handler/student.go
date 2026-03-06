package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

type UpdateBreakpointRequest struct {
	StudentID string `json:"studentId" binding:"required"`
	CourseID  string `json:"courseId" binding:"required"`
	LastPage  int    `json:"lastPage" binding:"required"`
}

// SaveNoteRequest 保存笔记请求
type SaveNoteRequest struct {
	StudentID string `json:"studentId" binding:"required"`
	CourseID  string `json:"courseId" binding:"required"`
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

	aiResp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{
		Question:    req.Question,
		CurrentPage: req.PageNum,
		Context:     context,
		Mode:        "llm",
	})
	if err != nil {
		logger.Errorf("调用AI问答失败: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "AI服务暂不可用，请稍后重试",
		})
		return
	}

	answer := aiResp.Answer

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
			"answer":             answer,
			"sourcePage":         aiResp.SourcePage,
			"sourceExcerpt":      aiResp.SourceExcerpt,
			"needReteach":        aiResp.Intent.NeedReteach,
			"followUpSuggestion": aiResp.FollowUpSuggestion,
			"aiUnavailable":      false,
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

	var coursePage model.CoursePage
	context := ""
	if err := h.db.Where("course_id = ? AND page_index = ?", req.CourseID, req.PageNum).First(&coursePage).Error; err == nil {
		context = coursePage.ScriptText
	}

	traceQuestion := req.Question + "（圈选坐标: " + strconv.FormatFloat(req.X, 'f', 2, 64) + "," + strconv.FormatFloat(req.Y, 'f', 2, 64) + "）"
	aiResp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{
		Question:    traceQuestion,
		CurrentPage: req.PageNum,
		Context:     context,
		Mode:        "llm",
	})
	if err != nil {
		logger.Errorf("调用AI溯源问答失败: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "AI服务暂不可用，请稍后重试",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": gin.H{
			"answer":        aiResp.Answer,
			"sourcePage":    aiResp.SourcePage,
			"sourceExcerpt": aiResp.SourceExcerpt,
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

// UpdateStudentBreakpoint 更新学习断点
// PUT /api/student/breakpoint
func (h *StudentHandler) UpdateStudentBreakpoint(c *gin.Context) {
	var req UpdateBreakpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必填参数: studentId, courseId 或 lastPage",
		})
		return
	}

	if req.LastPage < 1 {
		req.LastPage = 1
	}

	var progress model.UserProgress
	err := h.db.Where("user_id = ? AND course_id = ?", req.StudentID, req.CourseID).First(&progress).Error
	if err == nil {
		if err := h.db.Model(&progress).Update("last_page", req.LastPage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新断点失败"})
			return
		}
	} else {
		newProgress := model.UserProgress{UserID: req.StudentID, CourseID: req.CourseID, LastPage: req.LastPage}
		if err := h.db.Create(&newProgress).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存断点失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "断点更新成功",
		"data": gin.H{
			"lastPageNum": req.LastPage,
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

	note := model.StudentNote{
		UserID:   req.StudentID,
		CourseID: req.CourseID,
		PageNum:  req.PageNum,
		Note:     req.Note,
	}
	if err := h.db.Create(&note).Error; err != nil {
		logger.Errorf("保存笔记失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "保存成功",
		"data": gin.H{
			"status": "saved",
			"noteId": note.ID,
		},
	})
}
