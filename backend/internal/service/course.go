package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/pkg/logger"
	"smart-teaching-backend/pkg/oss"
	"sort"
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

	// 打开上传的文件用于上传
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

		// 优先尝试基于原始文件生成真实预览图，失败时回退到占位图
		pagePreviewURLs := map[int]string{}
		tmpPath, tmpErr := saveUploadedFileToTemp(file)
		if tmpErr == nil {
			defer os.Remove(tmpPath)
			generated, genErr := s.generatePagePreviewImages(ctx, tmpPath, course.ID, strings.ToLower(filepath.Ext(file.Filename)))
			if genErr == nil {
				pagePreviewURLs = generated
			} else {
				logger.Errorf("生成课件预览图失败，使用占位图: %v", genErr)
			}
		} else {
			logger.Errorf("保存课件临时文件失败，使用占位预览图: %v", tmpErr)
		}

		course.TotalPage = parsed.TotalPages
		if err := tx.Save(course).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("更新课件信息失败: %w", err)
		}

		pages := make([]model.CoursePage, 0, len(parsed.ParsedPages))
		for _, page := range parsed.ParsedPages {
			imageURL := strings.TrimSpace(pagePreviewURLs[page.Page])
			if imageURL == "" {
				// 回退：使用占位预览图，确保前端预览可用
				imageURL = fmt.Sprintf("https://picsum.photos/seed/%s_%d/800/600", course.ID, page.Page)
			}
			pages = append(pages, model.CoursePage{
				CourseID:   course.ID,
				PageIndex:  page.Page,
				ImageURL:   imageURL,
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
			CourseID:          courseID,
			NodeID:            node.NodeID,
			ChapterTitle:      chapterByNodeID[node.NodeID],
			PageIndex:         pageIndex,
			EstimatedDuration: node.EstimatedDuration,
			Title:             node.Title,
			Summary:           node.Summary,
			SourcePages:       string(sourcePages),
			CorePoints:        string(corePoints),
			Examples:          string(examples),
			CommonConfusions:  string(commonConfusions),
			SortOrder:         index,
		})
	}

	return tx.Create(&nodes).Error
}

// saveUploadedFileToTemp 将上传的课件文件保存到本地临时路径，供后续预览图生成使用
func saveUploadedFileToTemp(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	tmpFile, err := os.CreateTemp("", "course_*"+ext)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, src); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

// generatePagePreviewImages 使用本地渲染工具生成每页 PNG 预览，并上传到 MinIO，返回 pageIndex 到 URL 的映射
func (s *courseService) generatePagePreviewImages(ctx context.Context, localPath string, courseID string, ext string) (map[int]string, error) {
	// 目前仅支持 PDF/PPTX，其它类型直接跳过
	if ext != ".pdf" && ext != ".pptx" && ext != ".ppt" {
		return map[int]string{}, nil
	}

	pdfPath := localPath
	// PPT/PPTX 先转换为 PDF
	if ext == ".pptx" || ext == ".ppt" {
		var err error
		pdfPath, err = convertPptxToPdf(localPath)
		if err != nil {
			return nil, fmt.Errorf("PPT 转 PDF 失败: %w", err)
		}
		defer os.Remove(pdfPath)
	}

	imgDir, err := os.MkdirTemp("", "course_imgs_*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(imgDir)

	base := filepath.Join(imgDir, "page")
	cmd := exec.CommandContext(ctx, "pdftoppm", "-png", pdfPath, base)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("pdftoppm 生成预览失败: %w, output=%s", err, string(output))
	}

	files, err := filepath.Glob(filepath.Join(imgDir, "page-*.png"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return map[int]string{}, nil
	}
	sort.Strings(files)

	result := make(map[int]string, len(files))
	for idx, imgPath := range files {
		pageIndex := idx + 1

		f, err := os.Open(imgPath)
		if err != nil {
			continue
		}
		info, err := f.Stat()
		if err != nil {
			f.Close()
			continue
		}

		objectName := fmt.Sprintf("courses/%s/previews/page-%d.png", courseID, pageIndex)
		url, err := s.ossClient.UploadFile(ctx, objectName, f, info.Size(), "image/png")
		f.Close()
		if err != nil {
			logger.Errorf("上传预览图失败: %v", err)
			continue
		}
		result[pageIndex] = url
	}
	return result, nil
}

// convertPptxToPdf 使用本地 LibreOffice 将 PPT/PPTX 转换为 PDF
func convertPptxToPdf(pptPath string) (string, error) {
	outDir := filepath.Dir(pptPath)
	cmd := exec.Command("soffice", "--headless", "--convert-to", "pdf", "--outdir", outDir, pptPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("LibreOffice 转 PDF 失败: %w, output=%s", err, string(output))
	}
	pdfPath := strings.TrimSuffix(pptPath, filepath.Ext(pptPath)) + ".pdf"
	if _, err := os.Stat(pdfPath); err != nil {
		return "", fmt.Errorf("未找到转换后的 PDF 文件: %s", pdfPath)
	}
	return pdfPath, nil
}
