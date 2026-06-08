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

	// 首次运行时创建演示文章
	var articleCount int64
	db.Model(&Article{}).Count(&articleCount)
	if articleCount == 0 {
		createDemoArticle()
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

// createDemoArticle 首次运行时创建一篇演示文章
func createDemoArticle() {
	cate := Category{Name: "使用指南"}
	db.FirstOrCreate(&cate, Category{Name: "使用指南"})

	tagNames := []string{"Markdown", "KaTeX", "Mermaid", "Demo"}
	var tags []Tag
	for _, name := range tagNames {
		var tag Tag
		db.FirstOrCreate(&tag, Tag{Name: name})
		tags = append(tags, tag)
	}

	demo := Article{
		Title:     "YanBlog 功能演示",
		Cid:       int(cate.ID),
		Desc:      "演示所有 Markdown 特性：代码高亮、数学公式、流程图、表格、链接卡片等。",
		Content:   getDemoContent(),
		Img:       "/uploads/defaults/hero.jpg",
		Tags:      "Markdown,KaTeX,Mermaid,Demo",
		TagModels: tags,
	}

	if err := db.Create(&demo).Error; err != nil {
		fmt.Println("创建演示文章失败:", err)
	} else {
		fmt.Println("已创建演示文章: YanBlog 功能演示")
	}
}

func getDemoContent() string {
	return "## 文本样式\n\n支持 **粗体**、*斜体*、`行内代码`、~~删除线~~。\n\n## 列表\n\n- 第一项\n- 第二项\n  - 嵌套子项\n- 第三项\n\n1. 步骤一\n2. 步骤二\n3. 步骤三\n\n## 引用\n\n> 这是一段引用文字。\n> 可以有多行。\n\n## 链接卡片\n\n单独一行的链接会自动渲染为卡片样式：\n\nhttps://github.com\n\n## 表格\n\n| 功能 | 描述 | 状态 |\n|------|------|------|\n| 代码块 | Mac 风格 + 高亮 + 行号 | 已支持 |\n| 公式 | KaTeX 行内和块级 | 已支持 |\n| 图表 | Mermaid 流程图 | 已支持 |\n\n## 代码块\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, YanBlog!\")\n}\n```\n\n```javascript\nfor (let i = 1; i <= 100; i++) {\n  if (i % 15 === 0) console.log(\"FizzBuzz\")\n  else if (i % 3 === 0) console.log(\"Fizz\")\n  else if (i % 5 === 0) console.log(\"Buzz\")\n  else console.log(i)\n}\n```\n\n## 数学公式（KaTeX）\n\n行内：$E = mc^2$，$a^2 + b^2 = c^2$\n\n块级：\n\n$$\n\\int_{a}^{b} f(x) \\,dx = F(b) - F(a)\n$$\n\n$$\n\\sum_{n=1}^{\\infty} \\frac{1}{n^2} = \\frac{\\pi^2}{6}\n$$\n\n## 流程图（Mermaid）\n\n```mermaid\ngraph TD\n    A[开始] --> B{是否登录?}\n    B -->|是| C[进入后台]\n    B -->|否| D[跳转登录页]\n    C --> E[管理文章]\n    D --> H[输入账号密码]\n    H --> B\n```\n\n```mermaid\nsequenceDiagram\n    U->>F: 点击发布\n    F->>B: POST /api/v1/article/add\n    B->>D: INSERT\n    D-->>B: OK\n    B-->>F: {status:200}\n    F-->>U: 发布成功\n```\n\n---\n\n以上就是 YanBlog 支持的所有 Markdown 语法特性。"
}