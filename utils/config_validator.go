package utils

import (
	"fmt"
	"os"
	"strings"
)

// ValidateConfig 验证配置文件的完整性
// 只返回真正的错误（阻止启动），警告信息直接打印
func ValidateConfig() error {
	var errors []string
	var warnings []string

	// 验证数据库配置
	if ServerConfig.Database.DbUser == "" {
		errors = append(errors, "数据库用户名 (database.DbUser) 不能为空")
	}
	if ServerConfig.Database.DbPassWord == "rootpassword" {
		warnings = append(warnings, "⚠️  数据库密码仍为默认值 (rootpassword)，建议修改为强密码")
	}
	if ServerConfig.Database.DbName == "" {
		errors = append(errors, "数据库名称 (database.DbName) 不能为空")
	}

	// 验证 JWT 密钥
	jwtKey := ServerConfig.JwtKey
	if jwtKey == "" {
		errors = append(errors, "JWT 密钥 (JwtKey) 未设置，请运行 'openssl rand -hex 32' 生成并填入配置")
	} else if len(jwtKey) < 32 {
		warnings = append(warnings, fmt.Sprintf("⚠️  JWT 密钥长度不足（当前 %d 位），建议使用 64 位随机密钥", len(jwtKey)))
	}

	// 验证服务器配置
	if ServerConfig.Server.HttpPort == "" {
		errors = append(errors, "服务器端口 (server.HttpPort) 不能为空")
	}

	// 打印警告信息（不阻止启动）
	for _, w := range warnings {
		fmt.Println(w)
	}

	// 如果有错误，返回错误信息
	if len(errors) > 0 {
		return fmt.Errorf("配置验证失败：\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// PrintStartupInfo 打印启动信息
func PrintStartupInfo() {
	fmt.Println("\n===========================================")
	fmt.Println("🚀 博客系统启动中...")
	fmt.Println("===========================================")

	// 显示配置信息
	dbType := ServerConfig.Database.Db
	if dbType == "" {
		dbType = "MYSQL"
	}
	fmt.Printf("📝 运行模式：%s\n", ServerConfig.Server.AppMode)
	fmt.Printf("🌐 服务端口：%s\n", ServerConfig.Server.HttpPort)
	fmt.Printf("💾 数据库类型：%s\n", dbType)
	fmt.Printf("💾 数据库地址：%s@%s:%d/%s\n",
		ServerConfig.Database.DbUser,
		ServerConfig.Database.DbHost,
		ServerConfig.Database.DbPort,
		ServerConfig.Database.DbName)

	// 天气配置
	if ServerConfig.Weather.ApiKey != "" {
		fmt.Printf("🌤️  天气服务：已启用 (%s)\n", ServerConfig.Weather.DefaultCity)
	} else {
		fmt.Println("🌤️  天气服务：未启用")
	}

	fmt.Println("===========================================")
}

// CheckFileExists 检查文件是否存在
func CheckFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
