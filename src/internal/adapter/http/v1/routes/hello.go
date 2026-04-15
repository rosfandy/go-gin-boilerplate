package routes

import (
	"owner-api-proxy/internal/adapter/http/v1/controller"

	"github.com/gin-gonic/gin"
)

func RegisterHelloRoutes(group *gin.RouterGroup) {
	group.GET("/hello", controller.Hello)
}
