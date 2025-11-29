package handlers

import (
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/models/user"
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

// Register godoc
// @Summary     Register a new user
// @Description Register new user. Returns only status and message (no data).
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       body  body     models.RegisterRequest true "Register payload"
// @Success     201   {object} map[string]interface{} "Created (no data)"
// @Failure     400   {object} map[string]interface{} "Bad request"
// @Router      /auth/register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req user.RegisterRequest

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

	logrus.WithField("email", req.Email).Info("User registration successfully")

	response.Success(ctx, nil, "User registration successfully")
}

/* Activation Account */

func (h *AuthHandler) Activate(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		response.Error(ctx, "Missing token", 400, nil)
		return
	}

	if err := h.AuthService.ActivateAccount(ctx, token); err != nil {
		response.Error(ctx, err.Error(), 400, nil)
		return
	}

	response.Success(ctx, nil, "Account activated successfully, you can now login.")
}

func (h *AuthHandler) ResendActivation(ctx *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, "Invalid request body", 400, nil)
		return
	}

	if err := h.AuthService.ResendActivation(ctx, req.Email); err != nil {
		response.Error(ctx, err.Error(), 400, nil)
		return
	}

	response.Success(ctx, nil, "Activation link has been sent to your email.")
}

/* Activation Account */

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req user.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.WithField("error", err).Warn("Invalid JSON format in login")
		response.Error(ctx, "invalid request body", 400, nil)
		return
	}

	token, err := h.AuthService.Login(ctx.Request.Context(), &req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"email": req.Identifier,
			"error": err,
		}).Warn("Login failed")
		response.Error(ctx, err.Error(), 401, nil)
		return
	}

	logrus.WithField("Identifier", req.Identifier).Info("User login successfully")

	response.Success(ctx, gin.H{
		"token": token,
	}, "login successful")
}
