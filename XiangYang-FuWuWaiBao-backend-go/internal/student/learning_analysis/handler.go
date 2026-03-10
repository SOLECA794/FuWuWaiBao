package learning_analysis

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// 注意替换 module name
	"smart-learning/pkg/response"
)

// Handler 定义 HTTP 处理器结构体
type Handler struct {
	svc IService
}

func NewHandler(svc IService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetPersonalMicroStats(c *gin.Context) {
	courseID := c.Param("courseId")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, 400, "courseId 不能为空")
		return
	}

	data, err := h.svc.GetPersonalMicroStats(c.Request.Context(), courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "获取个人学情分析失败")
		return
	}
	response.Success(c, data)
}
