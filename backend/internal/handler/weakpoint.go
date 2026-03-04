package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WeakPointHandler struct {
	db *gorm.DB
}

func NewWeakPointHandler(db *gorm.DB) *WeakPointHandler {
	return &WeakPointHandler{db: db}
}

// ParseKnowledge 拆解知识点层级
// POST /api/ai/parseKnowledge
func (h *WeakPointHandler) ParseKnowledge(c *gin.Context) {
	var req struct {
		FileContent string `json:"fileContent" binding:"required"`
		FileType    string `json:"fileType" binding:"required"`
		StudentID   string `json:"studentId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// TODO: 调用AI服务解析知识点
	// 返回模拟数据
	data := gin.H{
		"structure": []gin.H{
			{
				"chapter": "第一章：数据清洗基础",
				"knowledgePoints": []gin.H{
					{"name": "缺失值处理", "subPoints": []string{"fillna()", "interpolate()", "dropna()"}},
					{"name": "异常值识别", "subPoints": []string{"Z-Score", "IQR"}},
				},
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// GetWeakPointList 获取薄弱点列表
// GET /api/weakPoint/getList?studentId=xxx&courseId=xxx
func (h *WeakPointHandler) GetWeakPointList(c *gin.Context) {
	studentId := c.Query("studentId")
	courseId := c.Query("courseId")

	if studentId == "" || courseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少参数: studentId 或 courseId",
		})
		return
	}

	// 从数据库查询或生成模拟数据
	weakPoints := []gin.H{
		{"name": "缺失值填充", "count": 5, "mastery": 60},
		{"name": "异常值识别", "count": 3, "mastery": 45},
		{"name": "重复值处理", "count": 2, "mastery": 80},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": weakPoints,
	})
}

// GetWeakPointExplain 获取薄弱点讲解
// POST /api/weakPoint/getExplain
func (h *WeakPointHandler) GetWeakPointExplain(c *gin.Context) {
	var req struct {
		WeakPointName string `json:"weakPointName" binding:"required"`
		StudentID     string `json:"studentId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 根据不同薄弱点返回不同讲解内容
	explains := map[string]gin.H{
		"缺失值填充": {
			"title": "缺失值填充 · 知识点讲解",
			"content": "缺失值是数据中为空的部分，常用方法：\n" +
				"1. fillna() 填充常数、均值、中位数\n" +
				"2. interpolate() 线性插值（适合时序）\n" +
				"3. dropna() 直接删除行/列",
			"examples": []string{
				"df.fillna(0)  # 用0填充",
				"df.interpolate()  # 线性插值",
				"df.dropna()  # 删除缺失值",
			},
		},
		"异常值识别": {
			"title": "异常值识别 · 知识点讲解",
			"content": "异常值识别常用方法：\n" +
				"1. Z-Score法：|Z|>3 视为异常\n" +
				"2. IQR法：超出Q1-1.5IQR或Q3+1.5IQR视为异常",
			"examples": []string{
				"from scipy import stats\nz_scores = stats.zscore(data)",
				"Q1 = df.quantile(0.25)\nQ3 = df.quantile(0.75)\nIQR = Q3 - Q1",
			},
		},
	}

	data, exists := explains[req.WeakPointName]
	if !exists {
		data = gin.H{
			"title":   req.WeakPointName + " · 知识点讲解",
			"content": "这是" + req.WeakPointName + "的详细讲解内容...",
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// GetWeakPointTest 生成检测习题
// POST /api/weakPoint/getTest
func (h *WeakPointHandler) GetWeakPointTest(c *gin.Context) {
	var req struct {
		WeakPointName string `json:"weakPointName" binding:"required"`
		StudentID     string `json:"studentId" binding:"required"`
		QuestionType  string `json:"questionType"` // single/multiple
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 根据不同薄弱点生成不同习题
	tests := map[string]gin.H{
		"缺失值填充": {
			"questionId": "Q001",
			"content":    "处理缺失值时，以下哪种方法最适合时间序列数据？",
			"type":       "single",
			"options": []string{
				"A. fillna(0) 用0填充",
				"B. interpolate() 线性插值",
				"C. dropna() 删除缺失值",
				"D. fillna(method='ffill') 向前填充",
			},
		},
	}

	data, exists := tests[req.WeakPointName]
	if !exists {
		data = gin.H{
			"questionId": "Q002",
			"content":    "关于" + req.WeakPointName + "的正确说法是？",
			"type":       "single",
			"options":    []string{"A. 选项1", "B. 选项2", "C. 选项3", "D. 选项4"},
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// CheckAnswer 校验答案
// POST /api/weakPoint/checkAnswer
func (h *WeakPointHandler) CheckAnswer(c *gin.Context) {
	var req struct {
		StudentID  string `json:"studentId" binding:"required"`
		QuestionID string `json:"questionId" binding:"required"`
		UserAnswer string `json:"userAnswer" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 校验答案（实际应从数据库查询正确答案）
	isCorrect := false
	correctAnswer := "B"
	explanation := ""

	if req.QuestionID == "Q001" {
		correctAnswer = "B"
		isCorrect = (req.UserAnswer == "B" || req.UserAnswer == "interpolate()")
		explanation = "interpolate() 线性插值适合时间序列数据，可以基于前后值推算缺失值。"
	}

	// 更新掌握度
	masteryDelta := 0
	if isCorrect {
		masteryDelta = 10
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"isCorrect":     isCorrect,
			"correctAnswer": correctAnswer,
			"explanation":   explanation,
			"masteryDelta":  masteryDelta,
		},
	})
}
