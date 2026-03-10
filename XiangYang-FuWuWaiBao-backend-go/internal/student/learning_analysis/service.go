package learning_analysis

import (
	"context"
)

// IService 定义学情分析的业务逻辑层接口
type IService interface {
	GetPersonalMicroStats(ctx context.Context, courseID string) (map[string]interface{}, error)
}

type service struct {
	repo IRepository
}

func NewService(repo IRepository) IService {
	return &service{repo: repo}
}

func (s *service) GetPersonalMicroStats(ctx context.Context, courseID string) (map[string]interface{}, error) {
	// 1. 调用数据访问层拉取聚合数据
	rawStats, err := s.repo.GetStudentCourseStats(ctx, courseID)
	if err != nil {
		return nil, err
	}

	// 2. 🚧 【后期实现位置：业务逻辑补充】
	// 如果需要结合大模型（如 Ark）对学生的整体学情给出一句综合性的 AI 评价，可以在这里发起调用并拼接到返回结果中。

	// 3. 组装最终供接口返回的结构
	result := map[string]interface{}{
		"courseId":      courseID,
		"totalTimeMins": rawStats["totalTimeMins"],
		"questionCount": rawStats["questionCount"],
		"noteCount":     rawStats["noteCount"],
		"masteryRadar":  rawStats["radarData"],
	}

	return result, nil
}
