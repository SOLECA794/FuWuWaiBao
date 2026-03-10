package weak_points

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	// 注意替换为你实际的 module name
	"smart-learning/pkg/response"
)

type GenerateTestRequest struct {
	QuestionType string `json:"questionType" binding:"required,oneof=single multiple"`
}

type CheckAnswerRequest struct {
	UserAnswer string `json:"userAnswer" binding:"required"`
}

// Handler 定义 HTTP 处理器结构体
type Handler struct {
	svc IService // 注入 Service
}

func NewHandler(svc IService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetWeakPointsList(c *gin.Context) {
	courseID := c.Param("courseId")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, 400, "courseId 不能为空")
		return
	}

	data, err := h.svc.GetWeakPointsList(c.Request.Context(), courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "获取薄弱点失败")
		return
	}
	response.Success(c, data)
}

func (h *Handler) ExplainWeakPoint(c *gin.Context) {
	startTime := time.Now()
	weakPointID := c.Param("weakPointId")
	if weakPointID == "" {
		response.Error(c, http.StatusBadRequest, 400, "weakPointId 不能为空")
		return
	}

	data, err := h.svc.ExplainWeakPoint(c.Request.Context(), weakPointID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "AI 讲解生成失败")
		return
	}

	processingTime := time.Since(startTime).Milliseconds()
	c.Header("X-Processing-Time", fmt.Sprintf("%dms", processingTime))
	response.Success(c, data)
}

func (h *Handler) GenerateTest(c *gin.Context) {
	startTime := time.Now()
	weakPointID := c.Param("weakPointId")
	if weakPointID == "" {
		response.Error(c, http.StatusBadRequest, 400, "weakPointId 不能为空")
		return
	}

	var req GenerateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	data, err := h.svc.GenerateTest(c.Request.Context(), weakPointID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "AI 题目生成失败")
		return
	}

	processingTime := time.Since(startTime).Milliseconds()
	c.Header("X-Processing-Time", fmt.Sprintf("%dms", processingTime))
	response.Success(c, data)
}

func (h *Handler) CheckTestAnswer(c *gin.Context) {
	startTime := time.Now()
	questionID := c.Param("questionId")
	if questionID == "" {
		response.Error(c, http.StatusBadRequest, 400, "questionId 不能为空")
		return
	}

	var req CheckAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	data, err := h.svc.CheckTestAnswer(c.Request.Context(), questionID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "答案校验失败")
		return
	}

	processingTime := time.Since(startTime).Milliseconds()
	c.Header("X-Processing-Time", fmt.Sprintf("%dms", processingTime))
	response.Success(c, data)
}
