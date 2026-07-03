package v1

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// FileStats 文件统计信息
type FileStats struct {
	TotalFiles    int64   `json:"total_files"`
	TotalDirs     int64   `json:"total_dirs"`
	TotalSize     int64   `json:"total_size"`
	TotalSizeMB   float64 `json:"total_size_mb"`
	ImageCount    int64   `json:"image_count"`
	DocumentCount int64   `json:"document_count"`
	ArchiveCount  int64   `json:"archive_count"`
	OtherCount    int64   `json:"other_count"`
	LargestFile   string  `json:"largest_file"`
	LargestSize   int64   `json:"largest_size"`
}

// FileSearchRequest 文件搜索请求
type FileSearchRequest struct {
	Keyword  string `json:"keyword"`   // 搜索关键词
	Path     string `json:"path"`      // 搜索路径
	Ext      string `json:"ext"`       // 文件扩展名过滤
	MinSize  int64  `json:"min_size"`  // 最小大小（字节）
	MaxSize  int64  `json:"max_size"`  // 最大大小（字节）
	SortBy   string `json:"sort_by"`   // 排序字段：name, size, time
	SortDesc bool   `json:"sort_desc"` // 是否降序
	Page     int    `json:"page"`      // 页码
	PageSize int    `json:"page_size"` // 每页数量
}

