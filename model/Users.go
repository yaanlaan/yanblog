package model

import (
	"yanblog/utils/errmsg"

	"log"
	"golang.org/x/crypto/scrypt"
	"encoding/base64"
	"gorm.io/gorm"
)

// User 用户模型结构体
type User struct {
	
	gorm.Model
	Username string `gorm:"type:varchar(20);not null " json:"username" validate:"required,min=4,max=12" label:"用户名"`  // 用户名（4-12位）
	Password string `gorm:"type:varchar(20);not null" json:"password" validate:"required,min=6,max=20" label:"密码"`     // 密码（6-20位）
	Role     int    `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`                     // 角色码（大于等于2）
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
    if err != nil || user.Username != name{
        // 用户不存在或查询出错
        return errmsg.ERROR_USER_WITH_WRONG_ID
    }
    
    // 比较用户名是否一致
    if user.Username == name {
        return errmsg.SUCCESS
    }
    
    return errmsg.ERROR
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
// 参数: pageSize - 每页数量，pageNum - 页码
// 返回: 用户列表和总记录数
func GetUsers(pageSize int, pageNum int) ([]User, int64) {
	var users []User
	var total int64
	var err error

	if pageSize == -1 || pageNum == -1 {
		err = db.Find(&users).Count(&total).Error
	} else {
		err = db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Count(&total).Error
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
	

	err = db.Model(&user).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// DeleteUser 删除用户（软删除）
// 参数: id - 要删除的用户ID
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
	u.Password = ScryptPw(u.Password)
	return
}

// ScryptPw 使用scrypt算法对密码进行加密
// 参数: password - 明文密码
// 返回: 加密后的密码
func ScryptPw(password string) string {
	const KeyLen = 10                     // 密钥长度
	salt := make([]byte, 8)               // 盐值
	salt = []byte{12, 32, 4, 6, 66, 22, 222, 11}  // 固定盐值

	// 使用scrypt算法生成密钥
	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
	}
	
	fpw := base64.StdEncoding.EncodeToString(HashPw)
	return fpw
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

	if ScryptPw(password) != user.Password {
		return errmsg.ERROR_PASSWORD_WRONG
	}

	// 检查用户权限
	if user.Role != 2 {
		return errmsg.ERROR_USER_NO_RIGHT
	}
	
	return errmsg.SUCCESS
}