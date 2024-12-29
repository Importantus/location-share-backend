package routers

import (
	"location-share-backend/initializers"

	"github.com/gin-gonic/gin"
)

func Info(router *gin.RouterGroup) {
	router.GET("", func(ctx *gin.Context) {
		config, err := initializers.LoadConfig(".")

		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{
			"public_registration": config.RegistrationSecret == "",
		})
	})
}
