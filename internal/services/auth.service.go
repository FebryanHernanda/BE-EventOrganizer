package services

import (
	"context"
	"errors"
	"time"

	"github.com/FebryanHernanda/BE-EventOrganizer/internal/models"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo  *repository.UserRepository
	JWTSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		JWTSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) error {
	logrus.WithField("email", req.Email).Info("Attempting user registration")

	email, err := s.UserRepo.FindByEmail(ctx, req.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		logrus.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err,
		}).Error("Failed to check email")
		return err
	}

	if email != nil {
		return errors.New("email already registered")
	}

	if _, err := s.UserRepo.FindByUsername(ctx, req.Username); err == nil {
		logrus.WithField("username", req.Username).Warn("Username already registered")
		return errors.New("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithField("error", err).Error("Failed to hash password")
		return err
	}

	user := models.User{
		ID:        primitive.NewObjectID(),
		FullName:  req.FullName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      "user",
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.UserRepo.CreateUser(ctx, &user); err != nil {
		logrus.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err,
		}).Error("Failed to create user")
		return err
	}

	logrus.WithField("email", req.Email).Info("User registered successfully")
	return nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (string, error) {
	logrus.WithField("email", req.Email).Info("Attempting user login")

	user, err := s.UserRepo.FindByEmail(ctx, req.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		logrus.WithField("email", req.Email).Warn("Email not registered")
		return "", errors.New("email not registered")
	}

	if user == nil {
		logrus.WithField("email", req.Email).Warn("User not found (nil result)")
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logrus.WithField("email", req.Email).Warn("Invalid password")
		return "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		logrus.WithError(err).Error("Failed to sign token")
		return "", err
	}

	logrus.WithField("email", req.Email).Info("Login successfully")
	return tokenString, nil
}
