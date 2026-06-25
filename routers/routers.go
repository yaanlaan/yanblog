package routers

import (
	"os"
	v1 "yanblog/api/v1"
	middleware "yanblog/middlewares"
	"yanblog/utils"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.ServerConfig.Server.AppMode)

	// 初始化路由
	r := gin.New()
	r.MaxMultipartMemory = 200 << 20 // 批量上传支持 200MB

	// 使用中间件
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression)) // 开启 gzip 压缩
	r.Use(middleware.Cors())

	// 确保必要的目录存在
	os.MkdirAll("./uploads", 0755)

	// 静态文件服务
	r.Static("/uploads", "./uploads")
	r.Static("/assets", "./web/frontend/public/assets")
	r.Static("/static", "./web/frontend/public/static")
	r.Static("/iconfont", "./web/frontend/public/iconfont")
	r.StaticFile("/favicon.ico", "./web/frontend/public/favicon.ico")
	
	r.StaticFile("/config.yaml", utils.GetFrontEndConfigPath())

	// 用户路由分组
	auth := r.Group("api/v1")
	auth.Use(middleware.JwtToken())

	// 管理员权限分组（仅超级管理员和管理员可操作）
	admin := r.Group("api/v1")
	admin.Use(middleware.JwtToken())
	admin.Use(middleware.AdminRequired())

	{
		// 用户模块的路由接口
		admin.POST("user/add", v1.AddUser) // 添加用户（需管理员权限）
		admin.PUT("user/:id", v1.EditUser)
		admin.DELETE("user/:id", v1.DeleteUser)
		auth.GET("users", v1.GetUsers)           // 查询用户列表（需认证）
		auth.GET("users/search", v1.SearchUsers) // 搜索用户（需认证）
		// 分类模块的路由接口
		admin.POST("category/add", v1.AddCategory)
		admin.PUT("category/:id", v1.EditCate)
		admin.DELETE("category/:id", v1.DeleteCate)
		// 文章模块的路由接口
		admin.POST("article/add", v1.AddArticle)
		admin.POST("article/zip", v1.UploadArticleZip)       // 上传单个ZIP发布文章
		admin.POST("article/zip/batch", v1.UploadArticleZipBatch) // 批量上传ZIP
		// 优化的ZIP上传（支持进度跟踪、断点续传、并发控制）
		admin.POST("article/zip/optimized", v1.UploadArticleZipOptimized)
		admin.POST("article/zip/batch/optimized", v1.UploadArticleZipBatchOptimized)
		admin.GET("article/upload/:id", v1.GetUploadProgress)  // 获取上传进度
		admin.DELETE("article/upload/:id", v1.CancelUpload)    // 取消上传任务
		// V2 增强版上传（WebSocket实时推送 + 断点续传 + 历史记录）
		admin.POST("article/zip/v2", v1.UploadArticleZipV2)           // V2上传
		admin.GET("article/upload/v2/:id", v1.GetUploadProgressV2)    // V2进度查询
		admin.DELETE("article/upload/v2/:id", v1.CancelUploadV2)      // V2取消
		admin.GET("article/upload/v2/:id/ws", v1.WebSocketProgress)   // WebSocket进度
		admin.POST("article/upload/v2/:id/retry", v1.RetryFailedUpload) // 重试失败文件
		admin.GET("article/upload/history", v1.GetUploadHistory)      // 上传历史
		admin.DELETE("article/upload/history", v1.EmptyRecycleBin) // 清空历史
		admin.PUT("article/:id", v1.EditArt)

		admin.DELETE("article/:id", v1.DeleteArt)
		admin.POST("article/batch-delete", v1.BatchDeleteArt)
		// 标签模块
		admin.POST("tags/add", v1.AddTag)
		admin.PUT("tags/:id", v1.EditTag)
		admin.DELETE("tags/:id", v1.DeleteTag)
		// 上传文件
		admin.POST("upload", v1.UpLoad)
		// 文件管理
		auth.GET("files", v1.GetFileList)
		admin.DELETE("files", v1.DeleteFile)
		admin.POST("files/folder", v1.CreateDir)     // 创建目录
		admin.PUT("files", v1.RenameFile)            // 重命名
		admin.POST("files/move", v1.MoveFile)        // 移动文件/目录
		admin.POST("files/copy", v1.CopyFile)        // 复制文件
		admin.POST("files/batch-delete", v1.BatchDeleteFiles) // 批量删除
		admin.POST("files/batch-upload", v1.BatchUploadFiles) // 批量上传
		auth.GET("files/stats", v1.GetStorageStats) // 获取存储统计
		// V2 增强文件管理
		auth.GET("files/v2/stats", v1.GetFileStats)              // 详细统计
		auth.GET("files/v2/search", v1.SearchFiles)              // 搜索文件
		admin.POST("files/v2/compress", v1.CompressFiles)        // 压缩
		admin.POST("files/v2/extract", v1.ExtractZip)            // 解压
		admin.POST("files/v2/recycle", v1.MoveToRecycleBin)      // 删除到回收站
		auth.GET("files/v2/recycle", v1.GetRecycleBin)           // 回收站列表
		admin.POST("files/v2/recycle/restore", v1.RestoreFromRecycleBin) // 恢复
		admin.DELETE("files/v2/recycle", v1.EmptyRecycleBin)     // 清空回收站
		auth.GET("files/v2/preview", v1.GetFilePreview)          // 文件预览
		admin.PUT("files/v2/metadata", v1.SaveFileMetadata)      // 保存元数据
		auth.GET("files/v2/metadata", v1.GetFileMetadata)        // 获取元数据
		// 前端配置管理（更新需要管理员权限）
		admin.PUT("frontend/config", v1.UpdateFrontEndConfig)
		// 后端配置管理
		admin.GET("backend/config", v1.GetBackendConfig)
		admin.PUT("backend/config", v1.UpdateBackendConfig)
		admin.POST("config/reload", v1.ReloadConfig)
		admin.GET("config/all", v1.GetAllConfig)
		admin.GET("system/status", v1.GetSystemStatus) // 获取系统状态信息（需管理员权限）
		// 关于页面内容管理
		admin.PUT("about", v1.UpdateAboutContent)
	}

	// 公共路由分组
	router := r.Group("api/v1")

	{
		router.GET("about", v1.GetAboutContent) // 获取关于页面内容 (公开接口，虽然前端直接读取静态文件，但提供 API 更统一)
		router.GET("category", v1.GetCate)
		router.GET("category/search", v1.SearchCate)    // 搜索分类
		router.GET("category/info/:id", v1.GetCateInfo) // 获取分类信息
		router.GET("article", v1.GetArt)
		router.GET("article/search", v1.SearchArt)          // 搜索文章
		router.GET("article/top", v1.GetTopArt)             // 获取置顶文章
		router.GET("article/hot", v1.GetHotArt)             // 获取热门文章
		router.GET("article/related/:id", v1.GetRelatedArt) // 获取相关文章
		router.GET("article/random", v1.GetRandomArt)       // 随机获取一篇文章
		router.GET("article/adjacent/:id", v1.GetAdjacentArt) // 获取相邻文章
		router.GET("article/archive", v1.GetArchive)        // 归档
		router.GET("article/list/:id", v1.GetCateArt)
		router.GET("article/info/:id", v1.GetArtInfo)
		router.GET("tags", v1.GetTags)                  // 获取标签列表
		router.GET("weather", v1.GetWeather)            // 获取天气信息
		router.GET("health", v1.HealthCheck)            // 健康检查（公开接口）
		router.GET("frontend/config", v1.GetFrontEndConfig) // 获取前端配置（公开接口）
		router.POST("login", middleware.LoginRateLimit(), v1.Login)
		router.GET("sitemap.xml", v1.GetSitemap) // 站点地图
	}

	_ = r.Run(utils.ServerConfig.Server.HttpPort) // 启动服务，监听端口
}
