package app

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/controllers/common"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/app/services"
	"github.com/jassue/jassue-gin/global"
	client "github.com/zinclabs/sdk-go-zincsearch"
	"gorm.io/gorm"
	"html/template"
	"log"
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

type ShareItem struct {
	Fids string `json:"Fids" gorm:"size:200;not null;comment:资源id"`
	Name string `json:"Name" gorm:"comment:资源名称"`
}

func DoShare(c *gin.Context) {
	var input ShareItem

	// Bind JSON payload to input
	if err := c.BindJSON(&input); err != nil {

		response.BusinessFail(c, err.Error())

		return
	}

	var data []models.ShareItem
	if err := json.Unmarshal([]byte(input.Fids), &data); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	var ids []string
	if len(data) > 0 {
		for i := range data {
			res := &data[i]
			ids = append(ids, res.ID)

		}

	}

	services.QuarkService.SaveResouceByUrl(ids, input.Name, data, 8)

	//url := services.QuarkService.ShareItem(ids, input.Name)
	//db := global.App.DB
	//if url != "" {
	//	fmt.Println(url)
	//
	//	if len(data) > 0 {
	//		for i := range data {
	//			res := &data[i]
	//			ids = append(ids, res.ID)
	//			//var tmp models.ResourceItem
	//			var items []models.NetDiskItem
	//
	//			// Create a new NetDiskItem
	//			newItem := models.NetDiskItem{
	//				Type: 2,
	//				Url:  url,
	//			}
	//
	//			// Append the new item to the slice
	//			items = append(items, newItem)
	//
	//			// Convert the slice to JSON
	//			jsonData, err := json.MarshalIndent(items, "", "  ")
	//			if err != nil {
	//				continue
	//
	//			}
	//
	//			err1 := db.Create(&models.ResourceItem{Views: 0, Title: res.Name, DiskItems: string(jsonData), CategoryId: 8, Status: 1, CoverImg: ""})
	//
	//			if err1 != nil {
	//				continue
	//			}
	//
	//		}
	//
	//	}
	//
	//}

	var res = gin.H{
		"message": "分享成功",
	}
	response.Success(c, res)
}

func GetFavList(c *gin.Context) {
	//keyword := c.DefaultQuery("keyword", "")
	categoryId := c.DefaultQuery("category_id", "0")
	pageSizeStr := c.DefaultQuery("pageSize", "50")

	pageStr := c.DefaultQuery("page", "1")
	page, err1 := strconv.Atoi(pageStr)
	if err1 != nil {
		page = 1
	}
	pageSize, err3 := strconv.Atoi(pageSizeStr)
	if err3 != nil {
		pageSize = 10
	}
	cid := 0
	cid, err2 := strconv.Atoi(categoryId)
	if err2 != nil {
		cid = 0
	}
	fmt.Println(cid)
	itemIDs, _ := services.AppUserService.GetFavIds(c)
	fmt.Println(itemIDs)
	totalCount := len(itemIDs)
	// 计算分页
	start := (page - 1) * pageSize
	end := int(math.Min(float64(start+pageSize), float64(totalCount)))
	if start >= totalCount {
		var res = gin.H{
			"list":  []models.ResourceItem{},
			"total": totalCount,
		}
		response.Success(c, res)
		return
	}

	// 分页获取 itemID
	pagedItemIDs := itemIDs[start:end]

	// 从数据库中查询收藏的物品
	var items []models.ResourceItem
	if err := global.App.DB.Where("id IN ?", pagedItemIDs).Find(&items).Error; err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	var res = gin.H{
		"list":  items,
		"total": totalCount,
	}
	response.Success(c, res)
	return
	//err, data, total := services.ResourceItemService.GetResList(page, pageSize, int32(cid), keyword)
	//if err != nil {
	//	response.BusinessFail(c, err.Error())
	//	return
	//}

	//for i := range data {
	//	res := &data[i]
	//	tm1 := time.Unix(res.CreatedAt, 0)
	//	tm2 := time.Unix(res.UpdatedAt, 0)
	//	//res.Description = ""
	//	res.CreateTimeStr = tm1.Format("2006-01-02 15:04:05")
	//	res.UpdateTimeStr = tm2.Format("2006-01-02 15:04:05")
	//}
	//
	//var res = gin.H{
	//	"list":  data,
	//	"total": total,
	//}
	//response.Success(c, res)
}

