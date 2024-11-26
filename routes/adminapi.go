package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/controllers/adminapi"
	"github.com/jassue/jassue-gin/app/middleware"
	"reflect"
)

// Define a struct to hold your resource routes
type ResourceRoutes struct {
	BasePath   string
	Router     *gin.RouterGroup
	Controller interface{}
}

type Controller interface {
	GetList(c *gin.Context)
	GetDetail(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

func NewResourceRoutes(basePath string, router *gin.RouterGroup, controller interface{}) *ResourceRoutes {
	return &ResourceRoutes{
		BasePath:   basePath,
		Router:     router,
		Controller: controller,
	}
}

func (r *ResourceRoutes) SetupRoutes() {
	resource := r.BasePath
	ctrlValue := reflect.ValueOf(r.Controller)

	// Map HTTP methods to controller methods
	methodMappings := map[string]string{
		"LIST":   "GetList",
		"GET":    "GetDetail",
		"POST":   "Create",
		"PUT":    "Update",
		"DELETE": "Delete",
	}

	for method, action := range methodMappings {
		methodFunc := ctrlValue.MethodByName(action)
		if !methodFunc.IsValid() {
			continue
		}

		switch method {
		case "LIST":
			r.Router.GET(resource, func(c *gin.Context) {
				methodFunc.Call([]reflect.Value{reflect.ValueOf(c)})
			}).Use(middleware.Cors())
		case "GET":
			r.Router.GET(resource+"/:id", func(c *gin.Context) {
				methodFunc.Call([]reflect.Value{reflect.ValueOf(c)})
			}).Use(middleware.Cors())
		case "POST":
			r.Router.POST(resource, func(c *gin.Context) {
				methodFunc.Call([]reflect.Value{reflect.ValueOf(c)})
			}).Use(middleware.Cors())
		case "PUT":
			r.Router.PUT(resource+"/:id", func(c *gin.Context) {
				methodFunc.Call([]reflect.Value{reflect.ValueOf(c)})
			}).Use(middleware.Cors())
		case "DELETE":
			r.Router.DELETE(resource+"/:id", func(c *gin.Context) {
				methodFunc.Call([]reflect.Value{reflect.ValueOf(c)})
			}).Use(middleware.Cors())
		}
	}
}

func SetAdminApiGroupRoutes(router *gin.RouterGroup) {
	userRoutes := NewResourceRoutes("/users", router, &adminapi.UserController{})
	userRoutes.SetupRoutes()

	commentRoutes := NewResourceRoutes("/comments", router, &adminapi.CommentController{})
	commentRoutes.SetupRoutes()

	appuserRoutes := NewResourceRoutes("/appusers", router, &adminapi.AppUserController{})
	appuserRoutes.SetupRoutes()

	cateRoutes := NewResourceRoutes("/category", router, &adminapi.CateController{})
	cateRoutes.SetupRoutes()

	menuRoutes := NewResourceRoutes("/menus", router, &adminapi.MenuController{})
	menuRoutes.SetupRoutes()

	resourceRoutes := NewResourceRoutes("/resource", router, &adminapi.ResourceController{})
	resourceRoutes.SetupRoutes()

	crawlRoutes := NewResourceRoutes("/crawl", router, &adminapi.CrawlController{})
	crawlRoutes.SetupRoutes()

	dictRoutes := NewResourceRoutes("/dict", router, &adminapi.DictController{})
	dictRoutes.SetupRoutes()

	logRoutes := NewResourceRoutes("/log", router, &adminapi.LogController{})
	logRoutes.SetupRoutes()

	feedbackRoutes := NewResourceRoutes("/feedbacks", router, &adminapi.FeedbackController{})
	feedbackRoutes.SetupRoutes()
}
