package handlers

import (
	"context"

	"github.com/ekosachev/logos/internal"
	"github.com/gin-gonic/gin"
)

type UserServicer interface {
	CreateUser(ctx context.Context, user *internal.User) error
}

type UserHandler struct {
	service UserServicer
}

func NewUserHandler(service UserServicer) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(group *gin.RouterGroup) {
	_ = group.Group("/user")
}
