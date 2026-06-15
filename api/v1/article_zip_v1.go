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

// UploadArticleZip 处理单个ZIP上传
func UploadArticleZip(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "文件上传失败",
		})
		return
	}

	tempDir := "./temp_zip"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		_ = os.MkdirAll(tempDir, 0755)
	}

	tempZipPath := filepath.Join(tempDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename))
	if err := c.SaveUploadedFile(file, tempZipPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "保存文件失败",
		})
		return
	}
	defer os.Remove(tempZipPath)

	unzipDir := strings.TrimSuffix(tempZipPath, ".zip") + "_extracted"
	if err := unzip(tempZipPath, unzipDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "解压失败: " + err.Error(),
		})
		return
	}
	defer os.RemoveAll(unzipDir)

	article, code := processZipArticle(unzipDir)
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

// parseDate 尝试多种日期格式解析
func parseDate(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05+08:00",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02",
		"2006-1-2 15:04:05",
		"2006-1-2",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析日期: %s", s)
}

// UploadArticleZipBatch 批量上传ZIP文章
func UploadArticleZipBatch(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "请上传文件",
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "未找到文件，请使用 files 字段上传",
		})
		return
	}

	type Result struct {
		FileName string `json:"file_name"`
		Title    string `json:"title"`
		Status   int    `json:"status"`
		Message  string `json:"message"`
	}

	results := make([]Result, 0, len(files))
	successCount := 0

	for _, fileHeader := range files {
		result := Result{FileName: fileHeader.Filename}

		// 保存临时文件
		tempDir := "./temp_zip"
		if _, err := os.Stat(tempDir); os.IsNotExist(err) {
			_ = os.MkdirAll(tempDir, 0755)
		}

		tempZipPath := filepath.Join(tempDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename))
		if err := c.SaveUploadedFile(fileHeader, tempZipPath); err != nil {
			result.Status = errmsg.ERROR
			result.Message = "保存文件失败"
			results = append(results, result)
			continue
		}

		// 解压
		unzipDir := strings.TrimSuffix(tempZipPath, ".zip") + "_extracted"
		if err := unzip(tempZipPath, unzipDir); err != nil {
			result.Status = errmsg.ERROR
			result.Message = "解压失败: " + err.Error()
			results = append(results, result)
			os.Remove(tempZipPath)
			os.RemoveAll(unzipDir)
			continue
		}

		// 处理
		article, code := processZipArticle(unzipDir)
		result.Status = code
		result.Message = errmsg.GetErrMsg(code)
		if article != nil {
			result.Title = article.Title
			if code == errmsg.SUCCESS {
				successCount++
			}
		}

		results = append(results, result)

		// 清理
		os.Remove(tempZipPath)
		os.RemoveAll(unzipDir)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"total":   len(files),
		"success": successCount,
		"results": results,
	})
}

// processZipArticle 处理单个解压后的ZIP目录，返回文章和状态码
func processZipArticle(unzipDir string) (*model.Article, int) {
	// 查找 Markdown 文件
	var mdPath string
	filepath.Walk(unzipDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			mdPath = path
			return io.EOF
		}
		return nil
	})

	if mdPath == "" {
		return nil, errmsg.ERROR
	}

	// 解析 Front Matter
	contentBytes, err := os.ReadFile(mdPath)
	if err != nil {
		return nil, errmsg.ERROR
	}
	contentStr := string(contentBytes)

	var frontMatter ArticleFrontMatter
	var bodyContent string

	if strings.HasPrefix(contentStr, "---") {
		parts := strings.SplitN(contentStr, "---", 3)
		if len(parts) >= 3 {
			if err := yaml.Unmarshal([]byte(parts[1]), &frontMatter); err != nil {
				bodyContent = contentStr
			} else {
				bodyContent = parts[2]
			}
		} else {
			bodyContent = contentStr
		}
	} else {
		bodyContent = contentStr
	}

	if frontMatter.Title == "" {
		base := filepath.Base(mdPath)
		frontMatter.Title = strings.TrimSuffix(base, filepath.Ext(base))
	}

	// 处理图片
	imgRegex := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	matches := imgRegex.FindAllStringSubmatch(bodyContent, -1)
	processedContent := bodyContent

	for _, match := range matches {
		originalPath := match[2]
		if strings.HasPrefix(originalPath, "http") || strings.HasPrefix(originalPath, "//") {
			continue
		}
		mdDir := filepath.Dir(mdPath)
		imageFullPath := filepath.Join(mdDir, originalPath)
		if _, err := os.Stat(imageFullPath); err == nil {
			newURL, err := uploadLocalFile(imageFullPath, "article")
			if err == nil {
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

	// 处理分类
	var cid int
	if frontMatter.Category != "" {
		cid = model.GetOrCreateCategory(frontMatter.Category)
	} else {
		cid = 1
	}

	article := &model.Article{
		Title:   frontMatter.Title,
		Cid:     cid,
		Desc:    frontMatter.Desc,
		Content: processedContent,
		Img:     frontMatter.Cover,
		Tags:    strings.Join(frontMatter.Tags, ","),
	}

	// 保留原始创建时间
	if frontMatter.Date != "" {
		if t, err := parseDate(frontMatter.Date); err == nil {
			article.CreatedAt = t
		}
	}

	code := model.CreateArt(article)
	return article, code
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
			os.MkdirAll(fpath, 0755)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
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
		_ = os.MkdirAll(targetDir, 0755)
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
