package handler

import (
	"net/http"
	"swap/api"
	"swap/apperrors"

	"github.com/gin-gonic/gin"
)

func ToFieldErrorResponse(c *gin.Context, field, message string) {

	c.JSON(http.StatusBadRequest, api.NewResponse(http.StatusBadRequest, "Bad Request", gin.H{
		"errors": []apperrors.FieldError{
			{
				Field:   field,
				Message: message,
			},
		},
	}))
}