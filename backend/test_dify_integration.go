//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"smart-teaching-backend/internal/service"
)

func main() {
	baseURL := envOr("APP_AI_DIFY_BASE_URL", "http://127.0.0.1:18001")
	apiKey := os.Getenv("APP_AI_DIFY_API_KEY")
	if strings.TrimSpace(apiKey) == "" {
		fmt.Println("FAIL: APP_AI_DIFY_API_KEY 未设置")
		os.Exit(1)
	}

	fallback := service.NewAIEngineClient(envOr("APP_AI_BASE_URL", "http://127.0.0.1:8000"), 0)
	client := service.NewDifyClientWithFallback(baseURL, apiKey, os.Getenv("APP_AI_DIFY_WORKFLOW_API_KEY"), fallback)

	req := service.AskWithContextRequest{
		Question:    "什么是梯度下降？",
		CurrentPage: 3,
		Context:     "机器学习基础课程，包含优化算法章节",
		Mode:        "llm",
	}

	fmt.Printf("调用 DifyClient: baseURL=%s\n", baseURL)
	resp, err := client.AskWithContext(context.Background(), req)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("OK: DifyClient 联调成功")
	fmt.Printf("Question: %s\n", resp.Question)
	fmt.Printf("Answer: %s\n", truncate(resp.Answer, 200))
	fmt.Printf("SourcePage: %d\n", resp.SourcePage)
	fmt.Printf("NeedReteach: %v\n", resp.Intent.NeedReteach)
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func truncate(s string, n int) string {
	if len([]rune(s)) <= n {
		return s
	}
	runes := []rune(s)
	return string(runes[:n]) + "..."
}
