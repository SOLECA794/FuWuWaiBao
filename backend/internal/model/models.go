package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 自动生成UUID
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// Course 课件表
type Course struct {
	BaseModel
	Title        string     `gorm:"size:200;not null;index" json:"title"`
	FileURL      string     `gorm:"size:500" json:"file_url"`
	FileType     string     `gorm:"size:20" json:"file_type"` // ppt, pdf, pptx
	TotalPage    int        `gorm:"default:0" json:"total_page"`
	IsPublished  bool       `gorm:"default:false" json:"is_published"`
	PublishScope string     `gorm:"size:50;default:'all'" json:"publish_scope"`
	PublishedAt  *time.Time `json:"published_at,omitempty"`

	// 关联
	Pages []CoursePage `gorm:"foreignKey:CourseID" json:"pages,omitempty"`
}

// CoursePage 课件页表
type CoursePage struct {
	BaseModel
	CourseID   string `gorm:"size:36;not null;index:idx_course_page,unique" json:"course_id"`
	PageIndex  int    `gorm:"not null;index:idx_course_page,unique" json:"page_index"`
	ImageURL   string `gorm:"size:500" json:"image_url"`
	ScriptText string `gorm:"type:text" json:"script_text"`
	AudioURL   string `gorm:"size:500" json:"audio_url"`
}

// UserProgress 用户进度表
type UserProgress struct {
	BaseModel
	UserID   string `gorm:"size:36;not null;index:idx_user_course,unique" json:"user_id"`
	CourseID string `gorm:"size:36;not null;index:idx_user_course,unique" json:"course_id"`
	LastPage int    `gorm:"default:0" json:"last_page"`
}

// QuestionLog 提问日志表
type QuestionLog struct {
	BaseModel
	UserID    string `gorm:"size:36;not null;index" json:"user_id"`
	CourseID  string `gorm:"size:36;not null;index" json:"course_id"`
	PageIndex int    `gorm:"not null" json:"page_index"`
	Question  string `gorm:"type:text;not null" json:"question"`
	Answer    string `gorm:"type:text" json:"answer"`
}

// TeacherEdit 教师编辑记录表
type TeacherEdit struct {
	BaseModel
	TeacherID string `gorm:"size:36;not null" json:"teacher_id"`
	CourseID  string `gorm:"size:36;not null" json:"course_id"`
	PageIndex int    `json:"page_index"`
	OldScript string `gorm:"type:text" json:"old_script"`
	NewScript string `gorm:"type:text" json:"new_script"`
}

// MindMapNode 思维导图节点表
type MindMapNode struct {
	BaseModel
	CourseID  string `gorm:"size:36;not null;index" json:"course_id"`
	ParentID  string `gorm:"size:36" json:"parent_id"`
	Title     string `gorm:"size:200;not null" json:"title"`
	Content   string `gorm:"type:text" json:"content"`
	PageIndex int    `gorm:"default:0" json:"page_index"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
}

// StudentNote 学生笔记表
type StudentNote struct {
	BaseModel
	UserID   string `gorm:"size:36;not null;index" json:"user_id"`
	CourseID string `gorm:"size:36;not null;index" json:"course_id"`
	PageNum  int    `gorm:"not null" json:"page_num"`
	Note     string `gorm:"type:text" json:"note"`
}

// StudentFavorite 收藏模型
type StudentFavorite struct {
	BaseModel
	UserID   string `gorm:"size:36;not null;index:idx_user_course_node,unique" json:"user_id"`
	CourseID string `gorm:"size:36;index:idx_user_course_node,unique" json:"course_id"`
	NodeID   string `gorm:"size:100;index:idx_user_course_node,unique" json:"node_id"`
	PageNum  int    `gorm:"default:0" json:"page_num"`
	Title    string `gorm:"size:300" json:"title"`
}

// PracticeTask 学生练习任务
type PracticeTask struct {
	BaseModel
	TaskID     string `gorm:"size:100;not null;index:idx_task_user,unique" json:"task_id"`
	UserID     string `gorm:"size:36;not null;index:idx_task_user,unique" json:"user_id"`
	CourseID   string `gorm:"size:36;index" json:"course_id"`
	NodeID     string `gorm:"size:100" json:"node_id"`
	PageNum    int    `json:"page_num"`
	Difficulty int    `json:"difficulty"`
	Count      int    `json:"count"`
	Questions  string `gorm:"type:text" json:"questions"` // JSON字符串
}

// PracticeAttempt 练习提交记录
type PracticeAttempt struct {
	BaseModel
	TaskID      string `gorm:"size:100;not null;index:idx_attempt_task_user,unique" json:"task_id"`
	UserID      string `gorm:"size:36;not null;index:idx_attempt_task_user,unique" json:"user_id"`
	Score       int    `json:"score"`
	Correct     int    `json:"correct"`
	Total       int    `json:"total"`
	Details     string `gorm:"type:text" json:"details"` // JSON字符串
}

// WeakPoint 薄弱点模型
type WeakPoint struct {
	BaseModel
	StudentID    string `gorm:"size:36;not null;index" json:"student_id"`
	CourseID     string `gorm:"size:36;not null;index" json:"course_id"`
	Name         string `gorm:"size:100;not null" json:"name"`
	Description  string `gorm:"type:text" json:"description"`
	Count        int    `gorm:"default:0" json:"count"`         // 出现次数
	MasteryLevel int    `gorm:"default:0" json:"mastery_level"` // 掌握程度 0-100
}

// KnowledgePoint 知识点模型
type KnowledgePoint struct {
	BaseModel
	CourseID string `gorm:"size:36;not null;index" json:"course_id"`
	ParentID string `gorm:"size:36" json:"parent_id"` // 父知识点ID
	Name     string `gorm:"size:100;not null" json:"name"`
	Level    int    `gorm:"default:1" json:"level"` // 1:章节 2:知识点 3:子知识点
	Content  string `gorm:"type:text" json:"content"`
	Examples string `gorm:"type:text" json:"examples"` // 示例代码/案例
}

// Question 习题模型
type Question struct {
	BaseModel
	WeakPointID  string `gorm:"size:36;index" json:"weak_point_id"`
	QuestionType string `gorm:"size:20" json:"question_type"` // single/multiple
	Content      string `gorm:"type:text;not null" json:"content"`
	Options      string `gorm:"type:text" json:"options"`     // JSON格式的选项
	Answer       string `gorm:"type:text;not null" json:"-"`  // 正确答案
	Explanation  string `gorm:"type:text" json:"explanation"` // 解析
	Difficulty   int    `gorm:"default:1" json:"difficulty"`  // 1-5
}

// AnswerRecord 答题记录
type AnswerRecord struct {
	BaseModel
	StudentID    string `gorm:"size:36;not null;index" json:"student_id"`
	QuestionID   string `gorm:"size:36;not null" json:"question_id"`
	UserAnswer   string `gorm:"type:text" json:"user_answer"`
	IsCorrect    bool   `json:"is_correct"`
	MasteryDelta int    `json:"mastery_delta"` // 掌握度变化
}
