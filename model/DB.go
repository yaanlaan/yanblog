package model

import (
	"fmt"
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

	db.AutoMigrate(&User{}, &Category{}, &Article{})

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
