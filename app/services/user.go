package services

import (
	"errors"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
	"github.com/jassue/jassue-gin/utils"
	"strconv"
)

type userService struct {
}

var UserService = new(userService)

// Register 注册
func (userService *userService) Register(params request.Register) (err error, user models.User) {
	var result = global.App.DB.Where("mobile = ? OR email = ?", params.UserName, params.Email).Select("id").First(&models.User{})
	if result.RowsAffected != 0 {
		err = errors.New("用户名或邮箱已存在")
		return
	}
	user = models.User{UserName: params.UserName, Name: params.UserName, Mobile: params.UserName, Email: params.Email, Password: utils.BcryptMake([]byte(params.Password))}
	err = global.App.DB.Create(&user).Error
	return
}

// Login 登录
func (userService *userService) Login(params request.Login) (err error, user *models.AdminUser) {
	err = global.App.DB.Where("name = ?", params.UserName).First(&user).Error
	if err != nil || !utils.BcryptMakeCheck([]byte(params.Password), user.Password) {
		err = errors.New("用户名不存在或密码错误:" + err.Error())
	}
	return
}

// GetUserInfo 获取用户信息
func (userService *userService) GetUserInfo(id string) (err error, user models.AdminUser) {
	intId, err := strconv.Atoi(id)
	err = global.App.DB.First(&user, intId).Error
	if err != nil {
		err = errors.New("数据不存在")
	}
	return
}

func (userService *userService) GetUserId(id string) (err error, user models.AdminUser) {
	intId, err := strconv.Atoi(id)
	err = global.App.DB.First(&user, intId).Error
	if err != nil {
		err = errors.New("数据不存在")
	}
	return
}
