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
	Pages         []CoursePage   `gorm:"foreignKey:CourseID" json:"pages,omitempty"`
	TeachingNodes []TeachingNode `gorm:"foreignKey:CourseID" json:"teaching_nodes,omitempty"`
}

// CoursePage 课件页表
type CoursePage struct {
	BaseModel
	CourseID         string `gorm:"size:36;not null;index:idx_course_page,unique" json:"course_id"`
	PageIndex        int    `gorm:"not null;index:idx_course_page,unique" json:"page_index"`
	ImageURL         string `gorm:"size:500" json:"image_url"`
	SourceText       string `gorm:"type:text" json:"source_text"`
	ScriptText       string `gorm:"type:text" json:"script_text"`
	AudioURL         string `gorm:"size:500" json:"audio_url"`
	AudioStatus      string `gorm:"size:30;default:'not_generated'" json:"audio_status"`
	AudioProvider    string `gorm:"size:60" json:"audio_provider"`
	AudioDurationSec int    `gorm:"default:0" json:"audio_duration_sec"`
}

// TeachingNode 教学节点表
type TeachingNode struct {
	BaseModel
	CourseID             string `gorm:"size:36;not null;index:idx_course_teaching_node,priority:1" json:"course_id"`
	NodeID               string `gorm:"size:100;not null;index:idx_course_teaching_node,priority:2,unique" json:"node_id"`
	ChapterTitle         string `gorm:"size:200" json:"chapter_title"`
	PageIndex            int    `gorm:"default:0;index" json:"page_index"`
	EstimatedDuration    int    `gorm:"default:0" json:"estimated_duration"`
	Title                string `gorm:"size:200;not null" json:"title"`
	Summary              string `gorm:"type:text" json:"summary"`
	SourcePages          string `gorm:"type:text" json:"source_pages"`
	CorePoints           string `gorm:"type:text" json:"core_points"`
	Examples             string `gorm:"type:text" json:"examples"`
	CommonConfusions     string `gorm:"type:text" json:"common_confusions"`
	ScriptText           string `gorm:"type:text" json:"script_text"`
	ReteachScript        string `gorm:"type:text" json:"reteach_script"`
	InteractiveQuestions string `gorm:"type:text" json:"interactive_questions"`
	TransitionText       string `gorm:"type:text" json:"transition_text"`
	MindmapMarkdown      string `gorm:"type:text" json:"mindmap_markdown"`
	AudioURL             string `gorm:"size:500" json:"audio_url"`
	AudioDurationSec     int    `gorm:"default:0" json:"audio_duration_sec"`
	AudioStartSec        int    `gorm:"default:0" json:"audio_start_sec"`
	AudioEndSec          int    `gorm:"default:0" json:"audio_end_sec"`
	TTSStatus            string `gorm:"size:30;default:'not_generated'" json:"tts_status"`
	VoiceProfile         string `gorm:"size:60" json:"voice_profile"`
	SortOrder            int    `gorm:"default:0" json:"sort_order"`
}

// UserProgress 用户进度表
type UserProgress struct {
	BaseModel
	UserID   string `gorm:"size:36;not null;index:idx_user_course,unique" json:"user_id"`
	CourseID string `gorm:"size:36;not null;index:idx_user_course,unique" json:"course_id"`
	LastPage int    `gorm:"default:0" json:"last_page"`
}

// DialogueSession 问答会话表
type DialogueSession struct {
	BaseModel
	UserID         string     `gorm:"size:36;index" json:"user_id"`
	CourseID       string     `gorm:"size:36;not null;index" json:"course_id"`
	CurrentPage    int        `gorm:"default:1" json:"current_page"`
	CurrentNodeID  string     `gorm:"size:100" json:"current_node_id"`
	CurrentTimeSec int        `gorm:"default:0" json:"current_time_sec"`
	PlaybackMode   string     `gorm:"size:30;default:'timeline'" json:"playback_mode"`
	LastAskedAt    *time.Time `json:"last_asked_at,omitempty"`
}

