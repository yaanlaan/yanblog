package middlewares

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type apiRateLimit struct {
	tokens     int
	maxTokens  int
	rate       time.Duration
	lastAccess time.Time
	mu         sync.Mutex
}

var apiRateLimits sync.Map

const (
	defaultAPIRateLimit    = 100
	defaultAPIRateInterval = time.Minute
)

var (
	shutdownCtx    context.Context
	shutdownCancel context.CancelFunc
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

	if err := db.AutoMigrate(&LoginAttempt{}); err != nil {
		return err
	}

	loginLimiter = &rateLimiter{
		db:       db,
		maxTries: 5,              // 最多尝试次数
		window:   15 * time.Minute, // 时间窗口
		banTime:  5 * time.Minute,  // 封禁时间（5分钟）
	}

	shutdownCtx, shutdownCancel = context.WithCancel(context.Background())

	go cleanupExpiredRecords()

	return nil
}

func cleanupExpiredRecords() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			loginLimiter.mu.Lock()
			cutoffTime := time.Now().Add(-loginLimiter.banTime)
			loginLimiter.db.Where("first_time < ?", cutoffTime).Delete(&LoginAttempt{})
			loginLimiter.mu.Unlock()
		case <-shutdownCtx.Done():
			return
		}
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
					"message": "登录尝试过于频繁，请5分钟后再试",
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

// APIRateLimit 通用 API 限流中间件（基于令牌桶算法）
// 默认每分钟 100 次请求限制
func APIRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		item, _ := apiRateLimits.LoadOrStore(ip, &apiRateLimit{
			tokens:     defaultAPIRateLimit,
			maxTokens:  defaultAPIRateLimit,
			rate:       defaultAPIRateInterval,
			lastAccess: time.Now(),
		})

		limit := item.(*apiRateLimit)
		limit.mu.Lock()
		defer limit.mu.Unlock()

		now := time.Now()
		elapsed := now.Sub(limit.lastAccess)

		if elapsed >= limit.rate {
			limit.tokens = limit.maxTokens
			limit.lastAccess = now
		}

		if limit.tokens > 0 {
			limit.tokens--
			limit.lastAccess = now
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
		}
	}
}

func CleanupAPIRateLimits() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			apiRateLimits.Range(func(key, value interface{}) bool {
				limit := value.(*apiRateLimit)
				limit.mu.Lock()
				if now.Sub(limit.lastAccess) > 10*time.Minute {
					limit.mu.Unlock()
					apiRateLimits.Delete(key)
				} else {
					limit.mu.Unlock()
				}
				return true
			})
		case <-shutdownCtx.Done():
			return
		}
	}
}

func Shutdown() {
	if shutdownCancel != nil {
		shutdownCancel()
	}
}
