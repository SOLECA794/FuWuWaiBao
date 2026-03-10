package dto

// TracePoint 通用的页面坐标信息，可用于 AI 溯源提问 (ai_partner)、做笔记 (learning_record) 等场景
type TracePoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
