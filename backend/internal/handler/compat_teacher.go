package handler

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"smart-teaching-backend/internal/model"
	"smart-teaching-backend/internal/service"
	"smart-teaching-backend/pkg/apiresp"
)

func (h *CompatibilityHandler) DeleteCoursewareV1(c *gin.Context) {
	courseID := c.Param("courseId")
	if err := h.courseService.DeleteCourse(courseID); err != nil {
		apiresp.Internal(c, "删除失败", "")
		return
	}
	apiresp.OKMessage(c, "删除成功")
}

func (h *CompatibilityHandler) GetTeacherScriptV1(c *gin.Context) {
	courseID := c.Param("courseId")
	pageStr := c.Param("pageNum")
	c.Params = append(filterParams(c.Params, "page"), gin.Param{Key: "page", Value: pageStr}, gin.Param{Key: "courseId", Value: courseID})
	NewTeacherHandler(h.db, h.aiClient).GetScript(c)
}

func (h *CompatibilityHandler) UpdateTeacherScriptV1(c *gin.Context) {
	courseID := c.Param("courseId")
	pageStr := c.Param("pageNum")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		apiresp.BadRequest(c, "页码错误", "")
		return
	}
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.BadRequest(c, "参数错误", "")
		return
	}
	var coursePage model.CoursePage
	err = h.db.Where("course_id = ? AND page_index = ?", courseID, page).First(&coursePage).Error
	if err != nil {
		coursePage = model.CoursePage{CourseID: courseID, PageIndex: page, ScriptText: req.Content}
		err = h.db.Create(&coursePage).Error
	} else {
		err = h.db.Model(&coursePage).Update("script_text", req.Content).Error
	}
	if err != nil {
		apiresp.Internal(c, "保存失败", "")
		return
	}
	apiresp.OKMessage(c, "保存成功")
}

func (h *CompatibilityHandler) GetTeacherNodesV1(c *gin.Context) {
	courseID := c.Param("courseId")
	pageStr := c.Param("pageNum")
	params := filterParams(filterParams(c.Params, "page"), "courseId")
	c.Params = append(params, gin.Param{Key: "page", Value: pageStr}, gin.Param{Key: "courseId", Value: courseID})
	NewTeacherHandler(h.db, h.aiClient).GetTeachingNodes(c)
}

func (h *CompatibilityHandler) UpdateTeacherNodesV1(c *gin.Context) {
	courseID := c.Param("courseId")
	pageStr := c.Param("pageNum")
	params := filterParams(filterParams(c.Params, "page"), "courseId")
	c.Params = append(params, gin.Param{Key: "page", Value: pageStr}, gin.Param{Key: "courseId", Value: courseID})
	NewTeacherHandler(h.db, h.aiClient).UpdateTeachingNodes(c)
}

func (h *CompatibilityHandler) GenerateTeacherAudioV1(c *gin.Context) {
	courseID := c.Param("courseId")
	pageStr := c.Param("pageNum")
	params := filterParams(filterParams(c.Params, "page"), "courseId")
	c.Params = append(params, gin.Param{Key: "page", Value: pageStr}, gin.Param{Key: "courseId", Value: courseID})
	NewTeacherHandler(h.db, h.aiClient).GeneratePageAudio(c)
}

