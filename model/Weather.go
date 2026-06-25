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

// OpenMeteo 天气 API 响应结构体
type OpenMeteoCurrentWeather struct {
	Time                string  `json:"time"`
	Temperature2m       float64 `json:"temperature_2m"`
	RelativeHumidity2m  int     `json:"relative_humidity_2m"`
	WeatherCode         int     `json:"weather_code"`
	WindSpeed10m        float64 `json:"wind_speed_10m"`
}

type OpenMeteoWeatherResponse struct {
	Current OpenMeteoCurrentWeather `json:"current"`
}

// OpenMeteo 地理编码 API 响应结构体
type OpenMeteoGeoResult struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	Admin1    string  `json:"admin1"`
}

type OpenMeteoGeoResponse struct {
	Results []OpenMeteoGeoResult `json:"results"`
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

// GetWeather 获取天气信息（使用 Open-Meteo 免费 API，无需密钥）
func GetWeather(city string) (*Weather, error) {
	// 如果没有指定城市，则使用配置文件中的默认城市
	if city == "" {
		city = utils.ServerConfig.Weather.DefaultCity
	}

	// 处理城市名称，如果在映射中则使用英文名
	cityName := strings.TrimSpace(city)
	if mappedName, exists := cityMapping[cityName]; exists {
		cityName = mappedName
	}

	// 创建带超时的 HTTP 客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 第一步：通过地理编码 API 获取城市经纬度
	lat, lon, displayName, err := geocodeCity(ctx, client, cityName)
	if err != nil {
		return nil, fmt.Errorf("地理编码失败: %v", err)
	}

	// 第二步：通过 Open-Meteo 天气 API 获取当前天气
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m",
		lat, lon,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建天气请求失败: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求天气 API 失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取天气响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("天气 API 请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var weatherResp OpenMeteoWeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return nil, fmt.Errorf("解析天气响应失败: %v", err)
	}

	current := weatherResp.Current
	weather := &Weather{
		City:        displayName,
		Temperature: current.Temperature2m,
		Description: wmoCodeToDescription(current.WeatherCode),
		Humidity:    current.RelativeHumidity2m,
		WindSpeed:   current.WindSpeed10m,
	}

	return weather, nil
}

// geocodeCity 通过 Open-Meteo 地理编码 API 获取城市坐标
func geocodeCity(ctx context.Context, client *http.Client, cityName string) (float64, float64, string, error) {
	url := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=zh", cityName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, 0, "", fmt.Errorf("创建地理编码请求失败: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, "", fmt.Errorf("请求地理编码 API 失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, "", fmt.Errorf("读取地理编码响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, 0, "", fmt.Errorf("地理编码 API 请求失败，状态码: %d", resp.StatusCode)
	}

	var geoResp OpenMeteoGeoResponse
	if err := json.Unmarshal(body, &geoResp); err != nil {
		return 0, 0, "", fmt.Errorf("解析地理编码响应失败: %v", err)
	}

	if len(geoResp.Results) == 0 {
		return 0, 0, "", fmt.Errorf("未找到城市: %s", cityName)
	}

	result := geoResp.Results[0]
	displayName := result.Name
	if result.Admin1 != "" && result.Admin1 != result.Name {
		displayName = result.Name
	}

	return result.Latitude, result.Longitude, displayName, nil
}

// wmoCodeToDescription 将 WMO 天气代码转换为中文描述
func wmoCodeToDescription(code int) string {
	switch {
	case code == 0:
		return "晴"
	case code == 1:
		return "大部晴朗"
	case code == 2:
		return "多云"
	case code == 3:
		return "阴"
	case code == 45 || code == 48:
		return "雾"
	case code >= 51 && code <= 57:
		return "毛毛雨"
	case code >= 61 && code <= 65:
		return "雨"
	case code == 66 || code == 67:
		return "冻雨"
	case code >= 71 && code <= 77:
		return "雪"
	case code >= 80 && code <= 82:
		return "阵雨"
	case code == 85 || code == 86:
		return "阵雪"
	case code >= 95:
		return "雷暴"
	default:
		return "未知"
	}
}

// getSimulatedWeather 无网络时返回模拟天气数据（保留作为降级方案）
func getSimulatedWeather(city string) *Weather {
	month := time.Now().Month()
	var temp, humidity int
	var desc string

	switch {
	case month >= 6 && month <= 8:
		temp = 28 + hashCity(city)%10
		humidity = 60 + hashCity(city)%25
		desc = "多云"
	case month >= 3 && month <= 5:
		temp = 18 + hashCity(city)%8
		humidity = 50 + hashCity(city)%20
		desc = "晴朗"
	case month >= 9 && month <= 11:
		temp = 15 + hashCity(city)%8
		humidity = 45 + hashCity(city)%20
		desc = "微风"
	default:
		temp = 3 + hashCity(city)%8
		humidity = 35 + hashCity(city)%15
		desc = "寒冷"
	}

	return &Weather{
		City:        city,
		Temperature: float64(temp),
		Description: desc,
		Humidity:    humidity,
		WindSpeed:   float64(2 + hashCity(city)%5),
	}
}

// hashCity 基于城市名生成一个伪随机种子
func hashCity(s string) int {
	h := 0
	for _, c := range s {
		h = h*31 + int(c)
	}
	if h < 0 {
		h = -h
	}
	return h
}

