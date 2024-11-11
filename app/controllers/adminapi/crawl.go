package adminapi

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/app/services"
	"github.com/jassue/jassue-gin/global"
	"github.com/jassue/jassue-gin/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CrawlController struct{}

func (uc *CrawlController) GetList(c *gin.Context) {
	var crawlItems []models.CrawlItem
	var totalcrawlItems int64
	pageStr := c.DefaultQuery("pageNum", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	keyword := c.DefaultQuery("keyword", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		fmt.Println(err)
		page = 1
	}
	//3a0dc2d67ac44660a839a119412a2c4e-github
	//err1, str := services.QuarkService.SaveShare("3a0dc2d67ac44660a839a119412a2c4e", "https://pan.quark.cn/s/433aff339be1")
	//fmt.Println(str)
	//fmt.Println(err1)
	//fmt.Println(page)
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	// Calculate offset and limit
	offset := (page - 1) * pageSize
	limit := pageSize
	db := global.App.DB

	query := db.Model(models.CrawlItem{})
	if keyword != "" {
		query.Where("mobile LIKE ?", "%"+keyword+"%").Or("name LIKE ?", "%"+keyword+"%")
	}

	query.Count(&totalcrawlItems).Limit(limit).Offset(offset).Order("id desc").Find(&crawlItems)
	for i, item := range crawlItems {
		crawlItems[i].CreateTimeStr = utils.TimestampToDateYmd(item.CreatedAt)
	}
	var res = gin.H{
		"list":  crawlItems,
		"total": totalcrawlItems,
	}
	response.Success(c, res)
}

// // Create handles POST requests to create a new user
func (uc *CrawlController) Create(c *gin.Context) {
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
	err1 := json.Unmarshal([]byte(input.DiskItems), &input.DiskItemsArray)
	if err1 != nil {
		fmt.Println("Error decoding JSON:", err1)
		return
	}
}

// Batch save to my netdisk
func (uc *CrawlController) BatchSaveToDisk(c *gin.Context) {

	var input request.TransSave

	// Bind JSON payload to input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 拆分字符串为切片
	idsSlice := strings.Split(input.Ids, ",")

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
	var items []models.CrawlItem
	if len(ids) > 0 {
		// 这里使用 IN 查询来批量删除这些 ID
		global.App.DB.Where("id IN ?", ids).Find(&items)
		for _, item := range items {
			services.QuarkService.SaveShare(input.Fid, item.Url)
			time.Sleep(3 * time.Second)
		}

	}

	response.Success(c, nil)
}

// Delete handles DELETE requests to delete a user
func (uc *CrawlController) Delete(c *gin.Context) {
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
		global.App.DB.Where("id IN ?", ids).Delete(models.CrawlItem{})
	}
	response.Success(c, nil)
}
