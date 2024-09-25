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
	router.GET("/duanju/list", app.GetResList)
	router.GET("/res/info", app.Info)
	router.GET("/res/search", app.ResSearch)

	router.GET("/category/list", app.CateList)

	router.GET("/admin/user/loginOut", adminapi.LogOut)
	authRouter := router.Group("").Use(middleware.JWTAuth(services.AppGuardName))
	{
		authRouter.POST("/comment/add", app.CreateComment)
		authRouter.POST("/auth/info", app.Info)
		authRouter.POST("/auth/logout", app.Logout)
		authRouter.POST("/image_upload", common.ImageUpload)
	}
}
