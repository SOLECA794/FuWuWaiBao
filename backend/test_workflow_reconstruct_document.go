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

	resp, err := client.ReconstructDocument(context.Background(), service.ReconstructDocumentRequest{
		ParsedDocument: map[string]any{
			"doc_id":       "doc_test_001",
			"doc_name":     "AI导论.pptx",
			"doc_type":     "pptx",
			"total_pages":  2,
			"parsed_pages": []map[string]any{
				{
					"page":           1,
					"content":        "第一章 机器学习概述\n- 机器学习是人工智能的分支\n- 监督学习与非监督学习",
					"content_length": 58,
				},
				{
					"page":           2,
					"content":        "第二章 监督学习\n- 分类与回归\n- 例如：垃圾邮件识别",
					"content_length": 42,
				},
			},
		},
		Mode: "hybrid",
	})
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("OK: ReconstructDocument via Dify\nnodes=%d chapters=%d\n%s\n",
		len(resp.TeachingNodes), len(resp.Chapters), string(b))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
