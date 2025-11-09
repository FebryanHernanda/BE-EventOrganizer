package routers

import (
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/handlers"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/middleware"
	"github.com/gin-gonic/gin"
)

type RouteHandler struct {
	Health *handlers.HealthHandler
	Auth   *handlers.AuthHandler
	User   *handlers.UserHandler
}

func SetupRouter(rh *RouteHandler, jwtSecret string) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())

	router.GET("/Health", rh.Health.HealthCheck)

	AuthRoutes(router, rh.Auth)
	UserRoutes(router, rh.User, jwtSecret)

	return router
}
