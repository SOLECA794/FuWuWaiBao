package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/pkg/logger"
)

type StudentHandler struct {
	db          *gorm.DB
	redisClient *redis.Client
	ctx         context.Context
}

func NewStudentHandler(db *gorm.DB, redisClient *redis.Client) *StudentHandler {
	return &StudentHandler{
		db:          db,
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

// 6.1 开始学习会话
// POST /api/v1/student/session/start
func (h *StudentHandler) StartSession(c *gin.Context) {
	var req struct {
		UserID   string `json:"userId" binding:"required"`
		CourseID string `json:"courseId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 生成会话ID
	sessionId := "sess_" + uuid.New().String()

	// 保存到 Redis
	sessionKey := "session:" + sessionId
	err := h.redisClient.HSet(h.ctx, sessionKey, map[string]interface{}{
		"user_id":         req.UserID,
		"course_id":       req.CourseID,
		"current_page":    1,
		"current_node_id": "",
		"interrupted":     false,
		"created_at":      time.Now().Unix(),
	}).Err()

	if err != nil {
		logger.Errorf("保存会话失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建会话失败",
		})
		return
	}

	h.redisClient.Expire(h.ctx, sessionKey, 24*time.Hour)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data": gin.H{
			"sessionId": sessionId,
			"courseId":  req.CourseID,
		},
	})
}

// 6.2 上报播放进度
// POST /api/v1/student/progress/update
func (h *StudentHandler) UpdateProgress(c *gin.Context) {
	var req struct {
		SessionID string `json:"sessionId" binding:"required"`
		CourseID  string `json:"courseId" binding:"required"`
		Page      int    `json:"page" binding:"required"`
		NodeID    string `json:"nodeId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 更新 Redis 中的进度
	sessionKey := "session:" + req.SessionID
	err := h.redisClient.HSet(h.ctx, sessionKey, map[string]interface{}{
		"current_page":    req.Page,
		"current_node_id": req.NodeID,
		"updated_at":      time.Now().Unix(),
	}).Err()

	if err != nil {
		logger.Errorf("更新进度失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新进度失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
	})
}

// 5.1 获取学习断点
// GET /api/v1/student/coursewares/:courseId/breakpoint
func (h *StudentHandler) GetBreakpoint(c *gin.Context) {
	courseId := c.Param("courseId")
	userId := c.Query("userId")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少参数: userId",
		})
		return
	}

	var progress model.UserProgress
	err := h.db.Where("user_id = ? AND course_id = ?", userId, courseId).First(&progress).Error

	lastPageNum := 1
	if err == nil {
		lastPageNum = progress.LastPage
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data": gin.H{
			"lastPageNum": lastPageNum,
		},
	})
}

// 5.1 更新学习断点
// PUT /api/v1/student/coursewares/:courseId/breakpoint
func (h *StudentHandler) UpdateBreakpoint(c *gin.Context) {
	courseId := c.Param("courseId")

	var req struct {
		UserID  string `json:"userId" binding:"required"`
		PageNum int    `json:"pageNum" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	var progress model.UserProgress
	err := h.db.Where("user_id = ? AND course_id = ?", req.UserID, courseId).First(&progress).Error

	if err != nil {
		progress = model.UserProgress{
			UserID:   req.UserID,
			CourseID: courseId,
			LastPage: req.PageNum,
		}
		h.db.Create(&progress)
	} else {
		h.db.Model(&progress).Update("last_page", req.PageNum)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
	})
}

// 5.2 保存随堂笔记
// POST /api/v1/student/coursewares/:courseId/notes
func (h *StudentHandler) SaveNote(c *gin.Context) {
	courseId := c.Param("courseId")

	var req struct {
		UserID  string `json:"userId" binding:"required"`
		PageNum int    `json:"pageNum" binding:"required"`
		Content string `json:"content" binding:"required"`
		X       int    `json:"x"`
		Y       int    `json:"y"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	logger.Infof("保存笔记: userId=%s, courseId=%s, pageNum=%d", req.UserID, courseId, req.PageNum)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "保存成功",
		"data": gin.H{
			"status": "saved",
		},
	})
}

