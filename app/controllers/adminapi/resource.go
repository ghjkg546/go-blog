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
	"github.com/jassue/jassue-gin/utils"
	client "github.com/zinclabs/sdk-go-zincsearch"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ResourceController struct{}

func (uc *ResourceController) GetList(c *gin.Context) {
	var ResourceItem []models.ResourceItem
	var totalResourceItem int64
	pageStr := c.DefaultQuery("pageNum", "1")
	disk_type_id_str := c.DefaultQuery("disk_type_id", "")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	keyword := c.DefaultQuery("keyword", "")

	categoryIdStr := c.DefaultQuery("category_id", "")
	startTime := c.DefaultQuery("createTime[0]", "")
	endTime := c.DefaultQuery("createTime[1]", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		fmt.Println(err)
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	cagetoryId, err2 := strconv.Atoi(categoryIdStr)
	if err2 != nil {
		cagetoryId = 0
	}
	// Calculate offset and limit
	offset := (page - 1) * pageSize
	limit := pageSize
	db := global.App.DB

	query := db.Model(models.ResourceItem{})
	if cagetoryId != 0 {
		query.Where("category_id= ?", cagetoryId)
	}
	if disk_type_id_str != "" {
		query.Where("disk_items LIKE ?", "%"+`"type": `+disk_type_id_str+`%`)
		query.Or("disk_items LIKE ?", "%"+`"type":`+disk_type_id_str+`%`)

	}
	if startTime != "" {
		startTimeStamp, endTimeStamp := services.CrudService.ParseStartEndTime(startTime, endTime)
		query.Where("created_at between ? and ?", startTimeStamp, endTimeStamp)
	}
	if keyword != "" {
		query.Where("title LIKE ?", "%"+keyword+"%")
	}

	query.Count(&totalResourceItem).Preload("Category").Limit(limit).Offset(offset).Order("id desc").Find(&ResourceItem)
	for i, item := range ResourceItem {
		ResourceItem[i].CreateTimeStr = utils.TimestampToDateYmd(item.CreatedAt)
	}

	var res = gin.H{
		"list":  ResourceItem,
		"total": totalResourceItem,
	}
	response.Success(c, res)
}

func (uc *ResourceController) SyncToSearch(c *gin.Context) {
	var ResourceItem []models.ResourceItem
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

	query.Limit(limit).Offset(offset).Find(&ResourceItem)
	err = services.SearchItemService.BatchSync(&ResourceItem)
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
		parts := strings.Split(line, ">>>")

		// Create a new NetDiskItem
		newItem := models.NetDiskItem{
			Type: input.DiskTypeId,
			Url:  parts[1],
		}

		//&models.ResourceItem{Views: 0, Title: parts[0], DiskItems: string(jsonData), CategoryId: input.CategoryId, Status: 1, CoverImg: ""}
		err1 := CreateOrUpdateResourceItem(models.ResourceItem{Views: 0, Title: parts[0], DiskItems: "",
			CategoryId: input.CategoryId, Status: 1, CoverImg: ""}, newItem, db)
		if err1 != nil {
			response.Fail(c, 500, err1.Error())
			fmt.Println("操作失败:", err1)
		} else {
			fmt.Println("操作成功")
		}

	}

	response.Success(c, nil)
}

func CreateOrUpdateResourceItem(input models.ResourceItem, newItem models.NetDiskItem, db *gorm.DB) error {
	// 去掉 title 前后空格
	trimmedTitle := strings.TrimSpace(input.Title)

	// 先尝试根据 title 查找记录
	var existingItem models.ResourceItem
	db.Where("title = ?", trimmedTitle).First(&existingItem)
	var items []models.NetDiskItem
	fmt.Println(existingItem.GetUid())
	if existingItem.GetUid() == "0" {
		// Create a new NetDiskItem
		items = append(items, newItem)
		jsonData, err := json.MarshalIndent(items, "", "  ")
		fmt.Println(jsonData)
		//if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果未找到记录，则创建新记录
		err = db.Create(&models.ResourceItem{
			Views:      0,
			Title:      trimmedTitle,
			DiskItems:  string(jsonData),
			CategoryId: input.CategoryId,
			Status:     1,
			CoverImg:   "",
		}).Error
		if err != nil {
			return fmt.Errorf("创建记录失败: %v", err)
		}
		fmt.Println("创建新记录成功")
		//} else {
		//	return fmt.Errorf("查找记录失败: %v", err)
		//}
	} else {
		// 如果找到记录，则更新
		err1 := json.Unmarshal([]byte(existingItem.DiskItems), &items)
		if err1 != nil {
			fmt.Println("Error decoding JSON:", err1)
			return nil
		}
		items = append(items, newItem)
		jsonData, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return fmt.Errorf("update时json_encode失败: %v", err)

		}
		fmt.Println(string(jsonData))
		err = db.Model(&existingItem).Updates(models.ResourceItem{
			Views:     existingItem.Views + 1, // 示例更新逻辑
			DiskItems: string(jsonData),

			Status: 1,
		}).Error
		if err != nil {
			return fmt.Errorf("更新记录失败: %v", err)
		}
		fmt.Println("更新记录成功")
	}

	return nil
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
	}
	response.Success(c, nil)
}

