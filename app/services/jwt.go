package services

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/global"
	"github.com/jassue/jassue-gin/utils"
	"strconv"
	"time"
)

type jwtService struct {
}

var JwtService = new(jwtService)

type JwtUser interface {
	GetUid() string
}

// CustomClaims 自定义 Claims
type CustomClaims struct {
	jwt.RegisteredClaims
}

const (
	TokenType    = "bearer"
	AppGuardName = "app"
)

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (jwtService *jwtService) CreateToken(GuardName string, user JwtUser) (tokenData TokenOutPut, err error, token *jwt.Token) {
	token = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.App.Config.Jwt.JwtTtl) * time.Second)),
				ID:        user.GetUid(),
				Issuer:    GuardName,
				NotBefore: jwt.NewNumericDate(time.Now().Add(-1000 * time.Second)),
			},
		},
	)

	tokenStr, err := token.SignedString([]byte(global.App.Config.Jwt.Secret))

	tokenData = TokenOutPut{
		tokenStr,
		int(global.App.Config.Jwt.JwtTtl),
		TokenType,
	}
	return
}

func (jwtService *jwtService) getBlackListKey(tokenStr string) string {
	return "jwt_black_list:" + utils.MD5([]byte(tokenStr))
}

func (jwtService *jwtService) JoinBlackList(token *jwt.Token) (err error) {
	nowUnix := time.Now().Unix()
	timer := time.Duration(token.Claims.(*CustomClaims).ExpiresAt.Unix()-nowUnix) * time.Second
	err = global.App.Redis.SetNX(context.Background(), jwtService.getBlackListKey(token.Raw), nowUnix, timer).Err()
	return
}

func (jwtService *jwtService) IsInBlacklist(tokenStr string) bool {
	return false
	//joinUnixStr, err := global.App.Redis.Get(context.Background(), jwtService.getBlackListKey(tokenStr)).Result()
	//joinUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)
	//if joinUnixStr == "" || err != nil {
	//	return false
	//}
	//if time.Now().Unix()-joinUnix < global.App.Config.Jwt.JwtBlacklistGracePeriod {
	//	return false
	//}
	//return true
}

func (jwtService *jwtService) GetUserInfo(GuardName string, id string) (err error, user JwtUser) {
	switch GuardName {
	case AppGuardName:
		return UserService.GetUserInfo(id)
	default:
		err = errors.New("guard " + GuardName + " does not exist")
	}
	return
}

func (jwtService *jwtService) GetUserId(c *gin.Context) int32 {
	uidsStr, ok := c.Get("id")
	if !ok {
		// Handle the case where the value is not a string
		response.Fail(c, 500, "获取uid出错")
	}

	str, ok := uidsStr.(string)
	if !ok {
		response.Fail(c, 500, "uid不是字符串")
		return 0
	}
	uidsInt, _ := stringToInt32(str)

	return uidsInt

}

func stringToInt32(s string) (int32, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}
