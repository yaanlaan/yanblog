package v1

import (
	"context"
	"net/http"
	"time"
	"yanblog/model"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// GetWeather 获取天气信息
func GetWeather(c *gin.Context) {
	// 从查询参数获取城市名，如果没有指定则使用配置中的默认城市
	city := c.Query("city")
	
	// 创建带超时的上下文，避免长时间等待
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// 创建channel用于接收结果
	resultChan := make(chan struct {
		weather *model.Weather
		err     error
	}, 1)
	
	// 在goroutine中执行天气API调用
	go func() {
		weather, err := model.GetWeather(city)
		resultChan <- struct {
			weather *model.Weather
			err     error
		}{weather, err}
	}()
	
	// 等待结果或超时
	select {
	case result := <-resultChan:
		if result.err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR,
				"data":    nil,
				"message": "获取天气信息失败: " + result.err.Error(),
			})
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.SUCCESS,
			"data":    result.weather,
			"message": errmsg.GetErrMsg(errmsg.SUCCESS),
		})
	case <-ctx.Done():
		// 超时处理
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"data":    nil,
			"message": "获取天气信息超时",
		})
	}
}