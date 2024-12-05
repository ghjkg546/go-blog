package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
	"github.com/jassue/jassue-gin/utils"
)

type resourceItemService struct {
}

var ResourceItemService = new(resourceItemService)

// 获取资源列表
func (resourceItemService *resourceItemService) GetResList(page int, pageSize int, categoryId int32, keyword string, orderBy string) (err error, data []models.ResourceItem, total int64) {

	offset := (page - 1) * pageSize
	query := global.App.DB.Model(models.ResourceItem{})
	if orderBy != "" {
		fmt.Println(orderBy)
		query.Order(orderBy + ",id DESC")
	} else {
		query.Order("id DESC")
	}
	if categoryId > 0 {
		query.Where("category_id=?", categoryId)
	}
	if keyword != "" {
		query.Where("title LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Offset(offset).Limit(pageSize).Find(&data).Error; err != nil {
		err = errors.New("数据为空")
	}

	for i, _ := range data {
		var items []models.NetDiskItem
		var str string
		str, _ = utils.TruncateString(data[i].Description, 100)
		data[i].Description = str
		data[i].Url = fmt.Sprintf("/archives/%s.html", data[i].GetUid())
		err1 := json.Unmarshal([]byte(data[i].DiskItems), &items)
		if err1 == nil {
			data[i].DiskItemsArray = items
		} else {
			err1 = errors.New("出错了")
		}

	}
	return
}

// 获取资源列表
func (resourceItemService *resourceItemService) GetRecommendList(page int, pageSize int, categoryId int32, keyword string, orderBy string) (err error, data []models.ResourceItem, total int64) {

	offset := (page - 1) * pageSize
	query := global.App.DB.Model(models.ResourceItem{}).Where("is_recommend=?", 1)
	if orderBy != "" {
		fmt.Println(orderBy)
		query.Order(orderBy + ",id DESC")
	} else {
		query.Order("id DESC")
	}
	if categoryId > 0 {
		query.Where("category_id=?", categoryId)
	}
	if keyword != "" {
		query.Where("title LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Offset(offset).Limit(pageSize).Find(&data).Error; err != nil {
		err = errors.New("数据为空")
	}

	for i, _ := range data {
		var items []models.NetDiskItem
		var str string
		str, _ = utils.TruncateString(data[i].Description, 100)
		data[i].Description = str
		data[i].Url = fmt.Sprintf("/archives/%s.html", data[i].GetUid())
		err1 := json.Unmarshal([]byte(data[i].DiskItems), &items)
		if err1 == nil {
			data[i].DiskItemsArray = items
		} else {
			err1 = errors.New("出错了")
		}

	}
	return
}

// 获取资源列表
func (resourceItemService *resourceItemService) GetNewResList(page int, pageSize int, category_id int32) (err error, data []models.ResourceItem) {
	pageSize = 15
	offset := (page - 1) * pageSize
	query := global.App.DB.Order("id DESC")
	if category_id > 0 {
		query.Where("category_id=?", category_id)
	}

	if err := query.Model(models.ResourceItem{}).Offset(offset).Limit(pageSize).Find(&data).Error; err != nil {
		err = errors.New("数据为空")
	}

	for i, _ := range data {
		var items []models.NetDiskItem
		var str string
		str, _ = utils.TruncateString(data[i].Description, 100)
		data[i].Description = str
		fmt.Println(data[i].GetUid())
		data[i].Url = fmt.Sprintf("/archives/%s.html", data[i].GetUid())
		fmt.Println(data[i].Url)
		err1 := json.Unmarshal([]byte(data[i].DiskItems), &items)
		if err1 == nil {
			data[i].DiskItemsArray = items
		} else {
			err1 = errors.New("出错了")
		}

	}
	return
}

// 获取资源信息
func (resourceItemService *resourceItemService) GetResInfo(id int) (err error, data models.ResourceItem) {
	err = global.App.DB.First(&data, id).Error
	if err != nil {
		err = errors.New("数据为空")
		return
	}
	var items []models.NetDiskItem
	err1 := json.Unmarshal([]byte(data.DiskItems), &items)
	if err1 != nil {
		err1 = errors.New("解析json出错")
		return
	}
	data.DiskItemsArray = items
	return
}
