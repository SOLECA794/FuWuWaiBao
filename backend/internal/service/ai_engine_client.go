package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AIEngine interface {
	GenerateScript(ctx context.Context, req GenerateScriptRequest) (*GenerateScriptResponse, error)
	AskWithContext(ctx context.Context, req AskWithContextRequest) (*AskWithContextResponse, error)
	ParseKnowledge(ctx context.Context, req ParseKnowledgeRequest) (*ParseKnowledgeResponse, error)
}

type aiEngineClient struct {
	baseURL string
	client  *http.Client
}

type GenerateScriptRequest struct {
	Page       int    `json:"page"`
	Content    string `json:"content"`
	CourseName string `json:"course_name"`
	Mode       string `json:"mode"`
}

type GenerateScriptResponse struct {
	Page            int    `json:"page"`
	Script          string `json:"script"`
	MindmapMarkdown string `json:"mindmap_markdown"`
}

type AskWithContextRequest struct {
	Question    string `json:"question"`
	CurrentPage int    `json:"current_page"`
	Context     string `json:"context"`
	Mode        string `json:"mode"`
}

type AskWithContextResponse struct {
	Question           string `json:"question"`
	SourcePage         int    `json:"source_page"`
	SourceExcerpt      string `json:"source_excerpt"`
	Answer             string `json:"answer"`
	ResumePage         int    `json:"resume_page"`
	FollowUpSuggestion string `json:"follow_up_suggestion"`
	Intent             struct {
		NeedReteach bool   `json:"need_reteach"`
		Reason      string `json:"reason"`
	} `json:"intent"`
}

type ParseKnowledgeRequest struct {
	Text string `json:"text"`
	Mode string `json:"mode"`
}

type KnowledgeNode struct {
	Name     string          `json:"name"`
	Children []KnowledgeNode `json:"children,omitempty"`
}

type ParseKnowledgeResponse struct {
	Structure []KnowledgeNode `json:"structure"`
}

func NewAIEngineClient(baseURL string, timeout time.Duration) AIEngine {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "http://127.0.0.1:8000"
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &aiEngineClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: timeout},
	}
}

func (c *aiEngineClient) GenerateScript(ctx context.Context, req GenerateScriptRequest) (*GenerateScriptResponse, error) {
	var result GenerateScriptResponse
	if err := c.postJSON(ctx, "/generate-script", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *aiEngineClient) AskWithContext(ctx context.Context, req AskWithContextRequest) (*AskWithContextResponse, error) {
	var result AskWithContextResponse
	if err := c.postJSON(ctx, "/ask-with-context", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *aiEngineClient) ParseKnowledge(ctx context.Context, req ParseKnowledgeRequest) (*ParseKnowledgeResponse, error) {
	var result ParseKnowledgeResponse
	if err := c.postJSON(ctx, "/parse-knowledge", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *aiEngineClient) postJSON(ctx context.Context, path string, reqBody any, out any) error {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("构造请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("调用AI引擎失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("AI引擎返回异常状态: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("解析AI引擎响应失败: %w", err)
	}
	return nil
}
