package main

import (
	"net/http"

	"geraldaddo.com/live-voting-system/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, cleanup := log.InitLog()
	defer cleanup()

	gin.SetMode(gin.ReleaseMode)

	server := gin.New()

	server.Use(gin.Recovery())
	server.Use(func(ctx *gin.Context) {
		log.SetupRequestTracking(ctx, logger)
	})

	server.GET("/", func(ctx *gin.Context) {
		requestId := ctx.GetString("requestId")
		logger.Info("testing root path", zap.String("request_id", requestId))
		ctx.JSON(http.StatusOK, gin.H{"message": "server is working well"})
	})

	logger.Info("Starting server")
	server.Run(":8080")
}