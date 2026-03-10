package ai_partner

import (
	"context"
	"fmt"
)

// IService 定义业务逻辑层接口
type IService interface {
	GetKnowledgeGraph(ctx context.Context, courseID string) (map[string]interface{}, error)
	GenerateAnswer(ctx context.Context, courseID string, req AskRequest) (map[string]interface{}, error)
}

// service 是 IService 接口的具体实现
type service struct {
	repo IRepository // 依赖注入 Repository 层
}

// NewService 是构造函数，接收 IRepository 作为依赖
func NewService(repo IRepository) IService {
	return &service{
		repo: repo,
	}
}

func (s *service) GetKnowledgeGraph(ctx context.Context, courseID string) (map[string]interface{}, error) {
	// 业务逻辑：调用 Repo 获取数据
	return s.repo.GetCourseKnowledgeGraph(ctx, courseID)
}

func (s *service) GenerateAnswer(ctx context.Context, courseID string, req AskRequest) (map[string]interface{}, error) {
	// 🚧 【后期实现位置：核心业务逻辑】
	// 1. 根据 courseID 和 req.PageNum，可能需要调用 repo 获取当前页的讲稿上下文
	// 2. 组装 Prompt
	// 3. 调用 Ark 大模型 API 进行多模态解答
	// 4. 将生成的回答通过 s.repo.SaveAskRecord() 异步保存到数据库

	// 返回符合格式的 Mock 数据
	mockAnswer := map[string]interface{}{
		"answer": fmt.Sprintf("这是 Ark 模型针对第 %d 页【%s】内容的解答...", req.PageNum, req.Question),
		"source": fmt.Sprintf("第 %d 页讲稿", req.PageNum),
	}

	return mockAnswer, nil
}
