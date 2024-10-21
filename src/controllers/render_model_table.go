package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pa-pe/luca-control/src/utils"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// Структура для хранения конфигурации
type modelConfig struct {
	Fields      []string          `json:"fields"`
	Headers     map[string]string `json:"headers"`
	Classes     map[string]string `json:"classes"`
	RelatedData map[string]string `json:"relatedData"`
}

// Конфигурация всех моделей
var modelConfigs map[string]modelConfig

func RenderModel(c *gin.Context, db *gorm.DB) {
	currentAuthUser := GetCurrentAuthUser(c)

	err := loadModelConfig("config/renderModelTable.json")
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
	}

	modelName := "DbChanges"
	htmlTable, err := RenderModelTable(db, modelName)
	if err != nil {
		fmt.Println("Ошибка:", err)
		c.String(http.StatusInternalServerError, "Error retrieving data from "+modelName)
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

func RenderModelTable(db *gorm.DB, modelName string) (string, error) {
	config, ok := modelConfigs[modelName]
	if !ok {
		return "", fmt.Errorf("configuration not found for model: %s", modelName)
	}

	tableName := utils.CamelToSnake(modelName)
	var records []map[string]interface{}
	if err := db.Debug().Table(tableName).Find(&records).Error; err != nil {
		return "", err
	}

	var htmlTable strings.Builder
	htmlTable.WriteString("<table class='table'>\n<thead><tr>")

	for _, field := range config.Fields {
		header := config.Headers[field]
		if header == "" {
			header = config.Headers[utils.InvertCaseStyle(field)]
		}
		if header == "" {
			header = field
		}
		htmlTable.WriteString(fmt.Sprintf("<th>%s</th>", header))
	}
	htmlTable.WriteString("</tr></thead>\n<tbody>")

	// Кэш для хранения полученных значений связанных данных
	relatedDataCache := make(map[string]string)

	for _, record := range records {
		htmlTable.WriteString("<tr>")
		for _, field := range config.Fields {
			//value := record[field]
			value, exists := record[field]
			if !exists || value == nil {
				value, exists = record[utils.InvertCaseStyle(field)]
				if !exists || value == nil {
					value = "" // Если значение отсутствует, выводим пустую строку
				}
			}

			classAttr := ""
			if class, ok := config.Classes[field]; ok {
				classAttr = fmt.Sprintf(" class='%s'", class)
			}

			relatedDataField, relatedExists := config.RelatedData[field]
			if relatedExists {
				//relatedDataField = utils.CamelToSnake(relatedDataField)
				cacheKey := fmt.Sprintf("%s_%v", relatedDataField, value)
				if cachedValue, found := relatedDataCache[cacheKey]; found {
					// Используем значение из кэша
					value = cachedValue
				} else {
					// Выполняем запрос к базе данных и добавляем в кэш
					var relatedValue string
					err := db.Debug().Table(strings.Split(relatedDataField, ".")[0]).
						Select(strings.Split(relatedDataField, ".")[1]).
						Where("id = ?", value).
						Row().Scan(&relatedValue)
					if err != nil {
						//return "", err
					} else {
						value = relatedValue
						relatedDataCache[cacheKey] = relatedValue
					}
				}
			}
			htmlTable.WriteString(fmt.Sprintf("<td%s>%v</td>", classAttr, value))
		}
		htmlTable.WriteString("</tr>\n")
	}
	htmlTable.WriteString("</tbody>\n</table>")

	return htmlTable.String(), nil
}
