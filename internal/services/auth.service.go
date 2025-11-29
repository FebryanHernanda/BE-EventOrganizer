package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
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
	AuthRepo  *repository.AuthRepository
	Mailer    *mailerService.MailerService
	JWTSecret string
}

func NewAuthService(userRepo *repository.UserRepository, authRepo *repository.AuthRepository, mailer *mailerService.MailerService, jwtSecret string) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		AuthRepo:  authRepo,
		Mailer:    mailer,
		JWTSecret: jwtSecret,
	}
}

func generateActivationToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

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

	_ = s.AuthRepo.DeleteActivationToken(ctx, user.ID.Hex())
	token, err := generateActivationToken()
	if err != nil {
		return err
	}

	activation := &modelAuth.ActivationToken{
		UserID:    user.ID.Hex(),
		Email:     user.Email,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
		CreatedAt: time.Now(),
	}
	if err := s.AuthRepo.SaveActivationToken(ctx, activation); err != nil {
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

func (s *AuthService) ActivateAccount(ctx context.Context, token string) error {
	at, err := s.AuthRepo.FindValidActivationToken(ctx, token)
	if err != nil {
		logrus.WithError(err).Error("Failed activation token")
		return errors.New("invalid activation token")
	}
	if at == nil {
		logrus.Warn("Failed activation, link expired or invalid")
		return errors.New("activation link expired or invalid")
	}

	if err := s.AuthRepo.ActivateUser(ctx, at.UserID); err != nil {
		return err
	}

	_ = s.AuthRepo.MarkTokenUsed(ctx, token)

	logrus.WithField("userID", at.UserID).Info("Account activated successfully")
	return nil
}

func (s *AuthService) ResendActivation(ctx context.Context, email string) error {
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		logrus.WithError(err).Error("Failed fetching user by email")
		return errors.New("email not registered")
	}

	if user.IsActive {
		logrus.Warn("Account already activated, skip resend")
		return errors.New("account already activate")
	}

	if err = s.AuthRepo.DeleteActivationToken(ctx, user.ID.Hex()); err != nil {
		logrus.WithError(err).Warn("failed deleting old activation token")
	}

	token, err := generateActivationToken()
	if err != nil {
		logrus.WithError(err).Error(("Failed generating activation token"))
		return errors.New("failed generating new token")
	}

	activation := &modelAuth.ActivationToken{
		UserID:    user.ID.Hex(),
		Email:     user.Email,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
		CreatedAt: time.Now(),
	}
	if err := s.AuthRepo.SaveActivationToken(ctx, activation); err != nil {
		logrus.WithError(err).Error("Failed saving activation token")
		return errors.New("failed saving activation token")
	}

	if err := s.Mailer.SendActivationEmail(
		user.FullName,
		user.Email,
		user.CreatedAt.Format("2006-01-02"),
		token,
	); err != nil {
		logrus.WithError(err).Error("Failed sending activation email")
		return errors.New("failed sending activation email")
	}

	logrus.Info("Resend activation email sent successfully")
	return nil
}

func (s *AuthService) Login(ctx context.Context, req *modelUser.LoginRequest) (string, error) {
	logrus.WithField("identifier", req.Identifier).Info("Attempting user login")

	var user *modelUser.User
	var err error

	if strings.Contains(req.Identifier, "@") {
		user, err = s.UserRepo.FindByEmail(ctx, req.Identifier)
	} else {

		user, err = s.UserRepo.FindByUsername(ctx, req.Identifier)
	}

	if err != nil && err != mongo.ErrNoDocuments {
		logrus.WithField("identifier", req.Identifier).Warn("Email/Username not registered")
		return "", errors.New("invalid credentials")
	}
	if user == nil {
		logrus.WithField("identifier", req.Identifier).Warn("User not found (nil result)")
		return "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		logrus.WithField("userID", user.ID.Hex()).Warn("User attempted login but account is not active")
		return "", errors.New("account not activated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logrus.WithField("id", user.ID.Hex()).Warn("Invalid password")
		return "", errors.New("invalid credentials")
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

	logrus.WithField("id", user.ID.Hex()).Info("Login successfully")
	return tokenString, nil
}
