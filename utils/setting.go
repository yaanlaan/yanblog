package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		AppMode  string `yaml:"AppMode" json:"appMode"`
		HttpPort string `yaml:"HttpPort" json:"httpPort"`
		SiteUrl  string `yaml:"SiteUrl" json:"siteUrl"`
	} `yaml:"server" json:"server"`

	Database struct {
		Db         string `yaml:"Db" json:"db"`
		DbHost     string `yaml:"DbHost" json:"dbHost"`
		DbPort     int    `yaml:"DbPort" json:"dbPort"`
		DbUser     string `yaml:"DbUser" json:"dbUser"`
		DbPassWord string `yaml:"DbPassWord" json:"dbPassWord"`
		DbName     string `yaml:"DbName" json:"dbName"`
	} `yaml:"database" json:"database"`

	JwtKey string `yaml:"JwtKey" json:"jwtKey"`

	Weather struct {
		DefaultCity string `yaml:"DefaultCity" json:"defaultCity"`
	} `yaml:"weather" json:"weather"`

	FrontEndConfigPath string `yaml:"FrontEndConfigPath" json:"frontEndConfigPath"`

	Cities []struct {
		Name  string `yaml:"Name" json:"name"`
		Alias string `yaml:"Alias" json:"alias"`
	} `yaml:"cities" json:"cities"`
}

var ServerConfig = Config{}
var configMutex sync.RWMutex

func init() {
	configPath := getConfigPath("config/backend/config.yaml")
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("未找到 config/backend/config.yaml，尝试旧路径...")
		configPath = "config/config.yaml"
		file, err = os.ReadFile(configPath)
		if err != nil {
			log.Printf("未找到 config/config.yaml，使用模板 config/config_template.yaml")
			configPath = "config/config_template.yaml"
			file, err = os.ReadFile(configPath)
			if err != nil {
				log.Fatalf("读取配置文件失败，请从 config/config_template.yaml 创建 config/backend/config.yaml。错误信息：%s", err)
			}
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
	if ServerConfig.Database.DbPassWord != "" && ServerConfig.Database.DbPassWord != "rootpassword" {
		fmt.Println("  密码状态: 已配置（非默认值）")
	} else {
		fmt.Println("⚠️  数据库密码仍为默认值或未配置，建议修改！")
	}
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
	// 确保目录存在（Docker 容器中可能没有 config/backend/ 子目录）
	_ = os.MkdirAll(filepath.Dir(configPath), 0755)
	return os.WriteFile(configPath, data, 0644)
}

func ReloadConfig() error {
	configPath := getConfigPath("config/backend/config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "config/config.yaml"
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configPath = "config/config_template.yaml"
		}
	}
	file, err := os.ReadFile(configPath)
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