package middlewares

import (
	"yanblog/utils"
	"yanblog/utils/errmsg"

	"net/http"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte(utils.ServerConfig.JwtKey)

type MyClaims struct {
	Username string `json:"username"`           // 用户名
	jwt.RegisteredClaims                          // 标准JWT声明（包含过期时间、签发者等）
}

// SetToken 生成JWT token
// 参数: username - 用户名
// 返回: token字符串和状态码
func SetToken(username string) (string, int) {
	// 设置token过期时间为10小时后
	expireTime := time.Now().Add(10 * time.Hour)
	
	SetClaims := MyClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime), // 过期时间
			Issuer:    "yanblog",                      // 签发者
		},
	}
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	
	token, err := reqClaim.SignedString(JwtKey)
	if err != nil {
		return "", errmsg.ERROR
	}
	
	return token, errmsg.SUCCESS
}

// CheckToken 验证JWT token的有效性
// 参数: tokenString - JWT token字符串
// 返回: 解析后的声明和状态码
func CheckToken(tokenString string) (*MyClaims, int) {
	// 解析并验证token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil // 提供验证密钥
	})

	if err != nil {
		return nil, errmsg.ERROR
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, errmsg.SUCCESS
	} else {
		return nil, errmsg.ERROR
	}
}

// JwtToken JWT认证中间件函数
// 返回: Gin中间件处理函数
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int // 声明响应状态码变量
		
		// 1. 从请求头获取Authorization字段
		tokenHeader := c.Request.Header.Get("Authorization")
		
		// 2. 检查token是否存在
		if tokenHeader == "" {
			code = errmsg.ERROR_TOKEN_EXIST
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort() // 终止请求处理
			return
		}

		// 3. 验证token格式是否为"Bearer token"
		checkToken := strings.SplitN(tokenHeader, " ", 2)
		if len(checkToken) != 2 || checkToken[0] != "Bearer" {
			code = errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		// 4. 验证token有效性
		key, Tcode := CheckToken(checkToken[1])
		if Tcode == errmsg.ERROR {
			code = errmsg.ERROR_TOKEN_WRONG
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		// 5. 检查token是否过期
		if key.RegisteredClaims.ExpiresAt == nil || time.Now().After(key.RegisteredClaims.ExpiresAt.Time) {
			code = errmsg.ERROR_TOKEN_RUNTIME
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		// 6. 验证通过，将用户名存入上下文供后续处理使用
		c.Set("username", key.Username)
		c.Next()
	}
}