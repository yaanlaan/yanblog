// cors.go - CORS跨域资源共享中间件
package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"    
)

// Cors 跨域资源共享中间件
// 允许前端应用从不同域访问后端API
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建CORS配置并应用
		cors.New(cors.Config{
			//AllowAllOrigins: true,
			AllowOrigins:  []string{"*"},
			AllowMethods:  []string{"*"}, 
			AllowHeaders:  []string{"Origin"},
			ExposeHeaders: []string{"Content-Length", "Authorization"},
			//AllowCredentials: true, 
			//AllowOriginFunc: func(origin string) bool {
			//	return origin == "https://github.com"
			//},            
			MaxAge: 12 * time.Hour, 
		})
	}
}