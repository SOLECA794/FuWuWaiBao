package service

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"smart-teaching-backend/internal/model"
)

// KnowledgeMapService 知识图谱服务
type KnowledgeMapService struct {
	db *gorm.DB
}

// NewKnowledgeMapService 创建知识图谱服务
func NewKnowledgeMapService(db *gorm.DB) *KnowledgeMapService {
	return &KnowledgeMapService{db: db}
}

// UpdateMasteryScore 更新知识点掌握度
func (kms *KnowledgeMapService) UpdateMasteryScore(studentID, knowledgePointID string, isCorrect bool, responseTime int) error {
	var knowledgeMap model.StudentKnowledgeMap

	// 查找或创建知识图谱记录
	err := kms.db.Where("student_id = ? AND knowledge_point_id = ?", studentID, knowledgePointID).
		First(&knowledgeMap).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		knowledgeMap = model.StudentKnowledgeMap{
			StudentID:        studentID,
			KnowledgePointID: knowledgePointID,
			MasteryScore:     0.5, // 初始掌握度
			ConfidenceLevel:  0.1, // 初始置信度
			AttemptCount:     0,
			CorrectCount:     0,
			LastUpdated:      time.Now(),
			LearningCurve:    "[]",
			StrengthAreas:    "[]",
			WeakAreas:        "[]",
		}
	}

	// 更新统计数据
	knowledgeMap.AttemptCount++
	if isCorrect {
		knowledgeMap.CorrectCount++
	}

	// 计算新的掌握度（使用Elo Rating类似算法）
	accuracy := float32(knowledgeMap.CorrectCount) / float32(knowledgeMap.AttemptCount)
	expectedScore := knowledgeMap.MasteryScore

	// 根据响应时间调整（响应越快，掌握度越高）
	timeBonus := float32(1.0)
	if responseTime < 5000 { // 5秒内回答
		timeBonus = 1.2
	} else if responseTime > 30000 { // 30秒以上
		timeBonus = 0.8
	}

	// 更新掌握度
	K := float32(0.3) // 学习率
	actualScore := accuracy * timeBonus
	knowledgeMap.MasteryScore = expectedScore + K*(actualScore-expectedScore)

	// 限制在0-1范围内
	if knowledgeMap.MasteryScore > 1.0 {
		knowledgeMap.MasteryScore = 1.0
	} else if knowledgeMap.MasteryScore < 0.0 {
		knowledgeMap.MasteryScore = 0.0
	}

	// 更新置信度（随着尝试次数增加）
	knowledgeMap.ConfidenceLevel = 1.0 - (1.0 / float32(knowledgeMap.AttemptCount+1))

	// 更新学习曲线
	learningPoint := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"score":     knowledgeMap.MasteryScore,
		"is_correct": isCorrect,
		"response_time": responseTime,
	}
	kms.updateLearningCurve(&knowledgeMap, learningPoint)

	// 分析强项和弱项
	kms.analyzeStrengthsAndWeaknesses(&knowledgeMap)

	knowledgeMap.LastUpdated = time.Now()

	if err == gorm.ErrRecordNotFound {
		return kms.db.Create(&knowledgeMap).Error
	}
	return kms.db.Save(&knowledgeMap).Error
}

// GetKnowledgeMap 获取学生的知识图谱
func (kms *KnowledgeMapService) GetKnowledgeMap(studentID string) ([]model.StudentKnowledgeMap, error) {
	var knowledgeMaps []model.StudentKnowledgeMap
	err := kms.db.Where("student_id = ?", studentID).
		Preload("KnowledgePoint").
		Find(&knowledgeMaps).Error
	return knowledgeMaps, err
}

// GetKnowledgePointMastery 获取特定知识点的掌握情况
func (kms *KnowledgeMapService) GetKnowledgePointMastery(studentID, knowledgePointID string) (*model.StudentKnowledgeMap, error) {
	var knowledgeMap model.StudentKnowledgeMap
	err := kms.db.Where("student_id = ? AND knowledge_point_id = ?", studentID, knowledgePointID).
		First(&knowledgeMap).Error
	if err != nil {
		return nil, err
	}
	return &knowledgeMap, nil
}

// GetWeakKnowledgePoints 获取薄弱知识点
func (kms *KnowledgeMapService) GetWeakKnowledgePoints(studentID string, threshold float32) ([]model.StudentKnowledgeMap, error) {
	var knowledgeMaps []model.StudentKnowledgeMap
	err := kms.db.Where("student_id = ? AND mastery_score < ?", studentID, threshold).
		Preload("KnowledgePoint").
		Order("mastery_score ASC").
		Find(&knowledgeMaps).Error
	return knowledgeMaps, err
}

// GetStrongKnowledgePoints 获取强项知识点
func (kms *KnowledgeMapService) GetStrongKnowledgePoints(studentID string, threshold float32) ([]model.StudentKnowledgeMap, error) {
	var knowledgeMaps []model.StudentKnowledgeMap
	err := kms.db.Where("student_id = ? AND mastery_score >= ?", studentID, threshold).
		Preload("KnowledgePoint").
		Order("mastery_score DESC").
		Find(&knowledgeMaps).Error
	return knowledgeMaps, err
}

