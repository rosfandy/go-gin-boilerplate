package routes

import "github.com/gin-gonic/gin"

func RegisterV1Routes(group *gin.RouterGroup) {
	RegisterHelloRoutes(group)
}
