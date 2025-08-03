package v1

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// 系统状态信息结构体
type SystemStatus struct {
	Status      string  `json:"status"`       // 状态 (online/offline)
	Uptime      string  `json:"uptime"`       // 运行时间
	MemoryUsage float64 `json:"memory_usage"` // 内存使用率
	CPUUsage    float64 `json:"cpu_usage"`    // CPU使用率
	Goroutines  int     `json:"goroutines"`   // Goroutines数量
	StartTime   int64   `json:"start_time"`   // 启动时间戳（毫秒）
}

// 获取系统启动时间
var startTime = time.Now()

// GetSystemStatus 获取系统状态信息
func GetSystemStatus(c *gin.Context) {
	// 获取内存信息
	_, err := mem.VirtualMemory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "获取内存信息失败",
			"data":    nil,
		})
		return
	}

	// 获取CPU信息
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "获取CPU信息失败",
			"data":    nil,
		})
		return
	}

	// 获取当前进程信息
	pid := os.Getpid()
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "获取进程信息失败",
			"data":    nil,
		})
		return
	}

	// 获取进程内存使用
	procMemPercent, err := proc.MemoryPercent()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "获取进程内存使用率失败",
			"data":    nil,
		})
		return
	}

	// 计算运行时间
	uptime := time.Since(startTime)
	uptimeStr := formatUptime(uptime)

	// 构造返回数据
	systemStatus := SystemStatus{
		Status:      "online",
		Uptime:      uptimeStr,
		MemoryUsage: float64(procMemPercent),
		CPUUsage:    cpuPercent[0],
		Goroutines:  runtime.NumGoroutine(),
		StartTime:   startTime.UnixNano() / 1000000, // 转换为毫秒时间戳
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "OK",
		"data":    systemStatus,
	})
}

// formatUptime 格式化运行时间
func formatUptime(duration time.Duration) string {
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟%d秒", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟%d秒", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%d分钟%d秒", minutes, seconds)
	} else {
		return fmt.Sprintf("%d秒", seconds)
	}
}