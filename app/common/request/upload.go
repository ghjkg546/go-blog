package request

import "mime/multipart"

type ImageUpload struct {
	Business string                `form:"business" json:"business" binding:"required"`
	Image    *multipart.FileHeader `form:"image" json:"image" binding:"required"`
}

type FileUpload struct {
	Business   string                `form:"business" json:"business" binding:"required"`
	Image      *multipart.FileHeader `form:"image" json:"image" binding:"required"`
	CategoryId uint                  `form:"category_id" json:"category_id"`
}

func (imageUpload ImageUpload) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"business.required": "业务类型不能为空",
		"image.required":    "请选择图片",
	}
}
