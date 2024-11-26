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

type CommentController struct{}

// GetList handles GET requests for Comments
func (uc *CommentController) GetList(c *gin.Context) {
	var Comments []models.Comment
	var totalComments int64
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	keyword := c.DefaultQuery("keyword", "")
	startTime := c.DefaultQuery("createTime[0]", "")
	endTime := c.DefaultQuery("createTime[1]", "")
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

	query := db.Model(models.Comment{}).Table("comment AS c")
	if keyword != "" {
		query = query.Where("c.content LIKE ? OR user.mobile LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
			Joins("LEFT JOIN user ON user.id = c.user_id")
	}
	if startTime != "" {
		startTimeStamp, endTimeStamp := services.CrudService.ParseStartEndTime(startTime, endTime)
		query.Where("created_at between ? and ?", startTimeStamp, endTimeStamp)
	}
	query.Count(&totalComments).Offset(offset).Limit(limit).Preload("User").Find(&Comments)
	for i, item := range Comments {
		Comments[i].CreatedAtStr = utils.TimestampToDateYmd(item.CreatedAt)
	}
	response.Success(c, gin.H{
		"list":  Comments,
		"total": totalComments,
	})
}

// GetDetail handles GET requests for Comment details
func (uc *CommentController) GetDetail(c *gin.Context) {
	id := c.Param("id")
	var Comment models.Comment

	db := global.App.DB
	db.Model(models.Comment{}).First(&Comment, id)
	//Comment.Password = ""
	response.Success(c, Comment)
}

// Create handles POST requests to create a new Comment
func (uc *CommentController) Create(c *gin.Context) {
	// Example response, replace with actual logic
	c.JSON(201, gin.H{"message": "Comment created"})
}

// Update handles PUT requests to update a Comment
func (uc *CommentController) Update(c *gin.Context) {
	var input models.Comment
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

// Delete handles DELETE requests to delete a Comment
func (uc *CommentController) Delete(c *gin.Context) {
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
		global.App.DB.Where("id IN ?", ids).Delete(models.Comment{})
	}
	response.Success(c, nil)
}
