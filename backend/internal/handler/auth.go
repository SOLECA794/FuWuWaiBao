package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/pkg/logger"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

type authRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"` // 可选：teacher / student
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	username := strings.TrimSpace(req.Username)
	if username == "" || len(req.Password) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名或密码不合法（密码至少4位）"})
		return
	}

	// 统一小写用户名，避免 teacher 和 TEACHER 视为不同账号
	username = strings.ToLower(username)

	var existing model.User
	if err := h.db.Where("username = ?", username).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名已存在"})
		return
	}

	role := strings.ToLower(strings.TrimSpace(req.Role))
	if role != "student" {
		role = "teacher"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("生成密码哈希失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册失败，请稍后重试"})
		return
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	}

	if err := h.db.Create(user).Error; err != nil {
		logger.Errorf("创建用户失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册失败，请稍后重试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	username := strings.TrimSpace(req.Username)
	if username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名或密码不能为空"})
		return
	}

	username = strings.ToLower(username)

	var user model.User
	if err := h.db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "账号或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "账号或密码错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

