package v1

import (
	"encoding/json"
	"io/ioutil"
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

	// 1. 读取现有 YAML 配置，通过 Config struct 解析（保证键名统一）
	var cfg utils.Config
	found := false
	for _, p := range []string{
		"config/backend/config.yaml",
		"config/config.yaml",
		"config/config_template.yaml",
	} {
		if raw, err := ioutil.ReadFile(p); err == nil {
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

	err = ioutil.WriteFile(configPath, data, 0644)
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

// deepMerge 递归合并 src 到 dst，保留 dst 中已有的字段
func deepMerge(dst, src map[string]interface{}) {
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			// 两边都是 map，递归合并
			srcMap, srcIsMap := srcVal.(map[string]interface{})
			dstMap, dstIsMap := dstVal.(map[string]interface{})
			if srcIsMap && dstIsMap {
				deepMerge(dstMap, srcMap)
				continue
			}
		}
		// 否则直接覆盖或新增
		dst[key] = srcVal
	}
}