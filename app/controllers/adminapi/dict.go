package adminapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
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

// Update handles PUT requests to update a dict
func (uc *DictController) Update(c *gin.Context) {
	var input models.SysDict
	db := global.App.DB
	// Bind JSON payload to input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, item := range input.DictItems {
		fmt.Printf("Sort: %d, Status: %d, Name: %s, Value: %s\n", item.Sort, item.Status, item.Name, item.Value)
		if item.ID != 0 {
			// Update existing item if ID is present
			if err := db.Model(&models.SysDictItem{}).Where("id = ?", item.ID).Updates(item).Error; err != nil {
				log.Println("Error updating item:", err)
				continue
			}
			fmt.Printf("Updated item with ID %d\n", item.ID)
		}
	} // Save data to database
	result := db.Save(&input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	response.Success(c, nil)

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
