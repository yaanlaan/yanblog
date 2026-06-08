package v1

import (
	"io/ioutil"
	"net/http"
	"yanblog/middlewares"
	"yanblog/utils"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func GetFrontEndConfig(c *gin.Context) {
	configPath := utils.GetFrontEndConfigPath()

	content, err := ioutil.ReadFile(configPath)
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

	err := ioutil.WriteFile(configPath, []byte(input.Content), 0644)
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
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"data":    config,
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

	configPath := utils.GetConfigPath("config/backend/config.yaml")
	data, err := yaml.Marshal(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "序列化配置失败: " + err.Error(),
		})
		return
	}

	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "保存配置文件失败: " + err.Error(),
		})
		return
	}

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
	
	configPath := utils.GetFrontEndConfigPath()
	frontendContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		frontendContent = []byte("{}")
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        errmsg.SUCCESS,
		"backend":       backendConfig,
		"frontend_yaml": string(frontendContent),
		"message":       errmsg.GetErrMsg(errmsg.SUCCESS),
	})
}