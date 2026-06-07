package v1

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"yanblog/model"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// ArticleFrontMatter 定义文章头部的元数据结构
type ArticleFrontMatter struct {
	Title    string   `yaml:"title"`
	Date     string   `yaml:"date"`
	Tags     []string `yaml:"tags"`
	Category string   `yaml:"category"` // string name
	Desc     string   `yaml:"desc"`
	Cover    string   `yaml:"cover"`
}

// UploadArticleZip 处理ZIP上传
func UploadArticleZip(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "文件上传失败",
		})
		return
	}

	// 1. 保存上传的ZIP文件到临时目录
	tempDir := "./temp_zip"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		_ = os.MkdirAll(tempDir, os.ModePerm)
	}

	tempZipPath := filepath.Join(tempDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename))
	if err := c.SaveUploadedFile(file, tempZipPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "保存文件失败",
		})
		return
	}
	defer os.Remove(tempZipPath) // 处理完后删除ZIP包

	// 2. 解压ZIP
	unzipDir := strings.TrimSuffix(tempZipPath, ".zip") + "_extracted"
	if err := unzip(tempZipPath, unzipDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "解压失败: " + err.Error(),
		})
		return
	}
	defer os.RemoveAll(unzipDir) // 确保清理解压目录

	// 3. 寻找 Markdown 文件
	var mdPath string
	err = filepath.Walk(unzipDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			mdPath = path
			return io.EOF // Found, stop walking
		}
		return nil
	})

	if mdPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "压缩包内未找到 Markdown 文件",
		})
		return
	}

	// 4. 解析 Front Matter 和 Content
	contentBytes, err := os.ReadFile(mdPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "读取 Markdown 文件失败",
		})
		return
	}
	contentStr := string(contentBytes)

	var frontMatter ArticleFrontMatter
	var bodyContent string

	// 检查是否有 YAML Front Matter
	if strings.HasPrefix(contentStr, "---") {
		parts := strings.SplitN(contentStr, "---", 3)
		if len(parts) >= 3 {
			// 解析元数据
			if err := yaml.Unmarshal([]byte(parts[1]), &frontMatter); err != nil {
				// 解析失败，忽略元数据，当作普通内容
				bodyContent = contentStr
			} else {
				// 解析成功，取正文
				bodyContent = parts[2]
			}
		} else {
			bodyContent = contentStr
		}
	} else {
		bodyContent = contentStr
	}

	// 默认值处理
	if frontMatter.Title == "" {
		// 使用文件名作为标题
		base := filepath.Base(mdPath)
		frontMatter.Title = strings.TrimSuffix(base, filepath.Ext(base))
	}

	// 5. 处理并上传图片
	// 正则匹配 ![alt](src) 和 <img src="src">
	// 简单起见，主要匹配 Markdown 图片语法
	imgRegex := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	matches := imgRegex.FindAllStringSubmatch(bodyContent, -1)

	processedContent := bodyContent

	for _, match := range matches {
		originalPath := match[2]
		// 忽略网络图片
		if strings.HasPrefix(originalPath, "http") || strings.HasPrefix(originalPath, "//") {
			continue
		}

		// 拼接绝对路径 (相对于md文件所在目录)
		mdDir := filepath.Dir(mdPath)
		imageFullPath := filepath.Join(mdDir, originalPath)

		// 检查图片是否存在
		if _, err := os.Stat(imageFullPath); err == nil {
			// 上传图片
			newURL, err := uploadLocalFile(imageFullPath, "article")
			if err == nil {
				// 替换内容中的路径
				processedContent = strings.Replace(processedContent, originalPath, newURL, -1)
			}
		}
	}

	// 处理封面图
	if frontMatter.Cover != "" && !strings.HasPrefix(frontMatter.Cover, "http") {
		mdDir := filepath.Dir(mdPath)
		coverFullPath := filepath.Join(mdDir, frontMatter.Cover)
		if _, err := os.Stat(coverFullPath); err == nil {
			newCoverURL, err := uploadLocalFile(coverFullPath, "cover")
			if err == nil {
				frontMatter.Cover = newCoverURL
			}
		}
	}

	// 6. 保存文章到数据库
	// 处理分类
	var cid int
	if frontMatter.Category != "" {
		// 检查分类是否存在
		// 我们假设 model 包里可以直接查，或者我们需要加个方法 GetCategoryByName
		// 暂时用 CheckCategory 逻辑变通一下
		// 这里需要直接操作 DB 查找 ID
		// 由于 model 包不可见 db 变量(它是私有的 in model)，我们需要 model 提供方法
		// 我们先假设 category 总是要存在的，如果不存在就创建
		cid = model.GetOrCreateCategory(frontMatter.Category)
	} else {
		// 默认分类?
		cid = 1 // 假设1是默认
	}

	article := model.Article{
		Title:   frontMatter.Title,
		Cid:     cid,
		Desc:    frontMatter.Desc,
		Content: processedContent,
		Img:     frontMatter.Cover,
		Tags:    strings.Join(frontMatter.Tags, ","), // 逗号分隔
	}

	// 创建文章
	code := model.CreateArt(&article)

	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"data":    article,
		"message": "上传成功",
	})
}

// 辅助函数：解压
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// 防止 Zip Slip 漏洞
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// 辅助函数：上传本地文件
func uploadLocalFile(path string, uploadType string) (string, error) {
	// 复用 model.UploadFile 逻辑比较麻烦，因为它需要 multipart
	// 这里重新实现一个简单的本地复制逻辑，或者修改 model 以支持
	// 为了不破坏现有结构，这里简单实现复制到 ./uploads/article/...

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 确定目标目录
	baseDir := "./uploads"
	targetDir := filepath.Join(baseDir, "article", "content", time.Now().Format("200601"))
	if uploadType == "cover" {
		targetDir = filepath.Join(baseDir, "article", "cover")
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		_ = os.MkdirAll(targetDir, os.ModePerm)
	}

	// 生成文件名
	ext := filepath.Ext(path)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	newFilePath := filepath.Join(targetDir, newFileName)

	out, err := os.Create(newFilePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 复制
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	// 返回 URL
	// 这里假设静态资源映射是 /uploads -> ./uploads
	// 或者是 /static/uploads -> ./uploads
	// 需要确认路由配置。通常是 /public/ or /uploads/
	// 检查 routers.go

	relPath, _ := filepath.Rel(".", newFilePath)
	url := "/" + filepath.ToSlash(relPath) // e.g. /uploads/article/content/202310/123.png

	return url, nil
}
