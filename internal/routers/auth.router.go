package routers

import (
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/handlers"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, authLimiter *middleware.RateLimiterStore) {
	auth := router.Group("/auth")
	auth.Use(authLimiter.RateLimiterMiddleware())

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.GET("/activate", authHandler.Activate)
	auth.POST("/resend-activation", authHandler.ResendActivation)

}