func GetResList(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")
	categoryId := c.DefaultQuery("category_id", "0")
	pageSizeStr := c.DefaultQuery("pageSize", "50")

	pageStr := c.DefaultQuery("page", "1")
	page, err1 := strconv.Atoi(pageStr)
	if err1 != nil {
		page = 1
	}
	pageSize, err3 := strconv.Atoi(pageSizeStr)
	if err3 != nil {
		pageSize = 10
	}
	cid := 0
	cid, err2 := strconv.Atoi(categoryId)
	if err2 != nil {
		cid = 0
	}
	fmt.Println(cid)
	err, data, total := services.ResourceItemService.GetResList(page, pageSize, int32(cid), keyword)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	for i := range data {
		res := &data[i]
		tm1 := time.Unix(res.CreatedAt, 0)
		tm2 := time.Unix(res.UpdatedAt, 0)
		//res.Description = ""
		res.CreateTimeStr = tm1.Format("2006-01-02 15:04:05")
		res.UpdateTimeStr = tm2.Format("2006-01-02 15:04:05")
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

// GetFrontReasouceItems returns the blog items
func (bc *ResController) GetFrontReasouceItems(c *gin.Context) {
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
	cateName := ""
	if cate.ID > 0 {
		cid = cate.ID
		cateName = cate.Name
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
	if len(data) < maxIndex {
		maxIndex = len(data)

	}

	if len(data) >= 1 {
		subItems = data[1:maxIndex]
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	c.HTML(http.StatusOK, "index.html", gin.H{
		"CategoryId":     slug,
		"Cates":          cates,
		"blogItems":      subItems,
		"cateName":       cateName,
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
		data.IsFavorite = false
	}
	global.App.DB.Model(&models.ResourceItem{}).Where("id = ?", intId).UpdateColumn("views", gorm.Expr("views + ?", 1))
	data.IsFavorite, _ = services.AppUserService.IsFavorite(c, id)
	tm1 := time.Unix(data.CreatedAt, 0)

	data.CreateTimeStr = tm1.Format("2006-01-02 15:04:05")
	//response.Success(c, gin.H{"info": data, "comments": comments})
	buttonTextMap := map[int]string{
		1: "百度网盘",
		2: "夸克网盘",
		3: "阿里网盘",
		4: "移动彩云",
	}
	c.HTML(http.StatusOK, "detail.html", gin.H{
		"Content":       template.HTML(data.Description),
		"ButtonTextMap": buttonTextMap,
		"blogItem":      data,
		"Cates":         cates,
	})

	c.AbortWithStatus(http.StatusNotFound)
}

func (bc *ResController) ApplyForReasource(c *gin.Context) {
	//idStr := c.Param("id")
	//id := strings.TrimSuffix(idStr, ".html")
	//intId, err := strconv.Atoi(id)

	//if err != nil {
	//	c.String(400, "Invalid ID")
	//	return
	//}
	//fmt.Println(intId)
	//err, data := services.ResourceItemService.GetResInfo(intId)
	//if err != nil {
	//	response.BusinessFail(c, err.Error())
	//	return
	//}
	err, cates, _ := services.CategoryService.GetList()
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	//var comments []models.Comment
	//err1 := global.App.DB.Preload("User").Where("resource_item_id=?", intId).Order("id desc").Limit(20).Find(&comments)
	//if err1 != nil {
	//	fmt.Println(err1.Error)
	//}
	//global.App.DB.Model(&models.ResourceItem{}).Where("id = ?", intId).UpdateColumn("views", gorm.Expr("views + ?", 1))
	//
	//tm1 := time.Unix(data.CreatedAt, 0)
	//
	//data.CreateTimeStr = tm1.Format("2006-01-02 15:04:05")
	//response.Success(c, gin.H{"info": data, "comments": comments})
	c.HTML(http.StatusOK, "qiu.html", gin.H{
		"Content": template.HTML("求资源"),
		//"blogItem": data,
		"Cates": cates,
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
		response.Fail(c, 500, "参数出错")
		return
	}
	if !captcha.VerifyString(input.CaptchaId, input.CaptchaValue) {
		response.BusinessFail(c, "验证码错误")
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
	data.IsFavorite, _ = services.AppUserService.IsFavorite(c, c.Query("id"))
	data.CreateTimeStr = tm1.Format("2006-01-02 15:04:05")
	response.Success(c, gin.H{"info": data, "comments": comments})
}

// WXTextMsg 微信文本消息结构体
type WXTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	MsgId        int64
}

// WXMsgReceive 微信消息接收
func WXMsgReceive(c *gin.Context) {
	var textMsg WXTextMsg
	err := c.ShouldBindXML(&textMsg)
	if err != nil {
		log.Printf("[消息接收] - XML数据包解析失败: %v\n", err)
		return
	}

	log.Printf("[消息接收] - 收到消息, 消息类型为: %s, 消息内容为: %s\n", textMsg.MsgType, textMsg.Content)
	fmt.Println(textMsg.Content)

	keyword := textMsg.Content

	db := global.App.DB
	var items []models.ResourceItem
	query := db.Model(models.ResourceItem{})
	if keyword != "" {
		query.Where("title LIKE ?", "%"+keyword+"%")
	}

	query.Limit(20).Order("id desc").Find(&items)
	text := ""
	for _, item := range items {
		input := &item
		err1 := json.Unmarshal([]byte(input.DiskItems), &input.DiskItemsArray)
		if err1 != nil {
			fmt.Println("Error decoding JSON:", err1)
			return
		}
		typeStr := ""

		for i := range input.DiskItemsArray {
			menu := input.DiskItemsArray[i]
			typeStr = typeStr + menu.Url
		}
		text += fmt.Sprintf("标题：%s\n网盘连接：%s\n", item.Title, typeStr)
	}

	if text != "" {
		WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName, text)
	} else {
		WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName, "抱歉，暂未收录")
	}
}

// WXRepTextMsg 微信回复文本消息结构体
type WXRepTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	// 若不标记XMLName, 则解析后的xml名为该结构体的名称
	XMLName xml.Name `xml:"xml"`
}

// WXMsgReply 微信消息回复
func WXMsgReply(c *gin.Context, fromUser, toUser string, content string) {
	repTextMsg := WXRepTextMsg{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      content,
	}

	msg, err := xml.Marshal(&repTextMsg)
	if err != nil {
		log.Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
		return
	}
	_, _ = c.Writer.Write(msg)
}

// 与填写的服务器配置中的Token一致
const Token = "ghg546"

func WXCheckSignature(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	ok := common.CheckSignature(signature, timestamp, nonce, Token)
	if !ok {
		log.Println("微信公众号接入校验失败!")
		return
	}

	log.Println("微信公众号接入校验成功!")
	_, _ = c.Writer.WriteString(echostr)
}
