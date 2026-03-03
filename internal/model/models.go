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
	Title     string `gorm:"size:200;not null;index" json:"title"`
	FileURL   string `gorm:"size:500" json:"file_url"`
	FileType  string `gorm:"size:20" json:"file_type"` // ppt, pdf, pptx
	TotalPage int    `gorm:"default:0" json:"total_page"`

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
