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

// Me godoc
// @Summary     Get logged-in user profile
// @Description Return profile of the current logged-in user (no password).
// @Tags        User
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Success     200 {object} models.User "user profile (password hidden by json:\"-\")"
// @Failure     401 {object} map[string]interface{} "Unauthorized"
// @Failure     500 {object} map[string]interface{} "Server error"
// @Router      /user/me [get]
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