// RecycleBinItem 回收站项目
type RecycleBinItem struct {
	OriginalPath string    `json:"original_path"`
	RecyclePath  string    `json:"recycle_path"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	DeletedAt    time.Time `json:"deleted_at"`
	IsDir        bool      `json:"is_dir"`
}

// GetFileStats 获取文件统计信息
func GetFileStats(c *gin.Context) {
	basePath := "uploads"
	
	stats := &FileStats{}
	var largestSize int64
	var largestFile string
	
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if info.IsDir() {
			stats.TotalDirs++
			return nil
		}
		
		stats.TotalFiles++
		stats.TotalSize += info.Size()
		
		// 统计最大文件
		if info.Size() > largestSize {
			largestSize = info.Size()
			largestFile = path
		}
		
		// 按类型统计
		ext := strings.ToLower(filepath.Ext(info.Name()))
		switch {
		case isImageFile(ext):
			stats.ImageCount++
		case isDocumentFile(ext):
			stats.DocumentCount++
		case isArchiveFile(ext):
			stats.ArchiveCount++
		default:
			stats.OtherCount++
		}
		
		return nil
	})
	
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "统计失败",
		})
		return
	}
	
	stats.TotalSizeMB = float64(stats.TotalSize) / 1024 / 1024
	stats.LargestFile = largestFile
	stats.LargestSize = largestSize
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"data":    stats,
	})
}

// SearchFiles 搜索文件
func SearchFiles(c *gin.Context) {
	var req FileSearchRequest
	
	// 解析查询参数
	req.Keyword = c.Query("keyword")
	req.Path = c.Query("path")
	req.Ext = c.Query("ext")
	
	if minSize := c.Query("min_size"); minSize != "" {
		req.MinSize, _ = strconv.ParseInt(minSize, 10, 64)
	}
	if maxSize := c.Query("max_size"); maxSize != "" {
		req.MaxSize, _ = strconv.ParseInt(maxSize, 10, 64)
	}
	
	req.SortBy = c.DefaultQuery("sort_by", "name")
	req.SortDesc = c.DefaultQuery("sort_desc", "false") == "true"
	req.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	req.PageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	if req.Path == "" {
		req.Path = "uploads"
	} else {
		var ok bool
		req.Path, ok = safeUploadPath(req.Path)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  errmsg.ERROR,
				"message": "非法路径",
			})
			return
		}
	}
	
	// 执行搜索
	var results []FileInfo
	filepath.Walk(req.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		
		// 关键词过滤
		if req.Keyword != "" && !strings.Contains(strings.ToLower(info.Name()), strings.ToLower(req.Keyword)) {
			return nil
		}
		
		// 扩展名过滤
		if req.Ext != "" && strings.ToLower(filepath.Ext(info.Name())) != strings.ToLower(req.Ext) {
			return nil
		}
		
		// 大小过滤
		if req.MinSize > 0 && info.Size() < req.MinSize {
			return nil
		}
		if req.MaxSize > 0 && info.Size() > req.MaxSize {
			return nil
		}
		
		relPath, _ := filepath.Rel("uploads", path)
		ext := filepath.Ext(info.Name())
		
		results = append(results, FileInfo{
			Name:    info.Name(),
			IsDir:   false,
			Path:    filepath.ToSlash(relPath),
			Size:    info.Size(),
			Ext:     ext,
			ModTime: info.ModTime(),
			IsImage: isImageFile(ext),
		})
		
		return nil
	})
	
	// 排序
	sort.Slice(results, func(i, j int) bool {
		switch req.SortBy {
		case "size":
			if req.SortDesc {
				return results[i].Size > results[j].Size
			}
			return results[i].Size < results[j].Size
		case "time":
			if req.SortDesc {
				return results[i].ModTime.After(results[j].ModTime)
			}
			return results[i].ModTime.Before(results[j].ModTime)
		default: // name
			if req.SortDesc {
				return results[i].Name > results[j].Name
			}
			return results[i].Name < results[j].Name
		}
	})
	
	// 分页
	total := len(results)
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize
	
	if start >= total {
		results = []FileInfo{}
	} else if end > total {
		results = results[start:]
	} else {
		results = results[start:end]
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"total":   total,
		"page":    req.Page,
		"size":    req.PageSize,
		"data":    results,
	})
}

// CompressFiles 压缩文件/目录
func CompressFiles(c *gin.Context) {
	var req struct {
		Paths    []string `json:"paths" binding:"required"`    // 要压缩的文件/目录
		ZipName  string   `json:"zip_name" binding:"required"` // ZIP文件名
		ZipPath  string   `json:"zip_path"`                    // ZIP文件保存路径
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}
	
	// 确定ZIP路径
	if req.ZipPath == "" {
		req.ZipPath = "uploads"
	}
	
	zipPath, ok := safeUploadPath(filepath.Join(req.ZipPath, req.ZipName))
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  errmsg.ERROR,
			"message": "非法路径",
		})
		return
	}
	
	// 创建ZIP文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "创建ZIP失败",
		})
		return
	}
	defer zipFile.Close()
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	// 添加文件到ZIP
	fileCount := 0
	for _, path := range req.Paths {
		safePath, ok := safeUploadPath(path)
		if !ok {
			continue
		}
		
		err := filepath.Walk(safePath, func(filePath string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			
			relPath, _ := filepath.Rel(filepath.Dir(safePath), filePath)
			
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = filepath.ToSlash(relPath)
			header.Method = zip.Deflate
			
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}
			
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
			
			fileCount++
			return nil
		})
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  errmsg.ERROR,
				"message": fmt.Sprintf("压缩失败: %v", err),
			})
			os.Remove(zipPath)
			return
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": fmt.Sprintf("压缩成功，共 %d 个文件", fileCount),
		"data": gin.H{
			"zip_path":    filepath.ToSlash(strings.TrimPrefix(zipPath, "uploads/")),
			"file_count":  fileCount,
			"zip_size":    getFileSize(zipPath),
		},
	})
}

// ExtractZip 解压ZIP文件
func ExtractZip(c *gin.Context) {
	var req struct {
		ZipPath   string `json:"zip_path" binding:"required"`   // ZIP文件路径
		ExtractTo string `json:"extract_to"`                    // 解压目标目录
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}
	
	safeZipPath, ok := safeUploadPath(req.ZipPath)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  errmsg.ERROR,
			"message": "非法路径",
		})
		return
	}
	
	if req.ExtractTo == "" {
		req.ExtractTo = filepath.Dir(safeZipPath)
	} else {
		var ok bool
		req.ExtractTo, ok = safeUploadPath(req.ExtractTo)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  errmsg.ERROR,
				"message": "非法路径",
			})
			return
		}
	}
	
	// 打开ZIP
	zipReader, err := zip.OpenReader(safeZipPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "打开ZIP失败",
		})
		return
	}
	defer zipReader.Close()
	
	// 解压
	fileCount := 0
	for _, file := range zipReader.File {
		destPath := filepath.Join(req.ExtractTo, file.Name)
		
		// 防止 Zip Slip
		if !strings.HasPrefix(destPath, filepath.Clean(req.ExtractTo)+string(os.PathSeparator)) {
			continue
		}
		
		if file.FileInfo().IsDir() {
			os.MkdirAll(destPath, 0755)
			continue
		}
		
		os.MkdirAll(filepath.Dir(destPath), 0755)
		
		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			continue
		}
		
		srcFile, err := file.Open()
		if err != nil {
			destFile.Close()
			continue
		}
		
		io.Copy(destFile, srcFile)
		destFile.Close()
		srcFile.Close()
		
		fileCount++
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": fmt.Sprintf("解压成功，共 %d 个文件", fileCount),
		"data": gin.H{
			"file_count": fileCount,
			"extract_to": filepath.ToSlash(strings.TrimPrefix(req.ExtractTo, "uploads/")),
		},
	})
}

// MoveToRecycleBin 移动到回收站
func MoveToRecycleBin(c *gin.Context) {
	var req struct {
		Paths []string `json:"paths" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}
	
	recycleDir := "uploads/.recycle"
	os.MkdirAll(recycleDir, 0755)
	
	movedItems := make([]RecycleBinItem, 0)
	
	for _, path := range req.Paths {
		safePath, ok := safeUploadPath(path)
		if !ok {
			continue
		}
		
		info, err := os.Stat(safePath)
		if err != nil {
			continue
		}
		
		// 生成回收站路径
		timestamp := time.Now().Format("20060102_150405")
		name := filepath.Base(safePath)
		recyclePath := filepath.Join(recycleDir, fmt.Sprintf("%s_%s", timestamp, name))
		
		// 移动到回收站
		if err := os.Rename(safePath, recyclePath); err != nil {
			continue
		}
		
		movedItems = append(movedItems, RecycleBinItem{
			OriginalPath: path,
			RecyclePath:  filepath.ToSlash(strings.TrimPrefix(recyclePath, "uploads/")),
			Name:         name,
			Size:         info.Size(),
			DeletedAt:    time.Now(),
			IsDir:        info.IsDir(),
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": fmt.Sprintf("已删除 %d 个项目", len(movedItems)),
		"data":    movedItems,
	})
}

