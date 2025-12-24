package v1

import (
	"yanblog/middlewares"
	"yanblog/model"
	"yanblog/utils/errmsg"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login 用户登录接口
// 处理用户登录请求，验证用户名和密码，成功后生成JWT token
func Login(c *gin.Context) {
	var data model.User

	c.ShouldBindJSON(&data)
	var token string
	var code int

	fmt.Println("Login attempt for user:", data.Username)

	code = model.CheckLogin(data.Username, data.Password)

	var role int
	if code == errmsg.SUCCESS {
		token, code = middlewares.SetToken(data.Username)
		role = model.GetUserRole(data.Username)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   code,
		"message":  errmsg.GetErrMsg(code),
		"token":    token, // JWT token（登录成功时返回）
		"username": data.Username,
		"role":     role,
	})
}
