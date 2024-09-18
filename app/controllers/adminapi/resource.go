package adminapi

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
	"net/http"
	"strconv"
	"strings"
)

type ResourceController struct{}

func (uc *ResourceController) GetList(c *gin.Context) {
	var users []models.ResourceItem
	var totalUsers int64
	pageStr := c.DefaultQuery("pageNum", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	keyword := c.DefaultQuery("keyword", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		fmt.Println(err)
		page = 1
	}
	fmt.Println(page)
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	// Calculate offset and limit
	offset := (page - 1) * pageSize
	limit := pageSize
	db := global.App.DB

	query := db.Model(models.ResourceItem{})
	if keyword != "" {
		query.Where("mobile LIKE ?", "%"+keyword+"%").Or("name LIKE ?", "%"+keyword+"%")
	}

	query.Count(&totalUsers).Limit(limit).Offset(offset).Find(&users)
	var res = gin.H{
		"list":  users,
		"total": totalUsers,
	}
	response.Success(c, res)
}

// GetDetail handles GET requests for user details
func (uc *ResourceController) GetDetail(c *gin.Context) {
	id := c.Param("id")
	var user models.ResourceItem
	db := global.App.DB
	db.Model(models.ResourceItem{}).First(&user, id)
	//var netItems models.NetDiskItem
	err := json.Unmarshal([]byte(user.DiskItems), &user.DiskItemsArray)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	response.Success(c, user)
}

// Create handles POST requests to create a new user
func (uc *ResourceController) Create(c *gin.Context) {
	var input models.ResourceItem
	db := global.App.DB
	// Bind JSON payload to input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.Views = 0
	// Save data to database
	result := db.Save(&input)
	if result.Error != nil {
		response.Fail(c, 500, result.Error.Error())

		return
	}
	response.Success(c, nil)
}

// Create handles POST requests to create a new user
func (uc *ResourceController) BatchCreate(c *gin.Context) {
	var input request.BatchSave
	db := global.App.DB
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lines := strings.Split(input.Content, "\n")
	for _, line := range lines {
		fmt.Println("line")
		parts := strings.Split(line, ">>>")

		var items []models.NetDiskItem

		// Create a new NetDiskItem
		newItem := models.NetDiskItem{
			Type: 2,
			Url:  parts[1],
		}

		// Append the new item to the slice
		items = append(items, newItem)

		// Convert the slice to JSON
		jsonData, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			continue

		}

		err1 := db.Create(&models.ResourceItem{Views: 0, Title: parts[0], DiskItems: string(jsonData), CategoryId: input.CategoryId, Status: 1, CoverImg: ""})

		if err1.Error != nil {
			response.Fail(c, 500, err1.Error.Error())

			return
		}

	}

	response.Success(c, nil)
}

// Update handles PUT requests to update a user
func (uc *ResourceController) Update(c *gin.Context) {
	var input models.ResourceItem
	db := global.App.DB
	// Bind JSON payload to input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save data to database
	result := db.Save(&input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	response.Success(c, nil)
}

// Delete handles DELETE requests to delete a user
func (uc *ResourceController) Delete(c *gin.Context) {
	id := c.Param("id")
	global.App.DB.Delete(&models.ResourceItem{}, id)
	response.Success(c, nil)
}
