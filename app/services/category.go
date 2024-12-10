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

func (categoryService *categoryService) GetList() (err error, parentCategories []models.Category, total int64) {

	//var parentCategories []models.Category
	db := global.App.DB
	// 查询顶层分类
	if err := db.Where("parent_id = ? AND status = ?", 0, 1).
		Order("id DESC").
		Find(&parentCategories).Error; err != nil {
		err = errors.New("数据为空")
	}

	// 查询每个顶层分类的子分类
	for i, parent := range parentCategories {
		var children []models.Category
		if err := db.Where("parent_id = ? AND status = ?", parent.ID, 1).
			Order("id DESC").
			Find(&children).Error; err != nil {
			// 如果某个分类没有子分类，可以继续操作（非致命错误）
			children = []models.Category{}
		}
		parentCategories[i].Children = children
	}

	return
}

func (categoryService *categoryService) GetWithChildList() (err error, parentCategories []models.Category, total int64) {
	db := global.App.DB
	// 查询顶层分类
	if err := db.Model(&parentCategories).Where("status = ?", 1).Count(&total).
		Order("id DESC").
		Find(&parentCategories).Error; err != nil {
		fmt.Println(err.Error())
		err = errors.New("数据为空")
	}

	return
}

func (categoryService *categoryService) GetListByPid(pid int32) (err error, data []models.Category) {

	if err := global.App.DB.Model(models.Category{}).Where("parent_id=? and status=?", pid, 1).
		Order("id DESC").Find(&data).Error; err != nil {
		err = errors.New("数据为空")

	}

	return
}
