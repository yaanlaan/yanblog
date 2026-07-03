package v1

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// UploadTaskV2 增强版上传任务
type UploadTaskV2 struct {
	ID              string           `json:"id"`
	FileName        string           `json:"file_name"`
	FileSize        int64            `json:"file_size"`
	TotalFiles      int              `json:"total_files"`
	Processed       int              `json:"processed"`
	Success         int              `json:"success"`
	Failed          int              `json:"failed"`
	Status          string           `json:"status"` // processing, completed, failed, cancelled, retrying
	Errors          []UploadError    `json:"errors,omitempty"`
	StartTime       time.Time        `json:"start_time"`
	EndTime         *time.Time       `json:"end_time,omitempty"`
	Cancelled       bool             `json:"-"`
	RetryCount      int              `json:"retry_count"`        // 重试次数
	MaxRetries      int              `json:"max_retries"`        // 最大重试次数
	ProcessedFiles  []string         `json:"processed_files"`    // 已处理文件列表
	FailedFiles     []string         `json:"failed_files"`       // 失败文件列表（用于重试）
	Speed           float64          `json:"speed"`              // 上传速度 (MB/s)
	ETA             string           `json:"eta"`                // 预计剩余时间
	Clients         []*websocket.Conn `json:"-"`                 // WebSocket 客户端
	ClientsMu       sync.Mutex       `json:"-"`
	mu              sync.Mutex
}

// UploadError 上传错误详情
type UploadError struct {
	FileName string `json:"file_name"`
	Error    string `json:"error"`
	Retried  bool   `json:"retried"`  // 是否已重试
}

