package handlers

import (
	"context"
	"net/http"

	"github.com/ekosachev/logos/internal"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserServicer interface {
	CreateUser(ctx context.Context, user *internal.User) error
}

type UserHandler struct {
	service UserServicer
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func NewUserHandler(service UserServicer) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(group *gin.RouterGroup) {
	userGroup := group.Group("/user")
	{
		userGroup.POST("/", h.createUser)
	}
}

func (h *UserHandler) createUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required;alphanum;lte=64;gte=4;"`
		Password string `json:"password" binding:"required;alphanum;lte=64;gte=8;"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, err)
		return
	}

	user := internal.User{
		Username:     req.Username,
		PasswordHash: req.Password,
	}

	if err := h.service.CreateUser(c, &user); err != nil {
		statusCode := MapErrorToStatus(err)
		sendError(c, statusCode, err)
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data: UserResponse{
			Username: user.Username,
			ID:       user.ID,
		},
	})
}
