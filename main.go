package main

import (
	"fmt"
	"os"
	"yanblog/model"
	"yanblog/routers"
	"yanblog/utils"
)

func main() {
	// 验证配置文件
	if err := utils.ValidateConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ 配置错误：%v\n\n", err)
		fmt.Println("请检查 config/config.yaml 配置文件，确保所有必填项已正确设置。")
		fmt.Println("详细说明请参考 QUICK_START.md 文档")
		os.Exit(1)
	}
	
	// 打印启动信息
	utils.PrintStartupInfo()
	
	// 初始化数据库
	model.InitDB()
	
	// 初始化路由
	routers.InitRouter()
}