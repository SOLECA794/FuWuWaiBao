package models

// Mistake represents the model for tracking student's mistake collections and review status.
type Mistake struct {
    ID             int    `json:"id"`
    StudentID      int    `json:"student_id"`
    ProblemID      int    `json:"problem_id"`
    MistakeContent string `json:"mistake_content"`
    ReviewStatus   string `json:"review_status"` // e.g., 'pending', 'reviewed'
    CreatedAt      string `json:"created_at"`
    UpdatedAt      string `json:"updated_at"`
}