// UploadHistory 上传历史记录
type UploadHistory struct {
	TaskID      string    `json:"task_id"`
	FileName    string    `json:"file_name"`
	TotalFiles  int       `json:"total_files"`
	Success     int       `json:"success"`
	Failed      int       `json:"failed"`
	Status      string    `json:"status"`
	Duration    string    `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
}

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境应该检查 Origin
	},
}

// 全局管理
var (
	uploadTasksV2   = make(map[string]*UploadTaskV2)
	tasksMuV2       sync.RWMutex
	maxConcurrentV2 = 3
	semaphoreV2     = make(chan struct{}, maxConcurrentV2)
	
	uploadHistory   = make([]UploadHistory, 0, 100) // 保留最近100条
	historyMu       sync.RWMutex
	historyFile     = "./data/upload_history.json"
)

// init 加载历史记录
func init() {
	loadUploadHistory()
}

// WebSocketProgress 通过 WebSocket 推送进度
func WebSocketProgress(c *gin.Context) {
	taskID := c.Param("id")
	
	tasksMuV2.RLock()
	task, exists := uploadTasksV2[taskID]
	tasksMuV2.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  errmsg.ERROR,
			"message": "任务不存在",
		})
		return
	}
	
	// 升级为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	
	// 添加客户端
	task.ClientsMu.Lock()
	task.Clients = append(task.Clients, conn)
	task.ClientsMu.Unlock()
	
	// 立即发送当前进度
	sendProgress(task, conn)
	
	// 保持连接
	defer func() {
		task.ClientsMu.Lock()
		for i, client := range task.Clients {
			if client == conn {
				task.Clients = append(task.Clients[:i], task.Clients[i+1:]...)
				break
			}
		}
		task.ClientsMu.Unlock()
		conn.Close()
	}()
	
	// 读取消息（可选，用于接收取消命令）
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// sendProgress 发送进度到 WebSocket 客户端
func sendProgress(task *UploadTaskV2, conn *websocket.Conn) {
	task.mu.Lock()
	progress := map[string]interface{}{
		"type":          "progress",
		"task_id":       task.ID,
		"file_name":     task.FileName,
		"total_files":   task.TotalFiles,
		"processed":     task.Processed,
		"success":       task.Success,
		"failed":        task.Failed,
		"progress":      float64(task.Processed) / float64(task.TotalFiles) * 100,
		"status":        task.Status,
		"speed":         task.Speed,
		"eta":           task.ETA,
		"retry_count":   task.RetryCount,
		"errors":        task.Errors,
		"start_time":    task.StartTime,
		"end_time":      task.EndTime,
	}
	task.mu.Unlock()
	
	conn.WriteJSON(progress)
}

// broadcastProgress 广播进度到所有客户端
func broadcastProgress(task *UploadTaskV2) {
	task.ClientsMu.Lock()
	clients := make([]*websocket.Conn, len(task.Clients))
	copy(clients, task.Clients)
	task.ClientsMu.Unlock()
	
	for _, conn := range clients {
		sendProgress(task, conn)
	}
}

// GetUploadProgressV2 获取上传进度（HTTP）
func GetUploadProgressV2(c *gin.Context) {
	taskID := c.Param("id")
	
	tasksMuV2.RLock()
	task, exists := uploadTasksV2[taskID]
	tasksMuV2.RUnlock()
	
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
		"file_size":   task.FileSize,
		"total_files": task.TotalFiles,
		"processed":   task.Processed,
		"success":     task.Success,
		"failed":      task.Failed,
		"progress":    float64(task.Processed) / float64(task.TotalFiles) * 100,
		"task_status": task.Status,
		"speed":       task.Speed,
		"eta":         task.ETA,
		"retry_count": task.RetryCount,
		"errors":      task.Errors,
		"start_time":  task.StartTime,
		"end_time":    task.EndTime,
	}
	task.mu.Unlock()
	
	c.JSON(http.StatusOK, progress)
}

// CancelUploadV2 取消上传任务
func CancelUploadV2(c *gin.Context) {
	taskID := c.Param("id")
	
	tasksMuV2.RLock()
	task, exists := uploadTasksV2[taskID]
	tasksMuV2.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  errmsg.ERROR,
			"message": "任务不存在",
		})
		return
	}
	
	task.mu.Lock()
	if task.Status == "processing" || task.Status == "retrying" {
		task.Cancelled = true
		task.Status = "cancelled"
		now := time.Now()
		task.EndTime = &now
	}
	task.mu.Unlock()
	
	// 通知所有客户端
	broadcastProgress(task)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "任务已取消",
	})
}

// RetryFailedUpload 重试失败的文件
func RetryFailedUpload(c *gin.Context) {
	taskID := c.Param("id")
	
	tasksMuV2.RLock()
	task, exists := uploadTasksV2[taskID]
	tasksMuV2.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  errmsg.ERROR,
			"message": "任务不存在",
		})
		return
	}
	
	task.mu.Lock()
	if task.Status != "completed" && task.Status != "failed" {
		task.mu.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "只能重试已完成或失败的任务",
		})
		return
	}
	
	if len(task.FailedFiles) == 0 {
		task.mu.Unlock()
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.SUCCESS,
			"message": "没有失败的文件需要重试",
		})
		return
	}
	
	// 重置状态
	task.Status = "retrying"
	task.RetryCount++
	task.Processed = task.TotalFiles - len(task.FailedFiles)
	task.Success = task.TotalFiles - len(task.FailedFiles)
	task.Failed = 0
	task.Errors = make([]UploadError, 0)
	task.mu.Unlock()
	
	// 通知客户端
	broadcastProgress(task)
	
	// 异步重试
	go func() {
		// TODO: 实现重试逻辑
		// 这需要保存原始的 ZIP 文件或者文件路径
	}()
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": fmt.Sprintf("开始重试 %d 个失败文件", len(task.FailedFiles)),
		"data": gin.H{
			"failed_files": task.FailedFiles,
			"retry_count":  task.RetryCount,
		},
	})
}

// GetUploadHistory 获取上传历史
func GetUploadHistory(c *gin.Context) {
	page := 1
	size := 20
	
	historyMu.RLock()
	total := len(uploadHistory)
	start := max(0, total-page*size)
	end := min(total, start+size)
	
	history := make([]UploadHistory, 0)
	if start < total {
		history = uploadHistory[start:end]
	}
	historyMu.RUnlock()
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"total":   total,
		"page":    page,
		"size":    size,
		"data":    history,
	})
}

// UploadArticleZipV2 增强版ZIP上传
func UploadArticleZipV2(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "文件上传失败",
		})
		return
	}

	// 创建任务
	taskID := fmt.Sprintf("upload_v2_%d", time.Now().UnixNano())
	task := &UploadTaskV2{
		ID:         taskID,
		FileName:   file.Filename,
		FileSize:   file.Size,
		Status:     "processing",
		StartTime:  time.Now(),
		Errors:     make([]UploadError, 0),
		MaxRetries: 3,
		Clients:    make([]*websocket.Conn, 0),
	}
	
	tasksMuV2.Lock()
	uploadTasksV2[taskID] = task
	tasksMuV2.Unlock()
	
	// 并发控制
	select {
	case semaphoreV2 <- struct{}{}:
		defer func() { <-semaphoreV2 }()
	default:
		c.JSON(http.StatusTooManyRequests, gin.H{
			"status":  errmsg.ERROR_UPLOAD_BUSY,
			"message": fmt.Sprintf("上传任务繁忙，最多支持 %d 个并发任务", maxConcurrentV2),
		})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		updateTaskFailed(task, "打开文件失败: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "打开文件失败",
		})
		return
	}
	defer src.Close()

	// 保存 ZIP 文件（用于重试）
	zipPath := fmt.Sprintf("./temp_zip/%s.zip", taskID)
	os.MkdirAll("./temp_zip", 0755)
	
	outFile, err := os.Create(zipPath)
	if err != nil {
		updateTaskFailed(task, "保存文件失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  errmsg.ERROR,
			"message": "保存文件失败",
		})
		return
	}
	
	// 复制文件内容
	io.Copy(outFile, src)
	outFile.Close()
	src.Close()
	
	// 重新打开用于处理
	src, _ = os.Open(zipPath)
	defer src.Close()
	
	tempDir := fmt.Sprintf("./temp_zip/%s", taskID)
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)
	
	// 流式处理
	startTime := time.Now()
	successCount, errors := processZipStreamV2(c.Request.Context(), src, file.Size, tempDir, task, startTime)
	
	// 更新任务状态
	task.mu.Lock()
	task.Processed = task.TotalFiles
	task.Success = successCount
	task.Failed = task.TotalFiles - successCount
	
	for _, err := range errors {
		task.Errors = append(task.Errors, UploadError{
			FileName: err.FileName,
			Error:    err.Error,
			Retried:  false,
		})
		if err.FileName != "" {
			task.FailedFiles = append(task.FailedFiles, err.FileName)
		}
	}
	
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
	
	// 保存到历史
	saveToHistory(task)
	
	// 通知所有客户端
	broadcastProgress(task)
	
	// 清理 ZIP 文件（如果成功）
	if task.Status == "completed" && task.Failed == 0 {
		os.Remove(zipPath)
	}
	
	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"task_id": taskID,
		"message": fmt.Sprintf("上传完成，成功 %d/%d", task.Success, task.TotalFiles),
		"data": gin.H{
			"total":       task.TotalFiles,
			"success":     task.Success,
			"failed":      task.Failed,
			"errors":      task.Errors,
			"zip_saved":   task.Failed > 0, // 如果有失败，保留了 ZIP 用于重试
		},
	})
}

// processZipStreamV2 增强版流式处理
func processZipStreamV2(ctx context.Context, reader io.ReaderAt, size int64, tempDir string, task *UploadTaskV2, startTime time.Time) (int, []UploadError) {
	successCount := 0
	errors := make([]UploadError, 0)
	
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		errors = append(errors, UploadError{Error: "ZIP文件损坏: " + err.Error()})
		return 0, errors
	}
	
	// 统计文件数
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
		errors = append(errors, UploadError{Error: "ZIP文件中未找到Markdown文件"})
		return 0, errors
	}
	
	// 处理每个文件
	for _, zipFile := range zipReader.File {
		// 检查取消
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
		
		if zipFile.FileInfo().IsDir() || !strings.HasSuffix(strings.ToLower(zipFile.Name), ".md") {
			continue
		}
		
		// 提取文件
		fileDir := filepath.Join(tempDir, filepath.Dir(zipFile.Name))
		os.MkdirAll(fileDir, 0755)
		
		destPath := filepath.Join(tempDir, zipFile.Name)
		if err := extractZipFile(zipFile, destPath); err != nil {
			task.mu.Lock()
			task.Processed++
			task.Failed++
			errors = append(errors, UploadError{
				FileName: zipFile.Name,
				Error:    "解压失败: " + err.Error(),
			})
			task.FailedFiles = append(task.FailedFiles, zipFile.Name)
			task.mu.Unlock()
			
			broadcastProgress(task)
			continue
		}
		
		// 处理文章（复用原有逻辑）
		articleDir := filepath.Join(tempDir, filepath.Dir(zipFile.Name))
		article, code := processZipArticle(articleDir)
		
		// 计算速度和 ETA
		task.mu.Lock()
		task.Processed++
		
		elapsed := time.Since(startTime).Seconds()
		if elapsed > 0 {
			task.Speed = float64(task.Processed) / elapsed
			remaining := task.TotalFiles - task.Processed
			if task.Speed > 0 {
				etaSeconds := float64(remaining) / task.Speed
				task.ETA = formatDuration(etaSeconds)
			}
		}
		
		if code == errmsg.SUCCESS && article != nil {
			successCount++
			task.Success++
			task.ProcessedFiles = append(task.ProcessedFiles, zipFile.Name)
		} else {
			task.Failed++
			errMsg := errmsg.GetErrMsg(code)
			if errMsg == "" {
				errMsg = "处理失败"
			}
			errors = append(errors, UploadError{
				FileName: zipFile.Name,
				Error:    errMsg,
			})
			task.FailedFiles = append(task.FailedFiles, zipFile.Name)
		}
		task.mu.Unlock()
		
		// 广播进度
		broadcastProgress(task)
		
		// 清理
		os.RemoveAll(articleDir)
	}
	
	return successCount, errors
}

// updateTaskFailed 更新任务为失败状态
func updateTaskFailed(task *UploadTaskV2, errMsg string) {
	task.mu.Lock()
	task.Status = "failed"
	task.Errors = append(task.Errors, UploadError{Error: errMsg})
	now := time.Now()
	task.EndTime = &now
	task.mu.Unlock()
	
	broadcastProgress(task)
}

// formatDuration 格式化持续时间
func formatDuration(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.0f秒", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%.0f分钟", seconds/60)
	}
	return fmt.Sprintf("%.1f小时", seconds/3600)
}

// saveToHistory 保存到历史记录
func saveToHistory(task *UploadTaskV2) {
	task.mu.Lock()
	duration := "进行中"
	if task.EndTime != nil {
		duration = task.EndTime.Sub(task.StartTime).String()
	}
	
	history := UploadHistory{
		TaskID:     task.ID,
		FileName:   task.FileName,
		TotalFiles: task.TotalFiles,
		Success:    task.Success,
		Failed:     task.Failed,
		Status:     task.Status,
		Duration:   duration,
		CreatedAt:  task.StartTime,
	}
	task.mu.Unlock()
	
	historyMu.Lock()
	uploadHistory = append([]UploadHistory{history}, uploadHistory...)
	if len(uploadHistory) > 100 {
		uploadHistory = uploadHistory[:100]
	}
	historyMu.Unlock()
	
	// 持久化到文件
	go persistHistory()
}

// persistHistory 持久化历史记录
func persistHistory() {
	historyMu.RLock()
	data, _ := json.Marshal(uploadHistory)
	historyMu.RUnlock()
	
	os.MkdirAll(filepath.Dir(historyFile), 0755)
	os.WriteFile(historyFile, data, 0644)
}

// loadUploadHistory 加载历史记录
func loadUploadHistory() {
	data, err := os.ReadFile(historyFile)
	if err != nil {
		return
	}
	
	historyMu.Lock()
	json.Unmarshal(data, &uploadHistory)
	historyMu.Unlock()
}

// ClearUploadHistory 清空历史记录
func ClearUploadHistory(c *gin.Context) {
	historyMu.Lock()
	uploadHistory = make([]UploadHistory, 0, 100)
	historyMu.Unlock()
	
	os.Remove(historyFile)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  errmsg.SUCCESS,
		"message": "历史记录已清空",
	})
}

// min/max 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
