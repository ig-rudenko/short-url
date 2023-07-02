package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidationError(errs validator.ValidationErrors) string {
	var messages []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			messages = append(messages, fmt.Sprintf("Обязательное поле %s", err.Field()))
		case "url":
			messages = append(messages, "Неверный URL")
		default:
			messages = append(messages, fmt.Sprintf("Поле %s неверное", err.Field()))

		}
	}
	return strings.Join(messages, ", ")
}
