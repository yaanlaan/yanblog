package v1

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"yanblog/model"
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}
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
		// c.Abort() // 移除 Abort，确保返回 JSON
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 删除分类
func DeleteCate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
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
