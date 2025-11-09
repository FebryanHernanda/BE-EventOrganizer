package services

import (
	"context"

	"github.com/FebryanHernanda/BE-EventOrganizer/internal/models"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/repository"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	UserRepo  *repository.UserRepository
	JWTSecret string
}

func NewUserService(userRepo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		UserRepo:  userRepo,
		JWTSecret: jwtSecret,
	}
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	logrus.WithField("userID", userID).Info("Attempting get profile")

	user, err := s.UserRepo.GetByID(ctx, userID)
	if err != nil {
		logrus.WithField("userID", userID).WithError(err).Error("GetProfile error")
		return nil, err
	}

	logrus.WithField("userID", userID).Info("GetProfile success")

	return user, nil
}
