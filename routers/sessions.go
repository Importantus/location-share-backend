package routers

import (
	"location-share-backend/logic/sessions"
	"location-share-backend/middleware"
	"location-share-backend/models"

	"github.com/gin-gonic/gin"
)

func Sessions(router *gin.RouterGroup) {
	router.POST("", func(ctx *gin.Context) {
		var json models.SessionCreate

		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		key, err := sessions.CreateSession(&json)

		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"token": key})
	})

	router.Use(middleware.ReadAuthRequired()).GET("/", func(ctx *gin.Context) {
		sessions, err := sessions.ListSessions(ctx.MustGet(middleware.SESSION_KEY).(models.Session))

		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, sessions)
	})

	router.Use(middleware.WriteAuthRequired()).DELETE("", func(ctx *gin.Context) {
		err := sessions.DestroySession(ctx.MustGet(middleware.SESSION_KEY).(models.Session))

		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{})
	})
}
