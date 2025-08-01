package model

import (
	"yanblog/utils/errmsg"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}

// CheckCategory 查询分类是否存在
// 参数: name - 分类名称
// 返回: 状态码
func CheckCategory(name string) (code int) {
	var cate Category
	db.Select("id").Where("name = ?", name).First(&cate)
	if cate.ID > 0 {
		return errmsg.ERROR_CATENAME_USED
	}
	return errmsg.SUCCESS
}

// CreateCate 新增分类
// 参数: data - 分类信息
// 返回: 状态码
func CreateCate(data *Category) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}


// GetCate 查询分类列表
// 参数: pageSize - 每页数量, pageNum - 页码
// 返回: 分类列表和总数
func GetCate(pageSize int, pageNum int) ([]Category, int64) {
	var cate []Category
	var total int64
	var err error

	if pageSize == -1 && pageNum == -1 {
		err = db.Find(&cate).Count(&total).Error
	} else {
		err = db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&cate).Count(&total).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	
	return cate, total
}

// EditCate 编辑分类信息
// 参数: id - 分类ID, data - 更新的分类信息
// 返回: 状态码
func EditCate(id int, data *Category) int {
	var cate Category
	var maps = make(map[string]interface{})
	maps["name"] = data.Name

	err = db.Model(&cate).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 删除分类
func DeleteCate(id int) int {
	var cate Category
	err = db.Where("id = ? ", id).Delete(&cate).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
