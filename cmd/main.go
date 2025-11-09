package main

import (
	"context"
	"fmt"
	"os"

	"github.com/FebryanHernanda/BE-EventOrganizer/config"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/handlers"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/repository"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/routers"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/services"
)

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

	jwtSecret := os.Getenv("JWT_SECRET")

	healthRepo := repository.NewHealthRepository(db)
	healthService := services.NewHealthService(healthRepo)
	healthHandler := handlers.NewHealthHandler(healthService)

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)
	userService := services.NewUserService(userRepo, jwtSecret)
	userHandler := handlers.NewUserHandler(userService)

	rh := &routers.RouteHandler{
		Health: healthHandler,
		Auth:   authHandler,
		User:   userHandler,
	}

	router := routers.SetupRouter(rh, jwtSecret)

	fmt.Println("Server running on port:", port)
	router.Run(":" + port)

}
