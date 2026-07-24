package services

import (
	"net/http"

	"github.com/ekosachev/logos/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type MiddlewareService struct {
	refreshTokenRepository RefreshTokenStorer
}

func NewMiddlewareService(refreshTokenRepository RefreshTokenStorer) *MiddlewareService {
	return &MiddlewareService{refreshTokenRepository: refreshTokenRepository}
}

func verifyRefreshToken(tokenString string) (bool, *uuid.UUID) {
	cfg := config.GetConfig()
	secretKey := cfg.JWTRefreshSecret

	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) { return []byte(secretKey), nil })

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, err := uuid.Parse(claims["sub"].(string)); err == nil {
			return true, &userId
		}
	}

	return false, nil
}

func verifyAccessToken(tokenString string) (bool, *uuid.UUID) {
	cfg := config.GetConfig()
	secretKey := cfg.JWTAccessSecret

	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) { return []byte(secretKey), nil })
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, err := uuid.Parse(claims["sub"].(string)); err == nil {
			return true, &userId
		}
	}

	return false, nil
}

func (s *MiddlewareService) RequiresRefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("refreshToken")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "refreshToken cookie required"})
			return
		}

		if ok, userID := verifyRefreshToken(tokenString); ok && userID != nil {
			validToken, err := s.refreshTokenRepository.GetActiveRefreshToken(ctx, *userID)
			if err != nil || validToken == nil {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "error": "Unauthorized"})
				return
			}

			if validToken.Token != tokenString {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "error": "Unauthorized"})
				return
			}

			ctx.Set("userID", userID.String())
			ctx.Set("tokenID", validToken.ID.String())
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Invalid token"})
		}
	}
}

func (s *MiddlewareService) RequiresAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Authorization header required"})
			return
		}

		tokenString := authHeader[7:] // remove "Bearer: " prefix

		if ok, userID := verifyAccessToken(tokenString); ok && userID != nil {
			ctx.Set("userID", userID.String())
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Invalid token"})
		}
	}
}
