//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"smart-teaching-backend/internal/service"
)

func main() {
	wfKey := os.Getenv("APP_AI_DIFY_WORKFLOW_API_KEY")
	if wfKey == "" {
		fmt.Println("FAIL: APP_AI_DIFY_WORKFLOW_API_KEY 未设置")
		os.Exit(1)
	}

	client := service.NewDifyClientWithFallback(
		getenv("APP_AI_DIFY_BASE_URL", "http://127.0.0.1:18001"),
		os.Getenv("APP_AI_DIFY_API_KEY"),
		wfKey,
		nil,
	)

	md := `# 机器学习基础

## 监督学习
用标注数据训练模型，如分类和回归。

## 无监督学习
从无标注数据中发现模式，如聚类。`

	resp, err := client.GenerateFromMarkdown(context.Background(), service.GenerateFromMarkdownRequest{
		Markdown:   md,
		CourseName: "AI导论",
		Mode:       "llm",
	})
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("OK: GenerateFromMarkdown via Dify\n")
	fmt.Printf("course=%s key_points=%d nodes=%d scripts=%d\n",
		resp.CourseName, len(resp.KeyPoints), len(resp.NodeTree.Nodes), len(resp.Scripts))

	b, _ := json.MarshalIndent(map[string]any{
		"key_points": resp.KeyPoints,
		"nodes":      resp.NodeTree.Nodes,
		"scripts":    len(resp.Scripts),
	}, "", "  ")
	fmt.Println(string(b))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
