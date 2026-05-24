package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	difySceneAskWithContext       = "ask_with_context"
	difySceneGenerateScript       = "generate_script"
	difySceneParseDocument        = "parse_document"
	difySceneReconstructDocument  = "reconstruct_document"
	difySceneGenerateNodeScript   = "generate_node_script"
	difySceneGenerateFromMarkdown = "generate_from_markdown"
	difySceneGenerateAudio        = "generate_audio"
	difySceneParseKnowledge       = "parse_knowledge"
)

// DifyClient 实现 AIEngine 接口，与 Dify API 交互。
type DifyClient struct {
	baseURL        string
	apiKey         string
	workflowAPIKey string
	client         *http.Client
	fallback       AIEngine
	debug          bool
}

// NewDifyClient 创建新的 Dify 客户端。
func NewDifyClient(baseURL, apiKey string) *DifyClient {
	return NewDifyClientWithFallback(baseURL, apiKey, apiKey, nil)
}

// NewDifyClientWithFallback 创建带 fallback 的 Dify 客户端。
func NewDifyClientWithFallback(baseURL, apiKey, workflowAPIKey string, fallback AIEngine) *DifyClient {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:18001"
	}
	baseURL = strings.TrimRight(baseURL, "/")
	if workflowAPIKey == "" {
		workflowAPIKey = apiKey
	}
	return &DifyClient{
		baseURL:        baseURL,
		apiKey:         apiKey,
		workflowAPIKey: workflowAPIKey,
		client:         &http.Client{Timeout: 60 * time.Second},
		fallback:       fallback,
	}
}

// SetDebug 设置是否输出调试信息。
func (c *DifyClient) SetDebug(debug bool) {
	c.debug = debug
}

func (c *DifyClient) AskWithContext(ctx context.Context, req AskWithContextRequest) (*AskWithContextResponse, error) {
	inputs := map[string]any{
		"scene":           difySceneAskWithContext,
		"current_page":    req.CurrentPage,
		"context":         req.Context,
		"mode":            req.Mode,
		"question":        req.Question,
		"session_id":      req.SessionID,
		"history_summary": req.HistorySummary,
	}
	if len(req.RecentTurns) > 0 {
		inputs["recent_turns"] = req.RecentTurns
	}

	body, err := c.postJSON(ctx, "/v1/chat-messages", map[string]any{
		"query":         req.Question,
		"inputs":        inputs,
		"response_mode": "blocking",
		"user":          "smart-teaching-backend",
	})
	if err != nil {
		if c.fallback != nil {
			return c.fallback.AskWithContext(ctx, req)
		}
		return nil, err
	}

	answer := parseChatAnswer(body)
	resp := &AskWithContextResponse{
		Question:           req.Question,
		Answer:             answer,
		SourcePage:         req.CurrentPage,
		ResumePage:         req.CurrentPage,
		SourceExcerpt:      req.Context,
		FollowUpSuggestion: "",
	}
	resp.Intent.NeedReteach = detectReteachIntent(req.Question)
	if strings.TrimSpace(answer) != "" {
		if parsed, ok := parseStructuredAskAnswer(answer); ok {
			if strings.TrimSpace(parsed.Question) == "" {
				parsed.Question = req.Question
			}
			if parsed.SourcePage == 0 {
				parsed.SourcePage = req.CurrentPage
			}
			if parsed.ResumePage == 0 {
				parsed.ResumePage = req.CurrentPage
			}
			if strings.TrimSpace(parsed.SourceExcerpt) == "" {
				parsed.SourceExcerpt = req.Context
			}
			if strings.TrimSpace(parsed.Answer) == "" {
				parsed.Answer = answer
			}
			if !parsed.Intent.NeedReteach {
				parsed.Intent.NeedReteach = detectReteachIntent(req.Question)
			}
			return parsed, nil
		}
	}
	return resp, nil
}

