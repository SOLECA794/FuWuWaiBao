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
	baseURL := getenv("APP_AI_DIFY_BASE_URL", "http://127.0.0.1:18001")
	chatKey := os.Getenv("APP_AI_DIFY_API_KEY")
	wfKey := os.Getenv("APP_AI_DIFY_WORKFLOW_API_KEY")
	if wfKey == "" {
		fmt.Println("FAIL: APP_AI_DIFY_WORKFLOW_API_KEY 未设置")
		os.Exit(1)
	}

	client := service.NewDifyClientWithFallback(baseURL, chatKey, wfKey, nil)
	resp, err := client.GenerateScript(context.Background(), service.GenerateScriptRequest{
		Page:       1,
		Content:    "人工智能是研究如何让机器具备智能的学科。",
		CourseName: "AI导论",
		Mode:       "llm",
	})
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("OK: GenerateScript via Dify Workflow")
	fmt.Printf("page=%d script_len=%d mindmap_len=%d\n", resp.Page, len(resp.Script), len(resp.MindmapMarkdown))
	if len(resp.Script) > 0 {
		s := resp.Script
		if len([]rune(s)) > 80 {
			s = string([]rune(s)[:80]) + "..."
		}
		fmt.Printf("script: %s\n", s)
	}
}

func getenv(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}
