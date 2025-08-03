package main

import (
	"fmt"
	"yanblog/model"
	"yanblog/utils"
)

func main() {
	// 打印配置信息
	fmt.Printf("Weather Provider: %s\n", utils.ServerConfig.Weather.Provider)
	fmt.Printf("API Key: %s\n", utils.ServerConfig.Weather.ApiKey)
	
	// 测试获取北京天气信息（使用中文城市名）
	fmt.Println("\n测试获取北京天气信息:")
	weather, err := model.GetWeather("北京")
	if err != nil {
		fmt.Printf("获取天气信息失败: %v\n", err)
	} else {
		fmt.Printf("城市: %s\n", weather.City)
		fmt.Printf("温度: %.1f°C\n", weather.Temperature)
		fmt.Printf("描述: %s\n", weather.Description)
		fmt.Printf("湿度: %d%%\n", weather.Humidity)
		fmt.Printf("风速: %.1f m/s\n", weather.WindSpeed)
	}
	
	// 测试获取上海天气信息（使用中文城市名）
	fmt.Println("\n测试获取上海天气信息:")
	weather, err = model.GetWeather("上海")
	if err != nil {
		fmt.Printf("获取天气信息失败: %v\n", err)
	} else {
		fmt.Printf("城市: %s\n", weather.City)
		fmt.Printf("温度: %.1f°C\n", weather.Temperature)
		fmt.Printf("描述: %s\n", weather.Description)
		fmt.Printf("湿度: %d%%\n", weather.Humidity)
		fmt.Printf("风速: %.1f m/s\n", weather.WindSpeed)
	}
	
	// 测试获取无效城市天气信息
	fmt.Println("\n测试获取无效城市天气信息:")
	weather, err = model.GetWeather("不存在的城市")
	if err != nil {
		fmt.Printf("获取天气信息失败: %v\n", err)
	} else {
		fmt.Printf("城市: %s\n", weather.City)
		fmt.Printf("温度: %.1f°C\n", weather.Temperature)
		fmt.Printf("描述: %s\n", weather.Description)
		fmt.Printf("湿度: %d%%\n", weather.Humidity)
		fmt.Printf("风速: %.1f m/s\n", weather.WindSpeed)
	}
}