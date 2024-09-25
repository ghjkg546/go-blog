package common

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/app/services"
	"github.com/jassue/jassue-gin/global"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"regexp"
	"strings"
)

func ImageUpload(c *gin.Context) {
	var form request.ImageUpload
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	outPut, err := services.MediaService.SaveImage(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, outPut)
}

type Record struct {
	Title string
	Path  string
	Url   string
}

func CsvUpload(c *gin.Context) {
	var form request.FileUpload
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	form.Business = "csv"
	outPut, err := services.MediaService.SaveFile(form)
	fmt.Println(outPut.Url)
	apath := "storage/app/public/" + outPut.Path
	file, err := ioutil.ReadFile(apath)
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := csv.NewReader(transform.NewReader(bytes.NewReader(file), simplifiedchinese.GBK.NewDecoder()))
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	var result []Record
	db := global.App.DB
	for _, record := range records[1:] { // Skip header
		var items []models.NetDiskItem
		items = append(items, models.NetDiskItem{Url: record[2]})
		b, err := json.Marshal(items)
		if err != nil {
			fmt.Printf("json.Marshal failed, err:%v\n", err)
			continue
		}

		content := record[1]
		newStr := strings.Replace(content, "https://p1.im0db.com/s_ratio_poster/public/", "https://api.shareziyuan.email/storage/local/cover_img/", -1)

		re := regexp.MustCompile(`https://pan\.quark\.cn/s/[a-zA-Z0-9]+`)
		newStr = re.ReplaceAllString(newStr, record[2])
		parts := strings.Split(record[3], ",")
		firstElement := parts[0]
		db.Create(&models.ResourceItem{CoverImg: "https://api.shareziyuan.email/storage/local/cover_img/" + firstElement, Views: 0, Title: record[0], DiskItems: string(b), Description: newStr, CategoryId: form.CategoryId, Status: 1})
		result = append(result, Record{Title: record[0], Path: record[1], Url: record[2]})
	}
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, outPut)
}
