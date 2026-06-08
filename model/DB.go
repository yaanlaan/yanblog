package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"yanblog/utils"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB
var err error

func GetDB() *gorm.DB {
	return db
}

func InitDB() {
	var dbErr error

	dbType := strings.ToUpper(utils.ServerConfig.Database.Db)
	
	if dbType == "SQLITE" {
		db, dbErr = initSQLite()
	} else {
		db, dbErr = initMySQL()
	}

	if dbErr != nil {
		panic(fmt.Sprintf("无法连接到数据库！Error: %s", dbErr.Error()))
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	db.AutoMigrate(&User{}, &Category{}, &Article{}, &Tag{})
	migrateTags()

	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		admin := User{
			Username: "admin",
			Password: "123456",
			Role:     1,
		}
		if err := db.Create(&admin).Error; err != nil {
			fmt.Println("创建默认超级管理员失败:", err)
		} else {
			fmt.Println("===========================================")
			fmt.Println("⚠️  安全警告")
			fmt.Println("===========================================")
			fmt.Println("已创建默认超级管理员账号:")
			fmt.Println("  用户名: admin")
			fmt.Println("  密码:   123456")
			fmt.Println("")
			fmt.Println("🔴 请立即登录后台并修改默认密码！")
			fmt.Println("===========================================")
		}
	}
}

func initMySQL() (*gorm.DB, error) {
	var dbErr error
	var dsn string

	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			utils.ServerConfig.Database.DbUser,
			utils.ServerConfig.Database.DbPassWord,
			utils.ServerConfig.Database.DbHost,
			utils.ServerConfig.Database.DbPort,
			utils.ServerConfig.Database.DbName,
		)
		db, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,
			SkipDefaultTransaction: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})

		if dbErr == nil {
			sqlDB, err := db.DB()
			if err == nil {
				if err = sqlDB.Ping(); err == nil {
					fmt.Println("MySQL数据库连接成功！")
					return db, nil
				} else {
					dbErr = err
				}
			} else {
				dbErr = err
			}
		}

		fmt.Printf("等待MySQL数据库启动... (%d/%d) Error: %s\n", i+1, maxRetries, dbErr)
		time.Sleep(2 * time.Second)
	}

	return nil, dbErr
}

func initSQLite() (*gorm.DB, error) {
	dbPath := utils.ServerConfig.Database.DbName
	if dbPath == "" {
		dbPath = "yanblog.db"
	}

	// 确保数据库文件所在目录存在
	if dir := filepath.Dir(dbPath); dir != "." {
		os.MkdirAll(dir, 0755)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("SQLite数据库连接成功！")
	return db, nil
}

func migrateTags() {
	var count int64
	err := db.Model(&Tag{}).Count(&count).Error
	if err != nil {
		fmt.Println("Check tag count failed:", err)
		return
	}
	if count > 0 {
		return
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
			err := db.Model(&art).Association("TagModels").Replace(newTags)
			if err != nil {
				fmt.Printf("Update article %d tags failed: %s\n", art.ID, err)
			}
		}
	}
	fmt.Println("Tag migration completed.")
}