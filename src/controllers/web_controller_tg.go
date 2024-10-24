package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pa-pe/luca-control/src/storage/model"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func ListTgUsers(c *gin.Context, db *gorm.DB) {
	currentAuthUser := GetCurrentAuthUser(c)

	var tgUsers []model.TgUser
	if err := db.Find(&tgUsers).Error; err != nil {
		log.Printf("Error retrieving Telegram users: %v", err)
		c.String(http.StatusInternalServerError, "Error retrieving Telegram users")
		return
	}

	var srvsEmployeesList []model.SrvsEmployeesList
	if err := db.Find(&srvsEmployeesList).Error; err != nil {
		log.Printf("Error retrieving Employees users: %v", err)
		c.String(http.StatusInternalServerError, "Error retrieving Employees users")
		return
	}

	c.HTML(http.StatusOK, "tg_users.tmpl", gin.H{
		"Title":             "Telegram Users",
		"CurrentUser":       currentAuthUser.Username,
		"tgUsers":           tgUsers,
		"srvsEmployeesList": srvsEmployeesList,
	})
}

func ListTgMsgsAll(c *gin.Context, db *gorm.DB) {
	currentAuthUser := GetCurrentAuthUser(c)

	//var tgMsgs []model.TgMsg
	//if err := db.Find(&tgMsgs).Error; err != nil {
	//	c.String(http.StatusInternalServerError, "Error retrieving Telegram users")
	//	return
	//}

	var tgMsgWithUserName []model.TgMsgWithUserName

	// Making a connection between TgMsg and TgUser
	if err := db.Table("tg_msgs").
		Select("tg_msgs.internal_id, tg_msgs.tg_id, tg_msgs.tg_user_id, tg_users.user_name, tg_msgs.chat_id, tg_msgs.reply_to_message_id, tg_msgs.is_outgoing, tg_msgs.text, tg_msgs.date, tg_msgs.added_timestamp").
		Joins("left join tg_users on tg_msgs.tg_user_id = tg_users.id").
		Order("internal_id DESC").
		Scan(&tgMsgWithUserName).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving Telegram users")
		return
	}

	c.HTML(http.StatusOK, "tg_msgs_all.tmpl", gin.H{
		"Title":       "All Telegram Messages",
		"CurrentUser": currentAuthUser.Username,
		"tgMsgs":      tgMsgWithUserName,
	})
}
