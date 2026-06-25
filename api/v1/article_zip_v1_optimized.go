package v1

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// UploadTask 上传任务状态
type UploadTask struct {
	ID          string     `json:"id"`
	FileName    string     `json:"file_name"`
	TotalFiles  int        `json:"total_files"`
	Processed   int        `json:"processed"`
	Success     int        `json:"success"`
	Failed      int        `json:"failed"`
	Status      string     `json:"status"` // processing, completed, failed, cancelled
	Errors      []string   `json:"errors,omitempty"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Cancelled   bool       `json:"-"`
	mu          sync.Mutex
}

// 全局上传任务管理
var (
	uploadTasks = make(map[string]*UploadTask)
	tasksMu     sync.RWMutex
	maxConcurrentUploads = 3 // 最大并发上传数
	uploadSemaphore = make(chan struct{}, maxConcurrentUploads)
)

// CancelUpload 取消上传任务
func CancelUpload(c *gin.Context) {
	taskID := c.Param("id")
	
	tasksMu.RLock()
	task, exists := uploadTasks[taskID]
	tasksMu.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  errmsg.ERROR,
			"message": "任务不存在",
		})
		return
	}
	
	task.mu.Lock()
	if task.Status == "processing" {
		task.Cancelled = true
		task.Status = "cancelled"
		now := time.Now()
		task.EndTime = &now
	}
	task.mu.Unlock()
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "任务已取消",
	})
}

// GetUploadProgress 获取上传进度
func GetUploadProgress(c *gin.Context) {
	taskID := c.Param("id")
	
	tasksMu.RLock()
	task, exists := uploadTasks[taskID]
	tasksMu.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  errmsg.ERROR,
			"message": "任务不存在",
		})
		return
	}
	
	task.mu.Lock()
	progress := gin.H{
		"status":      errmsg.SUCCESS,
		"task_id":     task.ID,
		"file_name":   task.FileName,
		"total_files": task.TotalFiles,
		"processed":   task.Processed,
		"success":     task.Success,
		"failed":      task.Failed,
		"progress":    float64(task.Processed) / float64(task.TotalFiles) * 100,
		"task_status": task.Status,
		"errors":      task.Errors,
		"start_time":  task.StartTime,
		"end_time":    task.EndTime,
	}
	task.mu.Unlock()
	
	c.JSON(http.StatusOK, progress)
}

// UploadArticleZipOptimized 优化的ZIP上传（流式处理 + 进度跟踪）
func UploadArticleZipOptimized(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "文件上传失败",
		})
		return
	}

	// 创建任务
	taskID := fmt.Sprintf("upload_%d", time.Now().UnixNano())
	task := &UploadTask{
		ID:        taskID,
		FileName:  file.Filename,
		Status:    "processing",
		StartTime: time.Now(),
		Errors:    make([]string, 0),
	}
	
	tasksMu.Lock()
	uploadTasks[taskID] = task
	tasksMu.Unlock()
	
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		task.mu.Lock()
		task.Status = "failed"
		task.Errors = append(task.Errors, "打开文件失败: "+err.Error())
		now := time.Now()
		task.EndTime = &now
		task.mu.Unlock()
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "打开文件失败",
		})
		return
	}
	defer src.Close()

	// 使用信号量控制并发
	select {
	case uploadSemaphore <- struct{}{}:
		defer func() { <-uploadSemaphore }()
	default:
		c.JSON(http.StatusTooManyRequests, gin.H{
			"status":  errmsg.ERROR_UPLOAD_BUSY,
			"message": fmt.Sprintf("上传任务繁忙，最多支持 %d 个并发任务", maxConcurrentUploads),
		})
		return
	}

	// 流式处理 ZIP
	tempDir := fmt.Sprintf("./temp_zip/%s", taskID)
	defer os.RemoveAll(tempDir)
	
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		task.mu.Lock()
		task.Status = "failed"
		task.Errors = append(task.Errors, "创建临时目录失败")
		now := time.Now()
		task.EndTime = &now
		task.mu.Unlock()
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "创建临时目录失败",
		})
		return
	}

	// 流式解压（边读边处理，不保存完整ZIP）
	successCount, errors := processZipStream(c.Request.Context(), src, file.Size, tempDir, task)
	
	// 更新任务状态
	task.mu.Lock()
	task.Processed = task.TotalFiles
	task.Success = successCount
	task.Failed = task.TotalFiles - successCount
	task.Errors = append(task.Errors, errors...)
	
	if task.Cancelled {
		task.Status = "cancelled"
	} else if task.Failed == task.TotalFiles {
		task.Status = "failed"
	} else {
		task.Status = "completed"
	}
	now := time.Now()
	task.EndTime = &now
	task.mu.Unlock()
	
	// 返回结果
	if task.Cancelled {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.SUCCESS,
			"task_id": taskID,
			"message": "上传已取消",
			"data": gin.H{
				"processed": task.Processed,
				"success":   task.Success,
				"cancelled": true,
			},
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"task_id": taskID,
		"message": fmt.Sprintf("上传完成，成功 %d/%d", task.Success, task.TotalFiles),
		"data": gin.H{
			"total":   task.TotalFiles,
			"success": task.Success,
			"failed":  task.Failed,
			"errors":  task.Errors,
		},
	})
}

// processZipStream 流式处理 ZIP 文件
func processZipStream(ctx context.Context, reader io.ReaderAt, size int64, tempDir string, task *UploadTask) (int, []string) {
	successCount := 0
	errors := make([]string, 0)
	
	// 创建 ZIP 读取器
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		errors = append(errors, "ZIP文件损坏: "+err.Error())
		return 0, errors
	}
	
	// 统计 MD 文件数量
	mdCount := 0
	for _, f := range zipReader.File {
		if !f.FileInfo().IsDir() && strings.HasSuffix(strings.ToLower(f.Name), ".md") {
			mdCount++
		}
	}
	
	task.mu.Lock()
	task.TotalFiles = mdCount
	task.mu.Unlock()
	
	if mdCount == 0 {
		errors = append(errors, "ZIP文件中未找到Markdown文件")
		return 0, errors
	}
	
	// 处理每个文件
	for _, zipFile := range zipReader.File {
		// 检查是否取消
		select {
		case <-ctx.Done():
			task.mu.Lock()
			task.Cancelled = true
			task.mu.Unlock()
			return successCount, errors
		default:
		}
		
		if task.Cancelled {
			return successCount, errors
		}
		
		// 只处理 .md 文件
		if zipFile.FileInfo().IsDir() || !strings.HasSuffix(strings.ToLower(zipFile.Name), ".md") {
			continue
		}
		
		// 提取单个文件到临时目录
		fileDir := filepath.Join(tempDir, filepath.Dir(zipFile.Name))
		os.MkdirAll(fileDir, 0755)
		
		destPath := filepath.Join(tempDir, zipFile.Name)
		if err := extractZipFile(zipFile, destPath); err != nil {
			task.mu.Lock()
			task.Processed++
			task.Failed++
			errors = append(errors, fmt.Sprintf("%s: 解压失败 - %v", zipFile.Name, err))
			task.mu.Unlock()
			continue
		}
		
		// 处理文章
		articleDir := filepath.Join(tempDir, filepath.Dir(zipFile.Name))
		article, code := processZipArticle(articleDir)
		
		task.mu.Lock()
		task.Processed++
		
		if code == errmsg.SUCCESS && article != nil {
			successCount++
			task.Success++
		} else {
			task.Failed++
			errMsg := errmsg.GetErrMsg(code)
			if errMsg == "" {
				errMsg = "处理失败"
			}
			errors = append(errors, fmt.Sprintf("%s: %s", zipFile.Name, errMsg))
		}
		task.mu.Unlock()
		
		// 清理已处理的文件
		os.RemoveAll(articleDir)
	}
	
	return successCount, errors
}

// extractZipFile 提取单个ZIP文件
func extractZipFile(zipFile *zip.File, destPath string) error {
	// 防止 Zip Slip 漏洞
	cleanDest := filepath.Clean(destPath)
	if !strings.HasPrefix(cleanDest, filepath.Clean(filepath.Dir(destPath))+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", cleanDest)
	}
	
	rc, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	
	outFile, err := os.OpenFile(cleanDest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()
	
	// 限制文件大小（100MB）
	const maxFileSize = 100 << 20
	written, err := io.CopyN(outFile, rc, maxFileSize)
	if err != nil && err != io.EOF {
		return err
	}
	
	if written >= maxFileSize {
		return fmt.Errorf("文件过大，超过 %d MB 限制", maxFileSize>>20)
	}
	
	return nil
}

// UploadArticleZipBatchOptimized 优化的批量ZIP上传
func UploadArticleZipBatchOptimized(c *gin.Context) {
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

	// 创建主任务
	taskID := fmt.Sprintf("batch_%d", time.Now().UnixNano())
	task := &UploadTask{
		ID:        taskID,
		FileName:  fmt.Sprintf("batch_%d_files", len(files)),
		Status:    "processing",
		StartTime: time.Now(),
		Errors:    make([]string, 0),
	}
	
	tasksMu.Lock()
	uploadTasks[taskID] = task
	tasksMu.Unlock()
	
	// 使用信号量控制并发
	select {
	case uploadSemaphore <- struct{}{}:
		defer func() { <-uploadSemaphore }()
	default:
		c.JSON(http.StatusTooManyRequests, gin.H{
			"status":  errmsg.ERROR_UPLOAD_BUSY,
			"message": fmt.Sprintf("上传任务繁忙，最多支持 %d 个并发任务", maxConcurrentUploads),
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
	var resultsMu sync.Mutex
	
	// 并发处理文件
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3) // 每个批次任务最多3个并发
	
	for _, fileHeader := range files {
		if task.Cancelled {
			break
		}
		
		wg.Add(1)
		go func(fh *multipart.FileHeader) {
			defer wg.Done()
			
			// 获取处理权限
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-c.Request.Context().Done():
				task.mu.Lock()
				task.Cancelled = true
				task.mu.Unlock()
				return
			}
			
			// 检查是否取消
			if task.Cancelled {
				return
			}
			
			result := Result{FileName: fh.Filename}
			
			// 打开文件
			src, err := fh.Open()
			if err != nil {
				result.Status = errmsg.ERROR
				result.Message = "打开文件失败"
				resultsMu.Lock()
				results = append(results, result)
				resultsMu.Unlock()
				return
			}
			defer src.Close()
			
			// 创建临时目录
			tempDir := fmt.Sprintf("./temp_zip/%s_%s", taskID, fh.Filename)
			defer os.RemoveAll(tempDir)
			
			if err := os.MkdirAll(tempDir, 0755); err != nil {
				result.Status = errmsg.ERROR
				result.Message = "创建临时目录失败"
				resultsMu.Lock()
				results = append(results, result)
				resultsMu.Unlock()
				return
			}
			
			// 读取文件到内存
			fileData, err := io.ReadAll(src)
			if err != nil {
				result.Status = errmsg.ERROR
				result.Message = "读取文件失败"
				resultsMu.Lock()
				results = append(results, result)
				resultsMu.Unlock()
				return
			}
			
			// 流式处理
			reader := strings.NewReader(string(fileData))
			successCount, errs := processZipStream(c.Request.Context(), 
				&readerAt{reader}, 
				int64(len(fileData)), 
				tempDir, 
				task)
			
			task.mu.Lock()
			task.Success += successCount
			task.Failed += len(errs)
			task.Errors = append(task.Errors, errs...)
			task.mu.Unlock()
			
			// 找到第一个成功处理的文章标题
			if successCount > 0 {
				// 重新处理获取标题
				tempZipPath := filepath.Join(tempDir, "temp.zip")
				os.WriteFile(tempZipPath, fileData, 0644)
				
				unzipDir := tempDir + "_extracted"
				if err := unzip(tempZipPath, unzipDir); err == nil {
					article, _ := processZipArticle(unzipDir)
					if article != nil {
						result.Title = article.Title
					}
					os.RemoveAll(unzipDir)
				}
				os.Remove(tempZipPath)
			}
			
			result.Status = errmsg.SUCCESS
			if len(errs) > 0 {
				result.Message = fmt.Sprintf("部分成功 (%d/%d)", successCount, successCount+len(errs))
			} else {
				result.Message = "上传成功"
			}
			
			resultsMu.Lock()
			results = append(results, result)
			resultsMu.Unlock()
			
		}(fileHeader)
	}
	
	// 等待所有任务完成
	wg.Wait()
	
	// 更新主任务状态
	task.mu.Lock()
	task.Processed = task.TotalFiles
	if task.Cancelled {
		task.Status = "cancelled"
	} else if task.Failed == task.TotalFiles {
		task.Status = "failed"
	} else {
		task.Status = "completed"
	}
	now := time.Now()
	task.EndTime = &now
	task.mu.Unlock()
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"task_id": taskID,
		"total":   len(files),
		"success": task.Success,
		"failed":  task.Failed,
		"results": results,
	})
}

// readerAt 适配器，将 io.Reader 转换为 io.ReaderAt
type readerAt struct {
	r io.ReaderAt
}

func (r *readerAt) ReadAt(p []byte, off int64) (n int, err error) {
	return r.r.ReadAt(p, off)
}