func (h *CompatibilityHandler) AIGenerateTeacherScriptV1(c *gin.Context) {
	courseID := c.Param("courseId")
	var req struct {
		PageNum int    `json:"pageNum" binding:"required"`
		Mode    string `json:"mode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.BadRequest(c, "参数错误", "")
		return
	}
	var course model.Course
	if err := h.db.First(&course, "id = ?", courseID).Error; err != nil {
		apiresp.NotFound(c, "课件不存在", "")
		return
	}
	var page model.CoursePage
	contextText := ""
	if err := h.db.Where("course_id = ? AND page_index = ?", courseID, req.PageNum).First(&page).Error; err == nil {
		contextText = pageContextText(page)
	}
	script := ""
	mindmapMarkdown := ""
	teachingNodes := loadTeachingNodesByPage(h.db, courseID, req.PageNum)
	if generatedScript, generatedMindmap, usedNodes, err := generateAndStoreTeachingNodeScripts(c.Request.Context(), h.db, h.aiClient, course.Title, defaultString(req.Mode, "llm"), teachingNodes); usedNodes && err == nil {
		script = generatedScript
		mindmapMarkdown = generatedMindmap
	} else {
		resp, err := h.aiClient.GenerateScript(c.Request.Context(), service.GenerateScriptRequest{Page: req.PageNum, Content: contextText, CourseName: course.Title, Mode: defaultString(req.Mode, "llm")})
		if err != nil {
			apiresp.Internal(c, "AI生成失败", "")
			return
		}
		script = resp.Script
		mindmapMarkdown = resp.MindmapMarkdown
	}
	if err := h.db.Where("course_id = ? AND page_index = ?", courseID, req.PageNum).First(&page).Error; err == nil {
		_ = h.db.Model(&page).Update("script_text", script).Error
	} else {
		_ = h.db.Create(&model.CoursePage{CourseID: courseID, PageIndex: req.PageNum, ScriptText: script}).Error
	}
	apiresp.OK(c, "请求成功", gin.H{"courseId": courseID, "pageNum": req.PageNum, "content": script, "mindmapMarkdown": mindmapMarkdown})
}

func (h *CompatibilityHandler) PublishCoursewareV1(c *gin.Context) {
	courseID := c.Param("courseId")
	var req struct {
		Scope string `json:"scope"`
	}
	_ = c.ShouldBindJSON(&req)
	if strings.TrimSpace(req.Scope) == "" {
		req.Scope = "all"
	}
	now := time.Now()
	updates := map[string]any{
		"is_published":  true,
		"publish_scope": req.Scope,
		"published_at":  now,
	}
	if err := h.db.Model(&model.Course{}).Where("id = ?", courseID).Updates(updates).Error; err != nil {
		apiresp.Internal(c, "发布失败", "")
		return
	}
	apiresp.OK(c, "发布成功", gin.H{"courseId": courseID, "scope": req.Scope, "publishedAt": now.Format("2006-01-02 15:04:05")})
}

func (h *CompatibilityHandler) GetCardDataV1(c *gin.Context) {
	courseID := c.Param("courseId")
	c.Params = append(filterParams(c.Params, "courseId"), gin.Param{Key: "courseId", Value: courseID})
	NewTeacherHandler(h.db, h.aiClient).GetCardData(c)
}

func (h *CompatibilityHandler) GetKnowledgeGraphV1(c *gin.Context) {
	courseID := c.Param("courseId")
	var kps []model.KnowledgePoint
	if err := h.db.Where("course_id = ?", courseID).Order("level asc, created_at asc").Find(&kps).Error; err != nil {
		apiresp.Internal(c, "查询知识图谱失败", err.Error())
		return
	}

	svc := service.NewKnowledgeGraphService(h.db)
	tnodes, err := svc.ListTeachingNodesForGraph(courseID)
	if err != nil {
		apiresp.Internal(c, "查询讲授节点失败", err.Error())
		return
	}
	relations, err := svc.ListTeachingNodeRelations(courseID)
	if err != nil {
		apiresp.Internal(c, "查询节点关联失败", err.Error())
		return
	}

	apiresp.OK(c, "请求成功", gin.H{
		"courseId":       courseID,
		"knowledgeGraph": buildKnowledgeTree(kps),
		"teachingNodes":  buildTeachingNodeGraphSummaries(tnodes),
		"relations":      buildRelationPayloads(relations),
		"schemaVersion":  1,
	})
}

// SyncCourseKnowledgeGraphV1 根据讲授节点 JSON 重算 teaching_node_relations（教师维护节点后调用）。
func (h *CompatibilityHandler) SyncCourseKnowledgeGraphV1(c *gin.Context) {
	courseID := strings.TrimSpace(c.Param("courseId"))
	if courseID == "" {
		apiresp.BadRequest(c, "缺少 courseId", "")
		return
	}
	svc := service.NewKnowledgeGraphService(h.db)
	if err := svc.RebuildTeachingNodeRelations(courseID); err != nil {
		if errors.Is(err, service.ErrTeachingNodePrerequisiteCycle) {
			apiresp.BadRequest(c, "知识节点前置依赖存在环路", err.Error())
			return
		}
		apiresp.Internal(c, "同步知识节点关联失败", err.Error())
		return
	}
	relations, _ := svc.ListTeachingNodeRelations(courseID)
	apiresp.OK(c, "同步成功", gin.H{
		"courseId":      courseID,
		"edgeCount":     len(relations),
		"relations":     buildRelationPayloads(relations),
		"schemaVersion": 1,
	})
}

// GetTeachingNodeReferenceHealthV1 扫描课件下各表对讲授节点业务 ID 的引用，列出 teaching_nodes 中不存在的 node_id（脏引用）。
func (h *CompatibilityHandler) GetTeachingNodeReferenceHealthV1(c *gin.Context) {
	courseID := strings.TrimSpace(c.Param("courseId"))
	if courseID == "" {
		apiresp.BadRequest(c, "缺少 courseId", "")
		return
	}
	svc := service.NewKnowledgeGraphService(h.db)
	report, err := svc.ScanOrphanTeachingNodeReferences(courseID)
	if err != nil {
		apiresp.Internal(c, "扫描节点引用失败", err.Error())
		return
	}
	apiresp.OKDataExtra(c, "请求成功", report, gin.H{"schemaVersion": 1})
}

// PostTeachingNodeReferenceRepairV1 清除孤儿 node 引用；请求体须含 "confirm": true。可选 "nodeIds" 仅修复子集。
func (h *CompatibilityHandler) PostTeachingNodeReferenceRepairV1(c *gin.Context) {
	courseID := strings.TrimSpace(c.Param("courseId"))
	if courseID == "" {
		apiresp.BadRequest(c, "缺少 courseId", "")
		return
	}
	var req struct {
		Confirm bool     `json:"confirm"`
		NodeIDs []string `json:"nodeIds"`
	}
	_ = c.ShouldBindJSON(&req)
	if !req.Confirm {
		apiresp.BadRequest(c, "请在请求体中设置 confirm 为 true 以执行修复", "")
		return
	}
	svc := service.NewKnowledgeGraphService(h.db)
	report, err := svc.RepairOrphanTeachingNodeReferences(courseID, req.NodeIDs)
	if err != nil {
		apiresp.Internal(c, "修复节点引用失败", err.Error())
		return
	}
	apiresp.OKDataExtra(c, "修复完成", report, gin.H{"schemaVersion": 1})
}

func buildTeachingNodeGraphSummaries(nodes []model.TeachingNode) []gin.H {
	out := make([]gin.H, 0, len(nodes))
	for _, n := range nodes {
		prereqs, pw, related, kw := parseTeachingNodeKnowledgeSurface(n.NodeID, n.KnowledgeNodesJSON)
		row := gin.H{
			"id":           n.ID,
			"nodeId":       n.NodeID,
			"nodeType":     strings.TrimSpace(n.NodeType),
			"title":        n.Title,
			"pageIndex":    n.PageIndex,
			"sortOrder":    n.SortOrder,
			"chapterTitle": n.ChapterTitle,
			"summary":      n.Summary,
		}
		if len(prereqs) > 0 {
			row["prerequisites"] = prereqs
		}
		if len(pw) > 0 {
			wh := gin.H{}
			for k, v := range pw {
				wh[k] = v
			}
			row["prerequisiteWeights"] = wh
		}
		if len(related) > 0 {
			row["relatedNodeIds"] = related
		}
		if kw > 0 {
			row["knowledgeWeight"] = kw
		}
		out = append(out, row)
	}
	return out
}

func buildRelationPayloads(relations []model.TeachingNodeRelation) []gin.H {
	out := make([]gin.H, 0, len(relations))
	for _, r := range relations {
		out = append(out, gin.H{
			"fromNodeId":   r.FromNodeID,
			"toNodeId":     r.ToNodeID,
			"relationType": r.RelationType,
			"weight":       r.Weight,
		})
	}
	return out
}

func buildKnowledgeTree(nodes []model.KnowledgePoint) []gin.H {
	childrenMap := map[string][]model.KnowledgePoint{}
	roots := make([]model.KnowledgePoint, 0)
	for _, node := range nodes {
		if strings.TrimSpace(node.ParentID) == "" {
			roots = append(roots, node)
			continue
		}
		childrenMap[node.ParentID] = append(childrenMap[node.ParentID], node)
	}
	for parentID := range childrenMap {
		sort.Slice(childrenMap[parentID], func(i, j int) bool {
			if childrenMap[parentID][i].Level == childrenMap[parentID][j].Level {
				return childrenMap[parentID][i].CreatedAt.Before(childrenMap[parentID][j].CreatedAt)
			}
			return childrenMap[parentID][i].Level < childrenMap[parentID][j].Level
		})
	}
	var convert func(model.KnowledgePoint) gin.H
	convert = func(node model.KnowledgePoint) gin.H {
		children := childrenMap[node.ID]
		result := gin.H{"id": node.ID, "name": node.Name, "level": node.Level, "content": node.Content}
		if strings.TrimSpace(node.SourceTeachingNodeID) != "" {
			result["sourceTeachingNodeId"] = strings.TrimSpace(node.SourceTeachingNodeID)
		}
		if len(children) > 0 {
			next := make([]gin.H, 0, len(children))
			for _, child := range children {
				next = append(next, convert(child))
			}
			result["children"] = next
		}
		return result
	}
	result := make([]gin.H, 0, len(roots))
	for _, root := range roots {
		result = append(result, convert(root))
	}
	return result
}

func (h *CompatibilityHandler) AskCoursewareV1(c *gin.Context) {
	courseID := c.Param("courseId")
	var req struct {
		PageNum    int    `json:"pageNum" binding:"required"`
		Type       string `json:"type"`
		Question   string `json:"question" binding:"required"`
		StudentID  string `json:"studentId"`
		TracePoint *struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"tracePoint"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.BadRequest(c, "参数错误", "")
		return
	}
	var coursePage model.CoursePage
	contextText := ""
	if err := h.db.Where("course_id = ? AND page_index = ?", courseID, req.PageNum).First(&coursePage).Error; err == nil {
		contextText = pageContextText(coursePage)
	}
	if strings.TrimSpace(contextText) == "" {
		contextText = buildPageContextFromTeachingNodes(loadTeachingNodesByPage(h.db, courseID, req.PageNum))
	}
	question := req.Question
	if req.TracePoint != nil {
		question = fmt.Sprintf("%s（定位坐标: %.2f, %.2f）", req.Question, req.TracePoint.X, req.TracePoint.Y)
	}
	resp, err := h.aiClient.AskWithContext(c.Request.Context(), service.AskWithContextRequest{Question: question, CurrentPage: req.PageNum, Context: contextText, Mode: "llm"})
	if err != nil {
		apiresp.ServiceUnavailable(c, "AI服务暂不可用", "")
		return
	}
	if strings.TrimSpace(req.StudentID) != "" {
		_ = h.db.Create(&model.QuestionLog{UserID: req.StudentID, CourseID: courseID, PageIndex: req.PageNum, NodeID: fmt.Sprintf("p%d_n1", req.PageNum), Question: req.Question, Answer: resp.Answer}).Error
	}
	apiresp.OK(c, "请求成功", gin.H{"answer": resp.Answer, "sourcePage": resp.SourcePage, "sourceExcerpt": resp.SourceExcerpt, "followUpSuggestion": resp.FollowUpSuggestion, "needReteach": resp.Intent.NeedReteach})
}

func (h *CompatibilityHandler) UploadCoursewareV1(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		title = "未命名课件"
	}
	file, err := c.FormFile("file")
	if err != nil {
		apiresp.BadRequest(c, "请选择文件", "")
		return
	}
	course, err := h.courseService.UploadCourse(c.Request.Context(), file, title)
	if err != nil {
		apiresp.Internal(c, "上传失败", "")
		return
	}
	apiresp.OK(c, "上传成功", course)
}