// 6.3 获取某页讲稿（学生播放用）
// GET /api/v1/student/coursewares/:courseId/pages/:pageNum
func (h *StudentHandler) GetCoursewarePage(c *gin.Context) {
	courseId := c.Param("courseId")
	pageNumStr := c.Param("pageNum")

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "页码必须是数字",
		})
		return
	}

	var course model.Course
	if err := h.db.First(&course, "id = ?", courseId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课件不存在",
		})
		return
	}

	var coursePage model.CoursePage
	err = h.db.Where("course_id = ? AND page_index = ?", courseId, pageNum).First(&coursePage).Error

	nodes := []gin.H{
		{"node_id": fmt.Sprintf("p%d_n1", pageNum), "type": "opening", "text": "欢迎学习本节课..."},
		{"node_id": fmt.Sprintf("p%d_n2", pageNum), "type": "explain", "text": coursePage.ScriptText},
		{"node_id": fmt.Sprintf("p%d_n3", pageNum), "type": "transition", "text": "接下来我们继续学习..."},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data": gin.H{
			"courseId":     courseId,
			"page":         pageNum,
			"nodes":        nodes,
			"page_summary": coursePage.ScriptText,
		},
	})
}

// 6.4 问答流式接口（核心）
// POST /api/v1/student/qa/stream
func (h *StudentHandler) QAStream(c *gin.Context) {
	var req struct {
		SessionID string `json:"sessionId" binding:"required"`
		CourseID  string `json:"courseId" binding:"required"`
		Page      int    `json:"page" binding:"required"`
		NodeID    string `json:"nodeId"`
		Question  string `json:"question" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 设置 SSE 头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// 标记会话为中断状态
	sessionKey := "session:" + req.SessionID
	h.redisClient.HSet(h.ctx, sessionKey, "interrupted", true)

	// 模拟流式输出
	answer := "这是一个关于" + req.Question + "的详细解答。"
	runes := []rune(answer)

	for i := 0; i < len(runes); i++ {
		if i%3 == 0 {
			token := string(runes[i])
			fmt.Fprintf(c.Writer, "event: token\ndata: {\"text\":\"%s\"}\n\n", token)
			c.Writer.Flush()
			time.Sleep(30 * time.Millisecond)
		}
	}

	fmt.Fprintf(c.Writer, "event: sentence\ndata: {\"text\":\"%s\"}\n\n", answer)
	c.Writer.Flush()

	final := map[string]interface{}{
		"need_reteach":   false,
		"source_page":    req.Page,
		"resume_page":    req.Page,
		"resume_node_id": req.NodeID,
	}
	finalJSON, _ := json.Marshal(final)
	fmt.Fprintf(c.Writer, "event: final\ndata: %s\n\n", finalJSON)
	c.Writer.Flush()

	// 更新 Redis
	h.redisClient.HSet(h.ctx, sessionKey, map[string]interface{}{
		"interrupted":    false,
		"resume_page":    req.Page,
		"resume_node_id": req.NodeID,
	})
}

// 6.3 学生端 - 个人微观学情
// GET /api/v1/student/coursewares/:courseId/stats
func (h *StudentHandler) GetPersonalStats(c *gin.Context) {
	courseId := c.Param("courseId")
	userId := c.Query("userId")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少参数: userId",
		})
		return
	}

	// 查询该学生的提问次数
	var totalQuestions int64
	h.db.Model(&model.QuestionLog{}).
		Where("user_id = ? AND course_id = ?", userId, courseId).
		Count(&totalQuestions)

	// 计算学习时长（模拟）
	studyHours := 2.5

	// 计算掌握情况（模拟）
	mastery := gin.H{
		"章节一": 85,
		"章节二": 70,
		"章节三": 90,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data": gin.H{
			"totalQuestions": totalQuestions,
			"studyHours":     studyHours,
			"mastery":        mastery,
		},
	})
}
