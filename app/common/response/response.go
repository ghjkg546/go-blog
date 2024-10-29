package response

import (
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
	"net/http"
	"os"
)

type Response struct {
	ErrorCode int         `json:"code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"msg"`
}

type DirResponse struct {
	Status int               `json:"status"`
	Code   int               `json:"code"`
	Data   []models.FileInfo `json:"data"`
	Total  int               `json:"total"`
}

func ServerError(c *gin.Context, err interface{}) {
	msg := "Internal Server Error"
	if global.App.Config.App.Env != "production" && os.Getenv(gin.EnvGinMode) != gin.ReleaseMode {
		if _, ok := err.(error); ok {
			msg = err.(error).Error()
		}
	}
	c.JSON(http.StatusInternalServerError, Response{
		http.StatusInternalServerError,
		nil,
		msg,
	})
	c.Abort()
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		0,
		data,
		"ok",
	})
}

func Fail(c *gin.Context, errorCode int, msg string) {
	c.JSON(http.StatusOK, Response{
		errorCode,
		nil,
		msg,
	})
}

func FailByError(c *gin.Context, error global.CustomError) {
	Fail(c, error.ErrorCode, error.ErrorMsg)
}

func ValidateFail(c *gin.Context, msg string) {
	Fail(c, global.Errors.ValidateError.ErrorCode, msg)
}

func BusinessFail(c *gin.Context, msg string) {
	Fail(c, global.Errors.BusinessError.ErrorCode, msg)
}

func TokenFail(c *gin.Context) {
	FailByError(c, global.Errors.TokenError)
}
