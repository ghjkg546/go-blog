package adminapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type DictController struct{}

// GetList handles GET requests for dict details
func (uc *DictController) GetList(c *gin.Context) {
	var menus []models.SysDict
	var totaldicts int64

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	keyword := c.DefaultQuery("keywords", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	db := global.App.DB

	offset := (page - 1) * pageSize
	limit := pageSize

	query := db.Model(models.SysDict{})
	if keyword != "" {
		query.Where("name LIKE ?", "%"+keyword+"%")
	}

	query.Count(&totaldicts).Limit(limit).Offset(offset).Order("id desc").Find(&menus)
	var res = gin.H{
		"list":  menus,
		"total": totaldicts,
	}
	response.Success(c, res)
}

// GetDetail handles GET requests for dict details
func (uc *DictController) GetDetail(c *gin.Context) {
	id := c.Param("id")
	var dict models.SysDict
	db := global.App.DB
	db.Model(models.SysDict{}).Preload("DictItems").First(&dict, id)

	response.Success(c, dict)
}

// Create handles POST requests to create a new dict
func (uc *DictController) Create(c *gin.Context) {
	var input models.SysDict
	db := global.App.DB
	// Bind JSON payload to input
	if err := c.BindJSON(&input); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	// Save data to database
	result := db.Save(&input)
	if result.Error != nil {
		response.Fail(c, 500, result.Error.Error())
		return
	}
	response.Success(c, nil)
}

func getOldIdsByDictId(db *gorm.DB, dictId int32) ([]int32, error) {
	var ids []int32
	// Get all IDs for the given dict_id
	err := db.Model(&models.SysDictItem{}).Where("dict_id = ?", dictId).Pluck("id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func calculateExcludeIds(oldIds []int32, newIds []int32) []int32 {
	var excludeIds []int32

	for _, oldId := range oldIds {
		found := false
		for _, newId := range newIds {
			if oldId == newId {
				found = true
				break
			}
		}
		if !found {
			excludeIds = append(excludeIds, oldId)
		}
	}
	return excludeIds
}

// Update handles PUT requests to update a dict
func (uc *DictController) Update(c *gin.Context) {
	var input models.SysDict
	db := global.App.DB
	// Bind JSON payload to input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var newIds []int32
	for _, item := range input.DictItems {
		fmt.Printf("Sort: %d, Status: %d, Name: %s, Value: %s\n", item.Sort, item.Status, item.Name, item.Value)
		if item.ID != 0 {
			newIds = append(newIds, item.ID)

		}
	}
	//excludeIds := []uint{6, 7} // IDs to exclude and delete if dict_id does not contain them
	err := deleteDictItemsIfNotMatching(db, input.ID, newIds)
	if err != nil {
		response.Fail(c, 500, "删数据出错")
	}
	for _, item := range input.DictItems {
		if item.ID != 0 {
			// Update existing item if ID is present
			if err1 := db.Model(&models.SysDictItem{}).Where("id = ?", item.ID).Updates(item).Error; err1 != nil {
				log.Println("Error updating item:", err)
				continue
			}
		}
	} // Save data to database
	result := db.Save(&input)
	if result.Error != nil {
		response.Fail(c, 500, result.Error.Error())
		return
	}
	response.Success(c, nil)

}

func deleteDictItemsIfNotMatching(db *gorm.DB, dictId int32, newIds []int32) error {
	// Delete records where dict_id is not the given dict_id and the id is in the excludeIds list

	oldIds, err := getOldIdsByDictId(db, dictId)
	if err != nil {
		return err
	}

	// Calculate excludeIds by comparing oldIds with newIds
	excludeIds := calculateExcludeIds(oldIds, newIds)
	fmt.Println(excludeIds)
	// If there are no IDs to delete, return early
	if len(excludeIds) == 0 {
		return nil
	}
	err2 := db.Where("dict_id = ? AND id IN ?", dictId, excludeIds).
		Delete(&models.SysDictItem{}).Error
	if err2 != nil {
		return err2
	}
	return nil
}

// Delete handles DELETE requests to delete a dict
func (uc *DictController) Delete(c *gin.Context) {
	id := c.Param("id")
	db := global.App.DB
	var menu models.SysDictItem
	db.Where("dict_id=?", id).First(&menu)
	if menu.ID > 0 {
		response.Fail(c, 500, "该字典下还有字典项,请先删除字典项")
		return
	}
	db.Delete(&models.SysDict{}, id)
	response.Success(c, nil)
}
