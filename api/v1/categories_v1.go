package v1

import (
	"yanblog/model"
	"yanblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 添加分类
func AddCategory(c *gin.Context) {
	var data model.Category
	var code int
	_ = c.ShouldBindJSON(&data)
	code = model.CheckCategory(data.Name)
	if code == errmsg.SUCCESS {
		// 仅设置top字段的默认值
		if data.Top < 0 {
			data.Top = 0
		}
		model.CreateCate(&data)
	}
	if code == errmsg.ERROR_CATENAME_USED {
		code = errmsg.ERROR_CATENAME_USED
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询分类信息
func GetCateInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var code int

	data, code := model.GetCateInfo(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询分类列表
func GetCate(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	var code int

	if pageSize == -1 && pageNum == -1 {
		// 查询所有分类
		data, total := model.GetCate(-1, -1)
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	data, total := model.GetCate(pageSize, pageNum)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// 搜索分类
func SearchCate(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	keyword := c.Query("keyword")
	var code int

	if pageSize <= 0 {
		pageSize = -1
	}
	if pageNum <= 0 {
		pageNum = -1
	}

	data, total := model.SearchCategory(keyword, pageSize, pageNum)

	code = errmsg.SUCCESS
	
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// 编辑分类信息
func EditCate(c *gin.Context) {
	var data model.Category
	var code int
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)
	
	code = model.CheckCategoryWithID(id, data.Name)

	if code == errmsg.SUCCESS {
		// 不再设置默认值，保持 img 字段为空或用户输入的值
		// 仅确保 top 字段有效
		if data.Top < 0 {
			data.Top = 0
		}
		code = model.EditCate(id, &data)
	}
	if code == errmsg.ERROR_CATENAME_USED {
		c.Abort()
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 删除用户
func DeleteCate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var code int

	code = model.DeleteCate(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}