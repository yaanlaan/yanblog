package model

import (
	"fmt"
	"strings"
	"time"
	"yanblog/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB
var err error

func GetDB() *gorm.DB {
	return db
}

// InitDB 初始化数据库连接
// 参数: 无
// 返回: 无
func InitDB() {
	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.ServerConfig.Database.DbUser,
		utils.ServerConfig.Database.DbPassWord,
		utils.ServerConfig.Database.DbHost,
		utils.ServerConfig.Database.DbPort,
		utils.ServerConfig.Database.DbName,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// gorm日志模式：silent
		Logger: logger.Default.LogMode(logger.Silent),
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})
	if err != nil {
		panic("无法连接到数据库，请检查配置！: " + err.Error())
	}

	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, _ := db.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	db.AutoMigrate(&User{}, &Category{}, &Article{}, &Tag{})

	// 检查并迁移旧的标签数据 (从 articles.tags 字符串迁移到 tag 表)
	migrateTags()

	// 检查是否存在用户，如果不存在则创建默认超级管理员
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		admin := User{
			Username: "admin",
			Password: "123456",
			Role:     1, // 超级管理员
		}
		// 这里会触发 BeforeSave 钩子，自动加密密码
		if err := db.Create(&admin).Error; err != nil {
			fmt.Println("创建默认超级管理员失败:", err)
		} else {
			fmt.Println("已创建默认超级管理员账号: admin, 密码: 123456")
		}
	}
}

// migrateTags 迁移旧标签数据
func migrateTags() {
	var count int64
	// table might not exist if logic runs too early? No, AutoMigrate runs before this.
	err := db.Model(&Tag{}).Count(&count).Error
	if err != nil {
		fmt.Println("Check tag count failed:", err)
		return
	}
	if count > 0 {
		return // 已有标签数据，不执行迁移
	}

	fmt.Println("Detected empty tags table, starting migration from articles...")

	var articles []Article
	db.Find(&articles)

	for _, art := range articles {
		if art.Tags == "" {
			continue
		}
		tagsStr := strings.ReplaceAll(art.Tags, "，", ",")
		tagNames := strings.Split(tagsStr, ",")

		var newTags []Tag
		for _, name := range tagNames {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			var tag Tag
			// 查找或创建标签
			if err := db.Where("name = ?", name).First(&tag).Error; err != nil {
				tag = Tag{Name: name}
				if err := db.Create(&tag).Error; err != nil {
					fmt.Printf("Create tag %s failed: %s\n", name, err)
					continue
				}
			}
			newTags = append(newTags, tag)
		}

		if len(newTags) > 0 {
			// 更新关联
			// 需要先设置主键
			// art 是 gorm.Find 出来的，主键应该有
			err := db.Model(&art).Association("TagModels").Replace(newTags)
			if err != nil {
				fmt.Printf("Update article %d tags failed: %s\n", art.ID, err)
			}
		}
	}
	fmt.Println("Tag migration completed.")
}
