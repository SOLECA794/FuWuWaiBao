package handler

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/service"
)

func ensureDialogueSession(db *gorm.DB, sessionID, userID, courseID string, page int, nodeID string, currentTimeSec int) string {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		sessionID = fmt.Sprintf("sess_%d", time.Now().UnixNano())
	}

	now := time.Now()
	var session model.DialogueSession
	err := db.First(&session, "id = ?", sessionID).Error
	if err == nil {
		updates := map[string]any{
			"user_id":          strings.TrimSpace(userID),
			"course_id":        strings.TrimSpace(courseID),
			"current_page":     page,
			"current_node_id":  strings.TrimSpace(nodeID),
			"current_time_sec": maxInt(currentTimeSec, 0),
			"last_asked_at":    now,
		}
		_ = db.Model(&session).Updates(updates).Error
		return sessionID
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return sessionID
	}

	created := model.DialogueSession{
		BaseModel:      model.BaseModel{ID: sessionID},
		UserID:         strings.TrimSpace(userID),
		CourseID:       strings.TrimSpace(courseID),
		CurrentPage:    page,
		CurrentNodeID:  strings.TrimSpace(nodeID),
		CurrentTimeSec: maxInt(currentTimeSec, 0),
		LastAskedAt:    &now,
	}
	_ = db.Create(&created).Error
	return sessionID
}

func syncDialogueSessionState(db *gorm.DB, sessionID, userID, courseID string, page int, nodeID string, currentTimeSec int) {
	if strings.TrimSpace(sessionID) == "" || strings.TrimSpace(courseID) == "" {
		return
	}
	ensureDialogueSession(db, sessionID, userID, courseID, page, nodeID, currentTimeSec)
}

func appendDialogueTurn(db *gorm.DB, sessionID, userID, courseID string, page int, nodeID, question, answer string, sourcePage int, needReteach bool, followUpSuggestion string) {
	if strings.TrimSpace(sessionID) == "" || strings.TrimSpace(courseID) == "" || strings.TrimSpace(question) == "" {
		return
	}

	var turnCount int64
	_ = db.Model(&model.DialogueTurn{}).Where("session_id = ?", sessionID).Count(&turnCount).Error
	turn := model.DialogueTurn{
		SessionID:          sessionID,
		CourseID:           strings.TrimSpace(courseID),
		UserID:             strings.TrimSpace(userID),
		TurnIndex:          int(turnCount) + 1,
		PageIndex:          page,
		NodeID:             strings.TrimSpace(nodeID),
		Question:           strings.TrimSpace(question),
		Answer:             strings.TrimSpace(answer),
		SourcePage:         sourcePage,
		NeedReteach:        needReteach,
		FollowUpSuggestion: strings.TrimSpace(followUpSuggestion),
	}
	_ = db.Create(&turn).Error
}

func buildDialogueContext(db *gorm.DB, sessionID string, limit int) (string, []service.ConversationTurn) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return "", nil
	}
	if limit <= 0 {
		limit = 4
	}

	var turns []model.DialogueTurn
	_ = db.Where("session_id = ?", sessionID).Order("turn_index desc").Limit(limit).Find(&turns).Error
	if len(turns) == 0 {
		return "", nil
	}

	for left, right := 0, len(turns)-1; left < right; left, right = left+1, right-1 {
		turns[left], turns[right] = turns[right], turns[left]
	}

	history := make([]service.ConversationTurn, 0, len(turns))
	parts := make([]string, 0, len(turns))
	for _, turn := range turns {
		history = append(history, service.ConversationTurn{
			Question: turn.Question,
			Answer:   turn.Answer,
			Page:     turn.PageIndex,
			NodeID:   turn.NodeID,
		})
		parts = append(parts, fmt.Sprintf("第%d轮 学生：%s AI：%s", turn.TurnIndex, trimForDialogueSummary(turn.Question, 60), trimForDialogueSummary(turn.Answer, 100)))
	}

	return strings.Join(parts, "\n"), history
}

func trimForDialogueSummary(value string, limit int) string {
	value = strings.TrimSpace(value)
	if value == "" || limit <= 0 {
		return value
	}
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit]) + "..."
}

func resolveResumeNodeID(currentNodeID string, page int, needReteach bool) string {
	if needReteach {
		return defaultStringValue(currentNodeID, fmt.Sprintf("p%d_n1", page))
	}
	return nextNodeID(currentNodeID, page)
}
