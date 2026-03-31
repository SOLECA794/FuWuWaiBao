package models

type Question struct {
    ID             int       `json:"id"`
    NodeID         int       `json:"node_id"`
    Content        string    `json:"content"`
    AnswerOptions  []string  `json:"answer_options"`
    CorrectAnswer  string    `json:"correct_answer"`
    KnowledgePoints []string  `json:"knowledge_points"`
    Difficulty     int       `json:"difficulty"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}