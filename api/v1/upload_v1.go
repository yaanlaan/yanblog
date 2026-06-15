package v1

import (
	"net/http"
	"path/filepath"
	"strings"
	"yanblog/model"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// 允许上传的文件扩展名白名单
var allowedExtensions = map[string]bool{
	".jpg":  true, ".jpeg": true, ".png": true, ".gif": true,
	".webp": true, ".svg": true, ".bmp": true, ".ico": true,
	".pdf":  true, ".doc": true, ".docx": true, ".xls": true,
	".xlsx": true, ".ppt": true, ".pptx": true, ".txt": true,
	".md":   true, ".zip": true, ".rar": true, ".7z": true,
	".mp3":  true, ".mp4": true, ".avi": true, ".mov": true,
	".json": true, ".csv": true, ".xml": true,
}

// 单文件最大大小（10MB）
const maxFileSize = 10 << 20 // 10MB

func UpLoad(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "文件上传失败: " + err.Error(),
			"url":     "",
		})
		return
	}
	defer file.Close()

	// 校验文件大小
	if fileHeader.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "文件过大，单文件最大支持 10MB",
			"url":     "",
		})
		return
	}

	// 校验文件扩展名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtensions[ext] && ext != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "不支持的文件类型: " + ext,
			"url":     "",
		})
		return
	}

	// 获取额外的参数
	uploadType := c.PostForm("type") // article, category, common
	key := c.PostForm("key")         // 文章标题或其他标识
	idStr := c.PostForm("id")        // 文章ID (可选)

	// 检查文章标题冲突 (仅针对文章上传)
	if uploadType == "article" && key != "" && key != "default" {
		// 检查标题是否存在
		code := model.CheckArtTitle(key)
		if code == errmsg.ERROR_ART_TITLE_USED {
			// 标题已存在，检查是否属于当前编辑的文章
			// 如果 idStr 为空或为 "0"，说明是新增文章，此时标题存在就是冲突
			// 如果 idStr 不为空，需要检查该标题对应的ID是否与 idStr 一致

			// 简单逻辑：如果标题存在，且不是当前正在编辑的文章（通过ID判断），则禁止上传
			if !model.CheckUploadPermission(key, idStr) {
				c.JSON(http.StatusOK, gin.H{
					"status":  errmsg.ERROR_ART_TITLE_USED,
					"message": "文章标题已存在，请更换标题后再上传图片",
					"url":     "",
				})
				return
			}
		}
	}

	url, code := model.UpLoadFile(file, fileHeader, uploadType, key)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"url":     url,
	})
}
