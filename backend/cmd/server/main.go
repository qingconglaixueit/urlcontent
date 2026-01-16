package main

import (
	"fmt"
	"log"
	"net/http"

	"urltocontent/backend/internal/config"
	"urltocontent/backend/internal/handlers"
)

func main() {
	fmt.Println("\n====================")
	fmt.Println("🤖 URL to Content API")
	fmt.Println("====================\n")

	// 加载配置
	cfg := config.Load()
	fmt.Printf("📋 服务器端口: %s\n", cfg.Port)
	fmt.Printf("🚀 飞书 App ID: %s\n", cfg.FeishuAppID)
	fmt.Printf("📚 飞书 Wiki ID: %s\n", cfg.FeishuWikiID)
	fmt.Println()

	// 创建处理器
	handler := handlers.NewHandler(cfg)

	// 设置路由
	http.HandleFunc("/health", handlers.CORSMiddleware(handler.HealthCheckHandler))
	http.HandleFunc("/api/parse", handlers.CORSMiddleware(handler.ParseURLHandler))
	http.HandleFunc("/api/write", handlers.CORSMiddleware(handler.WriteToFeishuHandler))

	// 启动服务器
	addr := ":" + cfg.Port
	fmt.Printf("✅ 服务器启动成功，监听地址: http://localhost%s\n\n", addr)
	fmt.Println("可用端点:")
	fmt.Println("  - GET  /health    - 健康检查")
	fmt.Println("  - POST /api/parse - URL 解析")
	fmt.Println("  - POST /api/write - 写入飞书")
	fmt.Println("\n按 Ctrl+C 停止服务器\n")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}
