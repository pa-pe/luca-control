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
	PageTitle         string                            `json:"pageTitle"`
	DbTable           string                            `json:"dbTable"`
	SqlWhere          string                            `json:"sqlWhere"`
	Fields            []string                          `json:"fields"`
	Headers           map[string]string                 `json:"headers"`
	Classes           map[string]string                 `json:"classes"`
	RelatedData       map[string]string                 `json:"relatedData"`
	AddableFields     []string                          `json:"addableFields"`
	RequiredFields    []string                          `json:"requiredFields"`
	NoZeroValueFields []string                          `json:"noZeroValueFields"`
	CountRelatedData  map[string]CountRelatedDataConfig `json:"countRelatedData"`
	Links             map[string]LinkConfig             `json:"links"`
	Parent            map[string]string                 `json:"parent"`
	ParentConfig      *modelConfig
}

type CountRelatedDataConfig struct {
	Table      string `json:"table"`
	ForeignKey string `json:"foreignKey"`
}

type LinkConfig struct {
	Template string `json:"template"`
}

func RenderModel(c *gin.Context, db *gorm.DB) {
	currentAuthUser := GetCurrentAuthUser(c)
	modelName := c.Param("modelName")

	config, err := loadModelConfig(c, modelName)
	if err != nil {
		log.Print("No configuration found for RenderModel: " + modelName)
		c.String(http.StatusNotFound, "No configuration found for RenderModel: "+modelName)
		return
	}

	htmlTable, err := RenderModelTable(db, modelName, config)
	if err != nil {
		fmt.Println("Ошибка:", err)
		errStr := fmt.Sprint("Ошибка:", err)
		c.String(http.StatusNotFound, "Error RenderModel "+modelName+": "+errStr)
		return
	}

	c.HTML(http.StatusOK, "render_model_table.tmpl", gin.H{
		"Title":       config.PageTitle,
		"CurrentUser": currentAuthUser.Username,
		"Content":     template.HTML(htmlTable),
	})
}

func loadModelConfig(c *gin.Context, modelName string) (*modelConfig, error) {
	configPath := "config/renderModelTable/" + modelName + ".json"

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config modelConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	if config.PageTitle == "" {
		config.PageTitle = modelName
	}
	if config.DbTable == "" {
		config.DbTable = utils.CamelToSnake(modelName)
	}

	if parentModelName, parentExists := config.Parent["modelName"]; parentExists {
		//fmt.Println(parentModelName)
		config.ParentConfig, err = loadModelConfig(c, parentModelName)
		if err != nil {
			log.Print("Can`t load ParentConfig: " + parentModelName)
		}
	}

	queryParams := c.Request.URL.Query()
	for key, values := range queryParams {
		for _, val := range values {
			//		fmt.Printf("Parameter: %s, Value: %s\n", key, value)
			placeholder := fmt.Sprintf("$%s$", key)
			config.SqlWhere = strings.ReplaceAll(config.SqlWhere, placeholder, fmt.Sprintf("%v", val))
		}
	}

	return &config, nil
}

func breadcrumbBuilder(config *modelConfig) string {
	breadcrumbStr := `<nav style="--bs-breadcrumb-divider: '>';" aria-label="breadcrumb">` + "\n"
	breadcrumbStr += `  <ol class="breadcrumb">` + "\n"
	breadcrumbStr += `    <li class="breadcrumb-item"><a href="/">Home</a></li>` + "\n"

	if parentModelName, parentExists := config.Parent["modelName"]; parentExists {
		//fmt.Println(parentModelName)
		breadcrumbStr += `    <li class="breadcrumb-item"><a href="/render_table/` + parentModelName + `"">` + config.ParentConfig.PageTitle + `</a></li>` + "\n"
	}

	breadcrumbStr += `    <li class="breadcrumb-item active" aria-current="page">` + config.PageTitle + `</li>` + "\n"
	breadcrumbStr += `  </ol>` + "\n"
	breadcrumbStr += `</nav>` + "\n"

	return breadcrumbStr
}

