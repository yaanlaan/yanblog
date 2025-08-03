package model

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"yanblog/utils"
)

// Weather 天气信息结构体
type Weather struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
}

// OpenWeatherMapResponse OpenWeatherMap API响应结构体
type OpenWeatherMapResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

// 城市名称映射，将中文城市名映射为英文名
var cityMapping = map[string]string{
	"北京": "Beijing",
	"上海": "Shanghai",
	"广州": "Guangzhou",
	"深圳": "Shenzhen",
	"杭州": "Hangzhou",
	"南京": "Nanjing",
	"成都": "Chengdu",
	"武汉": "Wuhan",
	"西安": "Xi'an",
	"重庆": "Chongqing",
	"天津": "Tianjin",
	"苏州": "Suzhou",
	"长沙": "Changsha",
	"合肥": "Hefei",
}

// GetWeather 获取天气信息
func GetWeather(city string) (*Weather, error) {
	// 从配置文件获取API密钥
	apiKey := utils.ServerConfig.Weather.ApiKey
	if apiKey == "" || apiKey == "YOUR_REAL_OPENWEATHER_API_KEY_HERE" {
		// 如果没有API密钥，返回错误
		return nil, fmt.Errorf("天气API密钥未配置，请在config.yaml中配置有效的API密钥")
	}

	// 如果没有指定城市，则使用配置文件中的默认城市
	if city == "" {
		city = utils.ServerConfig.Weather.DefaultCity
	}

	// 处理城市名称，如果在映射中则使用英文名
	cityName := strings.TrimSpace(city)
	if mappedName, exists := cityMapping[cityName]; exists {
		cityName = mappedName
	}

	// 构建API请求URL
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=zh_cn", cityName, apiKey)

	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置10秒超时
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建天气API请求失败: %v", err)
	}

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求天气API失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取天气API响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		// 特殊处理API密钥错误
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("天气API密钥无效，请检查config.yaml中的ApiKey配置")
		}
		return nil, fmt.Errorf("天气API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析JSON响应
	var apiResponse OpenWeatherMapResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("解析天气API响应失败: %v", err)
	}

	// 检查API返回的状态码
	if apiResponse.Cod != 200 {
		return nil, fmt.Errorf("天气API返回错误状态，cod: %d, 响应: %s", apiResponse.Cod, string(body))
	}

	// 构建返回的天气信息
	weather := &Weather{
		City:        apiResponse.Name,
		Temperature: apiResponse.Main.Temp,
		Description: "晴", // 默认值
		Humidity:    apiResponse.Main.Humidity,
		WindSpeed:   apiResponse.Wind.Speed,
	}

	// 如果有天气描述，使用API返回的描述
	if len(apiResponse.Weather) > 0 {
		weather.Description = apiResponse.Weather[0].Description
	}

	return weather, nil
}