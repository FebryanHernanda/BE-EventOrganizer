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

	repo := repository.NewHealthRepository(db)
	service := services.NewHealthService(repo)
	handler := handlers.NewHealthHandler(service)

	router := routers.SetupRouter(handler)

	fmt.Println("Server running on port:", port)
	router.Run(":" + port)

}
