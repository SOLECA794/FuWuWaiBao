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

type CourseHandler struct {
	courseService service.CourseService
	db            *gorm.DB
}

// NewCourseHandler 创建课件处理器
func NewCourseHandler(courseService service.CourseService, db *gorm.DB) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
		db:            db,
	}
}

// GetPagePreview 获取课件预览图片
// GET /api/v1/courseware/:courseId/page/:pageNum
func (h *CourseHandler) GetPagePreview(c *gin.Context) {
	courseId := c.Param("courseId")
	pageNumStr := c.Param("pageNum")

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil || pageNum < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "页码必须是正整数",
		})
		return
	}

	// 查询课件页面
	var coursePage model.CoursePage
	err = h.db.Where("course_id = ? AND page_index = ?", courseId, pageNum).First(&coursePage).Error

	if err != nil {
		// 如果查询出错，返回404
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "预览图不存在",
		})
		return
	}

	if coursePage.ImageURL == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "预览图不存在",
		})
		return
	}

	// 重定向到图片URL
	c.Redirect(http.StatusFound, coursePage.ImageURL)
}

// 上传课件
func (h *CourseHandler) UploadCourse(c *gin.Context) {
	// 这个方法保持不变，继续使用 h.courseService
	file, _ := c.FormFile("file")
	title := c.PostForm("title")

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

// 获取课件详情
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

// 获取课件页面列表
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

// 删除课件
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
