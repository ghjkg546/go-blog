package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/controllers/adminapi"
	"github.com/jassue/jassue-gin/app/controllers/app"
	"github.com/jassue/jassue-gin/app/controllers/common"
	"github.com/jassue/jassue-gin/app/middleware"
	"github.com/jassue/jassue-gin/app/services"
)

func SetApiGroupRoutes(router *gin.RouterGroup) {
	router.POST("/auth/register", app.Register)
	router.POST("/auth/login", app.Login)
	router.GET("/user/info", app.UserInfo)
	router.POST("/user/fav", app.Fav)
	router.POST("/user/sign", app.SignIn)
	router.GET("/user/signstatus", app.GetSignStatus)
	router.GET("/user/favlist", app.GetFavList)
	router.GET("/duanju/list", app.GetResList)

	router.GET("/res/info", app.Info)
	router.GET("/res/search", app.ResSearch)

	router.GET("/generate-captcha", app.GenerateCaptcha)
	router.GET("/captcha/:captchaID", app.CaptchaImage)
	router.GET("/verify-captcha", app.VerifyCaptcha)

	router.GET("/category/list", app.CateList)

	router.GET("/wx", app.WXCheckSignature)
	router.POST("/wx", app.WXMsgReceive)
	router.POST("/user/avatar_upload", common.AvatarUpload)
	router.GET("/admin/user/loginOut", adminapi.LogOut)
	authRouter := router.Group("").Use(middleware.JWTAuth(services.AppGuardName))
	{
		authRouter.POST("/comment/add", app.CreateComment)
		authRouter.POST("/auth/info", app.Info)
		authRouter.POST("/auth/logout", app.Logout)
		authRouter.POST("/image_upload", common.ImageUpload)
	}
}
