package model

import (
	"yanblog/utils/errmsg"

	"gorm.io/gorm"
)

type Tag struct {
	ID    uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name  string `gorm:"type:varchar(100);not null;unique" json:"name"`
	Count int    `gorm:"-" json:"count"` // 统计该标签下的文章数，不存库
}

// CheckTagExist 检查标签是否存在
func CheckTagExist(name string) int {
	var tag Tag
	db.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return errmsg.ERROR_TAG_EXIST
	}
	return errmsg.SUCCESS
}

// CheckTagWithID 检查标签名称是否被其他标签使用（排除指定ID）
func CheckTagWithID(id int, name string) int {
	var tag Tag
	db.Select("id").Where("name = ? AND id != ?", name, id).First(&tag)
	if tag.ID > 0 {
		return errmsg.ERROR_TAG_EXIST
	}
	return errmsg.SUCCESS
}

// CreateTag 新增标签
func CreateTag(data *Tag) int {
	err := db.Create(data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCESS
}

// GetTags 获取标签列表 (带文章计数)
// 注意：默认值处理已统一由 API 层的 ParsePageParams 处理
func GetTags(pageSize int, pageNum int) ([]Tag, int64) {
	var tags []Tag
	var total int64

	// 先计算总数
	db.Model(&Tag{}).Count(&total)

	var err error
	// 查询标签列表
	if pageSize == -1 && pageNum == -1 {
		err = db.Find(&tags).Error
	} else {
		err = db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&tags).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}

	// 一次性统计所有标签下的文章数
	type CountResult struct {
		TagID uint
		Count int
	}
	var counts []CountResult
	db.Table("article_tags").Select("tag_id, COUNT(*) as count").Group("tag_id").Scan(&counts)
	countMap := make(map[uint]int, len(counts))
	for _, cr := range counts {
		countMap[cr.TagID] = cr.Count
	}
	for i := range tags {
		tags[i].Count = countMap[tags[i].ID]
	}

	return tags, total
}

// EditTag 编辑标签
func EditTag(id int, data *Tag) int {
	var tag Tag
	var maps = make(map[string]interface{})
	maps["name"] = data.Name

	err := db.Model(&tag).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// DeleteTag 删除标签
func DeleteTag(id int) int {
	var tag Tag
	err := db.Where("id = ?", id).Delete(&tag).Error
	if err != nil {
		return errmsg.ERROR
	}
	// 同时删除中间表关联? Gorm 的级联删除需要配置
	// 手动清理 pivot 表 (article_tags)
	db.Exec("DELETE FROM article_tags WHERE tag_id = ?", id)
	return errmsg.SUCCESS
}
