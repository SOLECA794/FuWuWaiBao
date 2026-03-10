package ai_partner

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"smart-learning/internal/dto"
	"smart-learning/pkg/response"
)

// AskRequest 定义智能多模态答疑的请求体
type AskRequest struct {
	PageNum    int             `json:"pageNum" binding:"required"`
	Type       string          `json:"type" binding:"required,oneof=text audio image"`
	Question   string          `json:"question" binding:"required"`
	TracePoint *dto.TracePoint `json:"tracePoint,omitempty"`
}

// Handler 定义 AI 学伴的 HTTP 处理器结构体
type Handler struct {
	svc IService // 依赖注入 Service 层
}

// NewHandler 是构造函数
func NewHandler(svc IService) *Handler {
	return &Handler{
		svc: svc,
	}
}

// GetKnowledgeGraph 获取课件知识图谱
func (h *Handler) GetKnowledgeGraph(c *gin.Context) {
	courseID := c.Param("courseId")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, 400, "courseId 不能为空")
		return
	}

	// 调用 Service 层
	graphData, err := h.svc.GetKnowledgeGraph(c.Request.Context(), courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "获取知识图谱失败")
		return
	}

	response.Success(c, graphData)
}

// AskQuestion 智能多模态答疑
func (h *Handler) AskQuestion(c *gin.Context) {
	startTime := time.Now()
	courseID := c.Param("courseId")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, 400, "courseId 不能为空")
		return
	}

	var req AskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	// 调用 Service 层处理核心业务
	answerData, err := h.svc.GenerateAnswer(c.Request.Context(), courseID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "AI 生成解答失败")
		return
	}

	processingTime := time.Since(startTime).Milliseconds()
	c.Header("X-Processing-Time", fmt.Sprintf("%dms", processingTime))

	response.Success(c, answerData)
}
