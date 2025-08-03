package routers

import (
	_ "net/http"
	v1 "yanblog/api/v1"
	middleware "yanblog/middlewares"
	"yanblog/utils"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.ServerConfig.Server.AppMode)

	// 初始化路由
	r := gin.New()

	// 使用中间件
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	// 用户路由分组
	auth := r.Group("api/v1")
	auth.Use(middleware.JwtToken())

	{
		// 用户模块的路由接口
		auth.PUT("user/:id", v1.EditUser)
		auth.DELETE("user/:id", v1.DeleteUser)
		// 分类模块的路由接口
		auth.POST("category/add", v1.AddCategory)
		auth.PUT("category/:id", v1.EditCate)
		auth.DELETE("category/:id", v1.DeleteCate)
		// 文章模块的路由接口
		auth.POST("article/add", v1.AddArticle)
		auth.PUT("article/:id", v1.EditArt)
		auth.DELETE("article/:id", v1.DeleteArt)
		// 上传文件
		auth.POST("upload", v1.UpLoad)
	}

	// 公共路由分组
	router := r.Group("api/v1")
	
	{
		router.POST("user/add", v1.AddUser)
		router.GET("users", v1.GetUsers)
		router.GET("users/search", v1.SearchUsers) // 搜索用户
		router.GET("category", v1.GetCate)
		router.GET("category/search", v1.SearchCate) // 搜索分类
		router.GET("category/info/:id", v1.GetCateInfo) // 获取分类信息
		router.GET("article", v1.GetArt)
		router.GET("article/search", v1.SearchArt) // 搜索文章
		router.GET("article/top", v1.GetTopArt) // 获取置顶文章
		router.GET("article/list/:id", v1.GetCateArt)
		router.GET("article/info/:id", v1.GetArtInfo)
		router.GET("weather", v1.GetWeather) // 获取天气信息
		router.POST("login", v1.Login)
	}

	_ = r.Run(utils.ServerConfig.Server.HttpPort) // 启动服务，监听端口

}