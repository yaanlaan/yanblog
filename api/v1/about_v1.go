package v1

import (
	"net/http"
	"os"
	"path/filepath"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// getAboutReadFilePath 获取关于页面读取路径（优先使用持久化数据）
func getAboutReadFilePath() string {
	paths := []string{
		"data/about.md",
		"./web/frontend/public/static/about.md",
		"../web/frontend/public/static/about.md",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "data/about.md"
}

// GetAboutContent 获取关于页面内容
func GetAboutContent(c *gin.Context) {
	filePath := getAboutReadFilePath()
	content, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "读取关于页面内容失败",
			"error":   err.Error(),
		})
		return
	}

	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"data":    string(content),
		"message": errmsg.GetErrMsg(errmsg.SUCCESS),
	})
}

// UpdateAboutContent 更新关于页面内容
func UpdateAboutContent(c *gin.Context) {
	var data struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	// 优先持久化存储在 data/about.md
	dataPath := "data/about.md"
	if err := os.MkdirAll(filepath.Dir(dataPath), 0755); err == nil {
		if err := os.WriteFile(dataPath, []byte(data.Content), 0644); err == nil {
			// 同步写回静态目录以备回退（如果目录存在）
			staticPath := "./web/frontend/public/static/about.md"
			if _, statErr := os.Stat(filepath.Dir(staticPath)); statErr == nil {
				_ = os.WriteFile(staticPath, []byte(data.Content), 0644)
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.SUCCESS,
				"message": errmsg.GetErrMsg(errmsg.SUCCESS),
			})
			return
		}
	}

	// 降级写回静态文件
	fallbackPath := "./web/frontend/public/static/about.md"
	if err := os.WriteFile(fallbackPath, []byte(data.Content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "写入文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": errmsg.GetErrMsg(errmsg.SUCCESS),
	})
}
