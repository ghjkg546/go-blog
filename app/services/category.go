package services

import (
	"errors"
	"fmt"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
)

type categoryService struct {
}

var CategoryService = new(categoryService)

func (categoryService *categoryService) GetList() (err error, data []models.Category, total int64) {

	if err := global.App.DB.Model(models.Category{}).Count(&total).Where("parent_id=? and status=?", 0, 1).
		Order("id DESC").Find(&data).Error; err != nil {
		fmt.Println(err.Error())
		err = errors.New("数据为空")

	}

	return
}
