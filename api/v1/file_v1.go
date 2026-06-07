package v1

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// safeUploadPath 验证并返回安全的文件路径，防止路径遍历攻击
func safeUploadPath(userPath string) (string, bool) {
	cleaned := filepath.Clean(filepath.Join("uploads", userPath))
	// 确保路径在 uploads 目录内
	absBase, _ := filepath.Abs("uploads")
	absTarget, _ := filepath.Abs(cleaned)
	if !strings.HasPrefix(absTarget, absBase) {
		return "", false
	}
	return cleaned, true
}

// FileInfo 文件信息结构体
type FileInfo struct {
	Name      string    `json:"name"`
	IsDir     bool      `json:"isDir"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	Ext       string    `json:"ext"`
	ModTime   time.Time `json:"modTime"`
	IsImage   bool      `json:"isImage"`
	Thumbnail string    `json:"thumbnail,omitempty"`
}

// GetFileList 获取文件列表
func GetFileList(c *gin.Context) {
	reqPath := c.Query("path")
	if reqPath == "" {
		reqPath = "uploads"
	} else {
		var ok bool
		reqPath, ok = safeUploadPath(reqPath)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  errmsg.ERROR,
				"message": "非法路径",
				"data":    []FileInfo{},
			})
			return
		}
	}

	if _, err := os.Stat(reqPath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "目录不存在",
			"data":    []FileInfo{},
		})
		return
	}

	files, err := ioutil.ReadDir(reqPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "读取目录失败",
			"data":    []FileInfo{},
		})
		return
	}

	fileList := make([]FileInfo, 0)
	for _, f := range files {
		fullPath := filepath.Join(reqPath, f.Name())
		relPath, err := filepath.Rel("uploads", fullPath)
		if err != nil {
			continue
		}
		relPath = filepath.ToSlash(relPath)

		ext := filepath.Ext(f.Name())
		isImage := isImageFile(ext)

		fileInfo := FileInfo{
			Name:    f.Name(),
			IsDir:   f.IsDir(),
			Path:    relPath,
			Size:    f.Size(),
			Ext:     ext,
			ModTime: f.ModTime(),
			IsImage: isImage,
		}

		if isImage && !f.IsDir() {
			fileInfo.Thumbnail = "/uploads/" + relPath
		}

		fileList = append(fileList, fileInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": errmsg.GetErrMsg(errmsg.SUCCESS),
		"data":    fileList,
	})
}

// isImageFile 判断是否为图片文件
func isImageFile(ext string) bool {
	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
		".svg":  true,
	}
	return imageExts[ext]
}

// DeleteFile 删除文件或目录
func DeleteFile(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "路径不能为空",
		})
		return
	}

	targetPath := filepath.Join("uploads", path)
	targetPath = filepath.Clean(targetPath)

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "文件或目录不存在",
		})
		return
	}

	err := os.RemoveAll(targetPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "删除成功",
	})
}

// CreateDir 创建目录
func CreateDir(c *gin.Context) {
	var data struct {
		Path string `json:"path"`
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	data.Path = filepath.Clean(data.Path)
	if data.Path == "." {
		data.Path = ""
	}

	targetPath := filepath.Join("uploads", data.Path, data.Name)
	targetPath = filepath.Clean(targetPath)

	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "目录已存在",
		})
		return
	}

	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "创建目录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "创建成功",
	})
}

// RenameFile 重命名文件或目录
func RenameFile(c *gin.Context) {
	var data struct {
		Path    string `json:"path"`
		NewName string `json:"newName"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	if data.Path == "" || data.NewName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数不能为空",
		})
		return
	}

	oldPath := filepath.Join("uploads", data.Path)
	dir := filepath.Dir(oldPath)
	newPath := filepath.Join(dir, data.NewName)

	oldPath = filepath.Clean(oldPath)
	newPath = filepath.Clean(newPath)

	if err := os.Rename(oldPath, newPath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "重命名失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "重命名成功",
	})
}

