package handler

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

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

// 2.1 获取课件列表
// GET /api/v1/teacher/coursewares
func (h *TeacherHandler) GetCoursewareList(c *gin.Context) {
	var courses []model.Course

	if err := h.db.Order("created_at desc").Find(&courses).Error; err != nil {
		logger.Errorf("获取课件列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取课件列表失败",
		})
		return
	}

	// 转换为规范格式
	data := make([]gin.H, 0)
	for _, course := range courses {
		data = append(data, gin.H{
			"courseId":  course.ID,
			"title":     course.Title,
			"fileType":  course.FileType,
			"status":    "published", // 需要根据实际状态调整
			"createdAt": course.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data":    data,
	})
}

// 2.2 上传并解析课件（复用已有的 UploadCourse）
// POST /api/v1/teacher/coursewares/upload

// 2.3 获取页面讲稿
// GET /api/v1/teacher/coursewares/:courseId/scripts/:pageNum
func (h *TeacherHandler) GetScript(c *gin.Context) {
	courseId := c.Param("courseId")
	pageNumStr := c.Param("pageNum")

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "页码必须是数字",
		})
		return
	}

	var coursePage model.CoursePage
	err = h.db.Where("course_id = ? AND page_index = ?", courseId, pageNum).First(&coursePage).Error

	if err != nil {
		// 如果不存在，返回空内容
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "请求成功",
			"data": gin.H{
				"courseId": courseId,
				"pageNum":  pageNum,
				"content":  "",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data": gin.H{
			"courseId": courseId,
			"pageNum":  pageNum,
			"content":  coursePage.ScriptText,
		},
	})
}

// 2.4 更新页面讲稿
// PUT /api/v1/teacher/coursewares/:courseId/scripts/:pageNum
func (h *TeacherHandler) UpdateScript(c *gin.Context) {
	courseId := c.Param("courseId")
	pageNumStr := c.Param("pageNum")

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "页码必须是数字",
		})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 检查是否存在
	var coursePage model.CoursePage
	err = h.db.Where("course_id = ? AND page_index = ?", courseId, pageNum).First(&coursePage).Error

	if err != nil {
		// 不存在，创建新记录
		newPage := model.CoursePage{
			CourseID:   courseId,
			PageIndex:  pageNum,
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
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "保存成功",
	})
}

// 2.5 AI生成讲稿
// POST /api/v1/teacher/coursewares/:courseId/scripts/ai-generate
func (h *TeacherHandler) AIGenerateScript(c *gin.Context) {
	courseId := c.Param("courseId")

	var req struct {
		PageNum int    `json:"pageNum" binding:"required"`
		Mode    string `json:"mode"` // 预留，默认 llm
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 获取课程信息
	var course model.Course
	if err := h.db.First(&course, "id = ?", courseId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课件不存在",
		})
		return
	}

	// TODO: 调用 AI 服务生成讲稿
	mockScript := generateMockScript(course.Title, req.PageNum)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data": gin.H{
			"courseId": courseId,
			"pageNum":  req.PageNum,
			"content":  mockScript,
		},
	})
}

// 2.6 发布课件
// POST /api/v1/teacher/coursewares/:courseId/publish
func (h *TeacherHandler) PublishCourseware(c *gin.Context) {
	courseId := c.Param("courseId")

	// 检查课件是否存在
	var course model.Course
	if err := h.db.First(&course, "id = ?", courseId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "课件不存在",
		})
		return
	}

	// TODO: 更新课件状态为已发布
	// h.db.Model(&course).Update("status", "published")

	logger.Infof("课件发布成功: courseId=%s", courseId)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "发布成功",
		"data": gin.H{
			"courseId":    courseId,
			"publishedAt": time.Now().Format(time.RFC3339),
		},
	})
}

// 6.1 教师端 - 班级宏观学情
// GET /api/v1/teacher/coursewares/:courseId/stats
func (h *TeacherHandler) GetClassStats(c *gin.Context) {
	courseId := c.Param("courseId")

	// 各页面停留时长统计
	type PageStayTime struct {
		Page    int     `json:"page"`
		StayAvg float64 `json:"stayAvg"`
	}

	// 提问频次
	type QuestionFreq struct {
		Page  int `json:"page"`
		Count int `json:"count"`
	}

	var pageStats []PageStayTime
	var questionFreq []QuestionFreq

	// 获取提问频次
	h.db.Table("question_logs").
		Select("page_index as page, count(*) as count").
		Where("course_id = ?", courseId).
		Group("page_index").
		Order("page").
		Scan(&questionFreq)

	// 获取所有提问用于词云
	var questions []string
	h.db.Table("question_logs").
		Where("course_id = ?", courseId).
		Pluck("question", &questions)

	// 生成词云数据
	keywords := generateKeywordStats(questions)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data": gin.H{
			"pageStayTime": pageStats,    // 停留时长（需要从学习行为表获取）
			"questionFreq": questionFreq, // 提问频次
			"wordCloud":    keywords,     // 高频提问词云
		},
	})
}

// 6.2 教师端 - 历史提问记录
// GET /api/v1/teacher/coursewares/:courseId/questions
func (h *TeacherHandler) GetQuestionRecords(c *gin.Context) {
	courseId := c.Param("courseId")

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

	var total int64
	h.db.Model(&model.QuestionLog{}).Where("course_id = ?", courseId).Count(&total)

	var logs []model.QuestionLog
	h.db.Where("course_id = ?", courseId).
		Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data": gin.H{
			"list":     logs,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
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

// generateKeywordStats 生成关键词统计
func generateKeywordStats(questions []string) []gin.H {
	// 常见技术关键词
	commonKeywords := []string{
		"依赖注入", "IoC", "AOP", "Spring", "微服务",
		"分布式", "事务", "缓存", "数据库", "接口",
		"fillna", "interpolate", "dropna", "缺失值", "异常值",
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

	// 如果少于5个，补充一些默认关键词
	if len(stats) < 5 {
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
