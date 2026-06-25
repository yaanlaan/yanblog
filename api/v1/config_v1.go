package v1

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"yanblog/middlewares"
	"yanblog/utils"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func GetFrontEndConfig(c *gin.Context) {
	configPath := utils.GetFrontEndConfigPath()

	content, err := os.ReadFile(configPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "读取配置文件失败: " + err.Error(),
		})
		return
	}

	// 禁止缓存，确保修改后立即生效
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"data":    string(content),
		"message": errmsg.GetErrMsg(errmsg.SUCCESS),
	})
}

func UpdateFrontEndConfig(c *gin.Context) {
	var input struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	var temp interface{}
	if err := yaml.Unmarshal([]byte(input.Content), &temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "YAML 格式错误: " + err.Error(),
		})
		return
	}

	configPath := utils.GetFrontEndConfigPath()

	err := os.WriteFile(configPath, []byte(input.Content), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "保存配置文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "配置保存成功",
	})
}

func GetBackendConfig(c *gin.Context) {
	config := utils.GetConfig()

	// 过滤敏感字段，防止泄露到前端（使用结构体方式）
	safeConfig := map[string]interface{}{
		"server": map[string]interface{}{
			"AppMode":  config.Server.AppMode,
			"HttpPort": config.Server.HttpPort,
			"SiteUrl":  config.Server.SiteUrl,
		},
		"database": map[string]interface{}{
			"Db":     config.Database.Db,
			"DbHost": config.Database.DbHost,
			"DbPort": config.Database.DbPort,
			"DbUser": config.Database.DbUser,
			"DbName": config.Database.DbName,
			// 故意不返回 DbPassWord
		},
		"weather": config.Weather,
		// 故意不返回 JwtKey
		"FrontEndConfigPath": config.FrontEndConfigPath,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"data":    safeConfig,
		"message": errmsg.GetErrMsg(errmsg.SUCCESS),
	})
}

func UpdateBackendConfig(c *gin.Context) {
	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	// 输入验证：仅允许修改已知的配置字段
	allowedKeys := map[string]bool{
		"server": true, "database": true, "weather": true,
		"FrontEndConfigPath": true,
	}
	for key := range input {
		if !allowedKeys[key] {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  errmsg.ERROR,
				"message": "不允许修改的配置字段: " + key,
			})
			return
		}
	}

	// 验证 server 字段
	if serverInput, ok := input["server"].(map[string]interface{}); ok {
		allowedServerKeys := map[string]bool{"AppMode": true, "HttpPort": true, "SiteUrl": true}
		for key := range serverInput {
			if !allowedServerKeys[key] {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  errmsg.ERROR,
					"message": "不允许的服务器配置字段: " + key,
				})
				return
			}
		}
	}

	// 验证 database 字段（禁止通过此接口修改密码）
	if dbInput, ok := input["database"].(map[string]interface{}); ok {
		allowedDbKeys := map[string]bool{"Db": true, "DbHost": true, "DbPort": true, "DbUser": true, "DbName": true}
		for key := range dbInput {
			if !allowedDbKeys[key] {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  errmsg.ERROR,
					"message": "不允许的数据库配置字段: " + key,
				})
				return
			}
		}
	}

	// 1. 读取现有 YAML 配置，通过 Config struct 解析（保证键名统一）
	var cfg utils.Config
	found := false
	for _, p := range []string{
		"config/backend/config.yaml",
		"config/config.yaml",
		"config/config_template.yaml",
	} {
		if raw, err := os.ReadFile(p); err == nil {
			if yaml.Unmarshal(raw, &cfg) != nil {
				cfg = utils.Config{} // 重置
			} else {
				found = true
				break
			}
		}
	}
	_ = found

	// 2. 将前端的 JSON 输入（小写 key）合并到 Config struct
	// 通过 JSON 中间格式统一键名：input → JSON → Config struct
	if inputJson, err := json.Marshal(input); err == nil {
		// 用 json tag 把前端输入写入 cfg（只会覆盖匹配的字段）
		json.Unmarshal(inputJson, &cfg)
	}

	// 3. 写回 YAML（用 yaml tag 输出规范的 camelCase 键名）
	configPath := utils.GetConfigPath("config/backend/config.yaml")
	_ = os.MkdirAll(filepath.Dir(configPath), 0755)
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "序列化配置失败: " + err.Error(),
		})
		return
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "保存配置文件失败: " + err.Error(),
		})
		return
	}

	// 重新加载配置到内存，使修改即时生效
	_ = utils.ReloadConfig()
	middlewares.RefreshJwtKey()

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "后端配置保存成功",
	})
}

func ReloadConfig(c *gin.Context) {
	err := utils.ReloadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "重新加载配置失败: " + err.Error(),
		})
		return
	}

	// 刷新 JWT 密钥
	middlewares.RefreshJwtKey()

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "配置重新加载成功",
	})
}

func GetAllConfig(c *gin.Context) {
	backendConfig := utils.GetConfig()

	// 过滤敏感字段，防止密码和密钥泄露
	safeBackend := map[string]interface{}{
		"server": map[string]interface{}{
			"AppMode":  backendConfig.Server.AppMode,
			"HttpPort": backendConfig.Server.HttpPort,
			"SiteUrl":  backendConfig.Server.SiteUrl,
		},
		"database": map[string]interface{}{
			"Db":     backendConfig.Database.Db,
			"DbHost": backendConfig.Database.DbHost,
			"DbPort": backendConfig.Database.DbPort,
			"DbUser": backendConfig.Database.DbUser,
			"DbName": backendConfig.Database.DbName,
		},
		"weather":            backendConfig.Weather,
		"FrontEndConfigPath": backendConfig.FrontEndConfigPath,
	}
	
	configPath := utils.GetFrontEndConfigPath()
	frontendContent, err := os.ReadFile(configPath)
	if err != nil {
		frontendContent = []byte("{}")
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        errmsg.SUCCESS,
		"backend":       safeBackend,
		"frontend_yaml": string(frontendContent),
		"message":       errmsg.GetErrMsg(errmsg.SUCCESS),
	})
}
