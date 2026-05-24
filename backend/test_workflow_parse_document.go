//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"

	"smart-teaching-backend/internal/service"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run test_workflow_parse_document.go <文件路径.pdf|pptx>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	wfKey := os.Getenv("APP_AI_DIFY_WORKFLOW_API_KEY")
	if wfKey == "" {
		fmt.Println("FAIL: APP_AI_DIFY_WORKFLOW_API_KEY 未设置")
		os.Exit(1)
	}

	stat, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("FAIL: 读取文件失败: %v\n", err)
		os.Exit(1)
	}

	header := &multipart.FileHeader{
		Filename: stat.Name(),
		Size:     stat.Size(),
	}
	header.Header = make(map[string][]string)

	client := service.NewDifyClientWithFallback(
		getenv("APP_AI_DIFY_BASE_URL", "http://127.0.0.1:18001"),
		os.Getenv("APP_AI_DIFY_API_KEY"),
		wfKey,
		nil,
	)

	resp, err := client.ParseDocument(context.Background(), &localFileHeader{
		FileHeader: header,
		path:       filePath,
	})
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("OK: ParseDocument via Dify\npages=%d total=%d\n%s\n",
		len(resp.ParsedPages), resp.TotalPages, string(b))
}

type localFileHeader struct {
	*multipart.FileHeader
	path string
}

func (h *localFileHeader) Open() (multipart.File, error) {
	return os.Open(h.path)
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
