package model

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"yanblog/utils/errmsg"
)

// UpLoadFile 上传文件到本地
// 参数: file - 要上传的文件, fileHeader - 文件头信息, uploadType - 上传类型(article/category), key - 关键标识(文章标题/分类名)
// 返回: 文件访问URL和状态码
func UpLoadFile(file multipart.File, fileHeader *multipart.FileHeader, uploadType string, key string) (string, int) {
	// 1. 确定存储目录
	baseDir := "./uploads"
	targetDir := baseDir

	// 简单的文件名清理函数，防止路径遍历和非法字符
	// cleanKey := filepath.Clean(key) // 废弃: 不再依赖用户输入的 key 作为目录名，避免改名带来的路径失效问题

	switch uploadType {
	case "avatar":
		// 用户头像
		targetDir = filepath.Join(baseDir, "avatar")
	case "category":
		// 分类封面
		targetDir = filepath.Join(baseDir, "category")
	case "article", "markdown":
		// 文章内容图片 (Markdown中插入的)
		// 使用 年/月 格式，避免单文件夹文件过多
		now := time.Now()
		targetDir = filepath.Join(baseDir, "article", "content", fmt.Sprintf("%d%02d", now.Year(), now.Month()))
	case "cover":
		// 文章封面
		targetDir = filepath.Join(baseDir, "article", "cover")
	case "pdf":
		// PDF文件
		targetDir = filepath.Join(baseDir, "article", "pdf")
	case "system":
		// 系统配置图片 (Logo, 背景等)
		targetDir = filepath.Join(baseDir, "system")
	default:
		// 其他/通用
		targetDir = filepath.Join(baseDir, "common")
	}

	// 确保存储目录存在
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		_ = os.MkdirAll(targetDir, os.ModePerm)
	}

	// 2. 生成文件名
	ext := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(targetDir, fileName)

	// 3. 创建目标文件
	out, err := os.Create(filePath)
	if err != nil {
		return "", errmsg.ERROR
	}
	defer out.Close()

	// 4. 写入文件内容
	_, err = io.Copy(out, file)
	if err != nil {
		return "", errmsg.ERROR
	}

	// 5. 返回访问URL
	// 需要将文件路径转换为URL路径 (将反斜杠转换为斜杠)
	relPath, _ := filepath.Rel(".", filePath)
	url := "/" + filepath.ToSlash(relPath)

	return url, errmsg.SUCCESS
}
