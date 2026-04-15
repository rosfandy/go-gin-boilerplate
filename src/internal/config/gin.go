package config

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type GinConfig struct {
	Mode    string
	Address string
}

func LoadGinConfig() GinConfig {
	httpHost := ""
	httpPort := ":8080"

	if AppConfig != nil {
		httpHost = AppConfig.Server.Host

		if AppConfig.Server.Port != "" {
			httpPort = AppConfig.Server.Port
		}
	}

	address := httpPort
	if httpHost != "" {
		if strings.HasPrefix(httpPort, ":") {
			address = httpHost + httpPort
		} else {
			address = httpHost + ":" + httpPort
		}
	}

	return GinConfig{
		Mode:    gin.DebugMode,
		Address: address,
	}
}

func NewGinEngine(cfg GinConfig) *gin.Engine {
	gin.SetMode(cfg.Mode)

	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	return engine
}
