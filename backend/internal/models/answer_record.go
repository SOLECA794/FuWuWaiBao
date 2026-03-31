package models

// AnswerRecord represents a student's answer record.
type AnswerRecord struct {
    StudentID    string    `json:"student_id"`    // Student ID
    QuestionID   string    `json:"question_id"`   // Question ID
    AnswerContent string    `json:"answer_content"` // Answer content
    Score        int       `json:"score"`        // Score for the answer
    Explanation  string    `json:"explanation"`  // Explanation for the answer
    CreatedAt    time.Time `json:"created_at"`    // Timestamp when record was created
    UpdatedAt    time.Time `json:"updated_at"`    // Timestamp when record was last updated
}
