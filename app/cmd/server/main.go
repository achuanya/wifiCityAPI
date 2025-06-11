package main

import (
	"app/config"
	"app/internal/router"
	"app/pkg/database"
	"fmt"
	"log"
)

func main() {
	// 初始化数据库连接
	// 配置加载在 config 包的 init() 函数中自动完成
	// 所以我们在这里直接使用 database.Init()
	database.Init()

	// 设置并获取 Gin 路由引擎
	r := router.SetupRouter()

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", config.Cfg.Server.Port)
	log.Printf("服务器正在启动，监听地址: %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