func (c *DifyClient) ParseDocument(ctx context.Context, file *multipart.FileHeader) (*ParseDocumentResponse, error) {
	if file == nil {
		return nil, fmt.Errorf("文件不能为空")
	}
	if err := c.ensureFallbackAllowed("ParseDocument"); err != nil {
		return nil, err
	}

	uploadedID, err := c.uploadFile(ctx, file)
	if err != nil {
		if c.fallback != nil {
			return c.fallback.ParseDocument(ctx, file)
		}
		return nil, err
	}

	workflowResp, err := c.runWorkflow(ctx, map[string]any{
		"scene": difySceneParseDocument,
		"file": map[string]any{
			"transfer_method": "local_file",
			"upload_file_id":  uploadedID,
			"type":            "document",
		},
		"filename": file.Filename,
	})
	if err != nil {
		if c.fallback != nil {
			return c.fallback.ParseDocument(ctx, file)
		}
		return nil, err
	}

	var result ParseDocumentResponse
	if err := decodeWorkflowOutput(workflowResp, &result); err != nil {
		if c.fallback != nil {
			return c.fallback.ParseDocument(ctx, file)
		}
		return nil, err
	}
	return &result, nil
}

func (c *DifyClient) ReconstructDocument(ctx context.Context, req ReconstructDocumentRequest) (*ReconstructDocumentResponse, error) {
	workflowResp, err := c.runWorkflow(ctx, map[string]any{
		"scene":           difySceneReconstructDocument,
		"mode":            req.Mode,
		"parsed_document": mustJSONString(req.ParsedDocument),
	})
	if err != nil {
		if c.fallback != nil {
			return c.fallback.ReconstructDocument(ctx, req)
		}
		return nil, err
	}

	var result ReconstructDocumentResponse
	if err := decodeWorkflowOutput(workflowResp, &result); err != nil {
		if c.fallback != nil {
			return c.fallback.ReconstructDocument(ctx, req)
		}
		return nil, err
	}
	return &result, nil
}

func (c *DifyClient) GenerateNodeScript(ctx context.Context, req GenerateNodeScriptRequest) (*GenerateNodeScriptResponse, error) {
	workflowResp, err := c.runWorkflow(ctx, map[string]any{
		"scene":         difySceneGenerateNodeScript,
		"mode":          req.Mode,
		"course_name":   req.CourseName,
		"teaching_node": mustJSONString(req.TeachingNode),
	})
	if err != nil {
		if c.fallback != nil {
			return c.fallback.GenerateNodeScript(ctx, req)
		}
		return nil, err
	}

	var result GenerateNodeScriptResponse
	if err := decodeWorkflowOutput(workflowResp, &result); err != nil {
		if c.fallback != nil {
			return c.fallback.GenerateNodeScript(ctx, req)
		}
		return nil, err
	}
	return &result, nil
}

func (c *DifyClient) GenerateFromMarkdown(ctx context.Context, req GenerateFromMarkdownRequest) (*GenerateFromMarkdownResponse, error) {
	workflowResp, err := c.runWorkflow(ctx, map[string]any{
		"scene":       difySceneGenerateFromMarkdown,
		"mode":        req.Mode,
		"course_name": req.CourseName,
		"markdown":    req.Markdown,
	})
	if err != nil {
		if c.fallback != nil {
			return c.fallback.GenerateFromMarkdown(ctx, req)
		}
		return nil, err
	}

	var result GenerateFromMarkdownResponse
	if err := decodeWorkflowOutput(workflowResp, &result); err != nil {
		if c.fallback != nil {
			return c.fallback.GenerateFromMarkdown(ctx, req)
		}
		return nil, err
	}
	return &result, nil
}

func (c *DifyClient) GenerateScript(ctx context.Context, req GenerateScriptRequest) (*GenerateScriptResponse, error) {
	workflowResp, err := c.runWorkflow(ctx, map[string]any{
		"scene":       difySceneGenerateScript,
		"page":        req.Page,
		"content":     req.Content,
		"course_name": req.CourseName,
		"mode":        req.Mode,
	})
	if err != nil {
		if c.fallback != nil {
			return c.fallback.GenerateScript(ctx, req)
		}
		return nil, err
	}

	var result GenerateScriptResponse
	if err := decodeWorkflowOutput(workflowResp, &result); err != nil {
		if c.fallback != nil {
			return c.fallback.GenerateScript(ctx, req)
		}
		return nil, err
	}
	return &result, nil
}

