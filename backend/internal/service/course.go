package service

import (
	"context"
	"encoding/json"
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
	aiClient  AIEngine
}

func NewCourseService(db *gorm.DB, ossClient *oss.MinioClient, aiClient AIEngine) CourseService {
	return &courseService{
		db:        db,
		ossClient: ossClient,
		aiClient:  aiClient,
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

	if s.aiClient != nil {
		parsed, err := s.aiClient.ParseDocument(ctx, file)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("解析课件失败: %w", err)
		}

		course.TotalPage = parsed.TotalPages
		if err := tx.Save(course).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("更新课件信息失败: %w", err)
		}

		pages := make([]model.CoursePage, 0, len(parsed.ParsedPages))
		for _, page := range parsed.ParsedPages {
			pages = append(pages, model.CoursePage{
				CourseID:   course.ID,
				PageIndex:  page.Page,
				SourceText: strings.TrimSpace(page.Content),
			})
		}
		if len(pages) > 0 {
			if err := tx.Create(&pages).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("保存解析页面失败: %w", err)
			}
		}

		reconstructed, err := s.aiClient.ReconstructDocument(ctx, ReconstructDocumentRequest{
			ParsedDocument: map[string]any{
				"doc_id":       parsed.DocID,
				"doc_name":     parsed.DocName,
				"doc_type":     parsed.DocType,
				"total_pages":  parsed.TotalPages,
				"parsed_pages": parsed.ParsedPages,
			},
			Mode: "hybrid",
		})
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("重构课件内容失败: %w", err)
		}
		if err := s.saveTeachingNodes(tx, course.ID, reconstructed); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("保存教学节点失败: %w", err)
		}
	} else if err := tx.Save(course).Error; err != nil {
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

	if err := tx.Where("course_id = ?", id).Delete(&model.TeachingNode{}).Error; err != nil {
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

func (s *courseService) saveTeachingNodes(tx *gorm.DB, courseID string, reconstructed *ReconstructDocumentResponse) error {
	if reconstructed == nil || len(reconstructed.TeachingNodes) == 0 {
		return nil
	}

	chapterByNodeID := map[string]string{}
	for _, chapter := range reconstructed.Chapters {
		for _, nodeID := range chapter.NodeIDs {
			chapterByNodeID[nodeID] = chapter.Title
		}
	}

	nodes := make([]model.TeachingNode, 0, len(reconstructed.TeachingNodes))
	for index, node := range reconstructed.TeachingNodes {
		sourcePages, _ := json.Marshal(node.SourcePages)
		corePoints, _ := json.Marshal(node.CorePoints)
		examples, _ := json.Marshal(node.Examples)
		commonConfusions, _ := json.Marshal(node.CommonConfusions)
		pageIndex := 0
		if len(node.SourcePages) > 0 {
			pageIndex = node.SourcePages[0]
		}

		nodes = append(nodes, model.TeachingNode{
			CourseID:         courseID,
			NodeID:           node.NodeID,
			ChapterTitle:     chapterByNodeID[node.NodeID],
			PageIndex:        pageIndex,
			Title:            node.Title,
			Summary:          node.Summary,
			SourcePages:      string(sourcePages),
			CorePoints:       string(corePoints),
			Examples:         string(examples),
			CommonConfusions: string(commonConfusions),
			SortOrder:        index,
		})
	}

	return tx.Create(&nodes).Error
}
