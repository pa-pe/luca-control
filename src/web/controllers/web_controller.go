package controllers

import (
	tgmodels "github.com/pa-pe/luca-control/src/storage/model"
	"github.com/pa-pe/luca-control/src/utils"
	webmodels "github.com/pa-pe/luca-control/src/web/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WebController отвечает за обработку запросов веб-интерфейса
type WebController struct {
	DB *gorm.DB
}

// NewWebController создает новый экземпляр WebController
func NewWebController(db *gorm.DB) *WebController {
	return &WebController{DB: db}
}

// ShowTgUsers отображает список пользователей Telegram
func (wc *WebController) ShowTgUsers(c *gin.Context) {
	var users []tgmodels.TgUser
	result := wc.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных"})
		return
	}

	c.HTML(http.StatusOK, "tg_users.html", gin.H{"users": users})
}

func ShowMainMenu(c *gin.Context) {
	c.HTML(http.StatusOK, "main_menu.html", nil)
}

func ShowAddUserForm(c *gin.Context) {
	userRole, exists := c.Get("user_role")
	if !exists || userRole != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.HTML(http.StatusOK, "add_user.html", nil)
}

func AddWebUserHandler(c *gin.Context, db *gorm.DB) {
	userRole, exists := c.Get("user_role")
	if !exists || userRole != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	role := c.PostForm("role")

	// Хешируем пароль перед сохранением
	hashedPassword, err := utils.HashStr(password)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while hashing password")
		return
	}

	newUser := webmodels.WebUser{
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

func ListTgUsers(c *gin.Context, db *gorm.DB) {
	var tgUsers []tgmodels.TgUser
	if err := db.Find(&tgUsers).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving Telegram users")
		return
	}

	c.HTML(http.StatusOK, "tg_users.html", gin.H{
		"tgUsers": tgUsers,
	})
}

func ListWebUsers(c *gin.Context, db *gorm.DB) {
	var webUsers []webmodels.WebUser
	if err := db.Find(&webUsers).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving web users")
		return
	}

	c.HTML(http.StatusOK, "web_users.html", gin.H{
		"webUsers": webUsers,
	})
}
