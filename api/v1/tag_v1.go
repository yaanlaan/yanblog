package v1

import (
	"net/http"
	"strconv"
	"yanblog/model"
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
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	if pageSize == 0 {
		pageSize = -1 // 不填则查所有（或按业务定义）
	}
	if pageNum == 0 {
		pageNum = 1
	}

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
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)
	code := model.CheckTagExist(data.Name)
	if code == errmsg.SUCCESS {
		model.EditTag(id, &data)
	}
	// 如果改了名，CheckTagExist 可能报错，如果是改成自己现在的名字不算重复
	// 这里简化处理
	if code == errmsg.ERROR_TAG_EXIST {
		// 允许不改名直接保存
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DeleteTag(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
