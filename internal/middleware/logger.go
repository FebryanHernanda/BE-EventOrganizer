package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path

		entry := logrus.WithFields(logrus.Fields{
			"status":   status,
			"latency":  latency,
			"clientIP": clientIP,
			"method":   method,
			"path":     path,
		})

		switch {
		case status >= 500:
			entry.Error("Server error")
		case status >= 400:
			entry.Warn("Client error")
		default:
			entry.Info("Request completed")
		}
	}
}
