package main

import (
	"net/http"
	"os"
	"strconv"

	"geraldaddo.com/live-voting-system/db"
	"geraldaddo.com/live-voting-system/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, cleanup := log.InitLog()
	defer cleanup()
	
	var maxOpenConnections int64
	var maxIdleConnections int64
	var err error
	maxOpenConnections, err = strconv.ParseInt(os.Getenv("MAX_OPEN_CONN"), 10, 64)
	if err != nil {
		logger.Error("Could not parse max open connections")
		logger.Fatal(err.Error())
	}
	maxIdleConnections, err = strconv.ParseInt(os.Getenv("MAX_IDLE_CONN"), 10, 64)
	if err != nil {
		logger.Error("Could not parse max idle connections")
		logger.Fatal(err.Error())
	}

	db.InitDB(logger, os.Getenv("DB_URL"), int(maxOpenConnections), int(maxIdleConnections))

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