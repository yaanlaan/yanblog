package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// LoginAttempt 登录尝试记录（数据库模型）
type LoginAttempt struct {
	ID        uint      `gorm:"primarykey"`
	IP        string    `gorm:"type:varchar(45);uniqueIndex;not null"` // IPv6 最大 45 字符
	Count     int       `gorm:"not null;default:0"`
	FirstTime time.Time `gorm:"not null;index"`
	UpdatedAt time.Time `gorm:"not null"`
}

// TableName 指定表名
func (LoginAttempt) TableName() string {
	return "login_attempts"
}

// rateLimiter 登录频率限制器（基于 SQLite 持久化）
type rateLimiter struct {
	db       *gorm.DB
	mu       sync.Mutex
	maxTries int
	window   time.Duration
	banTime  time.Duration
}

var loginLimiter *rateLimiter

// InitRateLimiter 初始化限流器（必须在 main.go 中调用）
func InitRateLimiter(dbPath string) error {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(&LoginAttempt{}); err != nil {
		return err
	}

	loginLimiter = &rateLimiter{
		db:       db,
		maxTries: 5,               // 最多尝试次数
		window:   15 * time.Minute, // 时间窗口
		banTime:  30 * time.Minute, // 封禁时间
	}

	// 定期清理过期记录
	go cleanupExpiredRecords()

	return nil
}

// cleanupExpiredRecords 定期清理过期记录
func cleanupExpiredRecords() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		loginLimiter.mu.Lock()
		cutoffTime := time.Now().Add(-loginLimiter.banTime)
		loginLimiter.db.Where("first_time < ?", cutoffTime).Delete(&LoginAttempt{})
		loginLimiter.mu.Unlock()
	}
}

// LoginRateLimit 登录频率限制中间件
func LoginRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		loginLimiter.mu.Lock()
		defer loginLimiter.mu.Unlock()

		now := time.Now()
		var attempt LoginAttempt

		// 查询该 IP 的登录尝试记录
		result := loginLimiter.db.Where("ip = ?", ip).First(&attempt)

		if result.Error != nil {
			// 首次尝试，创建记录
			newAttempt := LoginAttempt{
				IP:        ip,
				Count:     1,
				FirstTime: now,
				UpdatedAt: now,
			}
			loginLimiter.db.Create(&newAttempt)
			c.Next()
			return
		}

		// 检查是否在封禁期内
		if attempt.Count >= loginLimiter.maxTries {
			if now.Sub(attempt.FirstTime) < loginLimiter.banTime {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"status":  429,
					"message": "登录尝试过于频繁，请30分钟后再试",
				})
				c.Abort()
				return
			}
			// 封禁期已过，重置
			loginLimiter.db.Model(&attempt).Updates(map[string]interface{}{
				"count":      0,
				"first_time": now,
				"updated_at": now,
			})
			attempt.Count = 0
			attempt.FirstTime = now
		}

		// 检查是否在时间窗口内
		if now.Sub(attempt.FirstTime) > loginLimiter.window {
			// 超过窗口，重置计数
			loginLimiter.db.Model(&attempt).Updates(map[string]interface{}{
				"count":      1,
				"first_time": now,
				"updated_at": now,
			})
		} else {
			// 窗口内，增加计数
			loginLimiter.db.Model(&attempt).Update("count", attempt.Count+1)
		}

		c.Next()
	}
}
