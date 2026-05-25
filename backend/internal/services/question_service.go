package services

// QuestionService struct for managing questions
type QuestionService struct{}

// Question struct represents a question entity
type Question struct {
    ID      int
    Title   string
    Content string
}

// GetQuestionByID retrieves a question by its ID
func (qs *QuestionService) GetQuestionByID(id int) (Question, error) {
    // 简单默认实现，后续替换为 DB 查询
    return Question{ID: id, Title: "", Content: ""}, nil
}

// ListQuestions lists all questions
func (qs *QuestionService) ListQuestions() ([]Question, error) {
    // 返回空列表作为占位实现
    return []Question{}, nil
}

// CreateQuestion creates a new question
func (qs *QuestionService) CreateQuestion(q Question) (Question, error) {
    // 占位：如果未提供 ID，保持 0
    return q, nil
}