// 等待分享列表
func (uc *ResourceController) WaitShareList(c *gin.Context) {
	pageStr := c.DefaultQuery("pageNum", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	fidStr := c.DefaultQuery("fid", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		fmt.Println(err)
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	var dirResp response.DirResponse
	dirResp = services.QuarkService.GetDirInfo(fidStr, page, pageSize)
	for i, item := range dirResp.Data {
		dirResp.Data[i].CreateTimeStr = utils.TimestampToDateYmd(item.CreatedAt / 1000)
	}
	var res = gin.H{
		"list":  dirResp.Data,
		"total": dirResp.Total,
	}
	response.Success(c, res)
}

// Fetch HTML content from a URL
func fetchListHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch data: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Extract href links using regex
func extractLinks(html string) []string {
	// Regular expression to find the href attribute within <li class="item tcod0"> <a> tags
	regex := `<h2><a href="(https://www\.ssyhb\.cn/\d+\.html)"\s*class="main">`
	re := regexp.MustCompile(regex)

	// Find all matches
	matches := re.FindAllStringSubmatch(html, -1)

	var links []string
	for _, match := range matches {
		if len(match) > 1 {
			links = append(links, match[1])
		}
	}

	return links
}

// 等待分享列表
func (uc *ResourceController) Crawl(c *gin.Context) {

	var input request.Crawl

	url := "https://www.ssyhb.cn/page_2.html"

	// Fetch HTML content
	html, err := fetchListHTML(url)
	if err != nil {
		fmt.Println("Error fetching HTML:", err)
		return
	}

	// Extract href links
	links := extractLinks(html)
	if len(links) == 0 {
		fmt.Println("No links found.")
		return
	}
	fmt.Println("Extracted links:")
	for _, link := range links {
		fmt.Println(link)
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(input)

	CrawlDetail(links, input)

	response.Success(c, gin.H{})
}

func CrawlDetail(links []string, input request.Crawl) {
	for _, link := range links {
		fmt.Println("正在抓取", link)
		input.DetailUrl = link
		name, link, err := CrawlHTML(input)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		db := global.App.DB
		var craw_item models.CrawlItem
		craw_item.Name = name
		craw_item.Url = link
		db.Save(&craw_item)
		time.Sleep(5 * time.Second) // 暂停 5 秒
	}

}

func CrawlHTML(params request.Crawl) (string, string, error) {
	// Step 1: Fetch HTML content from the URL
	resp, err := http.Get(params.DetailUrl)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read HTML content: %w", err)
	}

	// Step 2: Extract name and link using regex patterns
	nameRegex := regexp.MustCompile(params.NameRule)
	linkRegex := regexp.MustCompile(params.LinkRule)

	nameMatch := nameRegex.FindStringSubmatch(string(htmlData))
	linkMatch := linkRegex.FindStringSubmatch(string(htmlData))

	// Check if matches were found
	if len(nameMatch) < 2 || len(linkMatch) < 2 {
		return "", "", fmt.Errorf("failed to find matches with the provided regex rules")
	}

	// Return matched results
	return nameMatch[1], linkMatch[1], nil
}

func calculatePages(totalItems, itemsPerPage int) int {
	return int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))
}

// 批量分享
func (uc *ResourceController) BatchShare(c *gin.Context) {
	var input request.BatchShare

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var dirResp response.DirResponse
	dirResp = services.QuarkService.GetDirInfo(input.Fid, 1, 50)

	var ids []string
	if input.PageSize >= dirResp.Total {
		var chunks []models.ShareItem
		for _, item := range dirResp.Data {

			ids = append(ids, item.Fid)

			chunks = append(chunks, models.ShareItem{Name: item.FileName, ID: item.Fid})
		}
		services.QuarkService.SaveResouceByUrl(ids, "test", chunks, input.CategoryId)
		ids = []string{} // Clear the ids slice
		chunks = []models.ShareItem{}

	} else {
		pages := calculatePages(dirResp.Total, 50)
		for i := 0; i < pages; i++ {
			dirResp = services.QuarkService.GetDirInfo(input.Fid, i+1, 50)
			fmt.Printf("Processing page %d\n", i+1)
			var chunks []models.ShareItem
			for _, item := range dirResp.Data {

				ids = append(ids, item.Fid)

				chunks = append(chunks, models.ShareItem{Name: item.FileName, ID: item.Fid})

				if len(ids) >= input.PageSize {

					services.QuarkService.SaveResouceByUrl(ids, item.FileName, chunks, input.CategoryId)
					ids = []string{} // Clear the ids slice
					chunks = []models.ShareItem{}
					time.Sleep(5 * time.Second)
				}
			}
		}

	}

	response.Success(c, nil)
}
