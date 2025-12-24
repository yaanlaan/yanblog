package model

import (
	"strings"
	"yanblog/utils/errmsg"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型结构体
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null " json:"username" validate:"required,min=4,max=12" label:"用户名"` // 用户名（4-12位）
	Password string `gorm:"type:varchar(100);not null" json:"password" validate:"required,min=6,max=20" label:"密码"`  // 密码（6-20位）
	Role     int    `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=1" label:"角色码"`                    // 角色码（1:超级管理员, 2:管理员, 3:普通用户）
}

// CheckUser 检查用户名是否已存在
// 参数: name - 要检查的用户名
// 返回: 状态码（SUCCESS表示用户名可用，ERROR_USERNAME_USED表示用户名已被使用）
func CheckUser(name string) (code int) {
	var users User

	db.Select("id").Where("username = ?", name).First(&users)
	if users.ID > 0 {
		return errmsg.ERROR_USERNAME_USED
	}
	return errmsg.SUCCESS
}

// CheckUsernameWithID 检查指定ID的用户的用户名是否与传入的用户名一致
// 参数: id - 用户ID, name - 要比较的用户名
// 返回: 状态码（SUCCESS表示一致，ERROR表示不一致或用户不存在）
func CheckUserWithID(id int, name string) int {
	var user User

	// 查询指定ID的用户
	err := db.Select("username").Where("id = ?", id).First(&user).Error
	if err != nil || user.Username != name {
		// 用户不存在或查询出错
		return errmsg.ERROR_USER_WITH_WRONG_ID
	}

	// 比较用户名是否一致
	if user.Username == name {
		return errmsg.SUCCESS
	}

	return errmsg.ERROR
}

// SearchUser 搜索用户
// 参数: keyword - 搜索关键词, role - 角色筛选, pageSize - 每页数量, pageNum - 页码, currentRole - 当前用户角色
// 返回: 用户列表和总记录数
func SearchUser(keyword string, role int, pageSize int, pageNum int, currentRole int) ([]User, int64) {
	var users []User
	var total int64
	var err error

	// 构建查询条件
	query := db

	// 权限控制：低级用户无法看到高级用户
	if currentRole == 2 {
		// 管理员(2)只能看到管理员(2)和普通用户(3)
		query = query.Where("role >= ?", 2)
	} else if currentRole == 3 {
		// 普通用户(3)只能看到普通用户(3)
		query = query.Where("role = ?", 3)
	}
	// 超级管理员(1)可以看到所有用户

	// 如果有关键词，则添加用户名的模糊搜索
	if keyword != "" {
		searchTerm := "%" + strings.ToLower(keyword) + "%"
		query = query.Where("LOWER(username) LIKE ?", searchTerm)
	}

	// 如果有角色筛选条件，则添加角色筛选
	if role != 0 {
		query = query.Where("role = ?", role)
	}

	// 执行查询
	if pageSize == -1 || pageNum == -1 {
		err = query.Find(&users).Count(&total).Error
	} else {
		err = query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Count(&total).Error
	}

	// 处理错误情况
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}

	return users, total
}

// CreateUser 创建新用户
// 参数: data - 用户数据指针
// 返回: 状态码
func CreateUser(data *User) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// GetUsers 获取用户列表（支持分页）
// 参数: pageSize - 每页数量，pageNum - 页码, currentRole - 当前用户角色
// 返回: 用户列表和总记录数
func GetUsers(pageSize int, pageNum int, currentRole int) ([]User, int64) {
	var users []User
	var total int64
	var err error

	// 构建查询条件
	query := db

	// 权限控制：低级用户无法看到高级用户
	if currentRole == 2 {
		// 管理员(2)只能看到管理员(2)和普通用户(3)
		query = query.Where("role >= ?", 2)
	} else if currentRole == 3 {
		// 普通用户(3)只能看到普通用户(3)
		query = query.Where("role = ?", 3)
	}
	// 超级管理员(1)可以看到所有用户

	if pageSize == -1 || pageNum == -1 {
		err = query.Find(&users).Count(&total).Error
	} else {
		err = query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Count(&total).Error
	}

	// 处理错误情况
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0
	}

	return users, total
}

// EditUser 编辑用户信息
// 参数: id - 用户ID，data - 新的用户数据
// 返回: 状态码
func EditUser(id int, data *User) int {
	var user User

	var maps = make(map[string]interface{})
	maps["username"] = data.Username
	maps["role"] = data.Role

	// 如果密码不为空，则更新密码
	if data.Password != "" {
		// 密码加密在BeforeSave钩子中处理，但Updates方法不会触发BeforeSave钩子
		// 所以这里需要手动加密
		// fmt.Println("Updating password for user", id)
		maps["password"] = EncryptPassword(data.Password)
	} else {
		// fmt.Println("Password is empty, skipping update")
	}

	err = db.Model(&user).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// DeleteUser 删除用户
// 参数: id - 用户ID
// 返回: 状态码
func DeleteUser(id int) int {
	var user User
	err = db.Where("id = ? ", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// BeforeSave GORM钩子函数，在保存用户前自动执行
// 用于在保存用户前对密码进行加密
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.Password = EncryptPassword(u.Password)
	return
}

// EncryptPassword 使用bcrypt算法对密码进行加密
// 参数: password - 明文密码
// 返回: 加密后的密码
func EncryptPassword(password string) string {
	const cost = 10 // 负载值

	HashPw, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return password // 如果加密失败，返回原密码
	}

	return string(HashPw)
}

// CheckLogin 验证用户登录
// 参数: username - 用户名，password - 明文密码
// 返回: 状态码
func CheckLogin(username string, password string) int {
	var user User

	db.Where("username = ?", username).First(&user)

	if user.ID == 0 {
		return errmsg.ERROR_USER_NOT_EXIST
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errmsg.ERROR_PASSWORD_WRONG
	}

	// 检查用户权限
	if user.Role != 1 && user.Role != 2 {
		return errmsg.ERROR_USER_NO_RIGHT
	}

	return errmsg.SUCCESS
}

// GetUserRole 获取用户角色
func GetUserRole(username string) int {
	var user User
	db.Where("username = ?", username).First(&user)
	return user.Role
}
