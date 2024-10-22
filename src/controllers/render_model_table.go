package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pa-pe/luca-control/src/utils"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type modelConfig struct {
	PageTitle   string            `json:"pageTitle"`
	DbTable     string            `json:"dbTable"`
	Fields      []string          `json:"fields"`
	Headers     map[string]string `json:"headers"`
	Classes     map[string]string `json:"classes"`
	RelatedData map[string]string `json:"relatedData"`
}

func RenderModel(c *gin.Context, db *gorm.DB) {
	currentAuthUser := GetCurrentAuthUser(c)
	modelName := c.Param("modelName")

	config, err := loadModelConfig("config/renderModelTable/" + modelName + ".json")
	if err != nil {
		log.Print("No configuration found for RenderModel: " + modelName)
		c.String(http.StatusNotFound, "No configuration found for RenderModel: "+modelName)
		return
	}

	if config.PageTitle == "" {
		config.PageTitle = modelName
	}
	if config.DbTable == "" {
		config.DbTable = utils.CamelToSnake(modelName)
	}

	htmlTable, err := RenderModelTable(db, modelName, config)
	if err != nil {
		fmt.Println("Ошибка:", err)
		c.String(http.StatusNotFound, "RenderModel "+modelName+" not found")
		return
	}

	c.HTML(http.StatusOK, "render_model_table.tmpl", gin.H{
		"Title":       config.PageTitle,
		"CurrentUser": currentAuthUser.Username,
		"Content":     template.HTML(htmlTable),
	})
}

func loadModelConfig(configPath string) (*modelConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config modelConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}
	return &config, nil
}

func RenderModelTable(db *gorm.DB, modelName string, config *modelConfig) (string, error) {
	if config == nil || modelName == "" {
		log.Fatalf("configuration not found for model: %s", modelName)
	}

	//	tableName := utils.CamelToSnake(modelName)
	//	tableName := utils.CamelToSnake(config.DbTable)
	var records []map[string]interface{}
	if err := db.Debug().Table(config.DbTable).Find(&records).Error; err != nil {
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
