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

	//	router.SetHTMLTemplate(template.Must(template.ParseGlob("web/templates/*.html")))
	router.SetHTMLTemplate(template.Must(template.ParseGlob("web/templates/*.*")))

	//router.GET("/login", controllers.ShowLoginPage)
	router.GET("/login", func(c *gin.Context) { controllers.ShowLoginPage(c, db) })
	router.POST("/login", func(c *gin.Context) { controllers.HandleLogin(c, db) })
	router.GET("/logout", func(c *gin.Context) { controllers.HandleLogout(c, db) })

	// Auth users routes
	authorized := router.Group("/")
	authorized.Use(controllers.AuthRequired(db, &controllers.IsFirstRun))
	{
		authorized.GET("/", controllers.ShowAuthMain)
		//		authorized.GET("/logout", controllers.HandleLogout)
		authorized.GET("/web_users", func(c *gin.Context) { controllers.ListWebUsers(c, db) })
		authorized.GET("/web_users/add", controllers.ShowAddWebUserForm)
		authorized.POST("/web_users/add", func(c *gin.Context) { controllers.AddWebUserHandler(c, db) })
		authorized.GET("/tg_users", func(c *gin.Context) { controllers.ListTgUsers(c, db) })
		authorized.GET("/tg_msgs_all", func(c *gin.Context) { controllers.ListTgMsgsAll(c, db) })

		authorized.POST("/update_model", func(c *gin.Context) { controllers.UpdateModel(c, db) })
	}

	// Show initial setup page if it's the first run
	if controllers.IsFirstRun {
		router.GET("/initial-setup", controllers.ShowInitialSetupPage)
		router.POST("/initial-setup", func(c *gin.Context) { controllers.HandleInitialSetup(c, db) })
	}
}
