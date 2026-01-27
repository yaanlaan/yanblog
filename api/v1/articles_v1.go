package v1

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"yanblog/model"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// 添加文章
func AddArticle(c *gin.Context) {
	var data model.Article
	var code int
	_ = c.ShouldBindJSON(&data)

	// 检查标题是否重复
	code = model.CheckArtTitle(data.Title)
	if code == errmsg.SUCCESS {
		code = model.CreateArt(&data)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetArchive 获取归档文章
func GetArchive(c *gin.Context) {
	data, code := model.GetArchive()
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询分类下的所有文章
func GetCateArt(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	var code int

	if pageSize <= 0 {
		pageSize = -1
	}
	if pageNum <= 0 {
		pageNum = -1
	}

	data, code, total := model.GetCateArt(id, pageSize, pageNum)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询单个文章信息
func GetArtInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	// 增加阅读量
	model.IncrementArtViews(id)

	var code int
	data, code := model.GetArtInfo(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询文章列表
func GetArt(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	var code int

	if pageSize == 0 {
		pageSize = -1
	}
	if pageNum == 0 {
		pageNum = -1
	}
	data, code, total := model.GetArt(pageSize, pageNum)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// 搜索文章
func SearchArt(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	keyword := c.Query("keyword")
	cid, _ := strconv.Atoi(c.Query("cid"))
	var code int

	// 处理默认值
	if pageSize == 0 {
		pageSize = -1
	}
	if pageNum == 0 {
		pageNum = -1
	}

	data, code, total := model.SearchArticle(keyword, cid, pageSize, pageNum)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询置顶文章
func GetTopArt(c *gin.Context) {
	// 默认获取6篇置顶文章
	num, _ := strconv.Atoi(c.Query("num"))
	var code int
	if num <= 0 {
		num = 6
	}

	data, code := model.GetTopArt(num)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询热门文章
func GetHotArt(c *gin.Context) {
	num, _ := strconv.Atoi(c.Query("num"))
	var code int
	if num <= 0 {
		num = 5
	}

	data, code := model.GetHotArticles(num)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询相关文章
func GetRelatedArt(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	// 从查询参数获取 tags (通常前端不方便传太长，但这里可以接受 query param 或者 POST，这里选 query param 简单点，或者先获取文章详情再查？
	// 为了效率，让前端作为参数传过来比较好，避免后端多查一次数据库)
	// 或者，更安全的做法是后端先查 id 对应的 tags，再查相关文章。这样避免前端瞎传 tags。
	// 这里选择后端自查，虽然多一次 DB 交互，但逻辑更闭环。

	art, code := model.GetArtInfo(id)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	data, code := model.GetRelatedArticles(id, art.Tags)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// 编辑文章
func EditArt(c *gin.Context) {
	var data model.Article
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

	// 检查标题是否重复
	code = model.CheckArtTitleWithId(id, data.Title)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	// 获取旧文章信息，检查标题是否变更
	oldArt, code := model.GetArtInfo(id)
	if code == errmsg.SUCCESS && oldArt.Title != data.Title {
		// 标题变更，需要重命名文件夹并更新内容中的图片路径
		oldTitleClean := filepath.Clean(oldArt.Title)
		newTitleClean := filepath.Clean(data.Title)

		oldDir := filepath.Join("uploads", "articles", oldTitleClean)
		newDir := filepath.Join("uploads", "articles", newTitleClean)

		// 检查旧文件夹是否存在
		if _, err := os.Stat(oldDir); err == nil {
			// 重命名文件夹
			// 注意：如果新文件夹已存在（理论上不应该，因为检查了标题唯一性，但可能有残留），则可能失败
			// 这里简单处理，直接重命名
			err := os.Rename(oldDir, newDir)
			if err == nil {
				// 更新内容中的图片路径
				// 将 /uploads/articles/OldTitle/ 替换为 /uploads/articles/NewTitle/
				// 注意路径分隔符，URL中使用 /
				oldUrlPart := "/uploads/articles/" + filepath.ToSlash(oldTitleClean) + "/"
				newUrlPart := "/uploads/articles/" + filepath.ToSlash(newTitleClean) + "/"

				data.Content = strings.ReplaceAll(data.Content, oldUrlPart, newUrlPart)

				// 如果封面图也在该目录下，也需要更新
				if strings.Contains(data.Img, oldUrlPart) {
					data.Img = strings.ReplaceAll(data.Img, oldUrlPart, newUrlPart)
				}
			}
		}
	}

	code = model.EditArt(id, &data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 删除文章
func DeleteArt(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}
	var code int

	// 1. 获取文章信息以找到对应的文件夹
	data, code := model.GetArtInfo(id)
	if code == errmsg.SUCCESS {
		// 2. 删除文章对应的文件夹
		// 尝试删除 uploads/articles/{Title}
		if data.Title != "" {
			targetDir := filepath.Join("uploads", "articles", filepath.Clean(data.Title))
			_ = os.RemoveAll(targetDir)
		}
	}

	code = model.DeleteArt(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
