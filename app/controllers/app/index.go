package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/app/services"
	"github.com/jassue/jassue-gin/global"
	"gorm.io/gorm"
	"html/template"
	"math"
	"net/http"
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
	//intCategoryId, err := strconv.Atoi(categoryId)
	//if err != nil {
	//	intCategoryId = 0
	//}

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
	err, data, total := services.ResourceItemService.GetResList(page, pageSize, cate.ID, keyword)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	err, dataNew := services.ResourceItemService.GetNewResList(page, pageSize, cate.ID)
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
	fmt.Println(cates)
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	//totalPages = 100
	c.HTML(http.StatusOK, "index.html", gin.H{
		"CategoryId":     slug,
		"Cates":          cates,
		"blogItems":      data,
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
