package main

import (
	"log"

	"github.com/ekosachev/logos/internal/database"
	"github.com/ekosachev/logos/internal/handlers"
	"github.com/ekosachev/logos/internal/repositories"
	"github.com/ekosachev/logos/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.ConnectToDb()
	if err != nil {
		log.Fatal(err)
	}

	userRepository := repositories.NewUserRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(db)

	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(refreshTokenRepository, userRepository)
	middlewareService := services.NewMiddlewareService(refreshTokenRepository)

	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService, middlewareService)

	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "pong"}) })
	apiGroup := router.Group("/api/v1")
	{
		userHandler.RegisterRoutes(apiGroup)
		authHandler.RegisterRoutes(apiGroup)
	}

	if err = router.Run(); err != nil {
		log.Fatal(err)
	}
}
