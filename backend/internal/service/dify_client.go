package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// DifyClient 实现 AIEngine 接口，与 Dify API 交互
type DifyClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
	debug   bool
}

// NewDifyClient 创建新的 Dify 客户端
func NewDifyClient(baseURL, apiKey string) *DifyClient {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:18001"
	}
	// 移除末尾的斜杠
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &DifyClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		debug: false,
	}
}

// SetDebug 设置是否输出调试信息
func (c *DifyClient) SetDebug(debug bool) {
	c.debug = debug
}

// AskWithContext 调用 Dify 对话接口
func (c *DifyClient) AskWithContext(ctx context.Context, req AskWithContextRequest) (*AskWithContextResponse, error) {
	url := fmt.Sprintf("%s/api/v1/chat-messages", c.baseURL)

	// 构建请求体
	payload := map[string]interface{}{
		"query": req.Question,
		"user":  "smart-teaching-backend",
	}

	// 如果有上下文信息，添加到 inputs
	if strings.TrimSpace(req.Context) != "" {
		payload["inputs"] = map[string]interface{}{
			"context": req.Context,
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	if c.debug {
		fmt.Printf("[DifyClient] Request URL: %s\n", url)
		fmt.Printf("[DifyClient] Request Body: %s\n", string(body))
	}

	// 创建请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// 发送请求
	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Dify 服务调用失败: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	if c.debug {
		fmt.Printf("[DifyClient] Response Status: %d\n", httpResp.StatusCode)
		fmt.Printf("[DifyClient] Response Body: %s\n", string(respBody))
	}

	// 处理 Dify 流式响应
	// Dify 返回 SSE (Server-Sent Events) 格式，每行格式为 "data: {...}"
	lines := strings.Split(string(respBody), "\n")
	var answer string
	var lastEvent map[string]interface{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}

		// 移除 "data: " 前缀
		dataStr := strings.TrimPrefix(line, "data:")
		dataStr = strings.TrimSpace(dataStr)

		var event map[string]interface{}
		if err := json.Unmarshal([]byte(dataStr), &event); err != nil {
			continue
		}

		lastEvent = event

		// 从 "answer" 字段中提取答案
		if msg, ok := event["answer"].(string); ok {
			answer = msg
		}
	}

	// 如果是流式响应结束，从最后一条消息中提取最终答案
	if answer == "" && lastEvent != nil {
		if msg, ok := lastEvent["answer"].(string); ok {
			answer = msg
		}
	}

	result := &AskWithContextResponse{
		Question:           req.Question,
		SourcePage:         0,
		SourceExcerpt:      "",
		Answer:             answer,
		UsedFallback:       false,
		FallbackReason:     "",
		ResumePage:         0,
		FollowUpSuggestion: "",
	}
	result.Intent.NeedReteach = false
	result.Intent.Reason = ""

	return result, nil
}

// ParseDocument 解析文档
func (c *DifyClient) ParseDocument(ctx context.Context, file *multipart.FileHeader) (*ParseDocumentResponse, error) {
	return nil, nil
}

// ReconstructDocument 重建文档
func (c *DifyClient) ReconstructDocument(ctx context.Context, req ReconstructDocumentRequest) (*ReconstructDocumentResponse, error) {
	return nil, nil
}

// GenerateNodeScript 生成节点脚本
func (c *DifyClient) GenerateNodeScript(ctx context.Context, req GenerateNodeScriptRequest) (*GenerateNodeScriptResponse, error) {
	return nil, nil
}

// GenerateFromMarkdown 从 Markdown 生成
func (c *DifyClient) GenerateFromMarkdown(ctx context.Context, req GenerateFromMarkdownRequest) (*GenerateFromMarkdownResponse, error) {
	return nil, nil
}

// GenerateScript 生成脚本
func (c *DifyClient) GenerateScript(ctx context.Context, req GenerateScriptRequest) (*GenerateScriptResponse, error) {
	return nil, nil
}

// GenerateAudio 生成音频
func (c *DifyClient) GenerateAudio(ctx context.Context, req GenerateAudioRequest) (*GenerateAudioResponse, error) {
	return nil, nil
}

// ParseKnowledge 解析知识
func (c *DifyClient) ParseKnowledge(ctx context.Context, req ParseKnowledgeRequest) (*ParseKnowledgeResponse, error) {
	return nil, nil
}

