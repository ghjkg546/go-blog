package adminapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jassue/jassue-gin/app/common/response"
	"github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
	"net/http"
	"strconv"
)

type MenuController struct{}

func (uc *MenuController) Options(c *gin.Context) {
	var menus []models.SysMenu

	db := global.App.DB

	db.Model(models.SysMenu{}).Order("parent_id").Find(&menus)

	var rootRoutes []models.MenuOption
	for i := range menus {
		v := &menus[i]
		route := models.MenuOption{
			Value: v.ID,
			ID:    v.ID,
			//Path:      menu.RoutePath,
			ParentId: v.ParentID,
			Label:    v.Name,
			//Children:  child,
		}

		route.Children = []models.MenuOption{}

		rootRoutes = append(rootRoutes, route)
		//}

	}

	a := OptionListToTree(rootRoutes, 0, 1)
	response.Success(c, a)
}

// GetList handles GET requests for user details
func (uc *MenuController) GetList(c *gin.Context) {
	var menus []models.SysMenu
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

	db := global.App.DB

	query := db.Model(models.SysMenu{}).Order("parent_id").Count(&totalUsers)
	if keyword != "" {
		query.Where("name LIKE ?", "%"+keyword+"%")
	}
	query.Find(&menus)
	var rootRoutes []models.SysMenu
	for i := range menus {
		menu := &menus[i]
		route := models.SysMenu{
			ID: menu.ID,
			//Path:      menu.RoutePath,
			Component: menu.Component,
			Redirect:  menu.Redirect,
			Name:      menu.Name,
			ParentID:  menu.ParentID,
			Visible:   menu.Visible,
			Type:      menu.Type,
			//MetaStr:   "",
		}

		route.Children = []models.SysMenu{}

		rootRoutes = append(rootRoutes, route)
		//}

	}

	a := ListToTree(rootRoutes, 0, 1)
	response.Success(c, a)
}

func ListToTree(stuAll []models.SysMenu, pid uint, lev uint) []models.SysMenu {
	var goodArr []models.SysMenu
	for _, v := range stuAll {
		if v.ParentID == pid {
			// 这里可以理解为每次都从最原始的数据里面找出相对就的ID进行匹配，直到找不到就返回
			child := ListToTree(stuAll, v.ID, lev+1)
			if child == nil {
				child = []models.SysMenu{}
			}
			node := models.SysMenu{
				ID: v.ID,
				//Path:      menu.RoutePath,
				Component: v.Component,
				Redirect:  v.Redirect,
				Name:      v.Name,
				ParentID:  v.ParentID,
				Visible:   v.Visible,
				Type:      v.Type,
				Level:     lev,
				Children:  child,
			}

			goodArr = append(goodArr, node)
		}
	}
	return goodArr
}

func OptionListToTree(stuAll []models.MenuOption, pid uint, lev uint) []models.MenuOption {
	var goodArr []models.MenuOption
	for _, v := range stuAll {
		if v.ParentId == pid {
			// 这里可以理解为每次都从最原始的数据里面找出相对就的ID进行匹配，直到找不到就返回
			child := OptionListToTree(stuAll, v.ID, lev+1)
			if child == nil {
				child = []models.MenuOption{}
			}
			node := models.MenuOption{
				Value: v.ID,
				//Path:      menu.RoutePath,
				ParentId: v.ParentId,
				Label:    v.Label,
				Children: child,
			}

			goodArr = append(goodArr, node)
		}
	}
	return goodArr
}

// GetDetail handles GET requests for user details
func (uc *MenuController) GetDetail(c *gin.Context) {
	id := c.Param("id")
	var user models.SysMenu
	db := global.App.DB
	db.Model(models.SysMenu{}).First(&user, id)
	response.Success(c, user)
}

// Create handles POST requests to create a new user
func (uc *MenuController) Create(c *gin.Context) {
	var input models.SysMenu
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

// Update handles PUT requests to update a user
func (uc *MenuController) Update(c *gin.Context) {
	var input models.SysMenu
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
func (uc *MenuController) Delete(c *gin.Context) {
	id := c.Param("id")
	db := global.App.DB
	var menu models.SysMenu
	db.Where("parent_id=?", id).First(&menu)
	fmt.Println("id:")
	fmt.Println(id)
	if menu.ID > 0 {
		response.Fail(c, 500, "该分类下还有子菜单")
		return
	}
	db.Delete(&models.SysMenu{}, id)
	response.Success(c, nil)
}