// AudioAsset 音频资产表
type AudioAsset struct {
	BaseModel
	CourseID         string `gorm:"size:36;not null;index:idx_audio_asset,priority:1" json:"course_id"`
	PageIndex        int    `gorm:"default:1;index:idx_audio_asset,priority:2" json:"page_index"`
	NodeID           string `gorm:"size:100;index:idx_audio_asset,priority:3" json:"node_id"`
	Provider         string `gorm:"size:60" json:"provider"`
	VoiceType        string `gorm:"size:60" json:"voice_type"`
	Format           string `gorm:"size:20;default:'mp3'" json:"format"`
	Status           string `gorm:"size:30;default:'placeholder'" json:"status"`
	AudioURL         string `gorm:"size:500" json:"audio_url"`
	DurationSec      int    `gorm:"default:0" json:"duration_sec"`
	StartSec         int    `gorm:"default:0" json:"start_sec"`
	EndSec           int    `gorm:"default:0" json:"end_sec"`
	SourceScriptHash string `gorm:"size:64" json:"source_script_hash"`
}

// PlatformUser 平台用户表
type PlatformUser struct {
	BaseModel
	PlatformID      string `gorm:"size:64;index" json:"platform_id"`
	ExternalID      string `gorm:"size:64;uniqueIndex" json:"external_id"`
	Username        string `gorm:"size:80;index" json:"username"`
	DisplayName     string `gorm:"size:120" json:"display_name"`
	Email           string `gorm:"size:120;index" json:"email"`
	Phone           string `gorm:"size:32;index" json:"phone"`
	Role            string `gorm:"size:30;default:'student';index" json:"role"`
	Status          string `gorm:"size:30;default:'active'" json:"status"`
	OrgCode         string `gorm:"size:64;index" json:"org_code"`
	SchoolName      string `gorm:"size:160" json:"school_name"`
	Major           string `gorm:"size:120" json:"major"`
	Grade           string `gorm:"size:40" json:"grade"`
	ClassExternalID string `gorm:"size:64;index" json:"class_external_id"`
	ClassName       string `gorm:"size:160" json:"class_name"`
	AvatarURL       string `gorm:"size:500" json:"avatar_url"`
}