func RenderModelTable(db *gorm.DB, modelName string, config *modelConfig) (string, error) {
	if config == nil || modelName == "" {
		log.Fatalf("configuration not found for model: %s", modelName)
	}

	var records []map[string]interface{}
	if err := db.Debug().Table(config.DbTable).Where(config.SqlWhere).Find(&records).Error; err != nil {
		return "", err
	}

	var htmlTable strings.Builder
	htmlTable.WriteString(`<h2>` + config.PageTitle + `</h2>` + "\n")
	htmlTable.WriteString(breadcrumbBuilder(config))
	htmlTable.WriteString(RenderAddForm(config, modelName))

	htmlTable.WriteString("<table class='table mt-3'>\n<thead><tr>")

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

			if countConfig, countExists := config.CountRelatedData[field]; countExists {
				foreignKeyValue, ok := record[countConfig.ForeignKey]
				var count int64
				if ok {
					if err := db.Debug().Table(countConfig.Table).
						Where(fmt.Sprintf("%s = ?", countConfig.ForeignKey), foreignKeyValue).
						Count(&count).Error; err == nil {
					}
				}
				value = count
			}

			if linkConfig, linkExists := config.Links[field]; linkExists {
				link := linkConfig.Template
				for key, val := range record {
					placeholder := fmt.Sprintf("$%s$", key)
					link = strings.ReplaceAll(link, placeholder, fmt.Sprintf("%v", val))
				}
				value = fmt.Sprintf("<a href='%s'>%v</a>", link, value)
			}

			htmlTable.WriteString(fmt.Sprintf("<td%s>%v</td>", classAttr, value))
		}
		htmlTable.WriteString("</tr>\n")
	}
	htmlTable.WriteString("</tbody>\n</table>")

	return htmlTable.String(), nil
}

func HandleRenderTableAddRecord(c *gin.Context, db *gorm.DB) {
	currentAuthUser := GetCurrentAuthUser(c)
	if currentAuthUser.Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	modelName, ok := payload["modelName"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Model name is required"})
		return
	}

	config, err := loadModelConfig(c, modelName)
	if err != nil {
		log.Printf("No configuration found for model: %s", modelName)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No configuration found for model '" + modelName + "'"})
		return
	}

	insertData := make(map[string]interface{})
	for _, field := range config.AddableFields {
		if value, exists := payload[field]; exists {
			insertData[utils.CamelToSnake(field)] = value
		}
	}

	//	fmt.Println(insertData)

	// check RequiredFields
	for _, requiredField := range config.RequiredFields {
		if value, exists := payload[requiredField]; !exists || value == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Field '%s' is required", requiredField)})
			return
		}
	}

	// check NoZeroValueFields
	for _, noZeroField := range config.NoZeroValueFields {
		if value, exists := payload[noZeroField]; exists {
			if number, ok := value.(float64); ok && number == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Field '%s' cannot be zero", noZeroField)})
				return
			}
		}
	}

	if len(insertData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data to insert"})
		return
	}

	if err := db.Debug().Table(config.DbTable).Create(insertData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func RenderAddForm(config *modelConfig, modelName string) string {
	if config == nil || len(config.AddableFields) == 0 {
		return ""
	}

	var formBuilder strings.Builder
	formBuilder.WriteString(`<script src="/static/js/render_table_add.js"></script>
<link rel="stylesheet" href="/static/css/render_table_add.css">
	<div class="accordion" id="addFormAccordion">
        <div class="accordion-item">
            <h2 class="accordion-header" id="addFormHeading">
                <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#addFormCollapse" aria-expanded="false" aria-controls="addFormCollapse">
                    <i class="bi-plus-square"></i> &nbsp; Add New Record
                </button>
            </h2>
            <div id="addFormCollapse" class="accordion-collapse collapse" aria-labelledby="addFormHeading" data-bs-parent="#addFormAccordion">
                <div class="accordion-body" style="background: #e5eeff;">
`)

	formBuilder.WriteString(fmt.Sprintf(`<form id="addForm">
        <input type="hidden" name="modelName" value="%s">
`, modelName))

	for _, field := range config.AddableFields {
		label := config.Headers[field]
		if label == "" {
			label = field
		}

		requiredAttr := ""
		requiredLabel := ""
		for _, requiredField := range config.RequiredFields {
			if requiredField == field {
				requiredAttr = "required"
				requiredLabel = ` <span class="required-label">(required)</span>`
				break
			}
		}

		formBuilder.WriteString(fmt.Sprintf(`<div class="mb-3">
        <label for="%s" class="form-label">%s%s</label>
        <input type="text" class="form-control" id="%s" name="%s" %s>
    </div>`, field, label, requiredLabel, field, field, requiredAttr))

	}

	formBuilder.WriteString(`<button type="submit" class="btn btn-primary">Add</button></form>`)
	formBuilder.WriteString(`</div></div></div>`)
	return formBuilder.String()
}