func (c *DifyClient) GenerateAudio(ctx context.Context, req GenerateAudioRequest) (*GenerateAudioResponse, error) {
	workflowResp, err := c.runWorkflow(ctx, map[string]any{
		"scene":       difySceneGenerateAudio,
		"course_id":   req.CourseID,
		"page":        req.Page,
		"voice_type":  req.VoiceType,
		"format":      req.Format,
		"provider":    req.Provider,
		"playback_id": req.PlaybackID,
		"nodes":       mustJSONString(req.Nodes),
	})
	if err != nil {
		if c.fallback != nil {
			return c.fallback.GenerateAudio(ctx, req)
		}
		return nil, err
	}

	var result GenerateAudioResponse
	if err := decodeWorkflowOutput(workflowResp, &result); err != nil {
		if c.fallback != nil {
			return c.fallback.GenerateAudio(ctx, req)
		}
		return nil, err
	}
	return &result, nil
}

func (c *DifyClient) ParseKnowledge(ctx context.Context, req ParseKnowledgeRequest) (*ParseKnowledgeResponse, error) {
	workflowResp, err := c.runWorkflow(ctx, map[string]any{
		"scene": difySceneParseKnowledge,
		"text":  req.Text,
		"mode":  req.Mode,
	})
	if err != nil {
		if c.fallback != nil {
			return c.fallback.ParseKnowledge(ctx, req)
		}
		return nil, err
	}

	var result ParseKnowledgeResponse
	if err := decodeWorkflowOutput(workflowResp, &result); err != nil {
		if c.fallback != nil {
			return c.fallback.ParseKnowledge(ctx, req)
		}
		return nil, err
	}
	return &result, nil
}

func (c *DifyClient) uploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return "", fmt.Errorf("构造 multipart 请求失败: %w", err)
	}
	if _, err := io.Copy(part, src); err != nil {
		return "", fmt.Errorf("写入 multipart 文件失败: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("关闭 multipart writer 失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/files/upload", &body)
	if err != nil {
		return "", fmt.Errorf("构造文件上传请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("调用 Dify 文件上传接口失败: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取文件上传响应失败: %w", err)
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("Dify 文件上传接口返回异常状态: %d", resp.StatusCode)
	}

	var payload map[string]any
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return "", fmt.Errorf("解析文件上传响应失败: %w", err)
	}
	for _, key := range []string{"id", "file_id"} {
		if value, ok := payload[key].(string); ok && strings.TrimSpace(value) != "" {
			return value, nil
		}
	}
	return "", fmt.Errorf("Dify 文件上传响应缺少文件 ID")
}

func (c *DifyClient) runWorkflow(ctx context.Context, inputs map[string]any) (map[string]any, error) {
	requestBody := map[string]any{
		"inputs":        inputs,
		"response_mode": "blocking",
		"user":          "smart-teaching-backend",
	}
	body, err := c.postJSON(ctx, "/v1/workflows/run", requestBody)
	if err != nil {
		return nil, err
	}

	var envelope map[string]any
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("解析 workflow 响应失败: %w", err)
	}

	data, _ := asMap(envelope["data"])
	if data == nil {
		if status, ok := envelope["status"].(string); ok && strings.EqualFold(status, "succeeded") {
			return envelope, nil
		}
		return nil, fmt.Errorf("workflow 响应缺少 data 字段")
	}

	status, _ := data["status"].(string)
	if status != "" && !strings.EqualFold(status, "succeeded") {
		reason := "workflow 未成功执行"
		if msg, ok := data["error"].(string); ok && strings.TrimSpace(msg) != "" {
			reason = fmt.Sprintf("workflow 执行失败: %s", msg)
		}
		return nil, errors.New(reason)
	}

	outputs, _ := asMap(data["outputs"])
	if outputs == nil {
		outputs, _ = asMap(envelope["outputs"])
	}
	if outputs == nil {
		if len(data) > 0 {
			return data, nil
		}
		return envelope, nil
	}
	return outputs, nil
}

func (c *DifyClient) postJSON(ctx context.Context, path string, reqBody any) ([]byte, error) {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("构造请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用 Dify 接口失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 Dify 响应失败: %w", err)
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Dify 接口返回异常状态: %d", resp.StatusCode)
	}
	return respBody, nil
}

func parseChatAnswer(body []byte) string {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return ""
	}

	var direct map[string]any
	if err := json.Unmarshal(trimmed, &direct); err == nil {
		if answer, ok := direct["answer"].(string); ok {
			return answer
		}
		if data, ok := asMap(direct["data"]); ok {
			if answer, ok := data["answer"].(string); ok {
				return answer
			}
		}
	}

	return parseSSEAnswer(trimmed)
}

