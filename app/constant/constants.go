package constants

import (
	"fmt"
	"github.com/jassue/jassue-gin/global"
)

// 定义 Redis 收藏键的格式
const (
	FavKeyPattern = "%s:user:%s:fav" // %s 占位符将替换为 AppName 和 userID
)

// 获取 Redis 收藏键
func GetFavKey(userID string) string {
	return fmt.Sprintf(FavKeyPattern, global.App.Config.App.AppName, userID)
}

// 获取 Redis 签到键
func GetSignKey(userID string, week string) string {
	return fmt.Sprintf("%s:user_sign:%s:%s", global.App.Config.App.AppName, userID, week)
}
