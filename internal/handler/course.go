package handler

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"smart-teaching-backend/internal/service"
	"smart-teaching-backend/pkg/logger"
)

type CourseHandler struct {
	courseService service.CourseService
}

func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// UploadCourse 上传课件
// POST /api/v1/courses/upload
func (h *CourseHandler) UploadCourse(c *gin.Context) {
	// 获取表单数据
	title := c.PostForm("title")
	if title == "" {
		title = "未命名课件"
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		logger.Errorf("获取上传文件失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择要上传的文件",
		})
		return
	}

	// 验证文件类型
	ext := filepath.Ext(file.Filename)
	validExts := map[string]bool{".pdf": true, ".ppt": true, ".pptx": true}
	if !validExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "只支持PDF/PPT/PPTX格式",
		})
		return
	}

	// 上传课件
	course, err := h.courseService.UploadCourse(c.Request.Context(), file, title)
	if err != nil {
		logger.Errorf("上传课件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data":    course,
	})
}

// GetCourse 获取课件详情
// GET /api/v1/courses/:id
func (h *CourseHandler) GetCourse(c *gin.Context) {
	id := c.Param("id")

	course, err := h.courseService.GetCourse(id)
	if err != nil {
		logger.Errorf("获取课件失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课件不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": course,
	})
}

// GetCoursePages 获取课件页面列表
// GET /api/v1/courses/:id/pages
func (h *CourseHandler) GetCoursePages(c *gin.Context) {
	courseID := c.Param("id")

	pages, err := h.courseService.GetCoursePages(courseID)
	if err != nil {
		logger.Errorf("获取课件页面失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": pages,
	})
}

// UpdatePageScript 更新页面讲稿
// PUT /api/v1/pages/:id/script
func (h *CourseHandler) UpdatePageScript(c *gin.Context) {
	pageID := c.Param("id")

	var req struct {
		Script string `json:"script" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if err := h.courseService.UpdatePageScript(pageID, req.Script); err != nil {
		logger.Errorf("更新讲稿失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// DeleteCourse 删除课件
// DELETE /api/v1/courses/:id
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")

	if err := h.courseService.DeleteCourse(id); err != nil {
		logger.Errorf("删除课件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
