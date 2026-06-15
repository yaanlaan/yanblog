package v1

import (
	"yanblog/middlewares"
	"yanblog/model"
	"yanblog/utils/errmsg"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Login 用户登录接口
// 处理用户登录请求，验证用户名和密码，成功后生成JWT token
func Login(c *gin.Context) {
	var data model.User

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "请求参数错误",
		})
		return
	}

	code := model.CheckLogin(data.Username, data.Password)

	// 登录失败时直接返回
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}

	// 登录成功：生成 token
	token, tokenCode := middlewares.SetToken(data.Username)
	if tokenCode != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  tokenCode,
			"message": errmsg.GetErrMsg(tokenCode),
		})
		return
	}

	// 获取用户角色
	role := model.GetUserRole(data.Username)

	// 只返回一次成功响应（包含 token 和用户信息）
	c.JSON(http.StatusOK, gin.H{
		"status":   code,
		"message":  errmsg.GetErrMsg(code),
		"token":    token,
		"username": data.Username,
		"role":     role,
	})
}
