package log

import (
	"log"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog() (*zap.Logger, func()) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	cleanUp := func() {
		_ = logger.Sync()
	}

	return logger, cleanUp
}

func SetupRequestTracking(ctx *gin.Context) {
	start := time.Now()
	requestId := uuid.New().String()
	ctx.Header("X-Request-ID", requestId)
	ctx.Set("requestId", requestId)

	ctx.Next()

	logger.Info("http_request",
		zap.String("request_id", requestId),
		zap.Int("status", ctx.Writer.Status()),
		zap.String("method", ctx.Request.Method),
		zap.String("path", ctx.Request.URL.Path),
		zap.String("ip", ctx.ClientIP()),
		zap.Duration("latency", time.Since(start)),
		zap.String("user_agent", ctx.Request.UserAgent()),
	)
}