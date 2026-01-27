package v1

import (
	"io/ioutil"
	"net/http"
	"yanblog/utils"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// GetFrontEndConfig 获取前端配置
func GetFrontEndConfig(c *gin.Context) {
	configPath := utils.ServerConfig.FrontEndConfigPath
	if configPath == "" {
		configPath = "web/frontend/public/config.yaml" // 默认回退路径
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "读取配置文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"data":    string(content),
		"message": errmsg.GetErrMsg(errmsg.SUCCESS),
	})
}

// UpdateFrontEndConfig 更新前端配置
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

	// 验证 YAML 格式是否有效
	var temp interface{}
	if err := yaml.Unmarshal([]byte(input.Content), &temp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "YAML 格式错误: " + err.Error(),
		})
		return
	}

	configPath := utils.ServerConfig.FrontEndConfigPath
	if configPath == "" {
		configPath = "web/frontend/public/config.yaml"
	}

	// 写入文件
	err := ioutil.WriteFile(configPath, []byte(input.Content), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "保存配置文件失败: " + err.Error(),
		})
		return
	}

	// 如果是在 Docker 环境下运行，可能需要某种方式通知前端容器重载（不过Vue前端是静态文件，只要刷新浏览器一般即可，除非Vite开发服务器有缓存）
	// 对于静态 Nginx 部署，修改文件后下次请求就会获取新内容。

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "配置保存成功",
	})
}

// OpenFrontEndFile 打开前端文件(不仅仅是config，用于文件上传后的回调填入等，这里暂不需要)
func OpenFrontEndFile(c *gin.Context) {
	// 预留
}
