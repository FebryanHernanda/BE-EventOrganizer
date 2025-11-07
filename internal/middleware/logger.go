package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	logger := logrus.New()

	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()

		logger.WithFields(logrus.Fields{
			"method":  ctx.Request.Method,
			"path":    ctx.Request.URL.Path,
			"status":  status,
			"latency": latency,
		}).Info("request completed")
	}
}
