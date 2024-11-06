package adminapi

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
	client "github.com/zinclabs/sdk-go-zincsearch"
	"log"
	"net/http"
	"os"
	"time"
)

type LogController struct{}

// GetList handles GET requests for dict details
func (uc *LogController) GetList(c *gin.Context) {
	keyword := c.DefaultQuery("keywords", "")
	startime := c.DefaultQuery("createTime[0]", "")

	index := "share_log"

	query := *client.NewV1ZincQuery() // V1ZincQuery | Query

	params := *client.NewV1QueryParams()
	const layout = "2006-01-02T15:04:05.000-07:00"
	query.SetFrom(0)

	if startime != "" {
		startime1, err := time.Parse(layout, startime)
		if err == nil {
			params.SetStartTime(startime)
			// Add one day, then subtract one second
			endTime := startime1.Add(24*time.Hour - time.Second)
			params.SetEndTime(endTime.Format(layout))
		}
	}

	if keyword != "" {
		query.SetSearchType("match")
		params.SetTerm(keyword)
		params.SetField("request_data")
	} else {
		query.SetSearchType("match_all")
	}

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
		// 获取并格式化 timestamp 字段（假设时间字段是 "timestamp"）
		source := data.GetSource()
		if timestamp, exists := source["@timestamp"].(string); exists {
			parsedTime, err := time.Parse(time.RFC3339, timestamp)
			if err == nil {
				// 格式化为 "2006-01-02 15:04:05" 格式
				location, _ := time.LoadLocation("Asia/Shanghai")
				chinaTime := parsedTime.In(location)
				source["timestamp_str"] = chinaTime.Format("2006-01-02 15:04:05")
			} else {
				source["timestamp_str"] = "" // 解析失败时置空
			}
		}
		resList = append(resList, source)
	}
	var res = gin.H{
		"list":  resList,
		"total": resp.Hits.Total.Value,
	}
	response.Success(c, res)
}

// GetDetail handles GET requests for dict details
func (uc *LogController) GetDetail(c *gin.Context) {
	id := c.Param("id")
	var dict models.SysDict
	db := global.App.DB
	db.Model(models.SysDict{}).Preload("DictItems").First(&dict, id)

	response.Success(c, dict)
}

// Create handles POST requests to create a new dict
func (uc *LogController) Create(c *gin.Context) {
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
func (uc *LogController) Update(c *gin.Context) {
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
func (uc *LogController) Delete(c *gin.Context) {
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
