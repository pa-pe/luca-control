package src

import (
	"github.com/gin-gonic/gin"
	controllers2 "github.com/pa-pe/luca-control/src/controllers"
	"github.com/pa-pe/luca-control/src/storage/model"
	"gorm.io/gorm"
	"html/template"
	"log"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Check if there are any web users in the database
	var userCount int64
	if err := db.Model(&model.WebUser{}).Count(&userCount).Error; err != nil {
		log.Fatalf("Failed to check web users count: %v", err)
	}
	// If no users exist, it's the first run
	controllers2.IsFirstRun = userCount == 0

	// Configure router settings
	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalf("Could not set trusted proxies: %v", err)
	}

	//	router.SetHTMLTemplate(template.Must(template.ParseGlob("web/templates/*.html")))
	router.SetHTMLTemplate(template.Must(template.ParseGlob("web/templates/*.*")))

	//router.GET("/login", controllers.ShowLoginPage)
	router.GET("/login", func(c *gin.Context) { controllers2.ShowLoginPage(c, db) })
	router.POST("/login", func(c *gin.Context) { controllers2.HandleLogin(c, db) })
	router.GET("/logout", func(c *gin.Context) { controllers2.HandleLogout(c, db) })

	// Auth users routes
	authorized := router.Group("/")
	authorized.Use(controllers2.AuthRequired(db, &controllers2.IsFirstRun))
	{
		authorized.GET("/", controllers2.ShowAuthMain)
		//		authorized.GET("/logout", controllers.HandleLogout)
		authorized.GET("/web_users", func(c *gin.Context) { controllers2.ListWebUsers(c, db) })
		authorized.GET("/web_users/add", controllers2.ShowAddWebUserForm)
		authorized.POST("/web_users/add", func(c *gin.Context) { controllers2.AddWebUserHandler(c, db) })
		authorized.GET("/tg_users", func(c *gin.Context) { controllers2.ListTgUsers(c, db) })
		authorized.GET("/tg_msgs_all", func(c *gin.Context) { controllers2.ListTgMsgsAll(c, db) })

		authorized.POST("/update_model", func(c *gin.Context) { controllers2.UpdateModel(c, db) })
	}

	// Show initial setup page if it's the first run
	if controllers2.IsFirstRun {
		router.GET("/initial-setup", controllers2.ShowInitialSetupPage)
		router.POST("/initial-setup", func(c *gin.Context) { controllers2.HandleInitialSetup(c, db) })
	}
}
