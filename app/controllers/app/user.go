package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/app/services"
	"github.com/jassue/jassue-gin/global"
	"net/http"
	"strings"
)

func UserInfo(c *gin.Context) {
	tokenString := getTokenFromHeader(c)
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
		return
	}
	token, err := jwt.ParseWithClaims(tokenString, &services.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.App.Config.Jwt.Secret), nil
	})
	if err != nil {
		response.BusinessFail(c, "解析失败"+err.Error())
		return
	}
	// Check if the token is valid
	if claims, ok := token.Claims.(*services.CustomClaims); ok && token.Valid {
		// Access user info from claims
		userID := claims.ID
		db := global.App.DB

		query := db.Model(models.User{})
		var user models.User
		query.Where("id=?", userID).Find(&user)

		response.Success(c, gin.H{"user_id": userID, "username": user.Name})
		return
	}
	response.BusinessFail(c, "解析失败"+err.Error())

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

func Fav(c *gin.Context) {
	var form request.Fav
	if err := c.ShouldBindJSON(&form); err != nil {
		response.BusinessFail(c, "参数格式有误222")
		return
	}
	uid := services.AppUserService.GetUserId(c)
	if uid == "" {
		response.BusinessFail(c, "未登录或登录过期")
		return
	}
	res, _ := AddToFavorites(services.AppUserService.GetUserId(c), form.Id)
	response.Success(c, res)
}

// 添加商品到用户收藏
func AddToFavorites(userID string, itemID int) (bool, error) {
	key := fmt.Sprintf(global.App.Config.App.AppName+":user:%s:fav", userID)
	ctx := context.Background()

	// 检查 itemID 是否已经存在于集合中
	isMember, _ := global.App.Redis.SIsMember(ctx, key, itemID).Result()

	if isMember {
		// 如果已存在，则移除，并返回 false（取消收藏）
		err := global.App.Redis.SRem(ctx, key, itemID).Err()
		if err != nil {
			return false, err
		}
		return false, nil
	} else {
		// 如果不存在，则添加，并返回 true（已收藏）
		err := global.App.Redis.SAdd(ctx, key, itemID).Err()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	//key := fmt.Sprintf(global.App.Config.App.AppName+":user:%s:fav", userID)
	//return global.App.Redis.SAdd(context.Background(), key, itemID).Err()
}

// 检查商品是否被用户收藏
func IsFavorite(userID, itemID string) (bool, error) {
	key := fmt.Sprintf(global.App.Config.App.AppName+":user:%s:fav", userID)
	return global.App.Redis.SIsMember(context.Background(), key, itemID).Result()
}

// 获取用户的所有收藏商品
func GetFavoriteItems(redisClient *redis.Client, userID string) ([]string, error) {
	key := fmt.Sprintf(global.App.Config.App.AppName+":user:%s:fav", userID)
	return redisClient.SMembers(context.Background(), key).Result()
}
