package bootstrap

import (
	httpadapter "owner-api-proxy/internal/adapter/http"
	"owner-api-proxy/internal/container"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BuildHTTPEngine(engine *gin.Engine, db *gorm.DB) {
	appContainer := container.NewContainer(db)
	httpadapter.RegisterRoutes(engine, appContainer)
}
