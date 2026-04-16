package routes

import (
	"owner-api-proxy/internal/adapter/http/v1/controllers"

	"github.com/gin-gonic/gin"
)

func HelloRoutes(app *gin.RouterGroup) {
	app.GET("/hello", controllers.Hello)
}
