package main

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"smart-teaching-backend/internal/handler"
	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/repository"
	"smart-teaching-backend/internal/service"
	"smart-teaching-backend/pkg/config"
	applogger "smart-teaching-backend/pkg/logger"
	"smart-teaching-backend/pkg/oss"
)

func main() {
	// 获取项目根目录
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	// 加载配置
	cfg, err := config.LoadConfig(filepath.Join(projectRoot, "config"))
	if err != nil {
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
	_ = redisClient
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
	teacherHandler := handler.NewTeacherHandler(db)
	studentHandler := handler.NewStudentHandler(db)

	// 设置Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 gin 引擎
	r := gin.New()

	// 最重要的中间件：设置 UTF-8 响应头（放在最前面）
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	// 恢复中间件
	r.Use(gin.Recovery())

	// 请求体 UTF-8 检查中间件
	r.Use(func(c *gin.Context) {
		// 读取请求体
		bodyBytes, readErr := io.ReadAll(c.Request.Body)
		if readErr == nil && len(bodyBytes) > 0 {
			// 确保请求体是 UTF-8
			if !utf8.Valid(bodyBytes) {
				applogger.Warn("请求体包含非 UTF-8 字符")
			}
			// 重新设置请求体，因为已经被读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		c.Next()
	})

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
	api := r.Group("/api")
	{
		// 课件预览（公开）
		api.GET("/courseware/:courseId/page/:pageNum", teacherHandler.GetPagePreview)

		// 学生端接口
		student := api.Group("/student")
		{
			student.POST("/courseware/page", studentHandler.GetCoursewarePage)
			student.POST("/ai/question", studentHandler.AskAIQuestion)
			student.POST("/ai/traceQuestion", studentHandler.TraceAIQuestion)
			student.GET("/studyData", studentHandler.GetStudentStudyData)
			student.GET("/breakpoint", studentHandler.GetStudentBreakpoint)
			student.POST("/saveNote", studentHandler.SaveStudentNote)
		}

		// 教师端接口
		teacher := api.Group("/teacher")
		{
			// 1. 课件管理
			teacher.GET("/courseware-list", teacherHandler.GetCoursewareList)
			teacher.POST("/upload-courseware", courseHandler.UploadCourse)
			teacher.DELETE("/courseware/:courseId", courseHandler.DeleteCourse)
			teacher.POST("/publish-courseware", teacherHandler.PublishCourseware)

			// 2. 讲稿编辑
			teacher.GET("/script/:courseId/:page", teacherHandler.GetScript)
			teacher.POST("/script/save", teacherHandler.SaveScript)
			teacher.POST("/ai-generate-script", teacherHandler.AIGenerateScript)

			// 3. 学情分析
			teacher.GET("/student-stats/:courseId", teacherHandler.GetStudentStats)

			// 4. 提问记录
			teacher.GET("/question-records/:courseId", teacherHandler.GetQuestionRecords)
		}

		// 公开接口
		//api.GET("/courseware/:courseId/page/:pageNum", func(c *gin.Context) {
		//	c.JSON(200, gin.H{"message": "待实现"})
		//})
	}

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	applogger.Sugar.Infof("服务器启动成功，访问地址: http://localhost%s", addr)
	applogger.Sugar.Infof("健康检查: http://localhost%s/health", addr)

	if err := r.Run(addr); err != nil {
		applogger.Sugar.Fatalf("服务器启动失败: %v", err)
	}
}
