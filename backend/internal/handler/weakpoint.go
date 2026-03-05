package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WeakPointHandler struct {
	db *gorm.DB
}

func NewWeakPointHandler(db *gorm.DB) *WeakPointHandler {
	return &WeakPointHandler{db: db}
}

// 4.1 获取个人薄弱点列表
// GET /api/v1/student/coursewares/:courseId/weak-points
func (h *WeakPointHandler) GetWeakPointList(c *gin.Context) {
	// courseId := c.Param("courseId") // 使用 courseId 参数
	studentId := c.Query("studentId")

	if studentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少参数: studentId",
		})
		return
	}

	// 模拟数据
	weakPoints := []gin.H{
		{
			"id":          "wp_001",
			"name":        "缺失值填充",
			"description": "数据中空值的处理方法",
			"count":       5,
			"mastery":     60,
			"createdAt":   time.Now().AddDate(0, 0, -7),
		},
		{
			"id":          "wp_002",
			"name":        "异常值识别",
			"description": "识别数据中的异常值",
			"count":       3,
			"mastery":     45,
			"createdAt":   time.Now().AddDate(0, 0, -5),
		},
		{
			"id":          "wp_003",
			"name":        "重复值处理",
			"description": "处理数据中的重复记录",
			"count":       2,
			"mastery":     80,
			"createdAt":   time.Now().AddDate(0, 0, -3),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data":    weakPoints,
	})
}

// 4.2 薄弱点 AI 详细讲解
// GET /api/v1/student/weak-points/:weakPointId/explain
func (h *WeakPointHandler) GetWeakPointExplain(c *gin.Context) {
	weakPointId := c.Param("weakPointId")

	// 根据不同的薄弱点ID返回不同的讲解内容
	var data gin.H

	switch weakPointId {
	case "wp_001":
		data = gin.H{
			"id":    "wp_001",
			"name":  "缺失值填充",
			"title": "缺失值填充 · 知识点讲解",
			"content": "缺失值是数据中为空的部分，常用处理方法：\n\n" +
				"1. fillna() 填充常数、均值、中位数\n" +
				"   df.fillna(0)                 # 用0填充\n" +
				"   df.fillna(df.mean())         # 用均值填充\n" +
				"   df.fillna(df.median())       # 用中位数填充\n" +
				"   df.fillna(method='ffill')    # 用前一个值填充\n\n" +
				"2. interpolate() 线性插值（适合时序）\n" +
				"   df.interpolate()             # 线性插值\n" +
				"   df.interpolate(method='quadratic')  # 二次插值\n\n" +
				"3. dropna() 直接删除行/列\n" +
				"   df.dropna()                  # 删除任何包含缺失值的行\n" +
				"   df.dropna(axis=1)            # 删除任何包含缺失值的列\n" +
				"   df.dropna(thresh=2)          # 至少保留2个非空值的行\n\n" +
				"选择建议：\n" +
				"- 时序数据：优先使用 interpolate()\n" +
				"- 分类数据：可用众数填充或用'unknown'填充\n" +
				"- 连续数据：可用均值或中位数填充\n" +
				"- 缺失过多：考虑直接删除",
			"examples": []string{
				"df.fillna(0)",
				"df.interpolate()",
				"df.dropna()",
			},
			"relatedConcepts": []string{"数据清洗", "异常值处理", "数据预处理"},
		}
	case "wp_002":
		data = gin.H{
			"id":    "wp_002",
			"name":  "异常值识别",
			"title": "异常值识别 · 知识点讲解",
			"content": "异常值是数据中明显偏离其他观测值的点，常用识别方法：\n\n" +
				"1. Z-Score法\n" +
				"   from scipy import stats\n" +
				"   import numpy as np\n" +
				"   z_scores = np.abs(stats.zscore(data))\n" +
				"   threshold = 3\n" +
				"   outliers = np.where(z_scores > threshold)\n\n" +
				"2. IQR法（四分位距法）\n" +
				"   Q1 = df.quantile(0.25)\n" +
				"   Q3 = df.quantile(0.75)\n" +
				"   IQR = Q3 - Q1\n" +
				"   lower_bound = Q1 - 1.5 * IQR\n" +
				"   upper_bound = Q3 + 1.5 * IQR\n" +
				"   outliers = df[(df < lower_bound) | (df > upper_bound)]\n\n" +
				"3. 可视化方法\n" +
				"   - 箱线图（Box Plot）\n" +
				"   - 散点图（Scatter Plot）\n" +
				"   - 直方图（Histogram）\n\n" +
				"处理策略：\n" +
				"- 删除异常值（谨慎使用）\n" +
				"- 替换为边界值（盖帽法）\n" +
				"- 转换变量（如对数变换）",
			"examples": []string{
				"z_scores = np.abs(stats.zscore(data))",
				"Q1, Q3 = df.quantile([0.25, 0.75])",
				"df[~((df < lower_bound) | (df > upper_bound))]",
			},
			"relatedConcepts": []string{"数据清洗", "数据分布", "箱线图"},
		}
	case "wp_003":
		data = gin.H{
			"id":    "wp_003",
			"name":  "重复值处理",
			"title": "重复值处理 · 知识点讲解",
			"content": "重复值是指数据集中完全相同的记录，处理方法：\n\n" +
				"1. 识别重复值\n" +
				"   df.duplicated()              # 标记重复行\n" +
				"   df.duplicated(subset=['col1', 'col2'])  # 基于特定列\n\n" +
				"2. 查看重复值\n" +
				"   df[df.duplicated()]          # 查看重复行\n" +
				"   df.duplicated().sum()        # 统计重复数量\n\n" +
				"3. 删除重复值\n" +
				"   df.drop_duplicates()         # 删除完全重复的行\n" +
				"   df.drop_duplicates(subset=['col1'])  # 基于特定列删除\n" +
				"   df.drop_duplicates(keep='last')      # 保留最后一个\n\n" +
				"4. 聚合重复值\n" +
				"   df.groupby(df.columns.tolist()).agg({\n" +
				"       'col': 'first'\n" +
				"   })\n\n" +
				"注意事项：\n" +
				"- 先确认重复是否合理\n" +
				"- 时间序列数据可能需要保留\n" +
				"- 考虑数据来源的业务含义",
			"examples": []string{
				"df.duplicated().sum()",
				"df.drop_duplicates()",
				"df.drop_duplicates(subset=['id'])",
			},
			"relatedConcepts": []string{"数据去重", "数据聚合", "数据清洗"},
		}
	default:
		data = gin.H{
			"id":              weakPointId,
			"name":            "未知知识点",
			"title":           "知识点讲解",
			"content":         "这是该知识点的详细讲解内容...",
			"examples":        []string{},
			"relatedConcepts": []string{},
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data":    data,
	})
}

