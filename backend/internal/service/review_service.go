package service

import (
	"fmt"
	"time"

	"smart-teaching-backend/internal/model"

	"gorm.io/gorm"
)

type ReviewService struct {
	db *gorm.DB
}

func NewReviewService(db *gorm.DB) *ReviewService {
	return &ReviewService{db: db}
}

// GenerateReviewPackage 生成专项复习包
func (rs *ReviewService) GenerateReviewPackage(studentID, courseID string) (*model.ReviewPackage, error) {
	// 1. 查找薄弱知识点 (掌握度 < 0.6)
	var weakPoints []model.StudentKnowledgeMap
	rs.db.Where("student_id = ? AND mastery_score < ?", studentID, 0.6).
		Preload("KnowledgePoint").
		Find(&weakPoints)

	if len(weakPoints) == 0 {
		return nil, fmt.Errorf("当前学习状态良好，暂无薄弱知识点需要复习")
	}

	// 2. 创建复习包主记录
	pkg := model.ReviewPackage{
		StudentID:   studentID,
		Name:        fmt.Sprintf("%s 专项提升复习包", time.Now().Format("2006-01-02")),
		Description: "基于近期错题与薄弱点自动生成",
		Status:      "draft",
		GeneratedAt: time.Now(),
	}
	if err := rs.db.Create(&pkg).Error; err != nil {
		return nil, err
	}

	// 3. 循环处理每个薄弱点，抽取题目
	sortOrder := 0
	for _, wp := range weakPoints {
		kpID := wp.KnowledgePointID

		// A. 错题重做 (从 AnswerRecord 中找该知识点下做错的题)
		var mistakes []model.AnswerRecord
		rs.db.Where("student_id = ? AND is_correct = ? AND knowledge_points LIKE ?",
			studentID, false, "%"+kpID+"%").Limit(3).Find(&mistakes)

		for _, m := range mistakes {
			item := model.ReviewPackageItem{
				PackageID:       pkg.ID,
				ItemType:        "error_retake",
				ContentID:       m.QuestionID,
				KnowledgeNodeID: kpID,
				SortOrder:       sortOrder,
			}
			rs.db.Create(&item)
			sortOrder++
		}

		// B. 类似题目练习 (找同知识点下还没做过的题)
		var similarQuestions []model.Question
		rs.db.Where("knowledge_point_id = ? AND id NOT IN (?)", kpID,
			rs.db.Table("answer_records").Select("question_id").Where("student_id = ?", studentID)).
			Limit(2).Find(&similarQuestions)

		for _, q := range similarQuestions {
			item := model.ReviewPackageItem{
				PackageID:       pkg.ID,
				ItemType:        "similar_practice",
				ContentID:       q.ID,
				KnowledgeNodeID: kpID,
				SortOrder:       sortOrder,
			}
			rs.db.Create(&item)
			sortOrder++
		}
	}

	return &pkg, nil
}

// GetPackageDetail 获取复习包详情
func (rs *ReviewService) GetPackageDetail(packageID string) (map[string]interface{}, error) {
	var pkg model.ReviewPackage
	if err := rs.db.First(&pkg, "id = ?", packageID).Error; err != nil {
		return nil, err
	}

	var items []model.ReviewPackageItem
	rs.db.Where("package_id = ?", packageID).Order("sort_order ASC").Find(&items)

	// 组装返回数据，包含题目具体内容
	var resultItems []map[string]interface{}
	for _, item := range items {
		var question model.Question
		rs.db.First(&question, "id = ?", item.ContentID)

		resultItems = append(resultItems, map[string]interface{}{
			"id":          item.ID,
			"type":        item.ItemType,
			"is_selected": item.IsSelected,
			"question":    question,
			"node_id":     item.KnowledgeNodeID,
		})
	}

	return map[string]interface{}{
		"package": pkg,
		"items":   resultItems,
	}, nil
}

// UpdatePackageItems 更新复习包自定义内容
func (rs *ReviewService) UpdatePackageItems(packageID string, items []model.ReviewPackageItem) error {
	return rs.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&model.ReviewPackageItem{}).
				Where("id = ?", item.ID).
				Updates(map[string]interface{}{
					"is_selected": item.IsSelected,
					"sort_order":  item.SortOrder,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ExportPackage 导出复习包 (此处模拟生成链接，实际可对接 PDF 库)
func (rs *ReviewService) ExportPackage(packageID, format string) (string, error) {
	// TODO: 这里可以调用 gofpdf 或 wkhtmltopdf 生成真实文件
	// 暂时返回一个模拟的 MinIO 路径
	exportPath := fmt.Sprintf("exports/review_%s.%s", packageID, format)

	// 更新数据库中的导出状态
	rs.db.Model(&model.ReviewPackage{}).Where("id = ?", packageID).Update("export_url", exportPath)

	return exportPath, nil
}
