package model

import (
	"yanblog/utils/errmsg"
	"strings"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null" json:"name"`
	Img  string `gorm:"type:varchar(255)" json:"img"`
	Top  int    `gorm:"type:int;not null;default:0" json:"top"`
	// 添加文章计数字段（使用gorm:"-"标记，表示不直接映射到数据库字段，保证数据一致性）
	ArticleCount int `gorm:"-" json:"article_count"`
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

// CheckCategoryWithID 查询分类是否存在（排除指定ID的分类）
// 参数: id - 要排除的分类ID, name - 分类名称
// 返回: 状态码
func CheckCategoryWithID(id int, name string) (code int) {
	var cate Category
	db.Select("id").Where("name = ? AND id != ?", name, id).First(&cate)
	if cate.ID > 0 {
		return errmsg.ERROR_CATENAME_USED
	}
	return errmsg.SUCCESS
}

// CreateCate 新增分类
// 参数: data - 分类信息
// 返回: 状态码
func CreateCate(data *Category) int {
	// 保持img字段为空或用户输入的值，不再设置默认值
	
	// 确保top字段有效
	if data.Top < 0 {
		data.Top = 0
	}
	
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// SearchCategory 搜索分类
// 参数: keyword - 搜索关键词, pageSize - 每页数量, pageNum - 页码
// 返回: 分类列表和总数
func SearchCategory(keyword string, pageSize int, pageNum int) ([]Category, int64) {
	var cate []Category
	var total int64
	var err error

	// 构建查询条件
	query := db

	// 如果有关键词，则添加分类名的模糊搜索
	if keyword != "" {
		searchTerm := "%" + strings.ToLower(keyword) + "%"
		query = query.Where("LOWER(name) LIKE ?", searchTerm)
	}

	// 执行查询
	if pageSize == -1 && pageNum == -1 {
		err = query.Find(&cate).Count(&total).Error
	} else {
		err = query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&cate).Count(&total).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}
	
	return cate, total
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
	
	// 为每个分类获取文章数量
	for i := range cate {
		var count int64
		db.Model(&Article{}).Where("cid = ?", cate[i].ID).Count(&count)
		cate[i].ArticleCount = int(count)
	}
	
	return cate, total
}

// GetCateInfo 获取单个分类信息（包含文章数量）
// 参数: id - 分类ID
// 返回: 分类信息和状态码
func GetCateInfo(id int) (Category, int) {
	var cate Category
	err := db.Where("id = ?", id).First(&cate).Error
	if err != nil {
		return cate, errmsg.ERROR_CATE_NOT_EXIST
	}
	
	// 获取该分类下的文章数量
	var count int64
	db.Model(&Article{}).Where("cid = ?", id).Count(&count)
	cate.ArticleCount = int(count)
	
	return cate, errmsg.SUCCESS
}

// EditCate 编辑分类信息
// 参数: id - 分类ID, data - 更新的分类信息
// 返回: 状态码
func EditCate(id int, data *Category) int {
	var cate Category
	var maps = make(map[string]interface{})
	maps["name"] = data.Name
	maps["img"] = data.Img
	maps["top"] = data.Top

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