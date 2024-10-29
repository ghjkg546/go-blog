package services

import (
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
)

type dictService struct {
}

var DictService = new(dictService)

func (dictService *dictService) GetValueByDict(code string, name string) (value string) {
	db := global.App.DB
	var dict models.SysDict
	db.Model(models.SysDict{}).Preload("DictItems").Where("code = ?", code).First(&dict)
	for _, item := range dict.DictItems {
		if item.Name == name {
			return item.Value
		}
	}
	return ""
	return
}
