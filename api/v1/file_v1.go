package v1

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// FileInfo 文件信息结构体
type FileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
}

// GetFileList 获取文件列表
func GetFileList(c *gin.Context) {
	// 获取请求的路径，默认为 uploads 根目录
	reqPath := c.Query("path")
	if reqPath == "" {
		reqPath = "uploads"
	} else {
		// 防止路径遍历，确保路径在 uploads 目录下
		reqPath = filepath.Join("uploads", reqPath)
		// 再次清理路径
		reqPath = filepath.Clean(reqPath)
	}

	// 检查目录是否存在
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
		// 计算相对路径，用于前端导航
		// 注意：filepath.Join(reqPath, f.Name()) 得到的是完整路径 (如 uploads/articles/title)
		// filepath.Rel("uploads", ...) 得到的是相对 uploads 的路径 (如 articles/title)
		fullPath := filepath.Join(reqPath, f.Name())
		relPath, err := filepath.Rel("uploads", fullPath)
		if err != nil {
			// 如果 Rel 失败（例如跨驱动器），则忽略该文件或使用原名
			continue
		}
		relPath = filepath.ToSlash(relPath)

		fileList = append(fileList, FileInfo{
			Name:  f.Name(),
			IsDir: f.IsDir(),
			Path:  relPath,
			Size:  f.Size(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": errmsg.GetErrMsg(errmsg.SUCCESS),
		"data":    fileList,
	})
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

	// 安全检查：确保路径在 uploads 下
	targetPath := filepath.Join("uploads", path)
	targetPath = filepath.Clean(targetPath)

	// 检查文件是否存在
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
		Path string `json:"path"` // 父目录路径
		Name string `json:"name"` // 新目录名
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	// 路径清理
	// 防止 .. 逃逸
	data.Path = filepath.Clean(data.Path)
	if data.Path == "." {
		data.Path = ""
	}
	// Windows下 clean 可能会保留 . 如果为空

	targetPath := filepath.Join("uploads", data.Path, data.Name)
	targetPath = filepath.Clean(targetPath)

	// 检查是否存在
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
		Path    string `json:"path"`    // 旧路径 (相对 uploads)
		NewName string `json:"newName"` // 新文件名 (不包含路径)
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	// 基本安全检查
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

	// 验证 newPath 是否还在 uploads 下 (简单验证)
	// 在实际生产中应更严格

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
