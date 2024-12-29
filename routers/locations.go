package routers

import (
	"location-share-backend/customerrors"
	"location-share-backend/logic/locations"
	"location-share-backend/middleware"
	"location-share-backend/models"
	"location-share-backend/utils"

	"github.com/gin-gonic/gin"
)

func Locations(router *gin.RouterGroup) {
	router.Use(middleware.WriteAuthRequired()).POST("", func(ctx *gin.Context) {
		var json []models.LocationCreate

		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		session := ctx.MustGet(middleware.SESSION_KEY).(models.Session)

		var createdSessions []models.Location

		for _, location := range json {

			if location.UserID != session.UserID {
				ctx.JSON(400, customerrors.ErrInvalidUserID)
				return
			}

			createdLocation, err := locations.CreateLocation(session.ID, location)

			createdSessions = append(createdSessions, createdLocation)

			if err != customerrors.Success {
				ctx.JSON(err.Status, err)
				return
			}
		}

		ctx.JSON(200, createdSessions)
	})

	router.Use(middleware.ReadAuthRequired()).GET("", func(ctx *gin.Context) {
		from, err := utils.ParseTime(ctx.Query("from"))

		if err != nil {
			ctx.JSON(400, customerrors.ErrInvalidTime)
			return
		}

		to, err := utils.ParseTime(ctx.Query("to"))

		if err != nil {
			ctx.JSON(400, customerrors.ErrInvalidTime)
			return
		}

		locations, getError := locations.GetLocations(ctx.MustGet(middleware.SESSION_KEY).(models.Session).UserID, from, to)

		if getError != customerrors.Success {
			ctx.JSON(getError.Status, getError)
			return
		}

		ctx.JSON(200, locations)
	})

	router.Use(middleware.WriteAuthRequired()).DELETE("", func(ctx *gin.Context) {
		from, err := utils.ParseTime(ctx.Query("from"))

		if err != nil {
			ctx.JSON(400, customerrors.ErrInvalidTime)
			return
		}

		to, err := utils.ParseTime(ctx.Query("to"))

		if err != nil {
			ctx.JSON(400, customerrors.ErrInvalidTime)
			return
		}

		deletionError := locations.DeleteLocations(ctx.MustGet(middleware.SESSION_KEY).(models.Session).UserID, from, to)

		if deletionError != customerrors.Success {
			ctx.JSON(deletionError.Status, deletionError)
			return
		}

		ctx.JSON(200, gin.H{"message": "Location deleted"})
	})
}
