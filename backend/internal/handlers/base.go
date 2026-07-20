package handlers

import (
	"errors"
	"net/http"

	"github.com/ekosachev/logos/internal"
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data;omitempty"`
	Error   string `json:"error;omitempty"`
}

func sendError(c *gin.Context, status int, err error) {
	c.JSON(status, APIResponse{Success: false, Error: err.Error()})
}

func MapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, internal.ErrUsernameExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
