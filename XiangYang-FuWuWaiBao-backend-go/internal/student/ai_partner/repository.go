package ai_partner

import (
	"context"
)

// IRepository 定义数据访问层接口
type IRepository interface {
	// 获取课件解析后的图谱数据 (可能从数据库或 Redis 获取)
	GetCourseKnowledgeGraph(ctx context.Context, courseID string) (map[string]interface{}, error)
	// 保存提问记录
	SaveAskRecord(ctx context.Context, courseID string, question string, answer string) error
}

// repository 是 IRepository 接口的具体实现
type repository struct {
	// 实际开发中这里通常会注入 *gorm.DB 或 *redis.Client
	// db *gorm.DB
}

// NewRepository 是构造函数
func NewRepository() IRepository {
	return &repository{}
}

func (r *repository) GetCourseKnowledgeGraph(ctx context.Context, courseID string) (map[string]interface{}, error) {
	// 🚧 【后期实现位置】：编写 SQL 查询或 Redis 读取逻辑

	// 返回符合格式的 Mock 数据
	mockData := map[string]interface{}{
		"courseId": courseID,
		"chapters": []map[string]interface{}{
			{
				"title": "第一章：基础概念",
				"knowledgePoints": []map[string]interface{}{
					{"id": "kp-1", "name": "核心原理", "subPoints": []string{"概念A", "概念B"}},
				},
			},
		},
	}
	return mockData, nil
}

func (r *repository) SaveAskRecord(ctx context.Context, courseID string, question string, answer string) error {
	// 🚧 【后期实现位置】：将提问和回答写入数据库日志表
	return nil
}
