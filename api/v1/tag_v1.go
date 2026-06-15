package v1

import (
	"net/http"
	"yanblog/model"
	"yanblog/utils"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// AddTag 添加标签
func AddTag(c *gin.Context) {
	var data model.Tag
	_ = c.ShouldBindJSON(&data)
	code := model.CheckTagExist(data.Name)
	if code == errmsg.SUCCESS {
		model.CreateTag(&data)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetTags 获取标签列表
func GetTags(c *gin.Context) {
	pageSize, pageNum, _ := utils.ParsePageParams(c)

	data, total := model.GetTags(pageSize, pageNum)
	code := errmsg.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// EditTag 编辑标签
func EditTag(c *gin.Context) {
	var data model.Tag
	id, ok := utils.ParseIDParam(c)
	if !ok {
		return
	}
	_ = c.ShouldBindJSON(&data)

	// 检查标签名是否与其他标签冲突（排除自身）
	code := model.CheckTagWithID(id, data.Name)
	if code == errmsg.SUCCESS {
		code = model.EditTag(id, &data)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	id, ok := utils.ParseIDParam(c)
	if !ok {
		return
	}
	code := model.DeleteTag(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
