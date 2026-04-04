package main

import (
	"context"
	"fmt"
	"smart-teaching-backend/internal/service"
)

func main() {
	// 测试 DifyClient 创建和基本调用
	client := service.NewDifyClient("http://127.0.0.1:18001", "")
	
	// 测试请求
	req := service.AskWithContextRequest{
		Question:    "这是什么？",
		CurrentPage: 1,
		Context:     "这是一份关于数学的教学材料",
		Mode:        "llm",
	}
	
	fmt.Println("测试 DifyClient...")
	resp, err := client.AskWithContext(context.Background(), req)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	
	fmt.Printf("响应: %+v\n", resp)
}
