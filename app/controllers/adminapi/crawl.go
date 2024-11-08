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
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CrawlController struct{}

func (uc *CrawlController) GetList(c *gin.Context) {
	var users []models.CrawlItem
	var totalUsers int64
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

	query.Count(&totalUsers).Limit(limit).Offset(offset).Order("id desc").Find(&users)
	var res = gin.H{
		"list":  users,
		"total": totalUsers,
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
	//index := "resource_item" // string | Index
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

//
//// 等待分享列表
//func (uc *CrawlController) WaitShareList(c *gin.Context) {
//	pageStr := c.DefaultQuery("pageNum", "1")
//	pageSizeStr := c.DefaultQuery("pageSize", "10")
//	fidStr := c.DefaultQuery("fid", "")
//
//	page, err := strconv.Atoi(pageStr)
//	if err != nil || page < 1 {
//		fmt.Println(err)
//		page = 1
//	}
//
//	pageSize, err := strconv.Atoi(pageSizeStr)
//	if err != nil || pageSize < 1 {
//		pageSize = 10
//	}
//	var dirResp response.DirResponse
//	dirResp = services.QuarkService.GetDirInfo(fidStr, page, pageSize)
//
//	var res = gin.H{
//		"list":  dirResp.Data,
//		"total": dirResp.Total,
//	}
//	response.Success(c, res)
//}
//
//// 等待分享列表
//func (uc *CrawlController) Crawl(c *gin.Context) {
//
//	var input request.Crawl
//	//db := global.App.DB
//	if err := c.BindJSON(&input); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	fmt.Println(input)
//	//
//	//params := request.Crawl{
//	//	DetailUrl: "https://example.com",
//	//	NameRule:  `<h1>(.*?)</h1>`,  // Example regex for name
//	//	LinkRule:  `href="(https://example.com.*?)"`,  // Example regex for link
//	//}
//
//	name, link, err := CrawlHTML(input)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	fmt.Println("Name:", name)
//	fmt.Println("Link:", link)
//	//return
//	//page, err := strconv.Atoi(pageStr)
//	//if err != nil || page < 1 {
//	//	fmt.Println(err)
//	//	page = 1
//	//}
//	//
//	//pageSize, err := strconv.Atoi(pageSizeStr)
//	//if err != nil || pageSize < 1 {
//	//	pageSize = 10
//	//}
//	//var dirResp response.DirResponse
//	//dirResp = services.QuarkService.GetDirInfo(fidStr, page, pageSize)
//	//
//	//var res = gin.H{
//	//	"list":  dirResp.Data,
//	//	"total": dirResp.Total,
//	//}
//	response.Success(c, gin.H{})
//}
//
//func CrawlHTML(params request.Crawl) (string, string, error) {
//	// Step 1: Fetch HTML content from the URL
//	resp, err := http.Get(params.DetailUrl)
//	if err != nil {
//		return "", "", fmt.Errorf("failed to fetch URL: %w", err)
//	}
//	defer resp.Body.Close()
//
//	htmlData, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return "", "", fmt.Errorf("failed to read HTML content: %w", err)
//	}
//
//	// Step 2: Extract name and link using regex patterns
//	nameRegex := regexp.MustCompile(params.NameRule)
//	linkRegex := regexp.MustCompile(params.LinkRule)
//
//	nameMatch := nameRegex.FindStringSubmatch(string(htmlData))
//	linkMatch := linkRegex.FindStringSubmatch(string(htmlData))
//
//	// Check if matches were found
//	if len(nameMatch) < 2 || len(linkMatch) < 2 {
//		return "", "", fmt.Errorf("failed to find matches with the provided regex rules")
//	}
//
//	// Return matched results
//	return nameMatch[1], linkMatch[1], nil
//}
//
//func calculatePages(totalItems, itemsPerPage int) int {
//	return int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))
//}
//
//// 批量分享
//func (uc *CrawlController) BatchShare(c *gin.Context) {
//	var input request.BatchShare
//
//	if err := c.BindJSON(&input); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	var dirResp response.DirResponse
//	dirResp = services.QuarkService.GetDirInfo(input.Fid, 1, 50)
//	response.Success(c, nil)
//	var ids []string
//	if input.PageSize >= dirResp.Total {
//		var chunks []models.ShareItem
//		for _, item := range dirResp.Data {
//
//			ids = append(ids, item.Fid)
//
//			chunks = append(chunks, models.ShareItem{Name: item.FileName, ID: item.Fid})
//		}
//		services.QuarkService.SaveResouceByUrl(ids, "test", chunks, input.CategoryId)
//		ids = []string{} // Clear the ids slice
//		chunks = []models.ShareItem{}
//
//	} else {
//		pages := calculatePages(dirResp.Total, 50)
//		for i := 0; i < pages; i++ {
//			dirResp = services.QuarkService.GetDirInfo(input.Fid, i+1, 50)
//			fmt.Printf("Processing page %d\n", i+1)
//			var chunks []models.ShareItem
//			for _, item := range dirResp.Data {
//
//				ids = append(ids, item.Fid)
//
//				chunks = append(chunks, models.ShareItem{Name: item.FileName, ID: item.Fid})
//
//				if len(ids) >= input.PageSize {
//
//					services.QuarkService.SaveResouceByUrl(ids, item.FileName, chunks, input.CategoryId)
//					ids = []string{} // Clear the ids slice
//					chunks = []models.ShareItem{}
//					time.Sleep(5 * time.Second)
//				}
//			}
//		}
//
//	}
//
//	response.Success(c, nil)
//}
