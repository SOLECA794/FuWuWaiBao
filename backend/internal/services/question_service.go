package services

// QuestionService struct for managing questions
type QuestionService struct {}

// GetQuestionByID retrieves a question by its ID
func (qs *QuestionService) GetQuestionByID(id int) (Question, error) {
    // Implementation goes here
}

// ListQuestions lists all questions
func (qs *QuestionService) ListQuestions() ([]Question, error) {
    // Implementation goes here
}

// CreateQuestion creates a new question
func (qs *QuestionService) CreateQuestion(q Question) (Question, error) {
    // Implementation goes here
}

// Question struct represents a question entity
type Question struct {
    ID int
    Title string
    Content string
}