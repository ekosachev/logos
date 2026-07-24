package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthServicer interface {
	Login(ctx context.Context, username string, password string) (accessToken, refreshToken string, err error)
	Refresh(ctx context.Context, userID uuid.UUID) (accessToken, newRefreshToken string, err error)
}

type AuthHandler struct {
	service AuthServicer
}

func NewAuthHandler(service AuthServicer) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) RegisterRoutes(group *gin.RouterGroup) {
	authGroup := group.Group("/auth")
	{
		authGroup.POST("/login", h.login)
	}
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required;alphanum;lte=64;gte=4;"`
		Password string `json:"password" binding:"required;alphanum;lte=64;gte=8;"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, err)
		return
	}

	accessToken, refreshToken, err := h.service.Login(c, req.Username, req.Password)
	if err != nil {
		statusCode := MapErrorToStatus(err)
		sendError(c, statusCode, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

func getUUIDFromCxt(c *gin.Context, name string) (uuid.UUID, error) {
	uuidString := c.GetString(name)
	return uuid.Parse(uuidString)
}

func (h *AuthHandler) refresh(c *gin.Context) {
	userID, err := getUUIDFromCxt(c, "userID")
	if err != nil {
		sendError(c, http.StatusInternalServerError, err)
		return
	}

	accessToken, refreshToken, err := h.service.Refresh(c, userID)
	if err != nil {
		errorCode := MapErrorToStatus(err)
		sendError(c, errorCode, err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}
