package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/pkg/logger"
	"smart-teaching-backend/pkg/oss"
	"strings"

	"gorm.io/gorm"
)

type CourseService interface {
	UploadCourse(ctx context.Context, file *multipart.FileHeader, title string) (*model.Course, error)
	GetCourse(id string) (*model.Course, error)
	GetCoursePages(courseID string) ([]model.CoursePage, error)
	UpdatePageScript(pageID string, script string) error
	DeleteCourse(id string) error
}

type courseService struct {
	db        *gorm.DB
	ossClient *oss.MinioClient
}

func NewCourseService(db *gorm.DB, ossClient *oss.MinioClient) CourseService {
	return &courseService{
		db:        db,
		ossClient: ossClient,
	}
}

// UploadCourse 上传课件
func (s *courseService) UploadCourse(ctx context.Context, file *multipart.FileHeader, title string) (*model.Course, error) {
	// 开启事务
	tx := s.db.Begin()

	// 创建课件记录
	course := &model.Course{
		Title:    title,
		FileType: strings.TrimPrefix(filepath.Ext(file.Filename), "."),
	}

	if err := tx.Create(course).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建课件记录失败: %w", err)
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 生成对象名：courses/{course_id}/{filename}
	objectName := fmt.Sprintf("courses/%s/%s", course.ID, file.Filename)

	// 上传到MinIO
	fileURL, err := s.ossClient.UploadFile(ctx, objectName, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("上传文件到MinIO失败: %w", err)
	}

	// 更新课件记录的文件URL
	course.FileURL = fileURL
	if err := tx.Save(course).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新课件URL失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	logger.Infof("课件上传成功: %s, URL: %s", course.ID, fileURL)

	return course, nil
}

// GetCourse 获取课件信息
func (s *courseService) GetCourse(id string) (*model.Course, error) {
	var course model.Course
	err := s.db.Preload("Pages").First(&course, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("获取课件失败: %w", err)
	}
	return &course, nil
}

// GetCoursePages 获取课件所有页面
func (s *courseService) GetCoursePages(courseID string) ([]model.CoursePage, error) {
	var pages []model.CoursePage
	err := s.db.Where("course_id = ?", courseID).Order("page_index asc").Find(&pages).Error
	if err != nil {
		return nil, fmt.Errorf("获取课件页面失败: %w", err)
	}
	return pages, nil
}

// UpdatePageScript 更新页面讲稿
func (s *courseService) UpdatePageScript(pageID string, script string) error {
	return s.db.Model(&model.CoursePage{}).Where("id = ?", pageID).Update("script_text", script).Error
}

// DeleteCourse 删除课件
func (s *courseService) DeleteCourse(id string) error {
	// 开启事务
	tx := s.db.Begin()

	// 先删除关联的页面
	if err := tx.Where("course_id = ?", id).Delete(&model.CoursePage{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除课件
	if err := tx.Delete(&model.Course{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
