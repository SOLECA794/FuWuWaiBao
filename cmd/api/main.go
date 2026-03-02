package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger" // 重命名为 gormlogger

	"smart-teaching-backend/internal/handler"
	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/repository"
	"smart-teaching-backend/internal/service"
	"smart-teaching-backend/pkg/config"
	applogger "smart-teaching-backend/pkg/logger" // 重命名为 applogger
	"smart-teaching-backend/pkg/oss"
)

func main() {
	// 获取项目根目录
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	// 加载配置
	cfg, err := config.LoadConfig(filepath.Join(projectRoot, "config"))
	if err != nil {
		// 这里还不能用logger，直接用fmt
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}

	// 初始化日志
	err = applogger.InitLogger(
		cfg.Log.Level,
		filepath.Join(projectRoot, cfg.Log.Filename),
		cfg.Log.MaxSize,
		cfg.Log.MaxBackups,
		cfg.Log.MaxAge,
	)
	if err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		return
	}

	applogger.Info("系统启动中...",
		zap.String("port", fmt.Sprintf("%d", cfg.Server.Port)),
		zap.String("mode", cfg.Server.Mode),
	)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		applogger.Sugar.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(
		&model.Course{},
		&model.CoursePage{},
		&model.UserProgress{},
		&model.QuestionLog{},
		&model.TeacherEdit{},
		&model.MindMapNode{},
	)
	if err != nil {
		applogger.Sugar.Fatalf("数据库迁移失败: %v", err)
	}

	applogger.Info("数据库连接成功", zap.String("database", cfg.Database.DBName))

	// 连接Redis
	redisClient, err := repository.InitRedis(&cfg.Redis)
	if err != nil {
		applogger.Sugar.Fatalf("连接Redis失败: %v", err)
	}
	_ = redisClient // 暂时不使用
	applogger.Info("Redis连接成功")

	// 初始化MinIO客户端
	minioClient, err := oss.NewMinioClient(&cfg.OSS)
	if err != nil {
		applogger.Sugar.Fatalf("初始化MinIO失败: %v", err)
	}
	applogger.Info("MinIO连接成功")

	// 初始化服务
	courseService := service.NewCourseService(db, minioClient)

	// 初始化处理器
	courseHandler := handler.NewCourseHandler(courseService)

	// 设置Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	// 日志中间件
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		applogger.Info("HTTP请求",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// API路由
	api := r.Group("/api/v1")
	{
		// 课件相关
		api.POST("/courses/upload", courseHandler.UploadCourse)
		api.GET("/courses/:id", courseHandler.GetCourse)
		api.GET("/courses/:id/pages", courseHandler.GetCoursePages)
		api.PUT("/pages/:id/script", courseHandler.UpdatePageScript)
		api.DELETE("/courses/:id", courseHandler.DeleteCourse)
	}

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	applogger.Sugar.Infof("服务器启动成功，访问地址: http://localhost%s", addr)
	applogger.Sugar.Infof("健康检查: http://localhost%s/health", addr)

	if err := r.Run(addr); err != nil {
		applogger.Sugar.Fatalf("服务器启动失败: %v", err)
	}
}
