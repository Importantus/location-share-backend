package routers

import (
	"location-share-backend/customerrors"
	"location-share-backend/logic/locations"
	"location-share-backend/middleware"
	"location-share-backend/models"

	"github.com/gin-gonic/gin"
)

func SharedLocations(router *gin.RouterGroup) {
	router.Use(middleware.ReadAuthRequired()).GET("", func(ctx *gin.Context) {
		locations, getError := locations.GetSharedLocations(ctx.MustGet(middleware.SESSION_KEY).(models.Session).UserID)
		if getError != customerrors.Success {
			ctx.JSON(getError.Status, getError)
			return
		}
		ctx.JSON(200, locations)
	})
}
