package model

import (
	"strconv"
)

// GetArtIDByTitle 根据标题获取文章ID
func GetArtIDByTitle(title string) int {
	var art Article
	db.Select("id").Where("title = ?", title).First(&art)
	return int(art.ID)
}

// CheckUploadPermission 检查上传权限
// key: 文章标题
// idStr: 当前编辑的文章ID (字符串)
// 返回: true 表示允许上传, false 表示禁止
func CheckUploadPermission(key string, idStr string) bool {
	// 获取该标题对应的文章ID
	existingID := GetArtIDByTitle(key)

	if existingID == 0 {
		// 标题不存在，允许上传
		return true
	}

	// 标题存在
	currentID, _ := strconv.Atoi(idStr)

	// 如果是新增文章 (currentID == 0)，但标题已存在 -> 禁止
	if currentID == 0 {
		return false
	}

	// 如果是编辑文章，但ID不匹配 -> 禁止
	if currentID != existingID {
		return false
	}

	// ID匹配 -> 允许
	return true
}
