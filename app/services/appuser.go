package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jassue/jassue-gin/app/common/request"
	constants "github.com/jassue/jassue-gin/app/constant"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
	"github.com/jassue/jassue-gin/utils"
	"strconv"
	"strings"
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

func getTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// The header should be in the format "Bearer <token>"
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || splitToken[0] != "Bearer" {
		return ""
	}

	return splitToken[1]
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

// GetUserInfo 获取用户信息
func (appUserService *appUserService) GetUserId(c *gin.Context) string {
	tokenString := getTokenFromHeader(c)
	if tokenString == "" {
		return ""
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.App.Config.Jwt.Secret), nil
	})
	if err != nil {
		//response.BusinessFail(c, "解析失败"+err.Error())
		return ""
	}
	// Check if the token is valid
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// Access user info from claims
		userID := claims.ID
		return userID
	}
	return ""
}

// 检查是否被用户收藏
func (appUserService *appUserService) IsFavorite(c *gin.Context, itemID string) (bool, error) {
	userID := appUserService.GetUserId(c)

	return global.App.Redis.SIsMember(context.Background(), constants.GetFavKey(userID), itemID).Result()
}

// 获取所有收藏id
func (appUserService *appUserService) GetFavIds(c *gin.Context) ([]string, error) {
	userID := appUserService.GetUserId(c)
	if userID == "" {
		return []string{}, nil
	}
	ctx := context.Background()
	return global.App.Redis.SMembers(ctx, constants.GetFavKey(userID)).Result()
}
