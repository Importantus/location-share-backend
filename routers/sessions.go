package routers

import (
	"location-share-backend/customerrors"
	"location-share-backend/logic/sessions"
	"location-share-backend/logic/users"
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

		key, id, err := sessions.CreateSession(&json)

		if err != customerrors.Success {
			ctx.JSON(err.Status, err)
			return
		}

		session, error := sessions.GetSession(id)

		if error != customerrors.Success {
			ctx.JSON(error.Status, error)
			return
		}

		ctx.Set(middleware.SESSION_KEY, session)

		user, error := users.GetUser(session.UserID.String())

		if error != customerrors.Success {
			ctx.JSON(error.Status, error)
			return
		}

		ctx.JSON(200, gin.H{"token": key, "id": id, "user": models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}})
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

	router.Use(middleware.WriteAuthRequired()).POST("/register-fcm-token", func(ctx *gin.Context) {
		var json models.RegisterFCMToken

		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		_, error := sessions.RegisterFCMToken(ctx.MustGet(middleware.SESSION_KEY).(models.Session), json.Token)

		if error != customerrors.Success {
			ctx.JSON(error.Status, error)
			return
		}

		ctx.JSON(200, gin.H{})
	})
}