// AnalyzeLearningProgress 分析学习进度
func (kms *KnowledgeMapService) AnalyzeLearningProgress(studentID string, days int) (map[string]interface{}, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	var knowledgeMaps []model.StudentKnowledgeMap
	err := kms.db.Where("student_id = ? AND last_updated >= ?", studentID, startDate).
		Find(&knowledgeMaps).Error
	if err != nil {
		return nil, err
	}

	// 计算统计数据
	totalPoints := len(knowledgeMaps)
	weakPoints := 0
	strongPoints := 0
	averageMastery := float32(0)

	for _, km := range knowledgeMaps {
		averageMastery += km.MasteryScore
		if km.MasteryScore < 0.6 {
			weakPoints++
		} else if km.MasteryScore >= 0.8 {
			strongPoints++
		}
	}

	if totalPoints > 0 {
		averageMastery /= float32(totalPoints)
	}

	// 计算学习趋势（最近7天vs前7天）
	midDate := time.Now().AddDate(0, 0, -days/2)
	var recentMaps, oldMaps []model.StudentKnowledgeMap

	for _, km := range knowledgeMaps {
		if km.LastUpdated.After(midDate) {
			recentMaps = append(recentMaps, km)
		} else {
			oldMaps = append(oldMaps, km)
		}
	}

	recentAvg := float32(0)
	if len(recentMaps) > 0 {
		for _, km := range recentMaps {
			recentAvg += km.MasteryScore
		}
		recentAvg /= float32(len(recentMaps))
	}

	oldAvg := float32(0)
	if len(oldMaps) > 0 {
		for _, km := range oldMaps {
			oldAvg += km.MasteryScore
		}
		oldAvg /= float32(len(oldMaps))
	}

	trend := "stable"
	if recentAvg > oldAvg+0.05 {
		trend = "improving"
	} else if recentAvg < oldAvg-0.05 {
		trend = "declining"
	}

	return map[string]interface{}{
		"total_knowledge_points": totalPoints,
		"weak_points":           weakPoints,
		"strong_points":         strongPoints,
		"average_mastery":       averageMastery,
		"learning_trend":        trend,
		"recent_average":        recentAvg,
		"previous_average":      oldAvg,
	}, nil
}

// updateLearningCurve 更新学习曲线
func (kms *KnowledgeMapService) updateLearningCurve(km *model.StudentKnowledgeMap, newPoint map[string]interface{}) {
	var curve []map[string]interface{}
	if err := json.Unmarshal([]byte(km.LearningCurve), &curve); err != nil {
		curve = []map[string]interface{}{}
	}

	// 只保留最近50个数据点
	if len(curve) >= 50 {
		curve = curve[1:]
	}
	curve = append(curve, newPoint)

	curveJSON, _ := json.Marshal(curve)
	km.LearningCurve = string(curveJSON)
}

// analyzeStrengthsAndWeaknesses 分析强项和弱项
func (kms *KnowledgeMapService) analyzeStrengthsAndWeaknesses(km *model.StudentKnowledgeMap) {
	// 获取相关知识点
	var knowledgePoint model.KnowledgePoint
	kms.db.First(&knowledgePoint, "id = ?", km.KnowledgePointID)

	// 简单的强项弱项分析（可以根据具体需求扩展）
	strengths := []string{}
	weaknesses := []string{}

	if km.MasteryScore >= 0.8 {
		strengths = append(strengths, knowledgePoint.Name)
	} else if km.MasteryScore < 0.4 {
		weaknesses = append(weaknesses, knowledgePoint.Name)
	}

	// 如果准确率很高但掌握度不高，可能是运气或简单题目
	accuracy := float32(km.CorrectCount) / float32(km.AttemptCount)
	if accuracy > 0.8 && km.MasteryScore < 0.6 {
		weaknesses = append(weaknesses, knowledgePoint.Name+"(需要巩固)")
	}

	strengthsJSON, _ := json.Marshal(strengths)
	weaknessesJSON, _ := json.Marshal(weaknesses)

	km.StrengthAreas = string(strengthsJSON)
	km.WeakAreas = string(weaknessesJSON)
}

// RecommendNextStudy 推荐下一步学习内容
func (kms *KnowledgeMapService) RecommendNextStudy(studentID string) ([]model.KnowledgePoint, error) {
	// 获取薄弱知识点
	weakMaps, err := kms.GetWeakKnowledgePoints(studentID, 0.6)
	if err != nil {
		return nil, err
	}

	var knowledgePoints []model.KnowledgePoint
	for _, wm := range weakMaps {
		var kp model.KnowledgePoint
		if err := kms.db.First(&kp, "id = ?", wm.KnowledgePointID).Error; err == nil {
			knowledgePoints = append(knowledgePoints, kp)
		}
	}

	// 如果没有薄弱点，推荐难度适中的新知识点
	if len(knowledgePoints) == 0 {
		err := kms.db.Where("level = ?", 2).Limit(5).Find(&knowledgePoints).Error
		if err != nil {
			return nil, err
		}
	}

	return knowledgePoints, nil
}