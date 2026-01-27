package v1

import (
	"io/ioutil"
	"net/http"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// 关于页面内容文件路径
const AboutFilePath = "./web/frontend/public/static/about.md"

// GetAboutContent 获取关于页面内容
func GetAboutContent(c *gin.Context) {
	content, err := ioutil.ReadFile(AboutFilePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "Failed to read file",
			"error":   err.Error(),
		})
		return
	}

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
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "Invalid data",
		})
		return
	}

	// 写入文件
	err := ioutil.WriteFile(AboutFilePath, []byte(data.Content), 0644)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "Failed to write file",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": errmsg.GetErrMsg(errmsg.SUCCESS),
	})
}
