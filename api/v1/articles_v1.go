package v1

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"yanblog/model"
	"yanblog/utils"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// 添加文章
func AddArticle(c *gin.Context) {
	var input struct {
		Title     string `json:"title"`
		Cid       int    `json:"cid"`
		Desc      string `json:"desc"`
		Content   string `json:"content"`
		Img       string `json:"img"`
		Top       int    `json:"top"`
		Tags      string `json:"tags"`
		Type      int    `json:"type"`
		PdfUrl    string `json:"pdf_url"`
		CreatedAt string `json:"createdAt"`
	}
	var code int
	_ = c.ShouldBindJSON(&input)

	data := model.Article{
		Title:   input.Title,
		Cid:     input.Cid,
		Desc:    input.Desc,
		Content: input.Content,
		Img:     input.Img,
		Top:     input.Top,
		Tags:    input.Tags,
		Type:    input.Type,
		PdfUrl:  input.PdfUrl,
	}
	if input.CreatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", input.CreatedAt); err == nil {
			data.CreatedAt = t
		}
	}

	code = model.CheckArtTitle(data.Title)
	if code == errmsg.SUCCESS {
		code = model.CreateArt(&data)
	}

	utils.ErrorWithMessage(c, code, errmsg.GetErrMsg(code))
}

// GetArchive 获取归档文章
func GetArchive(c *gin.Context) {
	data, _ := model.GetArchive()
	utils.Success(c, data)
}

// 查询分类下的所有文章
func GetCateArt(c *gin.Context) {
	pageSize, pageNum, _ := utils.ParsePageParams(c)
	id, ok := utils.ParseIDParam(c)
	if !ok {
		return
	}

	data, _, total := model.GetCateArt(id, pageSize, pageNum)
	utils.SuccessWithTotal(c, data, total)
}

// 查询单个文章信息
func GetArtInfo(c *gin.Context) {
	id, ok := utils.ParseIDParam(c)
	if !ok {
		return
	}

	data, code := model.GetArtInfo(id)
	if code != errmsg.SUCCESS {
		utils.Error(c, code)
		return
	}

	// 仅在文章存在时增加阅读量
	model.IncrementArtViews(id)
	utils.Success(c, data)
}

// 查询文章列表
func GetArt(c *gin.Context) {
	pageSize, pageNum, _ := utils.ParsePageParams(c)
	excludeTop := c.Query("excludeTop") == "true"

	data, _, total := model.GetArt(pageSize, pageNum, excludeTop)
	utils.SuccessWithTotal(c, data, total)
}

// 搜索文章
func SearchArt(c *gin.Context) {
	pageSize, pageNum, _ := utils.ParsePageParams(c)
	keyword := c.Query("keyword")
	cid, _ := strconv.Atoi(c.Query("cid"))

	data, _, total := model.SearchArticle(keyword, cid, pageSize, pageNum)
	utils.SuccessWithTotal(c, data, total)
}

// 查询置顶文章
func GetTopArt(c *gin.Context) {
	num, _ := strconv.Atoi(c.Query("num"))
	if num <= 0 {
		num = 6
	}
	// 限制最大返回数量，防止恶意请求
	if num > utils.MaxPageSize {
		num = utils.MaxPageSize
	}
	data, _ := model.GetTopArt(num)
	utils.Success(c, data)
}

// 查询热门文章
func GetHotArt(c *gin.Context) {
	num, _ := strconv.Atoi(c.Query("num"))
	if num <= 0 {
		num = 5
	}
	// 限制最大返回数量，防止恶意请求
	if num > utils.MaxPageSize {
		num = utils.MaxPageSize
	}
	data, _ := model.GetHotArticles(num)
	utils.Success(c, data)
}

// 随机获取一篇文章
func GetRandomArt(c *gin.Context) {
	data, _ := model.GetRandomArticle()
	utils.Success(c, data)
}

// 获取相邻文章（上一篇/下一篇）
func GetAdjacentArt(c *gin.Context) {
	id, ok := utils.ParseIDParam(c)
	if !ok {
		return
	}

	prev, next, code := model.GetAdjacentArticle(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data": gin.H{
			"previous": prev,
			"next":     next,
		},
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询相关文章
func GetRelatedArt(c *gin.Context) {
	id, ok := utils.ParseIDParam(c)
	if !ok {
		return
	}

	art, code := model.GetArtInfo(id)
	if code != errmsg.SUCCESS {
		utils.Error(c, code)
		return
	}

	data, code := model.GetRelatedArticles(id, art.Tags)
	utils.Success(c, data)
}

// 编辑文章
func EditArt(c *gin.Context) {
	var data model.Article
	var code int
	id, ok := utils.ParseIDParam(c)
	if !ok {
		return
	}
	_ = c.ShouldBindJSON(&data)

	code = model.CheckArtTitleWithId(id, data.Title)
	if code != errmsg.SUCCESS {
		utils.Error(c, code)
		return
	}

	// 获取旧文章信息，检查标题是否变更
	oldArt, code := model.GetArtInfo(id)
	if code == errmsg.SUCCESS && oldArt.Title != data.Title {
		oldTitleClean := filepath.Clean(oldArt.Title)
		newTitleClean := filepath.Clean(data.Title)

		oldDir := filepath.Join("uploads", "articles", oldTitleClean)
		newDir := filepath.Join("uploads", "articles", newTitleClean)

		if _, err := os.Stat(oldDir); err == nil {
			err := os.Rename(oldDir, newDir)
			if err == nil {
				oldUrlPart := "/uploads/articles/" + filepath.ToSlash(oldTitleClean) + "/"
				newUrlPart := "/uploads/articles/" + filepath.ToSlash(newTitleClean) + "/"

				data.Content = strings.ReplaceAll(data.Content, oldUrlPart, newUrlPart)

				if strings.Contains(data.Img, oldUrlPart) {
					data.Img = strings.ReplaceAll(data.Img, oldUrlPart, newUrlPart)
				}
			}
		}
	}

	code = model.EditArt(id, &data)
	utils.Error(c, code)
}

// 删除文章
func DeleteArt(c *gin.Context) {
	id, ok := utils.ParseIDParam(c)
	if !ok {
		return
	}

	// 获取文章信息以找到对应的文件夹
	data, code := model.GetArtInfo(id)
	if code == errmsg.SUCCESS && data.Title != "" {
		targetDir := filepath.Join("uploads", "articles", filepath.Clean(data.Title))
		_ = os.RemoveAll(targetDir)
	}

	code = model.DeleteArt(id)
	utils.Error(c, code)
}

// 批量删除文章
func BatchDeleteArt(c *gin.Context) {
	var data struct {
		Ids []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&data); err != nil || len(data.Ids) == 0 {
		utils.BadRequest(c, "参数错误，需要 ids 数组")
		return
	}

	// 先清理各文章的关联文件夹
	for _, id := range data.Ids {
		art, code := model.GetArtInfo(id)
		if code == errmsg.SUCCESS && art.Title != "" {
			targetDir := filepath.Join("uploads", "articles", filepath.Clean(art.Title))
			_ = os.RemoveAll(targetDir)
		}
	}

	deleted, failed := model.BatchDeleteArts(data.Ids)

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"data": gin.H{
			"deleted": deleted,
			"failed":  failed,
			"total":   len(data.Ids),
		},
		"message": fmt.Sprintf("成功删除 %d 篇，失败 %d 篇", deleted, failed),
	})
}
