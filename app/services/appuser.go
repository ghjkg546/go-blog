package services

import (
	"errors"
	"fmt"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
	"github.com/jassue/jassue-gin/utils"
	"strconv"
)

type appUserService struct {
}

var AppUserService = new(appUserService)

// Register 注册
func (appUserService *appUserService) Register(params request.Register) (err error, user models.User) {
	var result = global.App.DB.Where("mobile = ?", params.UserName).Select("id").First(&models.User{})
	if result.RowsAffected != 0 {
		err = errors.New("用户名已存在")
		return
	}
	user = models.User{UserName: params.UserName, Name: params.UserName, Mobile: params.UserName, Password: utils.BcryptMake([]byte(params.Password))}
	err = global.App.DB.Create(&user).Error
	return
}

// Login 登录
func (appUserService *appUserService) Login(params request.Login) (err error, user *models.User) {
	err = global.App.DB.Where("username = ?", params.UserName).First(&user).Error
	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("数据库出错:" + err.Error())
	}
	fmt.Println(user.Password)
	fmt.Println(params)
	if !utils.BcryptMakeCheck([]byte(params.Password), user.Password) {
		err = errors.New("用户名不存在或密码错误")
	}

	return
}

// GetUserInfo 获取用户信息
func (appUserService *appUserService) GetUserInfo(id string) (err error, user models.AdminUser) {
	intId, err := strconv.Atoi(id)
	err = global.App.DB.First(&user, intId).Error
	if err != nil {
		err = errors.New("数据不存在")
	}
	return
}
