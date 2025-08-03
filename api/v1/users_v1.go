package v1

import (
	"yanblog/model"
	"yanblog/utils/errmsg"
	"yanblog/utils/validator"

	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 添加用户
func AddUser(c *gin.Context) {
	var data model.User
	var msg string
	var code int

	_ = c.ShouldBindJSON(&data)
	//fmt.Printf("data: %+v\n", data) 
	msg, code = validator.Validate(&data)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": msg,
		})
		return
	}

	code = model.CheckUser(data.Username)
	if code == errmsg.SUCCESS {
		model.CreateUser(&data)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		// "data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询用户列表
func GetUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	var code int

	if pageSize <=0 {
		pageSize = -1
	}
	if pageNum <= 0 {
		pageNum = -1
	}

	data, total := model.GetUsers(pageSize, pageNum)
	code = errmsg.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// 搜索用户
func SearchUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	keyword := c.Query("keyword")
	roleStr := c.Query("role")
	var role int
	var code int

	if roleStr != "" {
		role, _ = strconv.Atoi(roleStr)
	}

	if pageSize <=0 {
		pageSize = -1
	}
	if pageNum <= 0 {
		pageNum = -1
	}

	data, total := model.SearchUser(keyword, role, pageSize, pageNum)
	code = errmsg.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// 编辑用户
func EditUser(c *gin.Context) {
	var data model.User
	var code int
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)
	code = model.CheckUser(data.Username)
	if code == errmsg.SUCCESS {//用户名不在
		model.EditUser(id, &data)
	}
	if code == errmsg.ERROR_USERNAME_USED {//用户名已存在
		code = model.CheckUserWithID(id, data.Username)//判断是不是自己

		if code == errmsg.SUCCESS {//是自己
			model.EditUser(id, &data)
		}
		if code == errmsg.ERROR_USER_WITH_WRONG_ID{//不是自己	
			c.Abort()
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 删除用户
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var code int

	code = model.DeleteUser(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}