package model

import (
	// "fmt"
	// "os"
	"yanblog/utils/errmsg"
	"strings"

	"gorm.io/gorm"
)

type Article struct {
	Category Category `gorm:"foreignkey:Cid"`
	gorm.Model
	Title   string `gorm:"type:varchar(100);not null" json:"title"`
	Cid     int    `gorm:"type:int;not null" json:"cid"`
	Desc    string `gorm:"type:varchar(200)" json:"desc"`
	Content string `gorm:"type:longtext" json:"content"`
	Img     string `gorm:"type:varchar(100)" json:"img"`
	Top     int    `gorm:"type:int;default:0" json:"top"` // 0表示不置顶，其他数字1-6表示置顶等级，数字越小等级越高
}

// CreateArt 新增文章
// 参数: data - 文章信息
// 返回: 状态码
func CreateArt(data *Article) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCESS
}

// SearchArticle 搜索文章
// 参数: keyword - 搜索关键词, cid - 分类ID, pageSize - 每页数量, pageNum - 页码
// 返回: 文章列表、状态码和总数
func SearchArticle(keyword string, cid int, pageSize int, pageNum int) ([]Article, int, int64) {
	var articleList []Article
	var total int64
	var err error

	// 构建查询条件
	query := db.Preload("Category")

	// 如果有关键词，则添加标题和描述的模糊搜索
	if keyword != "" {
		searchTerm := "%" + strings.ToLower(keyword) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(desc) LIKE ?", searchTerm, searchTerm)
	}

	// 如果有分类ID，则添加分类筛选
	if cid > 0 {
		query = query.Where("cid = ?", cid)
	}

	// 添加排序
	query = query.Order("top ASC, created_at DESC")

	// 执行查询
	if pageSize == -1 && pageNum == -1 {
		// 查询所有
		err = query.Find(&articleList).Count(&total).Error
	} else {
		// 分页查询
		err = query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&articleList).Count(&total).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errmsg.ERROR, 0
	}
	return articleList, errmsg.SUCCESS, total
}

// GetCateArt 查询分类下的所有文章
// 参数: id - 分类ID, pageSize - 每页数量, pageNum - 页码
// 返回: 文章列表、状态码和总数
func GetCateArt(id int, pageSize int, pageNum int) ([]Article, int, int64) {
	var cateArtList []Article
	var total int64
	var err error
	
	if pageSize == -1 || pageNum == -1 {
		err = db.Preload("Category").Where("cid = ?", id).Order("top ASC, created_at DESC").Find(&cateArtList).Count(&total).Error
	}else{
		err = db.Preload("Category").Limit(pageSize).Offset((pageNum-1)*pageSize).Where("cid =?", id).Order("top ASC, created_at DESC").Find(&cateArtList).Count(&total).Error
	}
	
	if err != nil {
		return nil, errmsg.ERROR_CATE_NOT_EXIST, 0
	}
	return cateArtList, errmsg.SUCCESS, total
}

// GetArtInfo 查询单个文章
// 参数: id - 文章ID
// 返回: 文章信息和状态码
func GetArtInfo(id int) (Article, int) {
	var art Article
	err := db.Preload("Category").Where("id = ?", id).First(&art).Error
	if err != nil {
		return art, errmsg.ERROR_ART_NOT_EXIST
	}
	return art, errmsg.SUCCESS
}

// GetArt 查询文章列表
// 参数: pageSize - 每页数量, pageNum - 页码
// 返回: 文章列表、状态码和总数
func GetArt(pageSize int, pageNum int) ([]Article, int, int64) {
	var articleList []Article
	var err error
	var total int64
	if pageSize == -1 && pageNum == -1 { // 查询所有文章
		err = db.Preload("Category").Order("top ASC, created_at DESC").Find(&articleList).Count(&total).Error
	} else { // 分页查询	
		err = db.Preload("Category").Order("top ASC, created_at DESC").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&articleList).Count(&total).Error	
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errmsg.ERROR, 0
	}
	return articleList, errmsg.SUCCESS, total
}

// GetTopArt 查询置顶文章
// 参数: num - 查询置顶文章的数量
// 返回: 置顶文章列表和状态码
func GetTopArt(num int) ([]Article, int) {
	var topArtList []Article
	var err error
	
	err = db.Preload("Category").Where("top > 0").Order("top ASC").Limit(num).Find(&topArtList).Error
	if err != nil {
		return nil, errmsg.ERROR
	}
	return topArtList, errmsg.SUCCESS
}

// EditArt 编辑文章
// 参数: id - 文章ID, data - 更新的文章信息
// 返回: 状态码
func EditArt(id int, data *Article) int {
	var art Article
	var maps = make(map[string]interface{})
	maps["title"] = data.Title
	maps["cid"] = data.Cid
	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["img"] = data.Img
	maps["top"] = data.Top

	err = db.Model(&art).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// DeleteArt 删除文章
// 参数: id - 文章ID
// 返回: 状态码
func DeleteArt(id int) int {
	var art Article
	err = db.Where("id = ? ", id).Delete(&art).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}