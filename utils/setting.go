package utils

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// Config 定义配置文件的结构
type Config struct {
	Server struct { //服务器
		AppMode  string `yaml:"AppMode"`
		HttpPort string `yaml:"HttpPort"`
		SiteUrl  string `yaml:"SiteUrl"`
	} `yaml:"server"`

	Database struct { //数据库
		Db         string `yaml:"Db"`
		DbHost     string `yaml:"DbHost"`
		DbPort     int    `yaml:"DbPort"`
		DbUser     string `yaml:"DbUser"`
		DbPassWord string `yaml:"DbPassWord"`
		DbName     string `yaml:"DbName"`
	} `yaml:"database"`

	//JwtKey 用于JWT加密
	JwtKey string `yaml:"JwtKey"`

	// Weather 天气API配置
	Weather struct {
		Provider    string `yaml:"Provider"`
		ApiKey      string `yaml:"ApiKey"`
		DefaultCity string `yaml:"DefaultCity"`
	} `yaml:"weather"`

	FrontEndConfigPath string `yaml:"FrontEndConfigPath"`

	// Cities 城市配置
	Cities []struct {
		Name  string `yaml:"Name"`
		Alias string `yaml:"Alias"`
	} `yaml:"cities"`
}

var ServerConfig = Config{}

// init 初始化配置文件
// 参数: 无
// 返回: 无
func init() {
	file, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败，错误信息：%s", err)
	}
	LoadConfig(file)
	// fmt.Printf("完整配置加载成功: %+v\n", ServerConfig)
}

// LoadConfig 解析配置
// 参数: file - 配置文件内容
// 返回: 无
func LoadConfig(file []byte) {
	// 解析配置文件
	err := yaml.Unmarshal(file, &ServerConfig)
	if err != nil {
		log.Fatalf("解析数据库配置失败，错误信息：%s", err)
	}
	fmt.Printf("数据库配置加载成功: %+v\n", ServerConfig.Database)
}
