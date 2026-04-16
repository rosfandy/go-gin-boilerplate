package routes

import (
	"owner-api-proxy/internal/container"

	"github.com/gin-gonic/gin"
)

func RegisterV1Routes(group *gin.RouterGroup, dependencies *container.Container) {
	HelloRoutes(group)
	UserRoutes(group, dependencies)
}
