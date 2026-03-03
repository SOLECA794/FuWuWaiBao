package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/pkg/logger"
)

type TeacherHandler struct {
	db *gorm.DB
}

func NewTeacherHandler(db *gorm.DB) *TeacherHandler {
	return &TeacherHandler{db: db}
}

// ==================== 1. 课件管理模块 ====================

// GetCoursewareList 获取课件列表
// GET /api/teacher/courseware-list
func (h *TeacherHandler) GetCoursewareList(c *gin.Context) {
	var courses []model.Course

	// 查询所有课件，按创建时间倒序排列
	if err := h.db.Order("created_at desc").Find(&courses).Error; err != nil {
		logger.Errorf("获取课件列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取课件列表失败",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": courses,
	})
}

// PublishCourseware 发布课件
// POST /api/teacher/publish-courseware
func (h *TeacherHandler) PublishCourseware(c *gin.Context) {
	var req struct {
		CourseID string `json:"courseId" binding:"required"`
		Scope    string `json:"scope" binding:"required"` // all/class1/class2
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("发布课件参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 检查课件是否存在
	var course model.Course
	if err := h.db.First(&course, "id = ?", req.CourseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课件不存在",
		})
		return
	}

	// TODO: 这里可以添加发布逻辑
	// 比如：创建发布记录、推送通知到学生端等
	logger.Infof("课件发布成功: courseId=%s, scope=%s, title=%s",
		req.CourseID, req.Scope, course.Title)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "发布成功",
		"data": gin.H{
			"courseId":    req.CourseID,
			"scope":       req.Scope,
			"publishedAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// ==================== 2. 讲稿编辑模块 ====================

// GetScript 获取讲稿
// GET /api/teacher/script/:courseId/:page
func (h *TeacherHandler) GetScript(c *gin.Context) {
	courseId := c.Param("courseId")
	pageStr := c.Param("page")

	// 转换页码为整数
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "页码必须是数字",
		})
		return
	}

	// 查询指定课件指定页码的讲稿
	var coursePage model.CoursePage
	err = h.db.Where("course_id = ? AND page_index = ?", courseId, page).First(&coursePage).Error

	// 强制设置响应头为 UTF-8
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err != nil {
		// 如果不存在，返回空内容而不是错误（方便前端初始化）
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"courseId": courseId,
				"page":     page,
				"content":  "",
			},
		})
		return
	}

	// 返回讲稿内容
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"courseId": courseId,
			"page":     page,
			"content":  coursePage.ScriptText,
		},
	})
}