// TeachingCourse 教学课程表
type TeachingCourse struct {
	BaseModel
	PlatformID  string     `gorm:"size:64;index" json:"platform_id"`
	ExternalID  string     `gorm:"size:64;uniqueIndex" json:"external_id"`
	Code        string     `gorm:"size:80;index" json:"code"`
	Title       string     `gorm:"size:200;not null;index" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	TeacherID   string     `gorm:"size:36;index" json:"teacher_id"`
	OrgCode     string     `gorm:"size:64;index" json:"org_code"`
	SchoolName  string     `gorm:"size:160" json:"school_name"`
	Status      string     `gorm:"size:30;default:'draft';index" json:"status"`
	Semester    string     `gorm:"size:60" json:"semester"`
	Credit      float64    `gorm:"default:0" json:"credit"`
	Period      int        `gorm:"default:0" json:"period"`
	CoverURL    string     `gorm:"size:500" json:"cover_url"`
	StartsAt    *time.Time `json:"starts_at,omitempty"`
	EndsAt      *time.Time `json:"ends_at,omitempty"`
}

// CourseClass 教学班级表
type CourseClass struct {
	BaseModel
	PlatformID       string `gorm:"size:64;index" json:"platform_id"`
	ExternalID       string `gorm:"size:64;uniqueIndex" json:"external_id"`
	TeachingCourseID string `gorm:"size:36;not null;index" json:"teaching_course_id"`
	TeacherID        string `gorm:"size:36;index" json:"teacher_id"`
	ClassName        string `gorm:"size:160;not null" json:"class_name"`
	ClassCode        string `gorm:"size:80;index" json:"class_code"`
	Semester         string `gorm:"size:60" json:"semester"`
	Grade            string `gorm:"size:40" json:"grade"`
	Major            string `gorm:"size:120" json:"major"`
	Capacity         int    `gorm:"default:0" json:"capacity"`
	Status           string `gorm:"size:30;default:'active'" json:"status"`
}

// CourseEnrollment 选课关系表
type CourseEnrollment struct {
	BaseModel
	PlatformID       string     `gorm:"size:64;index" json:"platform_id"`
	ExternalID       string     `gorm:"size:64;index" json:"external_id"`
	TeachingCourseID string     `gorm:"size:36;not null;index:idx_enrollment_scope,priority:1" json:"teaching_course_id"`
	CourseClassID    string     `gorm:"size:36;index:idx_enrollment_scope,priority:2" json:"course_class_id"`
	UserID           string     `gorm:"size:36;not null;index:idx_enrollment_scope,priority:3" json:"user_id"`
	Role             string     `gorm:"size:30;default:'student'" json:"role"`
	Status           string     `gorm:"size:30;default:'active'" json:"status"`
	EnrolledAt       *time.Time `json:"enrolled_at,omitempty"`
}

// DialogueTurn 问答轮次表
type DialogueTurn struct {
	BaseModel
	SessionID          string `gorm:"size:36;not null;index" json:"session_id"`
	CourseID           string `gorm:"size:36;not null;index" json:"course_id"`
	UserID             string `gorm:"size:36;index" json:"user_id"`
	TurnIndex          int    `gorm:"default:1" json:"turn_index"`
	PageIndex          int    `gorm:"default:1;index" json:"page_index"`
	NodeID             string `gorm:"size:100" json:"node_id"`
	Question           string `gorm:"type:text;not null" json:"question"`
	Answer             string `gorm:"type:text" json:"answer"`
	SourcePage         int    `gorm:"default:1" json:"source_page"`
	NeedReteach        bool   `gorm:"default:false" json:"need_reteach"`
	FollowUpSuggestion string `gorm:"type:text" json:"follow_up_suggestion"`
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
	NodeID   string `gorm:"size:100;index" json:"node_id"`
	PageNum  int    `gorm:"not null" json:"page_num"`
	Note     string `gorm:"type:text" json:"note"`
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

// PracticeTask 练习任务
type PracticeTask struct {
	BaseModel
	UserID      string `gorm:"size:36;not null;index" json:"user_id"`
	CourseID    string `gorm:"size:36;not null;index" json:"course_id"`
	NodeID      string `gorm:"size:100;index" json:"node_id"`
	PageIndex   int    `gorm:"default:1;index" json:"page_index"`
	Difficulty  int    `gorm:"default:2" json:"difficulty"`
	QuestionIDs string `gorm:"type:text" json:"question_ids"` // JSON 数组
	Status      string `gorm:"size:30;default:'pending'" json:"status"`
}

// PracticeAttempt 练习作答记录
type PracticeAttempt struct {
	BaseModel
	TaskID       string `gorm:"size:36;not null;index" json:"task_id"`
	UserID       string `gorm:"size:36;not null;index" json:"user_id"`
	CourseID     string `gorm:"size:36;not null;index" json:"course_id"`
	NodeID       string `gorm:"size:100;index" json:"node_id"`
	TotalCount   int    `gorm:"default:0" json:"total_count"`
	CorrectCount int    `gorm:"default:0" json:"correct_count"`
	Score        int    `gorm:"default:0" json:"score"`
	AnswersJSON  string `gorm:"type:text" json:"answers_json"`
}

// NodeFavorite 节点/页面收藏
type NodeFavorite struct {
	BaseModel
	UserID    string `gorm:"size:36;not null;index:idx_favorite_scope,priority:1" json:"user_id"`
	CourseID  string `gorm:"size:36;not null;index:idx_favorite_scope,priority:2" json:"course_id"`
	NodeID    string `gorm:"size:100;index:idx_favorite_scope,priority:3" json:"node_id"`
	PageIndex int    `gorm:"default:1;index:idx_favorite_scope,priority:4" json:"page_index"`
	Title     string `gorm:"size:200" json:"title"`
}

// User 系统用户表（用于登录认证）
type User struct {
	BaseModel
	Username     string `gorm:"size:100;uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	// 角色：teacher / student
	Role string `gorm:"size:20;not null;default:'teacher'" json:"role"`
}
