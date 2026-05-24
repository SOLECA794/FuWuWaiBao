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

	resp, err := client.GenerateNodeScript(context.Background(), service.GenerateNodeScriptRequest{
		TeachingNode: map[string]any{
			"node_id":           "p1_n1",
			"title":             "监督学习",
			"summary":           "用标注数据训练模型，完成分类与回归任务。",
			"core_points":       []string{"需要标注样本", "分类与回归是典型任务"},
			"examples":          []string{"垃圾邮件识别", "房价预测"},
			"common_confusions": []string{"监督学习不等于必须有深度学习"},
		},
		CourseName: "AI导论",
		Mode:       "llm",
	})
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("OK: GenerateNodeScript via Dify\nnode_id=%s title=%s script_len=%d\n%s\n",
		resp.NodeID, resp.Title, len(resp.Script), string(b))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