// MoveFile 移动文件或目录
func MoveFile(c *gin.Context) {
	var data struct {
		SourcePath string `json:"sourcePath"`
		TargetPath string `json:"targetPath"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	if data.SourcePath == "" || data.TargetPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数不能为空",
		})
		return
	}

	source := filepath.Join("uploads", data.SourcePath)
	targetDir := filepath.Join("uploads", data.TargetPath)

	source = filepath.Clean(source)
	targetDir = filepath.Clean(targetDir)

	if _, err := os.Stat(source); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "源文件不存在",
		})
		return
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR,
				"message": "创建目标目录失败: " + err.Error(),
			})
			return
		}
	}

	fileName := filepath.Base(source)
	target := filepath.Join(targetDir, fileName)

	if _, err := os.Stat(target); !os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "目标文件已存在",
		})
		return
	}

	if err := os.Rename(source, target); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "移动失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "移动成功",
	})
}

// CopyFile 复制文件
func CopyFile(c *gin.Context) {
	var data struct {
		SourcePath string `json:"sourcePath"`
		TargetPath string `json:"targetPath"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	if data.SourcePath == "" || data.TargetPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数不能为空",
		})
		return
	}

	source := filepath.Join("uploads", data.SourcePath)
	targetDir := filepath.Join("uploads", data.TargetPath)

	source = filepath.Clean(source)
	targetDir = filepath.Clean(targetDir)

	sourceInfo, err := os.Stat(source)
	if os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "源文件不存在",
		})
		return
	}

	if sourceInfo.IsDir() {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "暂不支持复制目录",
		})
		return
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR,
				"message": "创建目标目录失败: " + err.Error(),
			})
			return
		}
	}

	fileName := filepath.Base(source)
	target := filepath.Join(targetDir, fileName)

	if _, err := os.Stat(target); !os.IsNotExist(err) {
		ext := filepath.Ext(fileName)
		name := fileName[:len(fileName)-len(ext)]
		target = filepath.Join(targetDir, name+"_copy"+ext)
	}

	srcFile, err := os.Open(source)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "打开源文件失败: " + err.Error(),
		})
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(target)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "创建目标文件失败: " + err.Error(),
		})
		return
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "复制失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "复制成功",
	})
}

// BatchDeleteFiles 批量删除文件
func BatchDeleteFiles(c *gin.Context) {
	var data struct {
		Paths []string `json:"paths"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	if len(data.Paths) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "请选择要删除的文件",
		})
		return
	}

	successCount := 0
	failCount := 0
	var failMessages []string

	for _, path := range data.Paths {
		targetPath := filepath.Join("uploads", path)
		targetPath = filepath.Clean(targetPath)

		if err := os.RemoveAll(targetPath); err != nil {
			failCount++
			failMessages = append(failMessages, path+": "+err.Error())
		} else {
			successCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       errmsg.SUCCESS,
		"message":      "批量删除完成",
		"successCount": successCount,
		"failCount":    failCount,
		"failDetails":  failMessages,
	})
}

// BatchUploadFiles 批量上传文件
func BatchUploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "解析表单失败: " + err.Error(),
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "请选择要上传的文件",
		})
		return
	}

	targetDir := c.PostForm("dir")
	if targetDir == "" {
		targetDir = ""
	}

	targetPath := filepath.Join("uploads", targetDir)
	targetPath = filepath.Clean(targetPath)

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR,
				"message": "创建目录失败: " + err.Error(),
			})
			return
		}
	}

	results := make([]map[string]interface{}, 0)
	successCount := 0
	failCount := 0

	for _, file := range files {
		fileName := file.Filename
		ext := filepath.Ext(fileName)
		newFileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
		dstPath := filepath.Join(targetPath, newFileName)

		if err := c.SaveUploadedFile(file, dstPath); err != nil {
			failCount++
			results = append(results, map[string]interface{}{
				"name":    fileName,
				"success": false,
				"error":   err.Error(),
			})
		} else {
			successCount++
			relPath, _ := filepath.Rel("uploads", dstPath)
			relPath = filepath.ToSlash(relPath)
			results = append(results, map[string]interface{}{
				"name":    fileName,
				"success": true,
				"url":     "/uploads/" + relPath,
				"path":    relPath,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       errmsg.SUCCESS,
		"message":      "批量上传完成",
		"successCount": successCount,
		"failCount":    failCount,
		"results":      results,
	})
}

// GetStorageStats 获取存储统计信息
func GetStorageStats(c *gin.Context) {
	rootDir := "uploads"
	var totalFiles int
	var totalSize int64
	var totalDirs int

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			totalDirs++
		} else {
			totalFiles++
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "获取统计信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     errmsg.SUCCESS,
		"message":    errmsg.GetErrMsg(errmsg.SUCCESS),
		"totalFiles": totalFiles,
		"totalDirs":  totalDirs,
		"totalSize":  totalSize,
	})
}