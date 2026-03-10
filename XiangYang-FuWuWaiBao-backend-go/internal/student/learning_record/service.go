package learning_record

import (
	"context"
	"time"
)

// IService 定义学习记录的业务逻辑层接口
type IService interface {
	GetBreakpoint(ctx context.Context, courseID string) (map[string]interface{}, error)
	UpdateBreakpoint(ctx context.Context, courseID string, req UpdateBreakpointRequest) error
	SaveNote(ctx context.Context, courseID string, req SaveNoteRequest) (map[string]interface{}, error)
}

type service struct {
	repo IRepository
}

func NewService(repo IRepository) IService {
	return &service{repo: repo}
}

func (s *service) GetBreakpoint(ctx context.Context, courseID string) (map[string]interface{}, error) {
	pageNum, err := s.repo.GetBreakpoint(ctx, courseID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"courseId": courseID,
		"pageNum":  pageNum,
	}, nil
}

func (s *service) UpdateBreakpoint(ctx context.Context, courseID string, req UpdateBreakpointRequest) error {
	// 可以在这里加业务校验：比如查一下课件总页数，如果 req.PageNum 大于总页数则报错
	return s.repo.UpdateBreakpoint(ctx, courseID, req.PageNum)
}

func (s *service) SaveNote(ctx context.Context, courseID string, req SaveNoteRequest) (map[string]interface{}, error) {
	noteData := map[string]interface{}{
		"pageNum": req.PageNum,
		"content": req.Content,
		"x":       req.X,
		"y":       req.Y,
	}

	noteID, createdAt, err := s.repo.CreateNote(ctx, courseID, noteData)
	if err != nil {
		return nil, err
	}

	// 组装返回给前端的数据
	return map[string]interface{}{
		"noteId":    noteID,
		"courseId":  courseID,
		"pageNum":   req.PageNum,
		"content":   req.Content,
		"x":         req.X,
		"y":         req.Y,
		"createdAt": createdAt.Format(time.RFC3339), // 格式化为标准时间字符串
	}, nil
}
