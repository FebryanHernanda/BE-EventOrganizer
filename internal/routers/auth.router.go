package routers

import (
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/handlers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.GET("/activate", authHandler.Activate)
}
