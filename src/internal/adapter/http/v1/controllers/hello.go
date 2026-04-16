package controllers

import (
	"net/http"
	"owner-api-proxy/internal/adapter/http/v1/presenters"

	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	msg := "Hello World"
	response := presenters.SuccessResponse{
		Success: true,
		Message: &msg,
	}

	c.JSON(http.StatusOK, response)
}
