package services

import (
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/models"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/repository"
)

type HealthService interface {
	GetHealth() models.Health
}

type healthService struct {
	repo repository.HealthRepository
}

func NewHealthService(repo repository.HealthRepository) HealthService {
	return &healthService{repo}
}

func (s *healthService) GetHealth() models.Health {
	return s.repo.CheckHealth()
}
