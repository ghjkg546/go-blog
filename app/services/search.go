package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
	client "github.com/zinclabs/sdk-go-zincsearch"
	"os"
	"strconv"
	"strings"
)

type searchItemService struct {
}

var SearchItemService = new(searchItemService)

// 批量同步资源到搜索引擎
func (searchItemService *searchItemService) BatchSync(items *[]models.ResourceItem) (err error) {
	index := "resource_item" // string | Index
	var query string
	for _, person := range *items {
		typeStr := ","
		var itemArray []models.NetDiskItem
		err1 := json.Unmarshal([]byte(person.DiskItems), &itemArray)
		if err1 != nil {
			fmt.Println("Error decoding JSON:", err1)
			continue
		}
		for i := range itemArray {
			menu := itemArray[i]
			typeStr = typeStr + strconv.Itoa(menu.Type)
		}
		cleanedData := strings.ReplaceAll(person.DiskItems, "\n", "")
		cleanedData = strings.TrimSpace(cleanedData)
		typeStr = typeStr + ","
		searchId := person.GetUid()
		document := map[string]interface{}{
			"_id":       searchId,
			"disk_type": typeStr,
			"title":     person.Title,
			"url":       itemArray,
		}
		jsonData, err := json.Marshal(document)
		if err != nil {
			fmt.Println(err)
			continue
		}
		global.App.DB.Model(&models.ResourceItem{}).Where("id = ?", searchId).Update("search_id", searchId)
		query += string(jsonData) + "\n"
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
	_, r, err := apiClient.Document.Multi(ctx, index).Query(query).Execute()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprintf(os.Stderr, "Error when calling `Document.Multi``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	return
}

func (searchItemService *searchItemService) CreateShareLog(requestData string, responseData string) (err error) {
	index := "share_log" // string | Index
	var query string
	fmt.Println("创建zin日志")
	document := map[string]interface{}{

		"request_data": requestData,
		"resonse_data": responseData,
	}
	jsonData, err := json.Marshal(document)
	if err != nil {
		fmt.Println(err)
		return
	}

	query += string(jsonData) + "\n"

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
	_, r, err := apiClient.Document.Multi(ctx, index).Query(query).Execute()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprintf(os.Stderr, "Error when calling `Document.Multi``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	return
}
