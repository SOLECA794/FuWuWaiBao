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

	resp, err := client.ParseKnowledge(context.Background(), service.ParseKnowledgeRequest{
		Text: "人工智能包括机器学习、深度学习、自然语言处理。机器学习分为监督学习和无监督学习。",
		Mode: "llm",
	})
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	b, _ := json.MarshalIndent(resp.Structure, "", "  ")
	fmt.Printf("OK: ParseKnowledge via Dify\nnodes=%d\n%s\n", len(resp.Structure), string(b))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
