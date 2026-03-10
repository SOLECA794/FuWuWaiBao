package router

import (
	"github.com/gin-gonic/gin"

	// 引入学生端的四大模块 (注意替换 smart-learning 为你自己的 module name)
	"smart-learning/internal/student/ai_partner"
	"smart-learning/internal/student/learning_analysis"
	"smart-learning/internal/student/learning_record"
	"smart-learning/internal/student/weak_points"
)

// SetupRouter 初始化并配置整个项目的 HTTP 路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 全局跨域等中间件可以在这里挂载
	// r.Use(middleware.Cors())

	v1 := r.Group("/api/v1")

	// 💡 预留：全局鉴权中间件
	// v1.Use(middleware.JWTAuth())

	// ==========================================
	// 1. 依赖注入 (Wire Dependencies)
	// 规则：Repository -> Service -> Handler
	// ==========================================

	// 1.1 AI 学伴模块
	aiRepo := ai_partner.NewRepository()
	aiSvc := ai_partner.NewService(aiRepo)
	aiHandler := ai_partner.NewHandler(aiSvc)

	// 1.2 薄弱点模块
	wpRepo := weak_points.NewRepository()
	wpSvc := weak_points.NewService(wpRepo)
	wpHandler := weak_points.NewHandler(wpSvc)

	// 1.3 学习记录模块
	recordRepo := learning_record.NewRepository()
	recordSvc := learning_record.NewService(recordRepo)
	recordHandler := learning_record.NewHandler(recordSvc)

	// 1.4 学情分析模块
	analysisRepo := learning_analysis.NewRepository()
	analysisSvc := learning_analysis.NewService(analysisRepo)
	analysisHandler := learning_analysis.NewHandler(analysisSvc)

	// ==========================================
	// 2. 注册路由 (Register Routes)
	// ==========================================

	// ---------------- 三、AI 学伴与互动答疑 ----------------
	aiGroup := v1.Group("/ai/coursewares")
	{
		aiGroup.GET("/:courseId/knowledge-graph", aiHandler.GetKnowledgeGraph)
		aiGroup.POST("/:courseId/ask", aiHandler.AskQuestion)
	}

	// ---------------- 学生端综合业务组 ----------------
	studentGroup := v1.Group("/student")
	{
		// 针对特定【课件】的操作集合
		coursewareGroup := studentGroup.Group("/coursewares/:courseId")
		{
			// 四、获取个人薄弱点列表
			coursewareGroup.GET("/weak-points", wpHandler.GetWeakPointsList)

			// 五、学习过程数据 (断点与笔记)
			coursewareGroup.GET("/breakpoint", recordHandler.GetBreakpoint)
			coursewareGroup.PUT("/breakpoint", recordHandler.UpdateBreakpoint)
			coursewareGroup.POST("/notes", recordHandler.SaveNote)

			// 六、个人微观学情分析
			coursewareGroup.GET("/stats", analysisHandler.GetPersonalMicroStats)
		}

		// 针对特定【薄弱点】的操作集合
		wpGroup := studentGroup.Group("/weak-points/:weakPointId")
		{
			wpGroup.GET("/explain", wpHandler.ExplainWeakPoint)
			wpGroup.POST("/generate-test", wpHandler.GenerateTest)
		}

		// 针对特定【题目】的操作集合
		studentGroup.POST("/tests/:questionId/check", wpHandler.CheckTestAnswer)
	}

	// 🚧 【预留】：教师端路由组 (Teacher Group) 待开发
	// teacherGroup := v1.Group("/teacher")
	// ...

	return r
}
