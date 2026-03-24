package service

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"smart-teaching-backend/internal/model"
)

// NotificationService 通知服务
type NotificationService struct {
	db *gorm.DB
}

// NewNotificationService 创建通知服务
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// CreateNotification 创建通知
func (ns *NotificationService) CreateNotification(studentID, title, content, notificationType, priority string, relatedID, relatedType string) (*model.Notification, error) {
	notification := &model.Notification{
		StudentID:   studentID,
		Title:       title,
		Content:     content,
		Type:        notificationType,
		Priority:    priority,
		Status:      "unread",
		RelatedID:   relatedID,
		RelatedType: relatedType,
		Channels:    `["app"]`, // 默认通过app推送
	}

	if err := ns.db.Create(notification).Error; err != nil {
		return nil, err
	}

	return notification, nil
}

// GetNotifications 获取用户通知
func (ns *NotificationService) GetNotifications(studentID string, status string, limit, offset int) ([]model.Notification, error) {
	var notifications []model.Notification
	query := ns.db.Where("student_id = ?", studentID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error

	return notifications, err
}

// MarkAsRead 标记通知为已读
func (ns *NotificationService) MarkAsRead(notificationID string) error {
	return ns.db.Model(&model.Notification{}).
		Where("id = ?", notificationID).
		Update("status", "read").Error
}

// MarkAllAsRead 标记所有通知为已读
func (ns *NotificationService) MarkAllAsRead(studentID string) error {
	return ns.db.Model(&model.Notification{}).
		Where("student_id = ? AND status = ?", studentID, "unread").
		Update("status", "read").Error
}

// GetUnreadCount 获取未读通知数量
func (ns *NotificationService) GetUnreadCount(studentID string) (int64, error) {
	var count int64
	err := ns.db.Model(&model.Notification{}).
		Where("student_id = ? AND status = ?", studentID, "unread").
		Count(&count).Error
	return count, err
}

// ScheduleNotification 调度定时通知
func (ns *NotificationService) ScheduleNotification(studentID, title, content, notificationType, priority string, scheduledAt time.Time, relatedID, relatedType string) error {
	notification := &model.Notification{
		StudentID:   studentID,
		Title:       title,
		Content:     content,
		Type:        notificationType,
		Priority:    priority,
		Status:      "scheduled",
		RelatedID:   relatedID,
		RelatedType: relatedType,
		ScheduledAt: &scheduledAt,
		Channels:    `["app"]`,
	}

	return ns.db.Create(notification).Error
}

// ProcessScheduledNotifications 处理定时通知
func (ns *NotificationService) ProcessScheduledNotifications() error {
	var notifications []model.Notification
	now := time.Now()

	err := ns.db.Where("status = ? AND scheduled_at <= ?", "scheduled", now).
		Find(&notifications).Error
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		// 标记为未读并设置发送时间
		sentAt := time.Now()
		updates := map[string]interface{}{
			"status":  "unread",
			"sent_at": sentAt,
		}

		ns.db.Model(&notification).Updates(updates)
	}

	return nil
}

// ArchiveOldNotifications 归档旧通知
func (ns *NotificationService) ArchiveOldNotifications(days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return ns.db.Model(&model.Notification{}).
		Where("status = ? AND created_at < ?", "read", cutoffDate).
		Update("status", "archived").Error
}

// GetNotificationStats 获取通知统计
func (ns *NotificationService) GetNotificationStats(studentID string) (map[string]interface{}, error) {
	var stats struct {
		Total     int64
		Unread    int64
		Read      int64
		Archived  int64
	}

	// 总数
	ns.db.Model(&model.Notification{}).Where("student_id = ?", studentID).Count(&stats.Total)

	// 未读数
	ns.db.Model(&model.Notification{}).Where("student_id = ? AND status = ?", studentID, "unread").Count(&stats.Unread)

	// 已读数
	ns.db.Model(&model.Notification{}).Where("student_id = ? AND status = ?", studentID, "read").Count(&stats.Read)

	// 归档数
	ns.db.Model(&model.Notification{}).Where("student_id = ? AND status = ?", studentID, "archived").Count(&stats.Archived)

	// 类型统计
	var typeStats []struct {
		Type  string
		Count int64
	}
	ns.db.Model(&model.Notification{}).
		Select("type, count(*) as count").
		Where("student_id = ?", studentID).
		Group("type").
		Find(&typeStats)

	typeCountMap := make(map[string]int64)
	for _, ts := range typeStats {
		typeCountMap[ts.Type] = ts.Count
	}

	return map[string]interface{}{
		"total":     stats.Total,
		"unread":    stats.Unread,
		"read":      stats.Read,
		"archived":  stats.Archived,
		"by_type":   typeCountMap,
	}, nil
}

// SendPushNotification 发送推送通知（模拟实现）
func (ns *NotificationService) SendPushNotification(notification *model.Notification) error {
	// 这里应该集成实际的推送服务，如：
	// - Firebase Cloud Messaging (FCM)
	// - Apple Push Notification Service (APNS)
	// - 华为推送
	// - 小米推送
	// - 等等

	// 暂时记录日志
	// log.Printf("Sending push notification to student %s: %s", notification.StudentID, notification.Title)

	// 标记为已发送
	now := time.Now()
	notification.SentAt = &now

	return ns.db.Save(notification).Error
}

// BatchSendNotifications 批量发送通知
func (ns *NotificationService) BatchSendNotifications() error {
	var notifications []model.Notification

	// 获取待发送的通知（未读且未发送）
	err := ns.db.Where("status = ? AND sent_at IS NULL", "unread").
		Find(&notifications).Error
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		if err := ns.SendPushNotification(&notification); err != nil {
			// 记录发送失败，但不中断其他通知
			continue
		}
	}

	return nil
}