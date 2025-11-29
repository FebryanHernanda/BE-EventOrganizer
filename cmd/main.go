package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/FebryanHernanda/BE-EventOrganizer/config"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/handlers"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/mailer/provider"
	mailerService "github.com/FebryanHernanda/BE-EventOrganizer/internal/mailer/services"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/middleware"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/repository"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/routers"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/services"
	"golang.org/x/time/rate"
)

// @title BE Event Organizer API
// @version 1.0
// @description API for BE-EventOrganizer
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
func main() {
	config.LoadEnv()
	config.ConnectMongo()

	client, db, err := config.ConnectMongo()
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	smtpConfig := config.LoadSMTPConfig()
	jwtSecret := os.Getenv("JWT_SECRET")

	smtpProvider := provider.NewSMTPProvider(smtpConfig)

	healthRepo := repository.NewHealthRepository(db)
	healthService := services.NewHealthService(healthRepo)
	healthHandler := handlers.NewHealthHandler(healthService)

	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)

	mailService := mailerService.NewMailerService(smtpProvider)

	authService := services.NewAuthService(userRepo, authRepo, mailService, jwtSecret)
	userService := services.NewUserService(authRepo, userRepo, jwtSecret)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	rh := &routers.RouteHandler{
		Health: healthHandler,
		Auth:   authHandler,
		User:   userHandler,
	}

	deps := &middleware.MiddlewareDeps{
		GlobalLimiter: middleware.NewRateLimiterStore(rate.Every(10*time.Second), 3),
		AuthLimiter:   middleware.NewRateLimiterStore(rate.Every(time.Second/5), 10),
	}

	router := routers.SetupRouter(rh, jwtSecret, deps)

	fmt.Println("Server running on port:", port)
	router.Run(":" + port)
}
