package main

import (
	"log"

	"location-share-backend/initializers"
	"location-share-backend/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load the configuration file")
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.Share{})
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Location{})
	initializers.DB.AutoMigrate(&models.Session{})

}
