package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/service"
)

type WeakPointHandler struct {
	db       *gorm.DB
	aiClient service.AIEngine
}

func NewWeakPointHandler(db *gorm.DB, aiClient service.AIEngine) *WeakPointHandler {
	return &WeakPointHandler{db: db, aiClient: aiClient}
}

// ParseKnowledge 拆解知识点层级
// POST /api/ai/parseKnowledge
func (h *WeakPointHandler) ParseKnowledge(c *gin.Context) {
	var req struct {
		FileContent string `json:"fileContent" binding:"required"`
		FileType    string `json:"fileType" binding:"required"`
		StudentID   string `json:"studentId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	text := strings.TrimSpace(req.FileContent)
	if text == "" {
		text = "请基于上传文件进行知识点层级拆解"
	}

	aiResp, err := h.aiClient.ParseKnowledge(c.Request.Context(), service.ParseKnowledgeRequest{
		Text: text,
		Mode: "llm",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "知识点拆解失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"structure": aiResp.Structure,
		},
	})
}

// GetWeakPointList 获取薄弱点列表
// GET /api/weakPoint/getList?studentId=xxx&courseId=xxx
func (h *WeakPointHandler) GetWeakPointList(c *gin.Context) {
	studentId := c.Query("studentId")
	courseId := c.Query("courseId")

	if studentId == "" || courseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少参数: studentId 或 courseId",
		})
		return
	}

	type pageStat struct {
		PageIndex int
		Count     int
	}

	var stats []pageStat
	h.db.Table("question_logs").
		Select("page_index, count(*) as count").
		Where("user_id = ? AND course_id = ?", studentId, courseId).
		Group("page_index").
		Order("count desc").
		Limit(5).
		Scan(&stats)

	weakPoints := make([]gin.H, 0, len(stats))
	for _, stat := range stats {
		mastery := 100 - stat.Count*12
		if mastery < 0 {
			mastery = 0
		}
		weakPoints = append(weakPoints, gin.H{
			"name":    "第" + strconv.Itoa(stat.PageIndex) + "页知识点",
			"count":   stat.Count,
			"mastery": mastery,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": weakPoints,
	})
}

// GetWeakPointExplain 获取薄弱点讲解
// POST /api/weakPoint/getExplain
func (h *WeakPointHandler) GetWeakPointExplain(c *gin.Context) {
	var req struct {
		WeakPointName string `json:"weakPointName" binding:"required"`
		StudentID     string `json:"studentId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	prompt := "请用教学口吻讲解知识点：" + req.WeakPointName + "。输出结构：先给定义，再给2个例子。"
	aiResp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{
		Question:    prompt,
		CurrentPage: 1,
		Context:     req.WeakPointName,
		Mode:        "llm",
	})
	if err != nil || strings.TrimSpace(aiResp.Answer) == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "AI讲解暂不可用，请稍后重试",
		})
		return
	}

	data := gin.H{
		"title":   req.WeakPointName + " · 知识点讲解",
		"content": aiResp.Answer,
		"examples": []string{
			"例1：结合本课件当前页核心概念进行复述",
			"例2：尝试将概念应用到一个真实场景中",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// GetWeakPointTest 生成检测习题
// POST /api/weakPoint/getTest
func (h *WeakPointHandler) GetWeakPointTest(c *gin.Context) {
	var req struct {
		WeakPointName string `json:"weakPointName" binding:"required"`
		StudentID     string `json:"studentId" binding:"required"`
		QuestionType  string `json:"questionType"` // single/multiple
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	prompt := "请为薄弱点“" + req.WeakPointName + "”生成一道单选题，返回严格JSON：{\"content\":\"题干\",\"options\":[\"A.xxx\",\"B.xxx\",\"C.xxx\",\"D.xxx\"],\"answer\":\"A\",\"explanation\":\"解析\"}"
	aiResp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{
		Question:    prompt,
		CurrentPage: 1,
		Context:     req.WeakPointName,
		Mode:        "llm",
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "AI习题生成暂不可用，请稍后重试",
		})
		return
	}

	testData, parseErr := parseWeakPointTest(aiResp.Answer)
	if parseErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "AI习题返回格式异常",
		})
		return
	}

	optionsJSON, _ := json.Marshal(testData.Options)
	questionRecord := model.Question{
		QuestionType: req.QuestionType,
		Content:      testData.Content,
		Options:      string(optionsJSON),
		Answer:       testData.Answer,
		Explanation:  testData.Explanation,
		Difficulty:   2,
	}
	if strings.TrimSpace(questionRecord.QuestionType) == "" {
		questionRecord.QuestionType = "single"
	}
	h.db.Create(&questionRecord)

	data := gin.H{
		"questionId": questionRecord.ID,
		"content":    testData.Content,
		"type":       questionRecord.QuestionType,
		"options":    testData.Options,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

type weakPointTestPayload struct {
	Content     string   `json:"content"`
	Options     []string `json:"options"`
	Answer      string   `json:"answer"`
	Explanation string   `json:"explanation"`
}

func parseWeakPointTest(raw string) (*weakPointTestPayload, error) {
	cleaned := strings.TrimSpace(raw)
	if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
		cleaned = strings.TrimPrefix(cleaned, "```")
		cleaned = strings.TrimSuffix(cleaned, "```")
		cleaned = strings.TrimSpace(cleaned)
	}

	var payload weakPointTestPayload
	if err := json.Unmarshal([]byte(cleaned), &payload); err != nil {
		return nil, err
	}
	if strings.TrimSpace(payload.Content) == "" || len(payload.Options) < 2 || strings.TrimSpace(payload.Answer) == "" {
		return nil, errInvalidAIResponse
	}
	if strings.TrimSpace(payload.Explanation) == "" {
		payload.Explanation = "请回顾本题涉及的知识点后再尝试一次。"
	}
	return &payload, nil
}

var errInvalidAIResponse = errors.New("invalid AI response payload")

// CheckAnswer 校验答案
// POST /api/weakPoint/checkAnswer
func (h *WeakPointHandler) CheckAnswer(c *gin.Context) {
	var req struct {
		StudentID  string `json:"studentId" binding:"required"`
		QuestionID string `json:"questionId" binding:"required"`
		UserAnswer string `json:"userAnswer" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	var question model.Question
	if err := h.db.First(&question, "id = ?", req.QuestionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "题目不存在",
		})
		return
	}

	correctAnswer := strings.TrimSpace(question.Answer)
	userAnswer := strings.TrimSpace(req.UserAnswer)
	isCorrect := strings.EqualFold(userAnswer, correctAnswer)
	if !isCorrect && len(correctAnswer) == 1 {
		isCorrect = strings.HasPrefix(strings.ToUpper(userAnswer), strings.ToUpper(correctAnswer))
	}
	explanation := question.Explanation

	// 更新掌握度
	masteryDelta := 0
	if isCorrect {
		masteryDelta = 10
	} else {
		masteryDelta = -5
	}

	_ = h.db.Create(&model.AnswerRecord{
		StudentID:    req.StudentID,
		QuestionID:   req.QuestionID,
		UserAnswer:   req.UserAnswer,
		IsCorrect:    isCorrect,
		MasteryDelta: masteryDelta,
	}).Error

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"isCorrect":     isCorrect,
			"correctAnswer": correctAnswer,
			"explanation":   explanation,
			"masteryDelta":  masteryDelta,
		},
	})
}
