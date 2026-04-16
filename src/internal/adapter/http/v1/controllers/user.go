package controllers

import (
	"net/http"
	"owner-api-proxy/internal/adapter/http/v1/presenters"
	"owner-api-proxy/internal/container"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context, dependcies *container.UserContainer) {
	userService := dependcies.UserService

	users, err := userService.GetUsers(c.Request.Context())
	if err != nil {
		errorMessage := "failed to get users: " + err.Error()
		presenters.SendErrorResponse(c, http.StatusInternalServerError, errorMessage)
		return
	}

	data := any(users)
	message := "users fetched"
	response := presenters.SuccessResponse{
		Success: true,
		Message: &message,
		Data:    &data,
	}

	presenters.SendSuccessResponse(c, response)
}
