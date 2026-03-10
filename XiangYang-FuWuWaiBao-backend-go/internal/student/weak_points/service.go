package weak_points

import (
	"context"
)

// IService 定义薄弱点业务逻辑层接口
type IService interface {
	GetWeakPointsList(ctx context.Context, courseID string) ([]map[string]interface{}, error)
	ExplainWeakPoint(ctx context.Context, weakPointID string) (map[string]interface{}, error)
	GenerateTest(ctx context.Context, weakPointID string, req GenerateTestRequest) (map[string]interface{}, error)
	CheckTestAnswer(ctx context.Context, questionID string, req CheckAnswerRequest) (map[string]interface{}, error)
}

type service struct {
	repo IRepository // 注入 Repository
}

func NewService(repo IRepository) IService {
	return &service{repo: repo}
}

func (s *service) GetWeakPointsList(ctx context.Context, courseID string) ([]map[string]interface{}, error) {
	return s.repo.GetWeakPointsByCourse(ctx, courseID)
}

func (s *service) ExplainWeakPoint(ctx context.Context, weakPointID string) (map[string]interface{}, error) {
	// 🚧 【后期实现位置：核心业务逻辑】
	// 1. 通过 s.repo.GetWeakPointDetail() 获取知识点上下文
	// 2. 调用 Ark 大模型进行 Reteach (重构与二次讲解)

	mockExplanation := map[string]interface{}{
		"weakPointId": weakPointID,
		"content":     "针对这个知识点，我们换个思路来理解。你可以把它想象成...",
	}
	return mockExplanation, nil
}

func (s *service) GenerateTest(ctx context.Context, weakPointID string, req GenerateTestRequest) (map[string]interface{}, error) {
	// 🚧 【后期实现位置：核心业务逻辑】
	// 1. 调用 Ark 大模型根据 req.QuestionType 动态生成单选/多选题目
	// 2. 验证大模型返回的 JSON 格式是否合法
	// 3. 调用 s.repo.SaveGeneratedQuestion() 保存题目到数据库

	mockQuestion := map[string]interface{}{
		"questionId":   "q-2048",
		"questionType": req.QuestionType,
		"content":      "关于该薄弱点，下列说法正确的是？",
		"options": []map[string]string{
			{"key": "A", "text": "描述 A"},
			{"key": "B", "text": "描述 B"},
		},
	}
	return mockQuestion, nil
}

func (s *service) CheckTestAnswer(ctx context.Context, questionID string, req CheckAnswerRequest) (map[string]interface{}, error) {
	// 🚧 【后期实现位置：核心业务逻辑】
	// 1. 获取题目正确答案，对比 req.UserAnswer
	// 2. 如果答错，可能需要调用大模型生成专属的错误解析
	// 3. 调用 s.repo.UpdateMasteryLevel() 更新掌握度

	isCorrect := req.UserAnswer == "A"
	updatedMastery, _ := s.repo.UpdateMasteryLevel(ctx, questionID, isCorrect)

	mockResult := map[string]interface{}{
		"questionId":     questionID,
		"isCorrect":      isCorrect,
		"correctAnswer":  "A",
		"explanation":    "因为选项 A 正确地指出了核心机制...",
		"updatedMastery": updatedMastery,
	}
	return mockResult, nil
}
