package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		AppMode  string `yaml:"AppMode"`
		HttpPort string `yaml:"HttpPort"`
		SiteUrl  string `yaml:"SiteUrl"`
	} `yaml:"server"`

	Database struct {
		Db         string `yaml:"Db"`
		DbHost     string `yaml:"DbHost"`
		DbPort     int    `yaml:"DbPort"`
		DbUser     string `yaml:"DbUser"`
		DbPassWord string `yaml:"DbPassWord"`
		DbName     string `yaml:"DbName"`
	} `yaml:"database"`

	JwtKey string `yaml:"JwtKey"`

	Weather struct {
		Provider    string `yaml:"Provider"`
		ApiKey      string `yaml:"ApiKey"`
		DefaultCity string `yaml:"DefaultCity"`
	} `yaml:"weather"`

	FrontEndConfigPath string `yaml:"FrontEndConfigPath"`

	Cities []struct {
		Name  string `yaml:"Name"`
		Alias string `yaml:"Alias"`
	} `yaml:"cities"`
}

var ServerConfig = Config{}
var configMutex sync.RWMutex

func init() {
	configPath := getConfigPath("config/backend/config.yaml")
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("尝试读取默认配置文件失败: %s", err)
		configPath = "config/config.yaml"
		file, err = ioutil.ReadFile(configPath)
		if err != nil {
			log.Fatalf("读取配置文件失败，错误信息：%s", err)
		}
	}
	LoadConfig(file)
}

func GetConfigPath(defaultPath string) string {
	return getConfigPath(defaultPath)
}

func getConfigPath(defaultPath string) string {
	if envPath := os.Getenv("YANBLOG_CONFIG_PATH"); envPath != "" {
		return envPath
	}
	return defaultPath
}

func LoadConfig(file []byte) {
	configMutex.Lock()
	defer configMutex.Unlock()

	content := string(file)
	content = replaceEnvVars(content)

	err := yaml.Unmarshal([]byte(content), &ServerConfig)
	if err != nil {
		log.Fatalf("解析数据库配置失败，错误信息：%s", err)
	}
	fmt.Printf("数据库配置加载成功: %s@%s:%d/%s\n",
		ServerConfig.Database.DbUser,
		ServerConfig.Database.DbHost,
		ServerConfig.Database.DbPort,
		ServerConfig.Database.DbName)
}

func replaceEnvVars(content string) string {
	re := regexp.MustCompile(`\$\{(\w+)(?::([^}]*))?\}`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		matches := re.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		envName := matches[1]
		defaultValue := ""
		if len(matches) > 2 {
			defaultValue = matches[2]
		}
		if value := os.Getenv(envName); value != "" {
			return value
		}
		return defaultValue
	})
}

func SaveConfig() error {
	configMutex.RLock()
	data, err := yaml.Marshal(&ServerConfig)
	configMutex.RUnlock()
	if err != nil {
		return err
	}
	configPath := getConfigPath("config/backend/config.yaml")
	return ioutil.WriteFile(configPath, data, 0644)
}

func ReloadConfig() error {
	configPath := getConfigPath("config/backend/config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "config/config.yaml"
	}
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	LoadConfig(file)
	// 通知 JWT 中间件刷新密钥（通过 model 包间接调用）
	// 注意：middlewares.RefreshJwtKey() 需要在路由初始化后调用
	return nil
}

// OnConfigReloaded 配置重新加载后的回调
var OnConfigReloaded func()

func GetConfig() Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return ServerConfig
}

func GetFrontEndConfigPath() string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	path := ServerConfig.FrontEndConfigPath
	if path == "" {
		if _, err := os.Stat("config/frontend/config.yaml"); err == nil {
			path = "config/frontend/config.yaml"
		} else {
			path = "web/frontend/public/config.yaml"
		}
	}
	return path
}