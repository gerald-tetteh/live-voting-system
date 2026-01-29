package main

import (
	"os"

	"geraldaddo.com/live-voting-system/log"
	"github.com/gin-gonic/gin"
)

func main() {
	logger, cleanup := log.InitLog()
	defer cleanup()

	server := gin.Default()

	server.Use(log.SetupRequestTracking)

	logger.Info("Starting server")
	server.Run(os.Getenv("PORT"))
}