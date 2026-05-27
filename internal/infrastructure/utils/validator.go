package utils

import (
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
            errors = append(errors, ValidError{
                Key:     e.Field(),
                Message: e.Error(),
            })
        }
    } else {
        errors = append(errors, ValidError{
            Key:     "request",
            Message: err.Error(),
        })
    }
    return nil
}