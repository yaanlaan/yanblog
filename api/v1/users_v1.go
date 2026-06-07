package v1

import (
	"yanblog/model"
	"yanblog/utils/errmsg"
	"yanblog/utils/validator"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 添加用户
func AddUser(c *gin.Context) {
	var data model.User
	var msg string
	var code int

	_ = c.ShouldBindJSON(&data)

	// 获取当前操作用户的用户名和角色
	currentUsername, _ := c.Get("username")
	currentUserRole := model.GetUserRole(currentUsername.(string))

	// 权限检查
	if currentUserRole == 1 {
		// 超级管理员可以创建管理员(2)和普通用户(3)
		if data.Role == 1 {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR_USER_NO_RIGHT,
				"message": "无权创建超级管理员",
			})
			return
		}
	} else if currentUserRole == 2 {
		// 管理员只能创建普通用户(3)
		if data.Role != 3 {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR_USER_NO_RIGHT,
				"message": "管理员只能创建普通用户",
			})
			return
		}
	} else {
		// 普通用户无权创建用户
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR_USER_NO_RIGHT,
			"message": "无权创建用户",
		})
		return
	}

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
		"status": code,
		// "data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询用户列表
func GetUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	var code int

	if pageSize <= 0 {
		pageSize = -1
	}
	if pageNum <= 0 {
		pageNum = -1
	}

	// 获取当前操作用户的用户名和角色
	currentUsername, _ := c.Get("username")
	currentUserRole := model.GetUserRole(currentUsername.(string))

	data, total := model.GetUsers(pageSize, pageNum, currentUserRole)
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

	if pageSize <= 0 {
		pageSize = -1
	}
	if pageNum <= 0 {
		pageNum = -1
	}

	// 获取当前操作用户的用户名和角色
	currentUsername, _ := c.Get("username")
	currentUserRole := model.GetUserRole(currentUsername.(string))

	data, total := model.SearchUser(keyword, role, pageSize, pageNum, currentUserRole)
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}
	_ = c.ShouldBindJSON(&data)

	// 调试日志
	// fmt.Printf("EditUser received data: %+v\n", data)

	// 获取当前操作用户的用户名和角色
	currentUsername, _ := c.Get("username")
	currentUserRole := model.GetUserRole(currentUsername.(string))

	// 获取目标用户的信息
	var targetUser model.User
	model.GetDB().Where("id = ?", id).First(&targetUser)

	// 权限检查
	// 1. 超级管理员(1)可以修改任何人，但不能修改自己的角色
	// 2. 管理员(2)可以修改自己和普通用户(3)，不能修改超级管理员(1)和其他管理员(2)
	// 3. 普通用户(3)只能修改自己，且不能修改角色

	// 禁止将用户角色修改为超级管理员(1)
	if data.Role == 1 && targetUser.Role != 1 {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR_USER_NO_RIGHT,
			"message": "无权将用户提升为超级管理员",
		})
		return
	}

	// 禁止修改其他超级管理员的角色（防止降权）
	if targetUser.Role == 1 && targetUser.Username != currentUsername.(string) {
		if data.Role != 1 {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR_USER_NO_RIGHT,
				"message": "无权修改其他超级管理员的角色",
			})
			return
		}
	}

	if currentUserRole == 1 {
		// 超级管理员
		if targetUser.Username == currentUsername.(string) {
			// 修改自己，强制保持超级管理员角色（防止降权）
			data.Role = 1
		}
		// 可以修改其他任何数据
	} else if currentUserRole == 2 {
		// 管理员
		if targetUser.Role == 1 {
			// 不能修改超级管理员
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR_USER_NO_RIGHT,
				"message": "无权修改超级管理员信息",
			})
			return
		}
		if targetUser.Role == 2 && targetUser.Username != currentUsername.(string) {
			// 不能修改其他管理员
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR_USER_NO_RIGHT,
				"message": "无权修改其他管理员信息",
			})
			return
		}
		// 修改自己或普通用户，不能修改角色
		data.Role = targetUser.Role
	} else {
		// 普通用户
		if targetUser.Username != currentUsername.(string) {
			// 不能修改他人
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR_USER_NO_RIGHT,
				"message": "无权修改他人信息",
			})
			return
		}
		// 修改自己，不能修改角色
		data.Role = targetUser.Role
	}

	code = model.CheckUser(data.Username)
	if code == errmsg.SUCCESS { //用户名不在
		model.EditUser(id, &data)
	}
	if code == errmsg.ERROR_USERNAME_USED { //用户名已存在
		code = model.CheckUserWithID(id, data.Username) //判断是不是自己

		if code == errmsg.SUCCESS { //是自己
			model.EditUser(id, &data)
		}
		if code == errmsg.ERROR_USER_WITH_WRONG_ID { //不是自己
			// 这里不需要Abort，直接返回JSON即可，因为Abort通常用于中间件
			// c.Abort()
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 删除用户
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  errmsg.ERROR,
			"message": "参数错误",
		})
		return
	}

	// 获取当前操作用户的用户名和角色
	currentUsername, _ := c.Get("username")
	currentUserRole := model.GetUserRole(currentUsername.(string))

	// 获取目标用户的信息
	var targetUser model.User
	model.GetDB().Where("id = ?", id).First(&targetUser)

	// 权限检查
	if currentUserRole == 1 {
		// 超级管理员不能删除自己
		if targetUser.Username == currentUsername.(string) {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR_USER_NO_RIGHT,
				"message": "超级管理员不能删除自己",
			})
			return
		}
		// 超级管理员可以删除其他任何人
	} else if currentUserRole == 2 {
		// 管理员只能删除普通用户
		if targetUser.Role != 3 {
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR_USER_NO_RIGHT,
				"message": "无权删除该用户",
			})
			return
		}
	} else {
		// 普通用户不能删除任何人
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR_USER_NO_RIGHT,
			"message": "无权执行删除操作",
		})
		return
	}

	var code int

	code = model.DeleteUser(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
