package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/service"
)

// NotificationHandler 通知处理器
type NotificationHandler struct {
	notificationService *service.NotificationService
}

// NewNotificationHandler 创建通知处理器
func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// CreateNotification 创建通知
func (nh *NotificationHandler) CreateNotification(c *gin.Context) {
	var req struct {
		UserID      uint   `json:"userId" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Content     string `json:"content" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Channels    []string `json:"channels"`
		ScheduledAt *int64  `json:"scheduledAt"` // Unix timestamp
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	notification := &model.Notification{
		UserID:   req.UserID,
		Title:    req.Title,
		Content:  req.Content,
		Type:     req.Type,
		Channels: req.Channels,
		Status:   "pending",
	}

	if req.ScheduledAt != nil {
		scheduledTime := time.Unix(*req.ScheduledAt, 0)
		notification.ScheduledAt = &scheduledTime
	}

	err := nh.notificationService.CreateNotification(notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": notification})
}

// GetNotifications 获取用户通知列表
func (nh *NotificationHandler) GetNotifications(c *gin.Context) {
	userIDStr := c.Query("userId")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")
	status := c.Query("status")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 userId"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "userId 参数无效"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "page 参数无效"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pageSize 参数无效"})
		return
	}

	notifications, total, err := nh.notificationService.GetNotifications(uint(userID), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取通知列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  notifications,
			"total": total,
			"page":  page,
			"pageSize": pageSize,
		},
	})
}

// GetNotification 获取单个通知详情
func (nh *NotificationHandler) GetNotification(c *gin.Context) {
	notificationIDStr := c.Param("id")

	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "通知ID无效"})
		return
	}

	notification, err := nh.notificationService.GetNotification(uint(notificationID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "通知不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": notification})
}

// MarkAsRead 标记通知为已读
func (nh *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationIDStr := c.Param("id")

	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "通知ID无效"})
		return
	}

	err = nh.notificationService.MarkAsRead(uint(notificationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "标记已读失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "标记成功"})
}

// MarkAllAsRead 标记所有通知为已读
func (nh *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userIDStr := c.Query("userId")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 userId"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "userId 参数无效"})
		return
	}

	err = nh.notificationService.MarkAllAsRead(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "标记已读失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "标记成功"})
}

// DeleteNotification 删除通知
func (nh *NotificationHandler) DeleteNotification(c *gin.Context) {
	notificationIDStr := c.Param("id")

	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "通知ID无效"})
		return
	}

	err = nh.notificationService.DeleteNotification(uint(notificationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

// GetUnreadCount 获取未读通知数量
func (nh *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userIDStr := c.Query("userId")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 userId"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "userId 参数无效"})
		return
	}

	count, err := nh.notificationService.GetUnreadCount(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取未读数量失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"unreadCount": count}})
}

// SendImmediateNotification 立即发送通知
func (nh *NotificationHandler) SendImmediateNotification(c *gin.Context) {
	var req struct {
		UserID   uint     `json:"userId" binding:"required"`
		Title    string   `json:"title" binding:"required"`
		Content  string   `json:"content" binding:"required"`
		Type     string   `json:"type" binding:"required"`
		Channels []string `json:"channels"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	err := nh.notificationService.SendImmediateNotification(req.UserID, req.Title, req.Content, req.Type, req.Channels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "发送通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "发送成功"})
}