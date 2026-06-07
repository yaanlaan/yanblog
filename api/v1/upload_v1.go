package v1

import (
	"net/http"
	"yanblog/model"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

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
