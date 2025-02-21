package main

import (
	"location-share-backend/initializers"
	"location-share-backend/logic/fcm"
	"location-share-backend/logic/ws"
	"location-share-backend/routers"
	"log"

	"github.com/gin-contrib/cors"
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

	hub := ws.GetHub()

	if err := fcm.InitializeFirebase(); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}

	server.Use(cors.New(corsConfig))

	v1 := server.Group("/v1")
	{
		routers.Sessions(v1.Group("/sessions"))
		routers.Users(v1.Group("/users"))
		routers.Shares(v1.Group("/shares"))
		routers.Locations(v1.Group("/locations"))
		routers.SharedLocations(v1.Group("/shared-locations"))
		routers.Info(v1.Group("/info"))
		routers.Websocket(v1.Group("/ws"), hub)
	}

	log.Fatal(server.Run(":" + config.ServerPort))
}
