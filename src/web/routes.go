package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pa-pe/luca-control/src/web/controllers"
	"github.com/pa-pe/luca-control/src/web/models"
	"gorm.io/gorm"
	"html/template"
	"log"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Check if there are any web users in the database
	var userCount int64
	if err := db.Model(&models.WebUser{}).Count(&userCount).Error; err != nil {
		log.Fatalf("Failed to check web users count: %v", err)
	}
	// If no users exist, it's the first run
	controllers.IsFirstRun = userCount == 0

	// Configure router settings
	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalf("Could not set trusted proxies: %v", err)
	}

	router.SetHTMLTemplate(template.Must(template.ParseGlob("web/templates/*.html")))

	router.GET("/login", controllers.ShowLoginPage)
	router.POST("/login", func(c *gin.Context) { controllers.HandleLogin(c, db) })

	// Auth users routes
	authorized := router.Group("/")
	authorized.Use(controllers.AuthRequired(db, &controllers.IsFirstRun))
	{
		authorized.GET("/", controllers.ShowMainMenu)
		authorized.GET("/tg_users", func(c *gin.Context) { controllers.ListTgUsers(c, db) })
		authorized.GET("/web_users", func(c *gin.Context) { controllers.ListWebUsers(c, db) })
		authorized.GET("/web_users/add", controllers.ShowAddUserForm)
		authorized.POST("/web_users/add", func(c *gin.Context) { controllers.AddWebUserHandler(c, db) })
	}

	// Show initial setup page if it's the first run
	if controllers.IsFirstRun {
		router.GET("/initial-setup", controllers.ShowInitialSetupPage)
		router.POST("/initial-setup", func(c *gin.Context) { controllers.HandleInitialSetup(c, db) })
	}
}
