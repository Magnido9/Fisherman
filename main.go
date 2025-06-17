package main

import (
	"Fisherman/cache"
	"Fisherman/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Init Redis
	if err := cache.InitRedis("localhost:6377"); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	router := gin.Default()
	router.GET("/readUrl", handlers.ReadHtmlApi)
	router.Run("localhost:8080")
}