func parseStructuredAskAnswer(answer string) (*AskWithContextResponse, bool) {
	trimmed := strings.TrimSpace(answer)
	if trimmed == "" || !(strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) {
		return nil, false
	}
	var resp AskWithContextResponse
	if err := json.Unmarshal([]byte(trimmed), &resp); err != nil {
		return nil, false
	}
	return &resp, true
}

func parseSSEAnswer(body []byte) string {
	lines := bytes.Split(body, []byte{'\n'})
	answer := ""
	for _, rawLine := range lines {
		line := strings.TrimSpace(string(rawLine))
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}
		dataStr := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if dataStr == "[DONE]" {
			continue
		}
		var event map[string]any
		if err := json.Unmarshal([]byte(dataStr), &event); err != nil {
			continue
		}
		if msg, ok := event["answer"].(string); ok {
			answer = msg
			continue
		}
		if msg, ok := event["content"].(string); ok {
			answer = msg
		}
	}
	return answer
}

func decodeWorkflowOutput(raw any, out any) error {
	if out == nil {
		return errors.New("输出目标不能为空")
	}
	payload, err := normalizeWorkflowPayload(raw)
	if err != nil {
		return err
	}
	if payload == nil {
		return errors.New("workflow 输出为空")
	}

	data, err := payloadToBytes(payload)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("解析 workflow 输出失败: %w", err)
	}
	return nil
}

func normalizeWorkflowPayload(raw any) (any, error) {
	if raw == nil {
		return nil, errors.New("workflow 输出为空")
	}
	switch value := raw.(type) {
	case []byte:
		var decoded any
		if err := json.Unmarshal(value, &decoded); err != nil {
			return nil, err
		}
		return normalizeWorkflowPayload(decoded)
	case string:
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return nil, errors.New("workflow 输出为空")
		}
		if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
			var decoded any
			if err := json.Unmarshal([]byte(trimmed), &decoded); err == nil {
				return normalizeWorkflowPayload(decoded)
			}
		}
		return value, nil
	case map[string]any:
		if result, ok := value["result"]; ok {
			return normalizeWorkflowPayload(result)
		}
		if outputs, ok := asMap(value["outputs"]); ok {
			if result, ok := outputs["result"]; ok {
				return normalizeWorkflowPayload(result)
			}
			return outputs, nil
		}
		if data, ok := asMap(value["data"]); ok {
			return normalizeWorkflowPayload(data)
		}
		if len(value) == 0 {
			return nil, errors.New("workflow 输出为空")
		}
		return value, nil
	default:
		return value, nil
	}
}

func payloadToBytes(payload any) ([]byte, error) {
	switch value := payload.(type) {
	case []byte:
		return value, nil
	case string:
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return nil, errors.New("workflow 输出为空")
		}
		if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
			return []byte(trimmed), nil
		}
		return json.Marshal(value)
	default:
		return json.Marshal(value)
	}
}

func detectReteachIntent(question string) bool {
	text := strings.TrimSpace(question)
	if text == "" {
		return false
	}
	patterns := []string{
		"听不懂",
		"没听懂",
		"不懂",
		"看不懂",
		"再讲一遍",
		"再讲一次",
		"重讲",
		"重新讲",
		"解释一下",
		"再解释",
		"换个角度",
	}
	for _, pattern := range patterns {
		if strings.Contains(text, pattern) {
			return true
		}
	}
	return false
}

func asMap(value any) (map[string]any, bool) {
	if value == nil {
		return nil, false
	}
	if result, ok := value.(map[string]any); ok {
		return result, true
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Map {
		return nil, false
	}
	result := make(map[string]any, rv.Len())
	for _, key := range rv.MapKeys() {
		if key.Kind() != reflect.String {
			continue
		}
		result[key.String()] = rv.MapIndex(key).Interface()
	}
	return result, true
}

func mustJSONString(value any) string {
	if value == nil {
		return ""
	}
	data, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(data)
}

func (c *DifyClient) ensureFallbackAllowed(method string) error {
	if c == nil {
		return errors.New("DifyClient 为空")
	}
	return nil
}
