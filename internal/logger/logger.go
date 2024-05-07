package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gin-gonic/gin"
)

type Logger struct {
	*zap.Logger
}

func New() *Logger {
	// Create a new file
	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		panic(err)
	}

	// Create a new logger
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(file),
		zap.InfoLevel,
	)
	logger := zap.New(core)

	gin.DefaultWriter = file

	return &Logger{logger}
}

func (l *Logger) GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		l.Info("Processed request",
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
			zap.String("client", c.ClientIP()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}
