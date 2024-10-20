package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pa-pe/luca-control/src/storage/model"
	"gorm.io/gorm"
	"net/http"
)

func ListTgUsers(c *gin.Context, db *gorm.DB) {
	currentAuthUser := GetCurrentAuthUser(c)

	var tgUsers []model.TgUser
	if err := db.Find(&tgUsers).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving Telegram users")
		return
	}

	c.HTML(http.StatusOK, "tg_users.tmpl", gin.H{
		"Title":       "Telegram Users",
		"CurrentUser": currentAuthUser.Username,
		"tgUsers":     tgUsers,
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

// UpdateModel updates the fields of a specified model based on the allowedFields
func UpdateModel(c *gin.Context, db *gorm.DB) {
	// Configuration for allowed fields
	var allowedFields = map[string][]string{
		"TgUser": {"chatbot_permit", "first_name", "last_name", "language_code", "shift_state"},
		"TgMsg":  {"text", "is_outgoing", "reply_to_message_id"},
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	modelName, ok := payload["model"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Model name is required"})
		return
	}

	id, ok := payload["id"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	// Retrieve allowed fields for the specified model
	allowed, ok := allowedFields[modelName]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Model is not allowed for updating"})
		return
	}

	// Prepare the map for updating
	updateData := make(map[string]interface{})
	for _, field := range allowed {
		if value, exists := payload[field]; exists {
			updateData[field] = value
		}
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	// Execute the update
	var result interface{}
	switch modelName {
	case "TgUser":
		result = &model.TgUser{ID: int64(id)}
	//case "TgMsg":
	//	result = &TgMsg{InternalID: int64(id)}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported model"})
		return
	}

	if err := db.Model(result).Debug().Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update model"})
		return
	}

	// client side js tests:
	//time.Sleep(1000 * time.Millisecond)
	//c.JSON(http.StatusBadRequest, gin.H{"error": "test fail"})
	//return

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Model updated successfully"})
}
