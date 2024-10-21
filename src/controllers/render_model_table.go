package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pa-pe/luca-control/src/storage/model"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"os"
	"reflect"
	"strings"
)

// Структура для хранения конфигурации
type modelConfig struct {
	Fields  []string          `json:"fields"`
	Headers map[string]string `json:"headers"`
}

// Конфигурация всех моделей
var modelConfigs map[string]modelConfig

func RenderModel(c *gin.Context, db *gorm.DB) {
	currentAuthUser := GetCurrentAuthUser(c)

	err := loadModelConfig("config/renderModelTable.json")
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
	}

	htmlTable, err := RenderModelTable(db, "DbChanges")
	if err != nil {
		fmt.Println("Ошибка:", err)
		c.String(http.StatusInternalServerError, "Error retrieving Telegram users")
		return
		//} else {
		//	fmt.Println(htmlTable)
	}

	c.HTML(http.StatusOK, "render_model_table.tmpl", gin.H{
		"Title":       "RenderModel",
		"CurrentUser": currentAuthUser.Username,
		"Content":     template.HTML(htmlTable),
	})
}

// Функция для загрузки конфигурации из JSON
func loadModelConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, &modelConfigs); err != nil {
		return fmt.Errorf("failed to parse config JSON: %w", err)
	}
	return nil
}

// RenderModelTable Функция для рендеринга HTML таблицы
func RenderModelTable(db *gorm.DB, modelName string) (string, error) {
	config, ok := modelConfigs[modelName]
	if !ok {
		return "", fmt.Errorf("model %s is not configured for rendering", modelName)
	}

	// Получаем данные из базы данных
	var records []interface{}
	switch modelName {
	case "DbChanges":
		var dbChanges []model.DbChanges
		if err := db.Find(&dbChanges).Error; err != nil {
			return "", fmt.Errorf("failed to query database: %w", err)
		}
		for _, record := range dbChanges {
			records = append(records, record)
		}
	// Можно добавить другие модели здесь
	default:
		return "", fmt.Errorf("unsupported model %s", modelName)
	}

	// Построение заголовка таблицы
	var sb strings.Builder
	sb.WriteString("<table class='table'>\n<thead>\n<tr>\n")
	for _, field := range config.Fields {
		header := field
		if h, exists := config.Headers[field]; exists {
			header = h
		}
		sb.WriteString("<th>" + header + "</th>\n")
	}
	sb.WriteString("</tr>\n</thead>\n<tbody>\n")

	// Получение значений для каждого поля в записях
	for _, record := range records {
		recordValue := reflect.ValueOf(record)

		sb.WriteString("<tr>\n")
		for _, field := range config.Fields {
			fieldValue := recordValue.FieldByName(field)
			sb.WriteString("<td>" + fmt.Sprintf("%v", fieldValue) + "</td>\n")
		}
		sb.WriteString("</tr>\n")
	}
	sb.WriteString("</tbody>\n</table>")

	return sb.String(), nil
}
