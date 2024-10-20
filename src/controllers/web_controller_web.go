package controllers

import (
	"github.com/pa-pe/luca-control/src/storage/model"
	"github.com/pa-pe/luca-control/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebController struct {
	DB *gorm.DB
}

//func NewWebController(db *gorm.DB) *WebController {
//	return &WebController{DB: db}
//}

func ShowAuthMain(c *gin.Context) {
	currentAuthUser := GetCurrentAuthUser(c)

	c.HTML(http.StatusOK, "auth_main.tmpl", gin.H{
		"Title":       "Main Menu",
		"CurrentUser": currentAuthUser.Username,
	})
}

func ListWebUsers(c *gin.Context, db *gorm.DB) {
	currentAuthUser := GetCurrentAuthUser(c)

	var webUsers []model.WebUser
	if err := db.Find(&webUsers).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving web users")
		return
	}

	c.HTML(http.StatusOK, "web_users.tmpl", gin.H{
		"Title":       "Web Users",
		"CurrentUser": currentAuthUser.Username,
		"webUsers":    webUsers,
	})
}

func ShowAddWebUserForm(c *gin.Context) {
	currentAuthUser := GetCurrentAuthUser(c)
	//userRole, exists := c.Get("user_role")
	//if !exists || userRole != "admin" {
	//	c.AbortWithStatus(http.StatusForbidden)
	//	return
	//}

	c.HTML(http.StatusOK, "web_user_add.tmpl", gin.H{
		"Title":       "Add New Web User",
		"CurrentUser": currentAuthUser.Username,
	})
}

func AddWebUserHandler(c *gin.Context, db *gorm.DB) {
	currentAuthUser := GetCurrentAuthUser(c)
	//	userRole, exists := c.Get("user_role")
	//	if !exists || userRole != "admin" {
	if currentAuthUser.Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	role := c.PostForm("role")

	hashedPassword, err := utils.HashStr(password)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while hashing password")
		return
	}

	newUser := model.WebUser{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error while creating user")
		return
	}

	c.Redirect(http.StatusSeeOther, "/web_users")
}
