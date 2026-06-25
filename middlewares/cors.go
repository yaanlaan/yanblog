// cors.go - CORS跨域资源共享中间件
package middlewares

import (
	"strings"
	"time"

	"yanblog/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 跨域资源共享中间件
// 允许前端应用从不同域访问后端API
func Cors() gin.HandlerFunc {
	// 从配置读取允许的域名，未配置时根据运行模式决定
	allowOrigins := []string{}
	if utils.ServerConfig.Server.SiteUrl != "" {
		// 生产环境：只允许配置的站点URL
		siteUrl := strings.TrimSuffix(utils.ServerConfig.Server.SiteUrl, "/")
		allowOrigins = []string{siteUrl}
	}

	// 如果没有配置站点URL或处于开发模式，允许所有来源
	isDev := utils.ServerConfig.Server.AppMode == "debug"
	allowAll := len(allowOrigins) == 0 || isDev

	return cors.New(cors.Config{
		// 根据环境决定是否允许所有来源
		AllowAllOrigins: allowAll,
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
}
