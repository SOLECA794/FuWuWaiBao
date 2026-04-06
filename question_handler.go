package question_handler

import (
	"net/http"	
	"github.com/gorilla/mux"
	"your_project/models" // replace with your actual models package
)

// GetQuestion handles the retrieval of a single question by ID.
func GetQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	questionID := vars["id"]

	// Fetch question from the database (pseudo-code)
	question, err := models.GetQuestionByID(questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Respond with the question data (pseudo-code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}

// ListQuestions handles the retrieval of a list of questions.
func ListQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := models.ListQuestions() // Fetch all questions (pseudo-code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the questions list
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}