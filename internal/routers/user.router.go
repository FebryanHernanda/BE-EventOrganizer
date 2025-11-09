package routers

import (
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/handlers"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userHandler *handlers.UserHandler, jwtSecret string) {
	user := router.Group("/user")
	user.Use(middleware.JWTAuthMiddleware(jwtSecret))
	user.GET("/me", userHandler.Me)
}
