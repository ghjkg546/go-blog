package adminapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/app/services"
	"github.com/jassue/jassue-gin/global"
	client "github.com/zinclabs/sdk-go-zincsearch"
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

	query.Count(&totalUsers).Limit(limit).Offset(offset).Order("id desc").Find(&users)
	var res = gin.H{
		"list":  users,
		"total": totalUsers,
	}
	response.Success(c, res)
}

func (uc *ResourceController) SyncToSearch(c *gin.Context) {
	var users []models.ResourceItem
	pageStr := c.DefaultQuery("pageNum", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "100000")
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

	query.Limit(limit).Offset(offset).Find(&users)
	err = services.SearchItemService.BatchSync(&users)
	if err != nil {
		response.BusinessFail(c, err.Error())
	}
	response.Success(c, "成功")
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
	index := "resource_item" // string | Index
	err1 := json.Unmarshal([]byte(input.DiskItems), &input.DiskItemsArray)
	if err1 != nil {
		fmt.Println("Error decoding JSON:", err1)
		return
	}

	typeStr := ","
	for i := range input.DiskItemsArray {
		menu := input.DiskItemsArray[i]
		typeStr = typeStr + strconv.Itoa(menu.Type)
	}
	typeStr = typeStr + ","
	document := map[string]interface{}{
		"_id":       input.GetUid(),
		"disk_type": typeStr,
		"title":     input.Title,
		"url":       input.DiskItems,
	}

	ctx := context.WithValue(context.Background(), client.ContextBasicAuth, client.BasicAuth{
		UserName: global.App.Config.Search.UserName,
		Password: global.App.Config.Search.Password,
	})

	configuration := client.NewConfiguration()
	configuration.Servers = client.ServerConfigurations{
		client.ServerConfiguration{
			URL: global.App.Config.Search.Url,
		},
	}

	apiClient := client.NewAPIClient(configuration)
	resp, _, err := apiClient.Document.Index(ctx, index).Document(document).Execute()
	if err != nil {
		fmt.Println(err)
		response.BusinessFail(c, err.Error())
	}

	global.App.DB.Model(&models.ResourceItem{}).Where("id = ?", input.GetUid()).Update("search_id", resp.GetId())
	response.Success(c, nil)
}

// Create handles POST requests to create a new reasouce
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
	idsString := c.Param("id")
	// 拆分字符串为切片
	idsSlice := strings.Split(idsString, ",")

	// 遍历切片并将字符串转换为整数
	var ids []int
	for _, idStr := range idsSlice {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Printf("转换错误: %s\n", err)
			continue
		}
		ids = append(ids, id)
	}
	// 使用 GORM 删除对应 ID 的记录
	if len(ids) > 0 {
		// 这里使用 IN 查询来批量删除这些 ID
		global.App.DB.Where("id IN ?", ids).Delete(models.ResourceItem{})
		fmt.Println("asdfasfdsaf")
	}
	response.Success(c, nil)
}
