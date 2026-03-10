package learning_record

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// 注意替换 module name
	"smart-learning/pkg/response"
)

type UpdateBreakpointRequest struct {
	PageNum int `json:"pageNum" binding:"required,min=1"`
}

type SaveNoteRequest struct {
	PageNum int     `json:"pageNum" binding:"required,min=1"`
	Content string  `json:"content" binding:"required"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
}

// Handler 定义 HTTP 处理器结构体
type Handler struct {
	svc IService // 注入 Service
}

func NewHandler(svc IService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetBreakpoint(c *gin.Context) {
	courseID := c.Param("courseId")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, 400, "courseId 不能为空")
		return
	}

	data, err := h.svc.GetBreakpoint(c.Request.Context(), courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "获取学习断点失败")
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateBreakpoint(c *gin.Context) {
	courseID := c.Param("courseId")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, 400, "courseId 不能为空")
		return
	}

	var req UpdateBreakpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	err := h.svc.UpdateBreakpoint(c.Request.Context(), courseID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "更新学习断点失败")
		return
	}

	// 更新操作成功，返回 nil 即可（对应文档要求：只需 code 200 和 success message）
	response.Success(c, nil)
}

func (h *Handler) SaveNote(c *gin.Context) {
	courseID := c.Param("courseId")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, 400, "courseId 不能为空")
		return
	}

	var req SaveNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	data, err := h.svc.SaveNote(c.Request.Context(), courseID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "保存笔记失败")
		return
	}
	response.Success(c, data)
}
