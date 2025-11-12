package routers

import (
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/handlers"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/middleware"
	"github.com/gin-gonic/gin"

	_ "github.com/FebryanHernanda/BE-EventOrganizer/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouteHandler struct {
	Health *handlers.HealthHandler
	Auth   *handlers.AuthHandler
	User   *handlers.UserHandler
}

func SetupRouter(rh *RouteHandler, jwtSecret string, deps *middleware.MiddlewareDeps) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())

	router.Use(deps.GlobalLimiter.RateLimiterMiddleware())

	// Routes
	router.GET("/Health", rh.Health.HealthCheck)
	AuthRoutes(router, rh.Auth, deps.AuthLimiter)
	UserRoutes(router, rh.User, jwtSecret)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
