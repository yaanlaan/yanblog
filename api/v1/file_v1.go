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
