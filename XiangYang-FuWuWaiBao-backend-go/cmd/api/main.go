package main

import (
	"log"

	"smart-learning/internal/router"
)

func main() {
	// 1. 初始化路由
	r := router.SetupRouter()

	// 2. 启动 HTTP 服务，默认监听 8080 端口
	log.Println("Server is starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
