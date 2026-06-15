package v1

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"yanblog/model"
	"yanblog/utils"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
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
	id, ok := utils.ParseIDParam(c)
	if !ok {
		return
	}

	data, code := model.GetCateInfo(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询分类列表
func GetCate(c *gin.Context) {
	pageSize, pageNum, _ := utils.ParsePageParams(c)
	data, total := model.GetCate(pageSize, pageNum)
	code := errmsg.SUCCESS

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// 搜索分类
func SearchCate(c *gin.Context) {
	pageSize, pageNum, _ := utils.ParsePageParams(c)
	keyword := c.Query("keyword")
	code := errmsg.SUCCESS

	data, total := model.SearchCategory(keyword, pageSize, pageNum)

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
	id, ok := utils.ParseIDParam(c)
	if !ok {
		return
	}
	_ = c.ShouldBindJSON(&data)

	code = model.CheckCategoryWithID(id, data.Name)

	if code == errmsg.SUCCESS {
		// 仅确保 top 字段有效
		if data.Top < 0 {
			data.Top = 0
		}
		code = model.EditCate(id, &data)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 删除分类
func DeleteCate(c *gin.Context) {
	id, ok := utils.ParseIDParam(c)
	if !ok {
		return
	}
	force, _ := strconv.ParseBool(c.Query("force")) // 是否强制删除关联文章

	var code int

	// 1. 获取分类信息（用于删除封面图）
	cate, code := model.GetCateInfo(id)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	// 2. 处理关联文章
	if force {
		// 获取该分类下所有文章
		articles, _, _ := model.GetCateArt(id, -1, -1)
		for _, art := range articles {
			// 删除文章文件夹
			if art.Title != "" {
				targetDir := filepath.Join("uploads", "articles", filepath.Clean(art.Title))
				_ = os.RemoveAll(targetDir)
			}
			// 删除文章记录
			model.DeleteArt(int(art.ID))
		}
	}

	// 3. 删除分类封面图 (如果是本地文件)
	if cate.Img != "" {
		// 简单判断是否在 uploads 下
		// 实际路径可能是 /uploads/category/123.jpg
		// os.Remove 需要相对路径 uploads/category/123.jpg
		// 去掉第一个 /
		if len(cate.Img) > 0 && cate.Img[0] == '/' {
			_ = os.Remove("." + cate.Img)
		} else {
			_ = os.Remove(cate.Img)
		}
	}

	code = model.DeleteCate(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
