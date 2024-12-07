package main

import (
	"location-share-backend/initializers"
	"location-share-backend/routers"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	routers.Sessions(server.Group("/sessions"))

	log.Fatal(server.Run(":" + config.ServerPort))
}
