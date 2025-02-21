package routers

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/logic/fcm"
	"location-share-backend/logic/users"
	"location-share-backend/middleware"
	"location-share-backend/models"

	"github.com/gin-gonic/gin"
)

func Users(router *gin.RouterGroup) {
	router.POST("", func(ctx *gin.Context) {
		var json models.UserCreate

		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		config, error := initializers.LoadConfig(".")

		if error != nil {
			ctx.JSON(500, gin.H{"error": error.Error()})
			return
		}

		if config.RegistrationSecret != "" && *json.RegistrationSecret != config.RegistrationSecret {
			ctx.JSON(403, gin.H{"error": "Invalid registration secret"})
			return
		}

		user, err := users.CreateUser(json)

		if err != customerrors.Success {
			ctx.JSON(err.Status, err)
			return
		}

		ctx.JSON(200, models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	})

	router.Use(middleware.ReadAuthRequired()).GET("", func(ctx *gin.Context) {
		userId := ctx.Query("id")

		users, err := users.ListUsers(userId)

		if err != customerrors.Success {
			ctx.JSON(err.Status, err)
			return
		}

		// Map users to UsersResponse
		usersResponse := make([]models.UserResponse, len(users))
		for i, user := range users {
			usersResponse[i] = models.UserResponse{
				ID:       user.ID,
				Username: user.Username,
				Name:     user.Name,
				Email: func() string {
					if user.ID == ctx.MustGet(middleware.SESSION_KEY).(models.Session).UserID {
						return user.Email
					}
					return ""
				}(),
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}
		}

		ctx.JSON(200, usersResponse)
	})

	router.Use(middleware.WriteAuthRequired()).DELETE("", func(ctx *gin.Context) {
		userId := ctx.Query("id")

		err := users.DeleteUser(userId)

		if err != customerrors.Success {
			ctx.JSON(err.Status, err)
			return
		}

		ctx.JSON(200, gin.H{})
	})

	router.Use(middleware.WriteAuthRequired()).PUT("", func(ctx *gin.Context) {
		var json models.UserUpdate

		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if json.ID != ctx.MustGet(middleware.SESSION_KEY).(models.Session).UserID {
			ctx.JSON(403, gin.H{"error": "You can only update your own user"})
			return
		}

		user, err := users.UpdateUser(json)

		if err != customerrors.Success {
			ctx.JSON(err.Status, err)
			return
		}

		ctx.JSON(200, user)
	})

	router.Use(middleware.ReadAuthRequired()).POST("wake-up", func(ctx *gin.Context) {
		fcm.WakeUpHandler(ctx)
	})
}
