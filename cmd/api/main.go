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
	applogger.Info("Redis连接成功")

	// 初始化MinIO客户端
	minioClient, err := oss.NewMinioClient(&cfg.OSS)
	if err != nil {
		applogger.Sugar.Fatalf("初始化MinIO失败: %v", err)
	}
	applogger.Info("MinIO连接成功")

	// 初始化服务
	courseService := service.NewCourseService(db, minioClient)

	// 初始化处理器 - 注意参数个数
	courseHandler := handler.NewCourseHandler(courseService) // 只传一个参数
	teacherHandler := handler.NewTeacherHandler(db)
	studentHandler := handler.NewStudentHandler(db, redisClient)
	aiHandler := handler.NewAIHandler(db)
	weakPointHandler := handler.NewWeakPointHandler(db)

	// 设置Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	// UTF-8中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	// 请求体UTF-8检查
	r.Use(func(c *gin.Context) {
		bodyBytes, readErr := io.ReadAll(c.Request.Body)
		if readErr == nil && len(bodyBytes) > 0 {
			if !utf8.Valid(bodyBytes) {
				applogger.Warn("请求体包含非 UTF-8 字符")
			}
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

	// API路由 - 使用 v1 版本
	api := r.Group("/api/v1")
	{
		// ==================== 教师端接口 ====================
		teacher := api.Group("/teacher/coursewares")
		{
			// 2.1 获取课件列表
			teacher.GET("/", teacherHandler.GetCoursewareList)

			// 2.2 上传并解析课件
			teacher.POST("/upload", courseHandler.UploadCourse)

			// 2.6 发布课件
			teacher.POST("/:courseId/publish", teacherHandler.PublishCourseware)

			// 1.3 删除课件
			teacher.DELETE("/:courseId", courseHandler.DeleteCourse)

			// 2.3 获取页面讲稿
			teacher.GET("/:courseId/scripts/:pageNum", teacherHandler.GetScript)

			// 2.4 更新页面讲稿
			teacher.PUT("/:courseId/scripts/:pageNum", teacherHandler.UpdateScript)

			// 2.5 AI生成讲稿
			teacher.POST("/:courseId/scripts/ai-generate", teacherHandler.AIGenerateScript)

			// 6.1 教师端 - 班级宏观学情
			teacher.GET("/:courseId/stats", teacherHandler.GetClassStats)

			// 6.2 教师端 - 历史提问记录
			teacher.GET("/:courseId/questions", teacherHandler.GetQuestionRecords)
		}

		// ==================== AI 学伴与互动答疑 ====================
		ai := api.Group("/ai/coursewares")
		{
			// 3.1 获取课件知识图谱
			ai.GET("/:courseId/knowledge-graph", aiHandler.GetKnowledgeGraph)

			// 3.2 智能多模态答疑
			ai.POST("/:courseId/ask", aiHandler.AskQuestion)
		}

		// ==================== 学生端接口 ====================
		student := api.Group("/student")
		{
			// 5.1 获取/更新学习断点
			student.GET("/coursewares/:courseId/breakpoint", studentHandler.GetBreakpoint)
			student.PUT("/coursewares/:courseId/breakpoint", studentHandler.UpdateBreakpoint)

			// 5.2 保存随堂笔记
			student.POST("/coursewares/:courseId/notes", studentHandler.SaveNote)

			// 6.3 学生端 - 个人微观学情
			student.GET("/coursewares/:courseId/stats", studentHandler.GetPersonalStats)

			// 4.1 获取个人薄弱点列表 - 暂时注释掉，等 weakPointHandler 实现后再启用
			student.GET("/coursewares/:courseId/weak-points", weakPointHandler.GetWeakPointList)

			// 4.2 薄弱点 AI 详细讲解 - 暂时注释掉
			student.GET("/weak-points/:weakPointId/explain", weakPointHandler.GetWeakPointExplain)

			// 4.3 生成随堂检测题 - 暂时注释掉
			student.POST("/weak-points/:weakPointId/generate-test", weakPointHandler.GenerateTest)

			// 4.4 提交并校验答案 - 暂时注释掉
			student.POST("/tests/:questionId/check", weakPointHandler.CheckAnswer)

			// 6.1 开始学习会话
			student.POST("/session/start", studentHandler.StartSession)

			// 6.2 上报播放进度
			student.POST("/progress/update", studentHandler.UpdateProgress)

			// 6.3 获取某页讲稿（学生播放用）
			student.GET("/coursewares/:courseId/pages/:pageNum", studentHandler.GetCoursewarePage)

			// 6.4 问答流式接口（核心）
			student.POST("/qa/stream", studentHandler.QAStream)
		}
	}

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	applogger.Sugar.Infof("服务器启动成功，访问地址: http://localhost%s", addr)
	applogger.Sugar.Infof("健康检查: http://localhost%s/health", addr)

	if err := r.Run(addr); err != nil {
		applogger.Sugar.Fatalf("服务器启动失败: %v", err)
	}
}
