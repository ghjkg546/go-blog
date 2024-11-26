package adminapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/app/services"
	"github.com/jassue/jassue-gin/global"
	"github.com/jassue/jassue-gin/utils"
	"net/http"
	"strconv"
	"strings"
)

type FeedbackController struct{}

// GetList handles GET requests
func (uc *FeedbackController) GetList(c *gin.Context) {
	var Feedbacks []models.Feedback
	var totalFeedbacks int64
	pageStr := c.DefaultQuery("pageNum", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	keyword := c.DefaultQuery("keywords", "")
	statusStr := c.DefaultQuery("status", "")
	startTime := c.DefaultQuery("createAt[0]", "")
	endTime := c.DefaultQuery("createAt[1]", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	// Calculate offset and limit
	offset := (page - 1) * pageSize
	limit := pageSize
	db := global.App.DB

	query := db.Model(models.Feedback{}).Table("feedback AS c")

	if keyword != "" {
		query = query.Where("c.content LIKE ? OR user.mobile LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
			Joins("LEFT JOIN user ON user.id = c.user_id")
	}
	if startTime != "" {
		startstamp, endstamp := services.CrudService.ParseStartEndTime(startTime, endTime)
		query.Where("created_at between ? and ?", startstamp, endstamp)
	}

	if (statusStr) != "" {
		status, err := strconv.Atoi(statusStr)
		if err != nil || pageSize < 1 {
			status = 10
		}
		query = query.Where("c.status = ?", status)

	}

	query.Count(&totalFeedbacks).Offset(offset).Limit(limit).Preload("User").Find(&Feedbacks)
	for i, item := range Feedbacks {
		Feedbacks[i].CreatedAtStr = utils.TimestampToDateYmd(item.CreatedAt)
	}
	response.Success(c, gin.H{
		"list":  Feedbacks,
		"total": totalFeedbacks,
	})
}

// GetDetail handles GET requests for dict details
func (uc *FeedbackController) GetDetail(c *gin.Context) {
	id := c.Param("id")
	var dict models.Feedback
	db := global.App.DB
	db.Model(models.Feedback{}).First(&dict, id)

	response.Success(c, dict)
}

// Create handles POST requests to create a new dict
func (uc *FeedbackController) Create(c *gin.Context) {
	var input models.Feedback
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
func (uc *FeedbackController) Update(c *gin.Context) {
	var input models.Feedback

	// Bind JSON payload to input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	global.App.DB.Where("id = ?", input.ID).Updates(models.Feedback{Status: 1, Message: input.Message})

	response.Success(c, nil)

}

// Delete handles DELETE requests to delete a dict
func (uc *FeedbackController) Delete(c *gin.Context) {
	idsString := c.Param("id")
	fmt.Printf(idsString)
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
		global.App.DB.Where("id IN ?", ids).Delete(models.Feedback{})
	}
	response.Success(c, nil)
}
