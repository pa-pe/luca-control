package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pa-pe/luca-control/src/storage/model"
	"github.com/pa-pe/luca-control/src/utils"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// Indicates whether it's the first run of the application
var IsFirstRun bool

// ShowInitialSetupPage displays the initial setup page for creating the first admin user
func ShowInitialSetupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "initial_setup.html", nil)
}

// HandleInitialSetup processes the form submission for creating the first admin user
func HandleInitialSetup(c *gin.Context, db *gorm.DB) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Check if the user already exists
	var existingUser model.WebUser
	if err := db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusConflict, "initial_setup.html", gin.H{"error": "User already exists"})
		return
	}

	// Hash the password and create the user
	hashedPassword, err := utils.HashStr(password)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "initial_setup.html", gin.H{"error": "Failed to hash password"})
		return
	}

	// Create the admin user
	newUser := model.WebUser{
		Username: username,
		Password: hashedPassword,
		Role:     "admin",
	}
	if err := db.Create(&newUser).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "initial_setup.html", gin.H{"error": "Failed to create admin user"})
		return
	}

	// Disable first run mode
	IsFirstRun = false

	log.Printf("InitialSetup finished with username=%s", username)

	// Redirect to the main page
	c.Redirect(http.StatusSeeOther, "/")
}
