package adminapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/request"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/app/services"
	"github.com/jassue/jassue-gin/global"
	"github.com/jassue/jassue-gin/utils"
	"net/http"
	"strconv"
)

type AppUserController struct{}

func (t *AppUserController) RoleList(c *gin.Context) {
	// 执行数据库操作
	db := global.App.DB

	// Fetch all routes
	var roles []models.SysRoleMenu
	db.Where("role_id = ?", 2).Find(&roles)
	var ids []uint
	for i := range roles {
		route := &roles[i]
		ids = append(ids, route.MenuID)

	}
	var menus []models.SysMenu

	db.Where("id in ?", ids).Order("parent_id asc").Find(&menus)

	routeMap := make(map[uint]*models.Route)
	var rootRoutes []*models.Route
	for i := range menus {
		menu := &menus[i]
		route := &models.Route{
			ID:        menu.ID,
			Path:      menu.RoutePath,
			Component: menu.Component,
			Redirect:  menu.Redirect,
			Name:      menu.Name,
			ParentID:  menu.ParentID,

			MetaStr: "",
		}
		meta := models.Meta{
			Title:      menu.Name,
			Icon:       menu.Icon,
			AlwaysShow: true,
		}
		route.Meta = meta
		route.Children = []models.Route{}
		routeMap[menu.ID] = route

		if route.ParentID == 0 {

			rootRoutes = append(rootRoutes, route)
		}

	}
	for _, route := range routeMap {
		if route.ParentID != 0 {
			parent, exists := routeMap[route.ParentID]
			if exists {
				parent.Children = append(parent.Children, *route)
			}
		}
	}
	//

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rootRoutes,
	})
}

func (uc *AppUserController) ChangePass(c *gin.Context) {
	var input request.ChangePass
	if err := c.BindJSON(&input); err != nil {
		response.Fail(c, 500, "参数出错")
		return
	}
	uid := services.JwtService.GetUserId(c)
	var user models.User
	fmt.Println("ui:d")
	fmt.Println(uid)
	global.App.DB.First(&user, uid)
	fmt.Println(user)
	if !utils.BcryptMakeCheck([]byte(input.OldPassword), user.Password) {
		response.Fail(c, 500, "原密码错误")
		return
	}
	//user.Password = utils.BcryptMake([]byte(input.NewPassword))
	global.App.DB.Model(&user).Where("id = ?", uid).Update("password", utils.BcryptMake([]byte(input.NewPassword)))

	response.Success(c, nil)
}

// GetList handles GET requests for user details
func (uc *AppUserController) GetList(c *gin.Context) {
	var users []models.User
	var totalUsers int64
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	keyword := c.DefaultQuery("keyword", "")

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

	query := db.Model(models.User{})
	if keyword != "" {
		query.Where("mobile LIKE ?", "%"+keyword+"%").Or("name LIKE ?", "%"+keyword+"%")
	}

	query.Find(&users).Count(&totalUsers).Limit(limit).Offset(offset)
	//for _, route := range routeMap
	for i := range users {
		user := &users[i]
		user.Password = ""
	}
	var res = gin.H{
		"list":  users,
		"total": totalUsers,
	}
	response.Success(c, res)
}

// GetDetail handles GET requests for user details
func (uc *AppUserController) GetDetail(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	db := global.App.DB
	db.Model(models.User{}).First(&user, id)
	user.Password = ""
	response.Success(c, user)
}

// Create handles POST requests to create a new user
func (uc *AppUserController) Create(c *gin.Context) {
	// Example response, replace with actual logic
	c.JSON(201, gin.H{"message": "User created"})
}

// Update handles PUT requests to update a user
func (uc *AppUserController) Update(c *gin.Context) {
	var input request.EditUser
	// Bind JSON payload to input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Password != "" {
		global.App.DB.Model(&models.User{}).Where("id = ?", input.Id).Update("password", utils.BcryptMake([]byte(input.Password)))
	}

	global.App.DB.Where("id = ?", input.Id).Updates(models.User{Mobile: input.Mobile, Name: input.Name})

	response.Success(c, nil)

}

// Delete handles DELETE requests to delete a user
func (uc *AppUserController) Delete(c *gin.Context) {
	id := c.Param("id")
	global.App.DB.Delete(&models.User{}, id)
	response.Success(c, nil)
}
