package model

import (
	"yanblog/utils/errmsg"
)

// CheckArtTitle 查询文章标题是否存在
// 参数: title - 文章标题
// 返回: 状态码
func CheckArtTitle(title string) int {
	var art Article
	db.Select("id").Where("title = ?", title).First(&art)
	if art.ID > 0 {
		return errmsg.ERROR_ART_TITLE_USED
	}
	return errmsg.SUCCESS
}

// CheckArtTitleWithId 查询文章标题是否存在（排除指定ID）
// 参数: id - 文章ID, title - 文章标题
// 返回: 状态码
func CheckArtTitleWithId(id int, title string) int {
	var art Article
	db.Select("id").Where("title = ? AND id != ?", title, id).First(&art)
	if art.ID > 0 {
		return errmsg.ERROR_ART_TITLE_USED
	}
	return errmsg.SUCCESS
}