// GetRecycleBin 获取回收站列表
func GetRecycleBin(c *gin.Context) {
	recycleDir := "uploads/.recycle"
	
	if _, err := os.Stat(recycleDir); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.SUCCESS,
			"data":    []RecycleBinItem{},
		})
		return
	}
	
	files, err := os.ReadDir(recycleDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": "读取回收站失败",
		})
		return
	}
	
	items := make([]RecycleBinItem, 0)
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			continue
		}
		
		items = append(items, RecycleBinItem{
			RecyclePath: filepath.ToSlash(filepath.Join(".recycle", f.Name())),
			Name:        f.Name(),
			Size:        info.Size(),
			DeletedAt:   info.ModTime(),
			IsDir:       f.IsDir(),
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"data":    items,
	})
}

// RestoreFromRecycleBin 从回收站恢复
func RestoreFromRecycleBin(c *gin.Context) {
	var req struct {
		RecyclePaths []string `json:"recycle_paths" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}
	
	restored := 0
	for _, recyclePath := range req.RecyclePaths {
		fullPath := filepath.Join("uploads", recyclePath)
		
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			continue
		}
		
		// 解析原始路径（从文件名中提取）
		name := filepath.Base(fullPath)
		parts := strings.SplitN(name, "_", 2)
		if len(parts) != 2 {
			continue
		}
		
		originalName := parts[1]
		originalPath := filepath.Join("uploads", originalName)
		
		// 处理冲突
		if _, err := os.Stat(originalPath); err == nil {
			// 文件已存在，添加后缀
			ext := filepath.Ext(originalName)
			base := strings.TrimSuffix(originalName, ext)
			originalPath = filepath.Join("uploads", fmt.Sprintf("%s_restored_%s", base, ext))
		}
		
		if err := os.Rename(fullPath, originalPath); err != nil {
			continue
		}
		
		restored++
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": fmt.Sprintf("已恢复 %d 个项目", restored),
	})
}

// EmptyRecycleBin 清空回收站
func EmptyRecycleBin(c *gin.Context) {
	recycleDir := "uploads/.recycle"
	
	if err := os.RemoveAll(recycleDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "清空回收站失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "回收站已清空",
	})
}

// GetFilePreview 获取文件预览
func GetFilePreview(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "路径不能为空",
		})
		return
	}
	
	safePath, ok := safeUploadPath(path)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  errmsg.ERROR,
			"message": "非法路径",
		})
		return
	}
	
	info, err := os.Stat(safePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  errmsg.ERROR,
			"message": "文件不存在",
		})
		return
	}
	
	if info.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "无法预览目录",
		})
		return
	}
	
	// 文件大小限制（1MB）
	if info.Size() > 1<<20 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "文件过大，无法预览",
		})
		return
	}
	
	ext := strings.ToLower(filepath.Ext(safePath))
	
	// 文本文件预览
	if isTextFile(ext) {
		content, err := os.ReadFile(safePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  errmsg.ERROR,
				"message": "读取文件失败",
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.SUCCESS,
			"type":    "text",
			"content": string(content),
			"ext":     ext,
		})
		return
	}
	
	// 图片预览
	if isImageFile(ext) {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.SUCCESS,
			"type":    "image",
			"url":     "/uploads/" + filepath.ToSlash(strings.TrimPrefix(safePath, "uploads/")),
			"size":    info.Size(),
		})
		return
	}
	
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  errmsg.ERROR,
		"message": "不支持预览的文件类型",
	})
}

// 辅助函数
func isDocumentFile(ext string) bool {
	docs := map[string]bool{
		".pdf":  true, ".doc": true, ".docx": true,
		".xls":  true, ".xlsx": true, ".ppt": true,
		".pptx": true, ".txt": true, ".md": true,
	}
	return docs[ext]
}

func isArchiveFile(ext string) bool {
	archives := map[string]bool{
		".zip": true, ".rar": true, ".7z": true,
		".tar": true, ".gz": true,
	}
	return archives[ext]
}

func isTextFile(ext string) bool {
	texts := map[string]bool{
		".txt":  true, ".md": true, ".json": true,
		".xml":  true, ".csv": true, ".log": true,
		".yaml": true, ".yml": true, ".html": true,
		".css":  true, ".js": true,
	}
	return texts[ext]
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// FormatSize 格式化文件大小
func FormatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes)/1024/1024)
	}
	return fmt.Sprintf("%.2f GB", float64(bytes)/1024/1024/1024)
}

// SaveFileMetadata 保存文件元数据
func SaveFileMetadata(c *gin.Context) {
	var req struct {
		Path     string            `json:"path" binding:"required"`
		Metadata map[string]string `json:"metadata"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}
	
	safePath, ok := safeUploadPath(req.Path)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  errmsg.ERROR,
			"message": "非法路径",
		})
		return
	}
	
	// 保存元数据到 .meta.json 文件
	metaPath := safePath + ".meta.json"
	data, _ := json.Marshal(req.Metadata)
	os.WriteFile(metaPath, data, 0644)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "元数据保存成功",
	})
}

// GetFileMetadata 获取文件元数据
func GetFileMetadata(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "路径不能为空",
		})
		return
	}
	
	safePath, ok := safeUploadPath(path)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  errmsg.ERROR,
			"message": "非法路径",
		})
		return
	}
	
	metaPath := safePath + ".meta.json"
	
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.SUCCESS,
			"data":    map[string]string{},
		})
		return
	}
	
	data, err := os.ReadFile(metaPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "读取元数据失败",
		})
		return
	}
	
	var metadata map[string]string
	json.Unmarshal(data, &metadata)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"data":    metadata,
	})
}
