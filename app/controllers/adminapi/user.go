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

type UserController struct{}

// 分类列表
func LogOut(c *gin.Context) {
	fmt.Println("logout")
	var res = gin.H{
		"list":  []interface{}{},
		"total": 0,
	}
	response.Success(c, res)
}

func (t *UserController) Login(c *gin.Context) {
	var form request.Login
	if err := c.ShouldBindJSON(&form); err != nil {
		fmt.Println(err.Error())
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if err, user := services.UserService.Login(form); err != nil {
		fmt.Println(err.Error())
		response.BusinessFail(c, err.Error())
	} else {
		tokenData, err, _ := services.JwtService.CreateToken(services.AppGuardName, user)
		if err != nil {
			fmt.Println(err.Error())
			response.BusinessFail(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "",
			"data": gin.H{
				"accessToken":  tokenData.AccessToken,
				"tokenType":    "Bearer",
				"refreshToken": nil,
				"username":     "admin",
				"password":     "admin",
				"role":         "admin",
				"roleId":       1,
				"permissions":  []string{"*.*.*"},
			},
		})
	}
}

func (t *UserController) RoleList(c *gin.Context) {
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

func (t *UserController) Me(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"userId":   2,
			"username": "admin",
			"nickname": "系统管理员",
			"avatar":   "https://foruda.gitee.com/images/1723603502796844527/03cdca2a_716974.gif",
			"roles":    []string{"ADMIN"},
			"perms":    []string{},
		},
		"msg": "一切ok",
	})
}

func (t *UserController) Logout(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "",
		"data":    "ok",
	})
}

func (uc *UserController) ChangePass(c *gin.Context) {
	var input request.ChangePass
	if err := c.BindJSON(&input); err != nil {
		response.Fail(c, 500, "参数出错")
		return
	}
	uid := services.JwtService.GetUserId(c)
	var user models.AdminUser
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
func (uc *UserController) GetList(c *gin.Context) {
	var users []models.AdminUser
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

	query := db.Model(models.AdminUser{})
	if keyword != "" {
		query.Where("mobile LIKE ?", "%"+keyword+"%").Or("name LIKE ?", "%"+keyword+"%")
	}

	query.Find(&users).Count(&totalUsers).Limit(limit).Offset(offset)
	//sql := fmt.Sprintf("SELECT * FROM user limit %d offset %d", limit, offset)
	var res = gin.H{
		"list":  users,
		"total": totalUsers,
	}
	response.Success(c, res)
}

// GetDetail handles GET requests for user details
func (uc *UserController) GetDetail(c *gin.Context) {
	id := c.Param("id")
	var user models.AdminUser

	db := global.App.DB
	db.Model(models.AdminUser{}).First(&user, id)
	user.Password = ""
	response.Success(c, user)
}

// Create handles POST requests to create a new user
func (uc *UserController) Create(c *gin.Context) {
	// Example response, replace with actual logic
	c.JSON(201, gin.H{"message": "User created"})
}

// Update handles PUT requests to update a user
func (uc *UserController) Update(c *gin.Context) {
	var input models.AdminUser
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
	// Example response, replace with actual logic
	response.Success(c, nil)

}

// Delete handles DELETE requests to delete a user
func (uc *UserController) Delete(c *gin.Context) {
	id := c.Param("id")
	global.App.DB.Delete(&models.AdminUser{}, id)
	c.JSON(200, gin.H{"message": "User deleted", "id": id})
}
