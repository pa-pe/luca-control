package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pa-pe/luca-control/src/storage/model"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// UpdateModel updates the fields of a specified model based on the allowedFields
func UpdateModel(c *gin.Context, db *gorm.DB) {
	currentAuthUser := GetCurrentAuthUser(c)
	if currentAuthUser.Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Configuration for allowed fields
	var allowedFields = map[string][]string{
		"TgUser": {"chatbot_permit", "SrvsEmployeesId"},
		//"TgMsg":  {"text", "is_outgoing", "reply_to_message_id"},
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

	fmt.Println(payload)

	id, ok := payload["id"].(float64)
	if !ok {
		idStr, ok := payload["id"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		idFromStr, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID recognize"})
			return
		}
		id = float64(idFromStr)
	}

	// Retrieve allowed fields for the specified model
	allowed, ok := allowedFields[modelName]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Model is not allowed for updating"})
		return
	}

	//Prepare the map for updating
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

	// Retrieve original values for fields to be updated
	originalData := make(map[string]interface{})
	if err := db.Model(&model.TgUser{}).Where("id = ?", int64(id)).Select(allowed).Take(&originalData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve original data"})
		return
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	if err := db.Model(result).Debug().Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update model"})
		return
	}

	// Log each change to DbChanges
	for field, newValue := range updateData {
		originalValue, exists := originalData[field]
		if exists && originalValue != newValue {
			dbChange := model.DbChanges{
				WebUserID: currentAuthUser.ID,
				ModelName: modelName + "." + field,
				DataFrom:  fmt.Sprintf("%v", originalValue),
				DataTo:    fmt.Sprintf("%v", newValue),
			}
			if err := db.Create(&dbChange).Error; err != nil {
				log.Printf("Failed to log change for field %s: %v", field, err)
			}
		}
	}

	// client side js tests:
	//time.Sleep(1000 * time.Millisecond)
	//c.JSON(http.StatusBadRequest, gin.H{"error": "test fail"})
	//return

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Model updated successfully"})
}
