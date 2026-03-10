package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

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

func NewCourseHandler(courseService service.CourseService, db *gorm.DB) *CourseHandler {
	return &CourseHandler{courseService: courseService, db: db}
}

func (h *CourseHandler) GetPagePreview(c *gin.Context) {
	courseID := c.Param("courseId")
	pageNum, err := strconv.Atoi(c.Param("pageNum"))
	if err != nil || pageNum < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "页码必须是正整数"})
		return
	}

	var coursePage model.CoursePage
	if err := h.db.Where("course_id = ? AND page_index = ?", courseID, pageNum).First(&coursePage).Error; err != nil || strings.TrimSpace(coursePage.ImageURL) == "" {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预览图不存在"})
		return
	}

	c.Redirect(http.StatusFound, coursePage.ImageURL)
}

func (h *CourseHandler) UploadCourse(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择要上传的文件"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExts := map[string]bool{".pdf": true, ".pptx": true}
	if !validExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只支持 PDF/PPTX 格式"})
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	if title == "" {
		title = strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
	}

	course, err := h.courseService.UploadCourse(c.Request.Context(), file, title)
	if err != nil {
		logger.Errorf("上传课件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "上传失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "上传成功", "data": course})
}

func (h *CourseHandler) GetCourse(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.Param("courseId")
	}

	course, err := h.courseService.GetCourse(id)
	if err != nil {
		logger.Errorf("获取课件失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "课件不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": course})
}

func (h *CourseHandler) GetCoursePages(c *gin.Context) {
	courseID := c.Param("id")
	if courseID == "" {
		courseID = c.Param("courseId")
	}

	pages, err := h.courseService.GetCoursePages(courseID)
	if err != nil {
		logger.Errorf("获取课件页面失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": pages})
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.Param("courseId")
	}

	if err := h.courseService.DeleteCourse(id); err != nil {
		logger.Errorf("删除课件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}
