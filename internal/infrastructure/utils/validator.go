package utils

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidError struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

func BindAndValidate(c *gin.Context, obj interface{}) []ValidError {
	if err := c.ShouldBind(obj); err != nil {
		return formatValidationError(err)
	}
	return nil
}

func formatValidationError(err error) []ValidError {
	var errors []ValidError
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			msg := buildErrorMessage(e)
			errors = append(errors, ValidError{
				Key:     e.Field(),
				Message: msg,
			})
		}
	} else {
		errors = append(errors, ValidError{
			Key:     "request",
			Message: err.Error(),
		})
	}
	return errors
}

func buildErrorMessage(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()
	param := e.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, strings.ReplaceAll(param, " ", ", "))
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, param)
	default:
		return e.Error()
	}
}
