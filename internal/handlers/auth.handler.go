package handlers

import (
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/models"
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/services"
	"github.com/FebryanHernanda/BE-EventOrganizer/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req models.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.WithField("error", err).Warn("Invalid JSON format in registration")
		response.Error(ctx, "Invalid JSON format", 400, nil)
		return
	}

	if err := h.AuthService.Register(ctx.Request.Context(), &req); err != nil {
		logrus.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err,
		}).Warn("Registration failed")
		response.Error(ctx, err.Error(), 400, nil)
		return
	}

	logrus.WithField("email", req.Email).Info("User registration validated successfully")

	response.Success(ctx, req, "User registration validated successfully")
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req models.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.WithField("error", err).Warn("Invalid JSON format in login")
		response.Error(ctx, "invalid request body", 400, nil)
		return
	}

	token, err := h.AuthService.Login(ctx.Request.Context(), &req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err,
		}).Warn("Login failed")
		response.Error(ctx, err.Error(), 401, nil)
		return
	}

	response.Success(ctx, gin.H{
		"token": token,
	}, "login successful")
}
