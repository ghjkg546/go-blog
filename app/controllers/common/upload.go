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
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func getTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// The header should be in the format "Bearer <token>"
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || splitToken[0] != "Bearer" {
		return ""
	}

	return splitToken[1]
}

func AvatarUpload(c *gin.Context) {
	var form request.ImageUpload
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	uid := services.AppUserService.GetUserId(c)
	if uid == "" {
		response.BusinessFail(c, "未登录或登录过期")
		return
	}
	outPut, err := services.MediaService.SaveImage(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	_ = UpdateuserAvater(uid, outPut.Url)
	response.Success(c, outPut)
}

// 更新用户分数，如果用户存在则加分
func UpdateuserAvater(uid string, url string) error {
	var user models.User

	// 根据 UID 查找用户
	result := global.App.DB.Where("id = ?", uid).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found")
		}
		return result.Error
	}

	user.Avatar = url

	// 保存更新后的用户
	if err := global.App.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

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

func InfoUpload(c *gin.Context) {
	// 获取表单字段
	id := c.PostForm("id")
	content := c.PostForm("content")
	fmt.Println(id, content)
	// 获取上传的图片
	file, header, err := c.Request.FormFile("img")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取上传的图片"})
		return
	}
	defer file.Close()

	// 确定图片保存路径
	saveDir := `./storage/app/public/local/cover_img`
	dbDir := `./storage/local/cover_img`
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建目录"})
		return
	}

	// 设置保存路径
	savePath := filepath.Join(saveDir, header.Filename)
	dbPath := filepath.Join(dbDir, header.Filename)

	// 保存图片
	out, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存图片"})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存图片失败"})
		return
	}
	// 原始 HTML 字符串
	originalHTML := content

	// 定义正则表达式匹配 img 标签中的 src
	re := regexp.MustCompile(`src="[^"]*"`) // 匹配 src 属性内容
	finalPath := global.App.Config.App.AppUrl + "/" + dbPath
	modifiedHTML := re.ReplaceAllString(originalHTML, fmt.Sprintf(`src="%s"`, finalPath))

	// 输出结果
	fmt.Println("修改后的HTML:", modifiedHTML)
	err1 := updateResourceItem(id, modifiedHTML, finalPath)
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	response.Success(c, "上传成功")

}

// 根据 ID 更新指定字段的内容
func updateResourceItem(id string, content string, imgPath string) error {
	// 使用 GORM 批量更新
	db := global.App.DB
	result := db.Model(&models.ResourceItem{}).Where("id = ?", id).Updates(map[string]interface{}{
		"description": content,
		"cover_img":   imgPath,
	})
	fmt.Println(result.RowsAffected)
	// 检查是否有影响的行
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到 ID 为 %d 的记录", id)
	}
	fmt.Println(result.Error)

	return nil
}

// 辅助函数：解析 ID
func parseID(id string) uint {
	var parsedID uint
	_, err := fmt.Sscanf(id, "%d", &parsedID)
	if err != nil {
		return 0
	}
	return parsedID
}
