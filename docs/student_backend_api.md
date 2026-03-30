# Student Backend API Documentation

## API Endpoints

### 1. Question Generation
   - **Endpoint:** `/api/generate-question`
   - **Method:** POST
   - **Description:** Generates a new question based on provided parameters.
   - **Request Body:**  
     ```json
     {
       "topic": "string",
       "difficulty": "string"
     }
     ```
   - **Response:**  
     ```json
     {
       "question": "string",
       "options": ["string"],
       "correctAnswer": "string"
     }
     ```

### 2. Answering Questions
   - **Endpoint:** `/api/answer-question`
   - **Method:** POST
   - **Description:** Submits an answer for grading.
   - **Request Body:**  
     ```json
     {
       "questionId": "string",
       "selectedOption": "string"
     }
     ```
   - **Response:**  
     ```json
     {
       "isCorrect": true,
       "feedback": "string"
     }
     ```

### 3. Grading Questions
   - **Endpoint:** `/api/grade-questions`
   - **Method:** POST
   - **Description:** Grades multiple answers submitted by the student.
   - **Request Body:**  
     ```json
     {
       "answers": [
         {
           "questionId": "string",
           "selectedOption": "string"
         }
       ]
     }
     ```
   - **Response:**  
     ```json
     {
       "results": [
         {
           "questionId": "string",
           "isCorrect": true
         }
       ]
     }
     ```

### 4. Mistake Collection
   - **Endpoint:** `/api/collect-mistakes`
   - **Method:** POST
   - **Description:** Collects mistakes made by students for further analysis.
   - **Request Body:**  
     ```json
     {
       "studentId": "string",
       "mistakes": [
         {
           "questionId": "string",
           "selectedOption": "string",
           "correctAnswer": "string"
         }
       ]
     }
     ```
   - **Response:**  
     ```json
     {
       "message": "Mistakes collected successfully"
     }
     ```

---

## Notes
- Ensure to validate inputs on both client and server sides.
- Handle error cases gracefully and provide meaningful feedback.