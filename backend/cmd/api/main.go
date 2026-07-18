package main

import (
	"log"

	"github.com/ekosachev/logos/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	_, err := database.ConnectToDb()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "pong"}) })

	if err = router.Run(); err != nil {
		log.Fatal(err)
	}
}
