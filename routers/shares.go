package routers

import (
	"location-share-backend/customerrors"
	"location-share-backend/logic/shares"
	"location-share-backend/middleware"
	"location-share-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Shares(router *gin.RouterGroup) {
	router.Use(middleware.WriteAuthRequired()).POST("", func(ctx *gin.Context) {
		var json models.ShareCreate

		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		share, err := shares.CreateShare(ctx.MustGet(middleware.SESSION_KEY).(models.Session).UserID, json.SharedWith, json.ValidUntil)

		if err != customerrors.Success {
			ctx.JSON(err.Status, err)
			return
		}

		ctx.JSON(200, share)
	})

	router.Use(middleware.ReadAuthRequired()).GET("", func(ctx *gin.Context) {
		shares, err := shares.GetShares(ctx.MustGet(middleware.SESSION_KEY).(models.Session).UserID)

		if err != customerrors.Success {
			ctx.JSON(err.Status, err)
			return
		}

		ctx.JSON(200, shares)
	})

	router.Use(middleware.WriteAuthRequired()).DELETE("", func(ctx *gin.Context) {
		shareId, err := uuid.Parse(ctx.Query("id"))

		if err != nil {
			ctx.JSON(400, customerrors.ErrInvalidUUID)
			return
		}

		deletionError := shares.DeleteShare(shareId, ctx.MustGet(middleware.SESSION_KEY).(models.Session).UserID)

		if deletionError != customerrors.Success {
			ctx.JSON(deletionError.Status, deletionError)
			return
		}

		ctx.JSON(200, gin.H{"message": "Share deleted"})
	})
}
