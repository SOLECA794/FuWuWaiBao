package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"smart-teaching-backend/pkg/apiresp"
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
		apiresp.BadRequest(c, "页码必须是正整数", "")
		return
	}

	// 优先使用预生成的图片预览
	var coursePage model.CoursePage
	if err := h.db.Where("course_id = ? AND page_index = ?", courseID, pageNum).First(&coursePage).Error; err == nil {
		if url := strings.TrimSpace(coursePage.ImageURL); url != "" {
			c.Redirect(http.StatusFound, url)
			return
		}
	}

	// 如果没有预览图，则回退到原始课件文件：
	// - 对于 PDF：直接跳转到 PDF 文件并附带 #page=N，依赖浏览器内置 PDF 预览能力。
	// - 其它类型（PPT/PPTX 等）：暂时仅返回 404，由前端自行兜底。
	var course model.Course
	if err := h.db.First(&course, "id = ?", courseID).Error; err == nil {
		fileURL := strings.TrimSpace(course.FileURL)
		if fileURL != "" && strings.EqualFold(course.FileType, "pdf") {
			redirectURL := fileURL
			// 附带页码信息，主流浏览器的 PDF 查看器支持 #page=N
			if !strings.Contains(fileURL, "#") {
				redirectURL = fileURL + "#page=" + strconv.Itoa(pageNum)
			}
			c.Redirect(http.StatusFound, redirectURL)
			return
		}
	}

	apiresp.NotFound(c, "预览图不存在", "")
}

func (h *CourseHandler) UploadCourse(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		apiresp.BadRequest(c, "请选择要上传的文件", "")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExts := map[string]bool{".pdf": true, ".pptx": true}
	if !validExts[ext] {
		apiresp.BadRequest(c, "只支持 PDF/PPTX 格式", "")
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	if title == "" {
		title = strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
	}

	course, err := h.courseService.UploadCourse(c.Request.Context(), file, title)
	if err != nil {
		logger.Errorf("上传课件失败: %v", err)
		apiresp.Internal(c, "上传失败", err.Error())
		return
	}

	apiresp.OK(c, "上传成功", course)
}

func (h *CourseHandler) GetCourse(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.Param("courseId")
	}

	course, err := h.courseService.GetCourse(id)
	if err != nil {
		logger.Errorf("获取课件失败: %v", err)
		apiresp.NotFound(c, "课件不存在", "")
		return
	}

	apiresp.OK(c, "请求成功", course)
}

func (h *CourseHandler) GetCoursePages(c *gin.Context) {
	courseID := c.Param("id")
	if courseID == "" {
		courseID = c.Param("courseId")
	}

	pages, err := h.courseService.GetCoursePages(courseID)
	if err != nil {
		logger.Errorf("获取课件页面失败: %v", err)
		apiresp.Internal(c, "获取失败", "")
		return
	}

	apiresp.OK(c, "请求成功", pages)
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.Param("courseId")
	}

	if err := h.courseService.DeleteCourse(id); err != nil {
		logger.Errorf("删除课件失败: %v", err)
		apiresp.Internal(c, "删除失败", "")
		return
	}

	apiresp.OKMessage(c, "删除成功")
}
