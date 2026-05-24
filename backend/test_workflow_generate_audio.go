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

	resp, err := client.GenerateAudio(context.Background(), service.GenerateAudioRequest{
		CourseID:   "course_demo_001",
		Page:       1,
		VoiceType:  getenv("TTS_VOICE", "Cherry"),
		Format:     "mp3",
		Provider:   getenv("TTS_PROVIDER", "dashscope"),
		PlaybackID: "audio_course_demo_001_1",
		Nodes: []service.GenerateAudioNode{
			{
				NodeID: "p1_n1",
				Title:  "监督学习",
				Text:   "同学们好，今天我们学习监督学习。它使用标注数据训练模型，典型任务包括分类和回归。",
			},
			{
				NodeID: "p1_n2",
				Title:  "无监督学习",
				Text:   "无监督学习不需要标注，常见应用是聚类和降维。",
			},
		},
	})
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("OK: GenerateAudio via Dify\n")
	fmt.Printf("audio_id=%s status=%s sections=%d total=%ds\n",
		resp.AudioID, resp.Status, len(resp.Sections), resp.TotalDuration)
	for _, s := range resp.Sections {
		fmt.Printf("  - %s audio_url=%s duration=%ds\n", s.NodeID, s.AudioURL, s.DurationSec)
	}

	b, _ := json.MarshalIndent(map[string]any{
		"audio_id":       resp.AudioID,
		"status":         resp.Status,
		"total_duration": resp.TotalDuration,
		"sections":       resp.Sections,
	}, "", "  ")
	fmt.Println(string(b))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
