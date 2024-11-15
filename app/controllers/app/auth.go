package app

import (
	"bytes"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/services"
	"net/http"
	"path"
	"strings"
	"time"
)

func Register(c *gin.Context) {
	var form request.Register
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if !captcha.VerifyString(form.CaptchaId, form.CaptchaValue) {
		response.BusinessFail(c, "验证码错误")
		return
	}
	if err, user := services.UserService.Register(form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		tokenData, err, _ := services.JwtService.CreateToken(services.AppGuardName, user)
		if err != nil {
			response.BusinessFail(c, err.Error())
			return
		}
		response.Success(c, tokenData)
	}
}

func Login(c *gin.Context) {
	var form request.Login
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err, user := services.AppUserService.Login(form); err != nil {
		fmt.Println("登录失败")
		response.BusinessFail(c, err.Error())
	} else {
		fmt.Println(user)
		tokenData, err, _ := services.JwtService.CreateToken(services.AppGuardName, user)
		if err != nil {
			response.Fail(c, 500, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "",
			"data": gin.H{
				"accessToken":  tokenData,
				"tokenType":    "Bearer",
				"refreshToken": nil,
				"username":     "admin",
				"role":         "user",
				"roleId":       1,
				"permissions":  []string{"*.*.*"},
			},
		})
	}
}

func Logout(c *gin.Context) {
	err := services.JwtService.JoinBlackList(c.Keys["token"].(*jwt.Token))
	if err != nil {
		response.BusinessFail(c, "登出失败")
		return
	}
	response.Success(c, nil)
}

type CaptchaResponse struct {
	CaptchaID string `json:"captcha_id"`
	ImageURL  string `json:"image_url"`
}

// 生成图形验证码
func GenerateCaptcha(c *gin.Context) {
	captchaID := captcha.NewLen(4)
	imageURL := "/captcha/" + captchaID + ".png"

	res := CaptchaResponse{
		CaptchaID: captchaID,
		ImageURL:  imageURL,
	}
	response.Success(c, res)

}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir, file := path.Split(r.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	fmt.Println("file : " + file)
	fmt.Println("ext : " + ext)
	fmt.Println("id : " + id)
	if ext == "" || id == "" {
		http.NotFound(w, r)
		return
	}

	if r.FormValue("reload") != "" {
		captcha.Reload(id)
	}
	lang := strings.ToLower(r.FormValue("lang"))
	download := path.Base(dir) == "download"
	if Serve(w, r, id, ext, lang, download, captcha.StdWidth, captcha.StdHeight) == captcha.ErrNotFound {
		http.NotFound(w, r)
	}
}

func Serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil

}

// 获取验证码图片
func CaptchaImage(c *gin.Context) {
	captchaId := c.Param("captchaID")
	fmt.Println("GetCaptchaPng : " + captchaId)

	ServeHTTP(c.Writer, c.Request)
	//captchaID := c.Param("captchaID")
	//if captchaID == "" {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Missing captcha ID"})
	//	return
	//}
	//
	//// 设置响应头
	//c.Header("Content-Type", "image/png")
	//err := captcha.WriteImage(c.Writer, captchaID, 120, 40)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate captcha image"})
	//}
}

// 验证用户输入的验证码
func VerifyCaptcha(c *gin.Context) {

	var form request.Captcha
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if form.CaptchaID == "" || form.CaptchaValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing captcha ID or value"})
		return
	}

	// 验证验证码
	if captcha.VerifyString(form.CaptchaID, form.CaptchaValue) {
		c.JSON(http.StatusOK, gin.H{"verified": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"verified": false})
	}
}
