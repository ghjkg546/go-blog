package request

type Register struct {
	UserName     string `form:"username" json:"username" binding:"required"`
	Password     string `form:"password" json:"password" binding:"required"`
	Email        string `form:"email" json:"email"`
	CaptchaId    string `form:"captcha_id" json:"captcha_id" `
	CaptchaValue string `form:"captcha_value" json:"captcha_value" binding:"required"`
}

func (register Register) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"username.required":      "用户名不能为空",
		"password.required":      "用户密码不能为空",
		"captcha_value.required": "验证码不能为空",
	}
}

//type Login struct {
//	UserName string `form:"username" json:"username" binding:"required,username"`
//	Password string `form:"password" json:"password" binding:"required"`
//}

type Captcha struct {
	CaptchaID    string `form:"captcha_id" json:"captcha_id"`
	CaptchaValue string `form:"captcha_value" json:"captcha_value"`
}

type Login struct {
	UserName string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type Fav struct {
	Id int `form:"id" json:"id"`
}

type PostComment struct {
	Content      string `form:"content" json:"content"`
	ResItemId    int32  `form:"resource_item_id" json:"resource_item_id"`
	CaptchaId    string `form:"captcha_id" json:"captcha_id" `
	CaptchaValue string `form:"captcha_value" json:"captcha_value" binding:"required"`
}

type ChangePass struct {
	ConfirmPassword string `form:"confirmPassword" json:"confirmPassword"`
	NewPassword     string `form:"newPassword" json:"newPassword"`
	OldPassword     string `form:"oldPassword" json:"oldPassword"`
}

type EditUser struct {
	Id       int    `form:"id" json:"id"`
	UserName string `form:"username" json:"username"`
	Name     string `form:"name" json:"name"`
	Password string `form:"password" json:"password"`
	Mobile   string `form:"mobile" json:"mobile"` // 添加手机号码字段
	Status   int    `form:"status" json:"status"` // 添加状态字段
}
