package app

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
