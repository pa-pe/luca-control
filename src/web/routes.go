package web

import (
	controllers2 "github.com/pa-pe/luca-control/src/web/controllers"
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalf("Could not set trusted proxies: %v", err)
	}

	router.SetHTMLTemplate(template.Must(template.ParseGlob("web/templates/*.html")))

	// Авторизация
	router.GET("/login", controllers2.ShowLoginPage)
	//    router.POST("/login", controllers.HandleLogin)
	router.POST("/login", func(c *gin.Context) { controllers2.HandleLogin(c, db) })

	// Маршруты для авторизованных пользователей
	authorized := router.Group("/")
	authorized.Use(controllers2.AuthRequired())
	{
		authorized.GET("/", controllers2.ShowMainMenu)
		authorized.GET("/tg_users", func(c *gin.Context) { controllers2.ListTgUsers(c, db) })
		authorized.GET("/web_users", func(c *gin.Context) { controllers2.ListWebUsers(c, db) })
		authorized.GET("/web_users/add", controllers2.ShowAddUserForm)
		authorized.POST("/web_users/add", func(c *gin.Context) { controllers2.AddWebUserHandler(c, db) })
	}
}
