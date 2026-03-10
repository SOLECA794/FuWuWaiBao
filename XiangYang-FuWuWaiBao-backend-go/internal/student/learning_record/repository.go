package learning_record

import (
	"context"
	"time"
)

// IRepository 定义学习记录的数据访问层接口
type IRepository interface {
	GetBreakpoint(ctx context.Context, courseID string) (int, error)
	UpdateBreakpoint(ctx context.Context, courseID string, pageNum int) error
	CreateNote(ctx context.Context, courseID string, noteData map[string]interface{}) (string, time.Time, error)
}

type repository struct {
	// 实际开发：
	// db *gorm.DB       // 用于存储笔记
	// redis *redis.Client // 用于存储高频更新的断点
}

func NewRepository() IRepository {
	return &repository{}
}

func (r *repository) GetBreakpoint(ctx context.Context, courseID string) (int, error) {
	// 🚧 【后期实现位置】：从 Redis 或数据库中读取该学生在该课件下的断点页码
	return 5, nil // Mock: 假设上次学到了第 5 页
}

func (r *repository) UpdateBreakpoint(ctx context.Context, courseID string, pageNum int) error {
	// 🚧 【后期实现位置】：将最新的页码覆盖写入 Redis 或数据库
	return nil
}

func (r *repository) CreateNote(ctx context.Context, courseID string, noteData map[string]interface{}) (string, time.Time, error) {
	// 🚧 【后期实现位置】：将笔记数据插入数据库的 note 表
	mockNoteID := "note-8848"
	mockCreatedAt := time.Now()
	return mockNoteID, mockCreatedAt, nil
}
