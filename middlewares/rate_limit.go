package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// loginAttempt 记录登录尝试
type loginAttempt struct {
	count     int
	firstTime time.Time
}

// rateLimiter 登录频率限制器
type rateLimiter struct {
	mu       sync.Mutex
	attempts map[string]*loginAttempt
	maxTries int
	window   time.Duration
	banTime  time.Duration
}

var loginLimiter = &rateLimiter{
	attempts: make(map[string]*loginAttempt),
	maxTries: 5,              // 最多尝试次数
	window:   15 * time.Minute, // 时间窗口
	banTime:  30 * time.Minute, // 封禁时间
}

// 定期清理过期记录
func init() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			loginLimiter.mu.Lock()
			now := time.Now()
			for ip, att := range loginLimiter.attempts {
				if now.Sub(att.firstTime) > loginLimiter.banTime {
					delete(loginLimiter.attempts, ip)
				}
			}
			loginLimiter.mu.Unlock()
		}
	}()
}

// LoginRateLimit 登录频率限制中间件
func LoginRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		loginLimiter.mu.Lock()
		att, exists := loginLimiter.attempts[ip]
		now := time.Now()

		if !exists {
			loginLimiter.attempts[ip] = &loginAttempt{
				count:     1,
				firstTime: now,
			}
			loginLimiter.mu.Unlock()
			c.Next()
			return
		}

		// 检查是否在封禁期内
		if att.count >= loginLimiter.maxTries {
			if now.Sub(att.firstTime) < loginLimiter.banTime {
				loginLimiter.mu.Unlock()
				c.JSON(http.StatusTooManyRequests, gin.H{
					"status":  429,
					"message": "登录尝试过于频繁，请30分钟后再试",
				})
				c.Abort()
				return
			}
			// 封禁期已过，重置
			att.count = 0
			att.firstTime = now
		}

		// 检查是否在时间窗口内
		if now.Sub(att.firstTime) > loginLimiter.window {
			// 超过窗口，重置计数
			att.count = 1
			att.firstTime = now
		} else {
			att.count++
		}

		loginLimiter.mu.Unlock()
		c.Next()
	}
}
