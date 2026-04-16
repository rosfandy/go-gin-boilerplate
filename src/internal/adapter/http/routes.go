package http

import (
	v1 "owner-api-proxy/internal/adapter/http/v1/routes"
	"owner-api-proxy/internal/container"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine, appContainer *container.Container) {
	apiV1 := engine.Group("/api/v1")
	v1.RegisterV1Routes(apiV1, appContainer)
}
