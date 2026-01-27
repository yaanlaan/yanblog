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
		return errmsg.ERROR_TAG_EXIST // 认为已存在（视业务逻辑而定，这里用于新增时的检查）
	}
	return errmsg.SUCCESS
}

// CreateTag 新增标签
func CreateTag(data *Tag) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCESS
}

// GetTags 获取标签列表 (带文章计数)
func GetTags(pageSize int, pageNum int) ([]Tag, int64) {
	var tags []Tag
	var total int64

	// 计算总数
	db.Model(&Tag{}).Count(&total)

	// 查询标签列表
	err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}

	// 统计每个标签下的文章数
	// 这里需要利用 association 或者手动 count
	// 由于 Count 是 gorm:"-"，我们需要手动填充
	for i, tag := range tags {
		var count int64
		// 使用 join 查询 article_tags 表
		db.Table("article_tags").Where("tag_id = ?", tag.ID).Count(&count)
		tags[i].Count = int(count)
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
