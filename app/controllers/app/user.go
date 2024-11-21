package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/common/response"
	constants "github.com/jassue/jassue-gin/app/constant"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/app/services"
	"github.com/jassue/jassue-gin/global"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
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

		response.Success(c, gin.H{"user_id": userID, "username": user.Name, "score": user.Score, "avatar": user.Avatar})
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

// 检查商品是否被用户
func IsFavorite(userID, itemID string) (bool, error) {
	key := fmt.Sprintf(global.App.Config.App.AppName+":user:%s:fav", userID)
	return global.App.Redis.SIsMember(context.Background(), key, itemID).Result()
}

// 获取用户的所有收藏
func GetFavoriteItems(redisClient *redis.Client, userID string) ([]string, error) {
	key := fmt.Sprintf(global.App.Config.App.AppName+":user:%s:fav", userID)
	return redisClient.SMembers(context.Background(), key).Result()
}

// 获取当前的周次 (year-week)
func getCurrentWeek() string {
	now := time.Now()
	_, week := now.ISOWeek()
	return fmt.Sprintf("%d-%02d", now.Year(), week)
}

// 用户签到接口
func SignIn(c *gin.Context) {
	week := getCurrentWeek()
	day := time.Now().Weekday()

	uid := services.AppUserService.GetUserId(c)
	if uid == "" {
		response.BusinessFail(c, "未登录或登录过期")
		return
	}
	exist, err := global.App.Redis.GetBit(context.Background(), constants.GetSignKey(uid, week), int64(day)).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取签到状态失败"})
		return
	}
	if exist > 0 {
		response.BusinessFail(c, "今天你已经签到过了")
		return
	}
	_, err1 := global.App.Redis.SetBit(context.Background(), constants.GetSignKey(uid, week), int64(day), 1).Result()
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "签到失败"})
		return
	}
	_ = UpdateUserScore(uid, 100)

	response.Success(c, "签到成功")

}

// 更新用户分数，如果用户存在则加分
func UpdateUserScore(uid string, increment int64) error {
	var user models.User

	// 根据 UID 查找用户
	result := global.App.DB.Where("id = ?", uid).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found")
		}
		return result.Error
	}
	// 用户存在，增加分数
	user.Score += increment

	// 保存更新后的用户
	if err := global.App.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// 获取用户签到状态接口
func GetSignStatus(c *gin.Context) {

	week := getCurrentWeek()

	uid := services.AppUserService.GetUserId(c)
	if uid == "" {
		response.BusinessFail(c, "未登录或登录过期")
		return
	}
	var status []int
	for i := 0; i < 7; i++ {
		bit, err := global.App.Redis.GetBit(context.Background(), constants.GetSignKey(uid, week), int64(i)).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取签到状态失败"})
			return
		}
		status = append(status, int(bit))
	}
	response.Success(c, status)
}