// SaveScript 保存讲稿
// POST /api/teacher/script/save
func (h *TeacherHandler) SaveScript(c *gin.Context) {
	var req struct {
		CourseID string `json:"courseId" binding:"required"`
		Page     int    `json:"page" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}

	// 读取原始请求体，检查编码
	bodyBytes, err := c.GetRawData()
	if err != nil {
		logger.Errorf("读取请求体失败: %v", err)
	} else {
		// 限制输出长度，避免日志过大
		maxLen := 50
		if len(bodyBytes) < maxLen {
			maxLen = len(bodyBytes)
		}
		logger.Infof("原始请求体(bytes): %v", bodyBytes[:maxLen])
		logger.Infof("原始请求体(string): %s", string(bodyBytes))
	}

	// 重新设置请求体
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("保存讲稿参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 检查内容是否为有效的 UTF-8
	contentBytes := []byte(req.Content)
	if !utf8.Valid(contentBytes) {
		logger.Warn("内容不是有效的 UTF-8，尝试转换")
		// 尝试修复编码
		validContent := make([]rune, 0, len(contentBytes))
		for i := 0; i < len(contentBytes); {
			r, size := utf8.DecodeRune(contentBytes[i:])
			if r != utf8.RuneError {
				validContent = append(validContent, r)
			}
			i += size
		}
		req.Content = string(validContent)
	}

	logger.Infof("接收到讲稿内容: courseId=%s, page=%d", req.CourseID, req.Page)
	logger.Infof("内容: %s", req.Content)
	logger.Infof("内容长度: %d", len(req.Content))

	// 检查课件是否存在
	var course model.Course
	if err := h.db.First(&course, "id = ?", req.CourseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课件不存在",
		})
		return
	}

	// 检查是否存在，存在则更新，不存在则创建
	var coursePage model.CoursePage
	err = h.db.Where("course_id = ? AND page_index = ?", req.CourseID, req.Page).First(&coursePage).Error

	if err != nil {
		// 不存在，创建新记录
		newPage := model.CoursePage{
			CourseID:   req.CourseID,
			PageIndex:  req.Page,
			ScriptText: req.Content,
		}
		if err := h.db.Create(&newPage).Error; err != nil {
			logger.Errorf("创建讲稿失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "保存失败",
			})
			return
		}
		logger.Infof("创建讲稿成功")
	} else {
		// 存在，更新
		if err := h.db.Model(&coursePage).Update("script_text", req.Content).Error; err != nil {
			logger.Errorf("更新讲稿失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "保存失败",
			})
			return
		}
		logger.Infof("更新讲稿成功")
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "保存成功",
	})
}

// AIGenerateScript AI生成讲稿
// POST /api/teacher/ai-generate-script
func (h *TeacherHandler) AIGenerateScript(c *gin.Context) {
	var req struct {
		CourseID   string `json:"courseId" binding:"required"`
		Page       int    `json:"page" binding:"required"`
		CourseName string `json:"courseName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("AI生成讲稿参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 检查课件是否存在
	var course model.Course
	if err := h.db.First(&course, "id = ?", req.CourseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课件不存在",
		})
		return
	}

	// 生成模拟讲稿内容
	mockScript := generateMockScript(req.CourseName, req.Page)

	logger.Infof("AI生成讲稿成功: courseId=%s, page=%d, courseName=%s",
		req.CourseID, req.Page, req.CourseName)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"courseId": req.CourseID,
			"page":     req.Page,
			"content":  mockScript,
		},
	})
}

// generateMockScript 生成模拟讲稿内容
func generateMockScript(courseName string, page int) string {
	templates := []string{
		`## %s 第%d页：课程导入
    
### 教学目标
- 了解本章节的核心概念
- 掌握基础知识框架
- 能够理解实际应用场景

### 重点内容
1. 概念引入：通过生活中的例子引入
2. 核心定义：准确定义本章节的核心概念
3. 基本原理：讲解基本原理和运作机制

### 案例分析
以一个具体的例子来说明本章节的内容：

示例：在实际开发中，我们经常会遇到...
    
### 小结
本章节主要介绍了...为后续学习打下基础。`,

		`## %s 第%d页：深入讲解
    
### 知识要点
- 要点一：深入理解核心机制
- 要点二：掌握常见的设计模式
- 要点三：学会实际应用技巧

### 代码示例

func main() {
    fmt.Println("Hello, " + courseName)
}

### 注意事项
- 注意边界条件的处理
- 注意性能优化要点
- 注意代码规范要求

### 练习题
1. 请实现一个简单的示例
2. 思考如何优化现有代码
3. 尝试解决实际问题`,

		`## %s 第%d页：实战应用
    
### 实战场景
在实际项目中，我们通常需要...

### 实现步骤
1. 第一步：环境准备
2. 第二步：核心代码实现
3. 第三步：测试验证
4. 第四步：优化改进

### 常见问题
- 问：遇到报错怎么办？
- 答：检查配置文件和环境变量
  
- 问：性能不够怎么办？
- 答：使用缓存和异步处理

### 扩展阅读
- 官方文档链接
- 相关技术博客
- 开源项目参考`,

		`## %s 第%d页：总结回顾
    
### 本章总结
知识点 | 重要程度 | 掌握情况
-------|----------|----------
基础概念 | ⭐⭐⭐ | 需要掌握
核心原理 | ⭐⭐⭐⭐ | 深入理解
实战应用 | ⭐⭐⭐⭐⭐ | 熟练运用

### 重点回顾
1. 核心概念：...
2. 关键技术：...
3. 最佳实践：...

### 下节预告
下一章我们将学习...

### 思考题
1. 如何将本章知识应用到实际项目？
2. 与之前学过的内容有什么联系？
3. 有没有更好的实现方式？`,
	}

	// 根据页码选择不同的模板（循环使用）
	template := templates[(page-1)%len(templates)]
	return fmt.Sprintf(template, courseName, page)
}

