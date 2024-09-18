package adminapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/app/services"
	"github.com/jassue/jassue-gin/global"
	"net/http"
	"strconv"
	"time"
)

type CateController struct{}

func (uc *CateController) AllCateList(c *gin.Context) {
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

func (uc *CateController) GetList(c *gin.Context) {
	var users []models.Category
	var totalUsers int64
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	keyword := c.DefaultQuery("keyword", "")
	status := c.DefaultQuery("status", "")

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

	query := db.Model(models.Category{})
	if keyword != "" {
		query.Or("name LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		query.Or("status = ?", status)
	}

	query.Find(&users).Count(&totalUsers).Limit(limit).Offset(offset)
	for i := range users {
		user := &users[i]
		tm1 := time.Unix(user.CreatedAt, 0)
		tm2 := time.Unix(user.UpdatedAt, 0)
		user.CreateTimeStr = tm1.Format("2006-01-02 15:04:05")
		user.UpdateTimeStr = tm2.Format("2006-01-02 15:04:05")
	}
	var res = gin.H{
		"list":  users,
		"total": totalUsers,
	}
	response.Success(c, res)
}

// GetDetail handles GET requests for user details
func (uc *CateController) GetDetail(c *gin.Context) {
	id := c.Param("id")
	var user models.Category
	db := global.App.DB
	db.Model(models.Category{}).First(&user, id)
	response.Success(c, user)
}

// Create handles POST requests to create a new user
func (uc *CateController) Create(c *gin.Context) {
	var input models.Category
	db := global.App.DB
	// Bind JSON payload to input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save data to database
	result := db.Save(&input)
	if result.Error != nil {
		response.Fail(c, 500, result.Error.Error())

		return
	}
	// Example response, replace with actual logic
	response.Success(c, nil)
}

// Update handles PUT requests to update a user
func (uc *CateController) Update(c *gin.Context) {
	var input models.Category
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
func (uc *CateController) Delete(c *gin.Context) {
	id := c.Param("id")
	db := global.App.DB

	db.Delete(&models.Category{}, id)
	response.Success(c, nil)
}
