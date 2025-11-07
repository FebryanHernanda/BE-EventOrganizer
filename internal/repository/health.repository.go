package repository

import (
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthRepository interface {
	CheckHealth() models.Health
}

type healthRepository struct {
	db *mongo.Database
}

func NewHealthRepository(db *mongo.Database) HealthRepository {
	return &healthRepository{db: db}
}

func (r *healthRepository) CheckHealth() models.Health {
	return models.Health{
		Status:  "OK",
		Message: "Service is running",
	}
}
