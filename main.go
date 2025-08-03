package main

import (
	"yanblog/model"
	"yanblog/routers"
)

func main() {
	// 初始化数据库
	model.InitDB()
	
	// 初始化路由
	routers.InitRouter()
}