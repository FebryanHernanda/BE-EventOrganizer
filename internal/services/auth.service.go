package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	mailerService "github.com/FebryanHernanda/BE-EventOrganizer/internal/mailer/services"
	modelAuth "github.com/FebryanHernanda/BE-EventOrganizer/internal/models/auth"
	modelUser "github.com/FebryanHernanda/BE-EventOrganizer/internal/models/user"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo  *repository.UserRepository
	Mailer    *mailerService.MailerService
	JWTSecret string
}

func NewAuthService(userRepo *repository.UserRepository, mailer *mailerService.MailerService, jwtSecret string) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		Mailer:    mailer,
		JWTSecret: jwtSecret,
	}
}

/* Activation Account */

func generateActivationToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) ActivateAccount(ctx context.Context, token string) error {
	at, err := s.UserRepo.FindActivationToken(ctx, token)
	if err != nil {
		return errors.New("invalid token")
	}
	if at == nil {
		return errors.New("token not found or already used")
	}

	if err := s.UserRepo.ActivateUser(ctx, at.UserID); err != nil {
		return err
	}

	if err := s.UserRepo.DeleteActivationToken(ctx, token); err != nil {
		logrus.WithError(err).Warn("Failed to delete used activation token")
	}

	logrus.WithField("userID", at.UserID).Info("Account activated successfully")
	return nil
}

/* Activation Account */

func (s *AuthService) Register(ctx context.Context, req *modelUser.RegisterRequest) error {
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

	user := modelUser.User{
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

	token, err := generateActivationToken()
	if err != nil {
		return err
	}

	activation := &modelAuth.ActivationToken{
		UserID:    user.ID.Hex(),
		Token:     token,
		CreatedAt: time.Now().Unix(),
	}
	if err := s.UserRepo.SaveActivationToken(ctx, activation); err != nil {
		return err
	}

	if err := s.Mailer.SendActivationEmail(
		user.FullName,
		user.Email,
		user.CreatedAt.Format("2006-01-02"),
		token,
	); err != nil {
		logrus.WithError(err).Error("Activation email failed")
	}

	logrus.WithField("email", req.Email).Info("User registered successfully")
	return nil
}

func (s *AuthService) Login(ctx context.Context, req *modelUser.LoginRequest) (string, error) {
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
