package model

import (
	// "fmt"
	// "os"
	"strings"
	"yanblog/utils/errmsg"

	"gorm.io/gorm"
)

type Article struct {
	Category Category `gorm:"foreignkey:Cid"`
	gorm.Model
	Title     string `gorm:"type:varchar(100);not null" json:"title"`
	Cid       int    `gorm:"type:int;not null" json:"cid"`
	Desc      string `gorm:"type:varchar(200)" json:"desc"`
	Content   string `gorm:"type:longtext" json:"content"`
	Img       string `gorm:"type:varchar(100)" json:"img"`
	Type      int    `gorm:"type:int;default:1" json:"type"` // 1: Markdown, 2: PDF
	PdfUrl    string `gorm:"type:varchar(200)" json:"pdf_url"`
	Top       int    `gorm:"type:int;default:0" json:"top"`   // 0表示不置顶，其他数字1-6表示置顶等级，数字越小等级越高
	Views     int    `gorm:"type:int;default:0" json:"views"` // 阅读量
	Tags      string `gorm:"type:varchar(200)" json:"tags"`
	TagModels []Tag  `gorm:"many2many:article_tags" json:"tag_models"`
}

// CreateArt 新增文章
// 参数: data - 文章信息
// 返回: 状态码
func CreateArt(data *Article) int {
	// 1. 处理标签逻辑：将 Tags 字符串解析为 Tag 模型
	if data.Tags != "" {
		tagNames := strings.Split(data.Tags, ",") // 假设前端传逗号分隔 e.g. "go,web"
		var tags []Tag
		for _, name := range tagNames {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			var tag Tag
			// FirstOrCreate: 如果不存在则创建，存在则查找
			db.FirstOrCreate(&tag, Tag{Name: name})
			tags = append(tags, tag)
		}
		data.TagModels = tags
	}

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
		query = query.Where("LOWER(title) LIKE ? OR LOWER(desc) LIKE ? OR LOWER(tags) LIKE ?", searchTerm, searchTerm, searchTerm)
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
	} else {
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

// IncrementArtViews 增加文章阅读量
func IncrementArtViews(id int) {
	// 使用 UpdateColumn 避免更新 UpdatedAt 字段，使用 gorm.Expr 实现原子增加
	db.Model(&Article{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + ?", 1))
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

// GetHotArticles 查询热门文章
func GetHotArticles(limit int) ([]Article, int) {
	var articleList []Article
	err := db.Preload("Category").Order("views DESC").Limit(limit).Find(&articleList).Error
	if err != nil {
		return nil, errmsg.ERROR
	}
	return articleList, errmsg.SUCCESS
}

// GetRelatedArticles 查询相关文章 (基于标签)
func GetRelatedArticles(id int, tags string) ([]Article, int) {
	var articleList []Article
	// if tags == "" {
	// 	// 如果没有标签，返回空的但成功的响应，或者随机推荐？这里暂且返回空
	// 	return []Article{}, errmsg.SUCCESS
	// }

	// 分割标签
	tagList := strings.Split(tags, ",") // 假设前端传逗号分隔 e.g. "go,web"

	tx := db.Preload("Category").Where("id != ?", id)

	// 构建 OR 查询条件
	var conditions []string
	var args []interface{}
	for _, tag := range tagList {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			conditions = append(conditions, "tags LIKE ?")
			args = append(args, "%"+tag+"%")
		}
	}

	if len(conditions) > 0 {
		tx = tx.Where(strings.Join(conditions, " OR "), args...)
	} else {
		// 如果没有有效标签，直接返回空
		return []Article{}, errmsg.SUCCESS
	}

	// 按最新发布排序，取前 5 条
	err := tx.Order("created_at DESC").Limit(5).Find(&articleList).Error
	if err != nil {
		return nil, errmsg.ERROR
	}
	return articleList, errmsg.SUCCESS
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

type Archive struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// GetArchive 获取文章归档
func GetArchive() ([]Archive, int) {
	var archives []Archive
	// MySQL 使用 DATE_FORMAT(created_at, '%Y-%m')
	err := db.Model(&Article{}).
		Select("DATE_FORMAT(created_at, '%Y-%m') as date, count(*) as count").
		Group("date").
		Order("date desc").
		Scan(&archives).Error

	if err != nil {
		return nil, errmsg.ERROR
	}
	return archives, errmsg.SUCCESS
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
	maps["tags"] = data.Tags
	// 补全缺失的更新字段
	maps["type"] = data.Type
	maps["pdf_url"] = data.PdfUrl

	// 处理标签更新
	var newTags []Tag
	if data.Tags != "" {
		tagNames := strings.Split(data.Tags, ",")
		for _, name := range tagNames {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			var tag Tag
			db.FirstOrCreate(&tag, Tag{Name: name})
			newTags = append(newTags, tag)
		}
	}

	err = db.Model(&art).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}

	// 更新关联
	art.ID = uint(id)
	db.Model(&art).Association("TagModels").Replace(newTags)

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

// GetSitemapData 获取站点地图所需数据
func GetSitemapData() ([]Article, int) {
	var articles []Article
	// 只查询 ID 和 UpdatedAt，减少数据量
	err := db.Select("id", "updated_at").Find(&articles).Error
	if err != nil {
		return nil, errmsg.ERROR
	}
	return articles, errmsg.SUCCESS
}
