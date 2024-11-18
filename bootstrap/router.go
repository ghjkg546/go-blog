package bootstrap

import (
	"context"
	"github.com/gin-gonic/gin"
	adminapi2 "github.com/jassue/jassue-gin/app/controllers/adminapi"
	"github.com/jassue/jassue-gin/app/controllers/app"
	"github.com/jassue/jassue-gin/app/controllers/common"
	"github.com/jassue/jassue-gin/app/middleware"
	"github.com/jassue/jassue-gin/app/services"
	"github.com/jassue/jassue-gin/global"
	"github.com/jassue/jassue-gin/routes"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//// Define a struct to hold your resource routes
//type ResourceRoutes struct {
//	BasePath   string
//	Router     *gin.Engine
//	Controller interface{}
//}

// Controller defines the methods that must be implemented by controllers
//
//	type Controller interface {
//		GetList(c *gin.Context)
//		GetDetail(c *gin.Context)
//		Create(c *gin.Context)
//		Update(c *gin.Context)
//		Delete(c *gin.Context)
//	}
//
//	func NewResourceRoutes(basePath string, router *gin.Engine, controller interface{}) *ResourceRoutes {
//		return &ResourceRoutes{
//			BasePath:   basePath,
//			Router:     router,
//			Controller: controller,
//		}
//	}
//
//	func (r *ResourceRoutes) SetupRoutes() {
//		resource := r.BasePath
//		//ctrlType := reflect.TypeOf(r.Controller)
//		ctrlValue := reflect.ValueOf(r.Controller)
//
//		// Map HTTP methods to controller methods
//		methodMappings := map[string]string{
//			"LIST":   "GetList",
//			"GET":    "GetDetail",
//			"POST":   "Create",
//			"PUT":    "Update",
//			"DELETE": "Delete",
//		}
//
//		for method, action := range methodMappings {
//			methodFunc := ctrlValue.MethodByName(action)
//			if !methodFunc.IsValid() {
//				continue
//			}
//
//			switch method {
//			case "LIST":
//				r.Router.GET(resource, func(c *gin.Context) {
//					methodFunc.Call([]reflect.Value{reflect.ValueOf(c)})
//				}).Use(middleware.Cors())
//			case "GET":
//				r.Router.GET(resource+"/:id", func(c *gin.Context) {
//					methodFunc.Call([]reflect.Value{reflect.ValueOf(c)})
//				}).Use(middleware.Cors())
//			case "POST":
//				r.Router.POST(resource, func(c *gin.Context) {
//					methodFunc.Call([]reflect.Value{reflect.ValueOf(c)})
//				}).Use(middleware.Cors())
//			case "PUT":
//				r.Router.PUT(resource+"/:id", func(c *gin.Context) {
//					methodFunc.Call([]reflect.Value{reflect.ValueOf(c)})
//				}).Use(middleware.Cors())
//			case "DELETE":
//				r.Router.DELETE(resource+"/:id", func(c *gin.Context) {
//					methodFunc.Call([]reflect.Value{reflect.ValueOf(c)})
//				}).Use(middleware.Cors())
//			}
//		}
//	}
func seq(start, end int) []int {
	var numbers []int
	for i := start; i <= end; i++ {
		numbers = append(numbers, i)
	}
	return numbers
}

func subtract(a, b int) int {
	return a - b
}

func add(a, b int) int {
	return a + b
}

func setupRouter() *gin.Engine {
	if global.App.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Logger(), middleware.CustomRecovery())
	router.SetFuncMap(template.FuncMap{
		"seq":      seq,
		"subtract": subtract,
		"add":      add,
	})
	router.LoadHTMLGlob("templates/*")
	router.Static("/css", "static/css")
	router.Static("/images", "static/images")
	router.Static("/js", "static/js")
	router.StaticFile("/sitemap.xml", "static/sitemap.xml")

	ResController := app.ResController{}
	router.GET("/", ResController.GetFrontReasouceItems)
	router.GET("/:cateid", ResController.GetFrontReasouceItems)
	router.GET("/category/:category_id/:page", ResController.GetFrontReasouceItems)
	router.GET("/category/:category_id", ResController.GetFrontReasouceItems)
	router.GET("/archives/:id", ResController.GetBlogDetail)
	router.GET("/qiu", ResController.ApplyForReasource)
	// 跨域处理
	router.Use(middleware.Cors())
	//userRoutes := NewResourceRoutes("/adminapi/users", router, &adminapi2.UserController{})
	//userRoutes.SetupRoutes()
	//
	//appuserRoutes := NewResourceRoutes("/adminapi/appusers", router, &adminapi2.AppUserController{})
	//appuserRoutes.SetupRoutes()
	//
	//cateRoutes := NewResourceRoutes("/adminapi/category", router, &adminapi2.CateController{})
	//cateRoutes.SetupRoutes()
	//
	//menuRoutes := NewResourceRoutes("/adminapi/menus", router, &adminapi2.MenuController{})
	//menuRoutes.SetupRoutes()
	//
	//resourceRoutes := NewResourceRoutes("/adminapi/resource", router, &adminapi2.ResourceController{})
	//resourceRoutes.SetupRoutes()
	//
	//crawlRoutes := NewResourceRoutes("/adminapi/crawl", router, &adminapi2.CrawlController{})
	//crawlRoutes.SetupRoutes()
	//
	//dictRoutes := NewResourceRoutes("/adminapi/dict", router, &adminapi2.DictController{})
	//dictRoutes.SetupRoutes()
	//
	//logRoutes := NewResourceRoutes("/adminapi/log", router, &adminapi2.LogController{})
	//logRoutes.SetupRoutes()

	ResourceController := adminapi2.ResourceController{}
	router.GET("/adminapi/share/waitlist", ResourceController.WaitShareList)
	router.POST("/adminapi/share/doshare", app.DoShare)
	router.POST("/adminapi/share/crawl", ResourceController.Crawl)

	crawController := adminapi2.CrawlController{}

	router.POST("/adminapi/crawl/savetodisk", crawController.BatchSaveToDisk)
	// 前端项目静态资源
	router.StaticFile("/staticfile", "./static/dist/index.html")
	router.Static("/assets", "./static/dist/assets")
	router.StaticFile("/favicon.ico", "./static/dist/favicon.ico")
	// 其他静态资源
	router.Static("/public", "./static")
	router.Static("/storage", "./storage/app/public")

	// 注册 api 分组路由
	apiGroup := router.Group("/api")
	apiGroup.Use(middleware.Cors())
	routes.SetApiGroupRoutes(apiGroup)

	// 注册 adminapi 分组路由
	admimApiGroup := router.Group("/adminapi")
	admimApiGroup.Use(middleware.Cors())
	routes.SetAdminApiGroupRoutes(admimApiGroup)

	//upload route
	authRouter := router.Group("/adminapi").Use(middleware.JWTAuth(services.AppGuardName))
	{
		authRouter.POST("/image_upload", common.ImageUpload)
		authRouter.POST("/csv_upload", common.CsvUpload)
		authRouter.POST("/user/password", func(ctx *gin.Context) {
			hello := adminapi2.UserController{}
			hello.ChangePass(ctx)
		})
	}

	// 简单的路由组: v1
	adminapi := router.Group("/adminapi/")
	{
		adminapi.GET("user/logout",
			//middleware.CorsMiddleware(),
			func(ctx *gin.Context) {
				hello := adminapi2.UserController{}
				hello.Logout(ctx)
			},
		)

		adminapi.GET("menu/options",
			//middleware.CorsMiddleware(),
			func(ctx *gin.Context) {
				hello := adminapi2.MenuController{}
				hello.Options(ctx)
			},
		)

		adminapi.GET("category/alllist",
			//middleware.CorsMiddleware(),
			func(ctx *gin.Context) {
				hello := adminapi2.CateController{}
				hello.AllCateList(ctx)
			},
		)

		adminapi.POST("user/login",
			//middleware.CorsMiddleware(),
			func(ctx *gin.Context) {
				hello := adminapi2.UserController{}
				hello.Login(ctx)
			},
		)
		adminapi.GET("user/me",
			//middleware.CorsMiddleware(),
			func(ctx *gin.Context) {
				hello := adminapi2.UserController{}
				hello.Me(ctx)
			},
		)

		adminapi.GET("role/list",
			//middleware.CorsMiddleware(),
			func(ctx *gin.Context) {
				hello := adminapi2.UserController{}
				hello.RoleList(ctx)
			},
		)

		adminapi.POST("resource/batchCreate",
			//middleware.CorsMiddleware(),
			func(ctx *gin.Context) {
				resouce := adminapi2.ResourceController{}
				resouce.BatchCreate(ctx)
			},
		)

		adminapi.POST("resource/batchShare",
			func(ctx *gin.Context) {
				resource := adminapi2.ResourceController{}
				resource.BatchShare(ctx)
			},
		)

		adminapi.GET("resource/syncToSearch",
			//middleware.CorsMiddleware(),
			func(ctx *gin.Context) {
				hello := adminapi2.ResourceController{}
				hello.SyncToSearch(ctx)
			},
		)

	}
	adminapi.Use(middleware.Cors())

	return router
}

func RunServer() {
	r := setupRouter()

	srv := &http.Server{
		Addr:    ":" + global.App.Config.App.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
