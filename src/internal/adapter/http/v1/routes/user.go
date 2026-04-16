package routes

import (
	"owner-api-proxy/internal/adapter/http/v1/controllers"
	"owner-api-proxy/internal/container"

	"github.com/gin-gonic/gin"
)

func UserRoutes(app *gin.RouterGroup, dependencies *container.Container) {
	api := app.Group("/users")

	api.GET("/", func(c *gin.Context) {
		controllers.GetUsers(c, dependencies.User)
	})
}