// ==================== 3. 学情分析模块 ====================

// GetStudentStats 获取学情分析数据
// GET /api/teacher/student-stats/:courseId
func (h *TeacherHandler) GetStudentStats(c *gin.Context) {
	courseId := c.Param("courseId")

	// 1. 统计每页提问次数
	type PageStats struct {
		PageIndex int `json:"page"`
		Count     int `json:"count"`
	}

	var pageStats []PageStats
	h.db.Table("question_logs").
		Select("page_index, count(*) as count").
		Where("course_id = ?", courseId).
		Group("page_index").
		Order("page_index").
		Scan(&pageStats)

	// 2. 获取所有提问记录用于关键词分析
	var questions []string
	h.db.Table("question_logs").
		Where("course_id = ?", courseId).
		Pluck("question", &questions)

	// 3. 简单的关键词统计（后续可以集成jieba分词）
	keywordStats := generateKeywordStats(questions)

	// 4. 获取总提问数
	var totalQuestions int64
	h.db.Model(&model.QuestionLog{}).Where("course_id = ?", courseId).Count(&totalQuestions)

	// 5. 获取活跃用户数（去重）
	var activeUsers int64
	h.db.Table("question_logs").
		Where("course_id = ?", courseId).
		Distinct("user_id").
		Count(&activeUsers)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"pageStats":      pageStats,
			"keywords":       keywordStats,
			"totalQuestions": totalQuestions,
			"activeUsers":    activeUsers,
		},
	})
}

// generateKeywordStats 生成关键词统计（简单实现）
func generateKeywordStats(questions []string) []gin.H {
	// 常见技术关键词
	commonKeywords := []string{
		"依赖注入", "IoC", "AOP", "Spring", "微服务",
		"分布式", "事务", "缓存", "数据库", "接口",
	}

	var stats []gin.H
	for _, keyword := range commonKeywords {
		count := 0
		for _, q := range questions {
			if strings.Contains(q, keyword) {
				count++
			}
		}
		if count > 0 {
			stats = append(stats, gin.H{
				"word":  keyword,
				"count": count,
			})
		}
	}

	// 按次数排序
	sort.Slice(stats, func(i, j int) bool {
		return stats[i]["count"].(int) > stats[j]["count"].(int)
	})

	// 如果少于10个，补充一些默认关键词
	if len(stats) < 10 {
		defaultKeywords := []string{"概念", "原理", "实现", "应用", "区别"}
		for _, kw := range defaultKeywords {
			stats = append(stats, gin.H{
				"word":  kw,
				"count": 1,
			})
		}
	}

	return stats
}

// ==================== 4. 提问记录模块 ====================

// GetQuestionRecords 获取提问记录
// GET /api/teacher/question-records/:courseId?page=1&pageSize=20
func (h *TeacherHandler) GetQuestionRecords(c *gin.Context) {
	courseId := c.Param("courseId")

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 查询总数
	var total int64
	h.db.Model(&model.QuestionLog{}).Where("course_id = ?", courseId).Count(&total)

	// 分页查询记录
	var logs []model.QuestionLog
	h.db.Where("course_id = ?", courseId).
		Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"pageSize":  pageSize,
			"totalPage": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// ==================== 5. 辅助接口 ====================

// GetPagePreview 获取课件预览图片
// GET /api/courseware/:courseId/page/:pageNum
func (h *TeacherHandler) GetPagePreview(c *gin.Context) {
	courseId := c.Param("courseId")
	pageNumStr := c.Param("pageNum")

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil || pageNum < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "页码必须是正整数",
		})
		return
	}

	// 查询课件页面
	var coursePage model.CoursePage
	err = h.db.Where("course_id = ? AND page_index = ?", courseId, pageNum).First(&coursePage).Error

	if err != nil || coursePage.ImageURL == "" {
		// 如果没有预览图，返回默认图片或提示
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "预览图不存在",
		})
		return
	}

	// 重定向到图片URL
	c.Redirect(http.StatusFound, coursePage.ImageURL)
}
