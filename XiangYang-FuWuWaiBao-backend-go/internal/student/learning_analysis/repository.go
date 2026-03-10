package learning_analysis

import (
	"context"
)

// IRepository 定义学情分析的数据访问层接口
type IRepository interface {
	GetStudentCourseStats(ctx context.Context, courseID string) (map[string]interface{}, error)
}

type repository struct {
	// db *gorm.DB // 用于执行聚合查询的数据库连接
}

func NewRepository() IRepository {
	return &repository{}
}

func (r *repository) GetStudentCourseStats(ctx context.Context, courseID string) (map[string]interface{}, error) {
	// 🚧 【后期实现位置】：编写 SQL 聚合查询
	// SELECT SUM(duration) FROM learning_logs...
	// SELECT COUNT(*) FROM ask_records...

	// 返回 Mock 的底层聚合数据
	mockData := map[string]interface{}{
		"totalTimeMins": 360,
		"questionCount": 42,
		"noteCount":     18,
		"radarData": []map[string]interface{}{
			{"dimension": "基础语法", "score": 0.85},
			{"dimension": "并发编程", "score": 0.90},
			{"dimension": "网络通信", "score": 0.80},
			{"dimension": "设计模式", "score": 0.75},
			{"dimension": "算法基础", "score": 0.88},
		},
	}
	return mockData, nil
}