// 4.3 生成随堂检测题
// POST /api/v1/student/weak-points/:weakPointId/generate-test
func (h *WeakPointHandler) GenerateTest(c *gin.Context) {
	weakPointId := c.Param("weakPointId")

	var req struct {
		QuestionType string `json:"questionType"` // single/multiple
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.QuestionType = "single"
	}

	// 根据薄弱点ID生成不同的题目
	questionId := "q_" + uuid.New().String()[:8]

	var testData gin.H

	switch weakPointId {
	case "wp_001":
		if req.QuestionType == "single" {
			testData = gin.H{
				"questionId": questionId,
				"type":       "single",
				"content":    "处理缺失值时，以下哪种方法最适合时间序列数据？",
				"options": []gin.H{
					{"key": "A", "value": "fillna(0) 用0填充"},
					{"key": "B", "value": "interpolate() 线性插值"},
					{"key": "C", "value": "dropna() 删除缺失值"},
					{"key": "D", "value": "fillna(method='ffill') 向前填充"},
				},
			}
		} else {
			testData = gin.H{
				"questionId": questionId,
				"type":       "multiple",
				"content":    "以下哪些是处理缺失值的常用方法？",
				"options": []gin.H{
					{"key": "A", "value": "fillna()"},
					{"key": "B", "value": "interpolate()"},
					{"key": "C", "value": "dropna()"},
					{"key": "D", "value": "replace()"},
				},
			}
		}
	case "wp_002":
		testData = gin.H{
			"questionId": questionId,
			"type":       "single",
			"content":    "在IQR法中，异常值的判断标准是什么？",
			"options": []gin.H{
				{"key": "A", "value": "超出 Q1-1.5*IQR 或 Q3+1.5*IQR"},
				{"key": "B", "value": "超出 Q1-2*IQR 或 Q3+2*IQR"},
				{"key": "C", "value": "Z-Score > 2"},
				{"key": "D", "value": "Z-Score > 1"},
			},
		}
	case "wp_003":
		testData = gin.H{
			"questionId": questionId,
			"type":       "single",
			"content":    "以下哪个函数用于删除重复值？",
			"options": []gin.H{
				{"key": "A", "value": "df.unique()"},
				{"key": "B", "value": "df.drop_duplicates()"},
				{"key": "C", "value": "df.duplicated()"},
				{"key": "D", "value": "df.nunique()"},
			},
		}
	default:
		testData = gin.H{
			"questionId": questionId,
			"type":       "single",
			"content":    "关于该知识点的正确说法是？",
			"options": []gin.H{
				{"key": "A", "value": "选项A"},
				{"key": "B", "value": "选项B"},
				{"key": "C", "value": "选项C"},
				{"key": "D", "value": "选项D"},
			},
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data":    testData,
	})
}

// 4.4 提交并校验答案
// POST /api/v1/student/tests/:questionId/check
func (h *WeakPointHandler) CheckAnswer(c *gin.Context) {
	questionId := c.Param("questionId")

	var req struct {
		UserAnswer string `json:"userAnswer" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少参数: userAnswer",
		})
		return
	}

	// 根据题目ID判断正确答案
	var isCorrect bool
	var correctAnswer string
	var explanation string
	var masteryDelta int

	// 简单判断，实际应该从数据库查询
	if len(questionId) > 5 {
		if req.UserAnswer == "B" || req.UserAnswer == "interpolate()" {
			isCorrect = true
			correctAnswer = "B"
			explanation = "interpolate() 线性插值适合时间序列数据，可以基于前后值推算缺失值。"
			masteryDelta = 10
		} else if req.UserAnswer == "A" {
			isCorrect = true
			correctAnswer = "A"
			explanation = "IQR法中，超出Q1-1.5*IQR或Q3+1.5*IQR的值被视为异常值。"
			masteryDelta = 10
		} else {
			isCorrect = false
			correctAnswer = "B"
			explanation = "答案不正确，请重新学习相关知识。"
			masteryDelta = -5
		}
	} else {
		// 默认答案
		correctAnswer = "A"
		isCorrect = (req.UserAnswer == "A")
		explanation = "这是正确答案的解析。"
		masteryDelta = 5
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "请求成功",
		"data": gin.H{
			"isCorrect":     isCorrect,
			"correctAnswer": correctAnswer,
			"explanation":   explanation,
			"masteryDelta":  masteryDelta,
			"newMastery":    65 + masteryDelta,
		},
	})
}
