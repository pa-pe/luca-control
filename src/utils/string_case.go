package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strings"
	"unicode"
)

// CamelToSnake convert from camelCase to snake_case
func CamelToSnake(input string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(input, "${1}_${2}")
	return strings.ToLower(snake)
}

//func CamelToSnake(input string) string {
//	var result []rune
//	for i, r := range input {
//		if unicode.IsUpper(r) {
//			if i > 0 {
//				result = append(result, '_')
//			}
//			result = append(result, unicode.ToLower(r))
//		} else {
//			result = append(result, r)
//		}
//	}
//	return string(result)
//}

// SnakeToCamel convert from snake_case to camelCase
func SnakeToCamel(input string) string {
	parts := strings.Split(input, "_")
	for i, part := range parts {
		parts[i] = cases.Title(language.Und).String(part)
	}
	return strings.Join(parts, "")
}

//func SnakeToCamel(input string) string {
//	parts := strings.Split(input, "_")
//	for i := range parts {
//		if i > 0 {
//			parts[i] = strings.Title(parts[i])
//		}
//	}
//	return strings.Join(parts, "")
//}

// InvertCaseStyle Функция для определения стиля строки и выполнения обратного преобразования
func InvertCaseStyle(input string) string {
	if IsCamelCase(input) {
		return CamelToSnake(input)
	}
	return SnakeToCamel(input)
}

// IsCamelCase Проверяет, является ли строка в стиле camelCase
func IsCamelCase(input string) bool {
	for _, r := range input {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}
