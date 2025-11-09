package handlers

import (
	"github.com/FebryanHernanda/BE-EventOrganizer/internal/services"
	"github.com/FebryanHernanda/BE-EventOrganizer/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) Me(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		logrus.Warn("Me unauthorized: missing userID in context")
		response.Error(ctx, "Unauthorized", 401, nil)
		return
	}

	userID := userIDVal.(string)

	user, err := h.UserService.GetProfile(ctx.Request.Context(), userID)
	if err != nil {
		logrus.WithField("userID", userID).Errorf("Me handler error: %v", err)
		response.Error(ctx, "failed to fetch profile", 500, nil)
		return
	}
	logrus.WithField("email", user.Email).Info("User get profile successfully")

	response.Success(ctx, user, "success get profile")
}
