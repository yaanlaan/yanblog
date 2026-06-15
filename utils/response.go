package utils

import (
	"net/http"
	"strconv"

	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// 分页相关常量
const (
	MaxPageSize = 100 // 单页最大记录数上限
)

// Response 统一 API 响应封装，消除各 handler 中重复的 c.JSON(gin.H{...}) 模式

// Success 返回成功响应（带数据）
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"data":    data,
		"message": errmsg.GetErrMsg(200),
	})
}

// SuccessWithTotal 返回成功响应（带数据和总数）
func SuccessWithTotal(c *gin.Context, data interface{}, total int64) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(200),
	})
}

// Error 返回错误响应
func Error(c *gin.Context, code int) {
	c.JSON(http.StatusOK, gin.H{ // 使用 200 + status 字段，保持前后端一致
		"status":  code,
		"data":    nil,
		"message": errmsg.GetErrMsg(code),
	})
}

// ErrorWithMessage 返回错误响应（自定义消息）
func ErrorWithMessage(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    nil,
		"message": message,
	})
}

// BadRequest 参数错误
func BadRequest(c *gin.Context, message string) {
	ErrorWithMessage(c, errmsg.ERROR, message)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	ErrorWithMessage(c, errmsg.ERROR_ART_NOT_EXIST, message)
}

// ParsePageParams 从查询参数解析分页信息
// 返回 pageSize, pageNum, 以及是否为"请求全部"模式
// 边界校验：pageSize 和 pageNum <=0 时设为 -1（查全部），pageSize 有上限限制
func ParsePageParams(c *gin.Context) (pageSize int, pageNum int, isAll bool) {
	pageSize, _ = strconv.Atoi(c.Query("pagesize"))
	pageNum, _ = strconv.Atoi(c.Query("pagenum"))

	if pageSize <= 0 {
		pageSize = -1
	}
	if pageNum <= 0 {
		pageNum = -1
	}

	// 限制单页最大数量，防止恶意请求
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	isAll = (pageSize == -1 && pageNum == -1)
	return
}

// ParseIDParam 从 URL 参数解析整数 ID
// 成功返回 id 和 true；解析失败返回 false 并自动写入错误响应
func ParseIDParam(c *gin.Context) (int, bool) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(c, "无效的 ID 参数")
		return 0, false
	}
	return id, true
}
