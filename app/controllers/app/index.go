package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/app/services"
	"github.com/jassue/jassue-gin/global"
	client "github.com/zinclabs/sdk-go-zincsearch"
	"gorm.io/gorm"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// 分类列表
func CateList(c *gin.Context) {
	err, data, total := services.CategoryService.GetList()
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	var res = gin.H{
		"list":  data,
		"total": total,
	}
	response.Success(c, res)
}

// 分类列表
func ResSearch(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")
	index := "resource_item"          // string | Index
	query := *client.NewV1ZincQuery() // V1ZincQuery | Query
	query.SetSearchType("match")
	params := *client.NewV1QueryParams()
	params.SetTerm(keyword)
	params.SetField("title")
	query.SetQuery(params)
	query.SetSortFields([]string{"-@timestamp"})
	query.SetMaxResults(20)
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
	resp, r, err := apiClient.Search.SearchV1(ctx, index).Query(query).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SearchApi.SearchV1``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	fmt.Fprintf(os.Stdout, "Response from `SearchApi.SearchV1`: %v\n", resp)
	var resList []interface{}
	for _, data := range resp.Hits.Hits {
		resList = append(resList, data.GetSource())
	}
	var res = gin.H{
		"list":  resList,
		"total": resp.Hits.Total.Value,
	}
	response.Success(c, res)
}

type ResController struct{}

func seq(start, end int) []int {
	var numbers []int
	for i := start; i <= end; i++ {
		numbers = append(numbers, i)
	}
	return numbers
}

// GetBlogItems returns the blog items
func (bc *ResController) GetBlogItems(c *gin.Context) {
	//categoryId := c.Param("category_id")
	slug := c.Param("category_id")
	keyword := c.Param("keyword")

	err, cates, total := services.CategoryService.GetList()
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	pageSize := 10 // Number of items per page
	pageStr := c.Param("page")
	page, err1 := strconv.Atoi(pageStr)
	if err1 != nil {
		page = 1
	}
	var cate models.Category
	global.App.DB.Where("slug=?", slug).First(&cate)
	var cid int32 = 0
	if cate.ID > 0 {
		cid = cate.ID
	}
	err, data, total := services.ResourceItemService.GetResList(page, pageSize, cid, keyword)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	err, dataNew := services.ResourceItemService.GetNewResList(page, pageSize, cid)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	for i := range data {
		res := &data[i]
		tm1 := time.Unix(res.CreatedAt, 0)
		tm2 := time.Unix(res.UpdatedAt, 0)
		res.Description = ""
		res.CreateTimeStr = tm1.Format("2006-01-02 15:04:05")
		res.UpdateTimeStr = tm2.Format("2006-01-02 15:04:05")
	}
	maxIndex := 10
	var subItems []models.ResourceItem
	var topItem models.ResourceItem
	if len(data) < maxIndex {
		maxIndex = len(data)

	}

	if len(data) >= 1 {
		subItems = data[1:maxIndex]
		topItem = data[0]
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	c.HTML(http.StatusOK, "index.html", gin.H{
		"CategoryId":     slug,
		"Cates":          cates,
		"blogItems":      subItems,
		"topItem":        topItem,
		"blogNew":        dataNew,
		"CurrentPage":    page,
		"TotalPages":     totalPages,
		"TotalPageArray": generatePages(totalPages, page),
		"PrevPage":       page - 1,
		"NextPage":       page + 1,
	})
}

func generatePages(totalPages, currentPage int) []string {
	pages := []string{}

	for i := 1; i <= totalPages; i++ {
		if i < 3 || i > totalPages-2 || (i >= currentPage-1 && i <= currentPage+1) {
			pages = append(pages, fmt.Sprintf("%d", i))
		} else if len(pages) == 0 || pages[len(pages)-1] != "..." {
			pages = append(pages, "...")
		}
	}

	return pages
}

func (bc *ResController) GetBlogDetail(c *gin.Context) {
	idStr := c.Param("id")
	id := strings.TrimSuffix(idStr, ".html")
	intId, err := strconv.Atoi(id)

	if err != nil {
		c.String(400, "Invalid ID")
		return
	}
	fmt.Println(intId)
	err, data := services.ResourceItemService.GetResInfo(intId)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	err, cates, _ := services.CategoryService.GetList()
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	var comments []models.Comment
	err1 := global.App.DB.Preload("User").Where("resource_item_id=?", intId).Order("id desc").Limit(20).Find(&comments)
	if err1 != nil {
		fmt.Println(err1.Error)
	}
	global.App.DB.Model(&models.ResourceItem{}).Where("id = ?", intId).UpdateColumn("views", gorm.Expr("views + ?", 1))

	tm1 := time.Unix(data.CreatedAt, 0)

	data.CreateTimeStr = tm1.Format("2006-01-02 15:04:05")
	//response.Success(c, gin.H{"info": data, "comments": comments})
	c.HTML(http.StatusOK, "detail.html", gin.H{
		"Content":  template.HTML(data.Description),
		"blogItem": data,
		"Cates":    cates,
	})

	c.AbortWithStatus(http.StatusNotFound)
}

func stringToInt32(s string) (int32, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

// Create handles POST requests to create a new user
func CreateComment(c *gin.Context) {
	var input request.PostComment
	if err := c.BindJSON(&input); err != nil {
		fmt.Println("出错000")
		response.Fail(c, 500, "参数出错")
		return
	}

	uidsStr, ok := c.Get("id")
	if !ok {
		// Handle the case where the value is not a string
		response.Fail(c, 500, "获取uid出错")
	}

	str, ok := uidsStr.(string)
	if !ok {
		response.Fail(c, 500, "uid不是字符串")
		return
	}
	uidsInt, _ := stringToInt32(str)

	fmt.Println("str")
	fmt.Println(uidsStr)

	comment := models.Comment{ParentID: 0, Content: input.Content, ResourceItemId: input.ResItemId, UserID: uidsInt}
	global.App.DB.Create(&comment)
	var comments []models.Comment
	global.App.DB.Where("resource_item_id=?", input.ResItemId).Model(&comments).Preload("User").Order("id desc").Limit(20).Find(&comments)

	// Example response, replace with actual logic
	response.Success(c, comments)
}

// 资源详情
func Info(c *gin.Context) {
	intId, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		c.String(400, "Invalid ID")
		return
	}
	err, data := services.ResourceItemService.GetResInfo(intId)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	var comments []models.Comment
	err1 := global.App.DB.Preload("User").Where("resource_item_id=?", intId).Order("id desc").Limit(20).Find(&comments)
	if err1 != nil {
		fmt.Println(err1.Error)
	}
	global.App.DB.Model(&models.ResourceItem{}).Where("id = ?", intId).UpdateColumn("views", gorm.Expr("views + ?", 1))

	tm1 := time.Unix(data.CreatedAt, 0)

	data.CreateTimeStr = tm1.Format("2006-01-02 15:04:05")
	response.Success(c, gin.H{"info": data, "comments": comments})
}
