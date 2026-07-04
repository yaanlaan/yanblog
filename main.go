package main

import (
	"fmt"
	"os"
	"yanblog/middlewares"
	"yanblog/model"
	"yanblog/routers"
	"yanblog/utils"
)

func main() {
	// 验证配置文件
	if err := utils.ValidateConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ 配置错误：%v\n\n", err)
		fmt.Println("请检查 config/config.yaml 配置文件，确保所有必填项已正确设置。")
		os.Exit(1)
	}

	// 初始化 JWT 密钥
	middlewares.InitJwtKey(utils.ServerConfig.JwtKey)
	// 注册配置重载回调
	utils.OnConfigReloaded = middlewares.RefreshJwtKey

	// 打印启动信息
	utils.PrintStartupInfo()

	// 初始化数据库
	model.InitDB()

	// 初始化登录限流器（基于 SQLite 持久化）
	if err := middlewares.InitRateLimiter("./data/rate_limit.db"); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  限流器初始化失败: %v\n", err)
		fmt.Println("将继续运行，但登录限流功能不可用")
	}

	// 启动 API 限流清理
	go middlewares.CleanupAPIRateLimits()

	// 初始化路由
	routers.InitRouter()
}