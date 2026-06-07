// cors.go - CORS跨域资源共享中间件
package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 跨域资源共享中间件
// 允许前端应用从不同域访问后端API
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: false, // 使用 AllowAllOrigins 时必须为 false（CORS 规范要求）
		MaxAge:           12 * time.Hour,
	})
}
