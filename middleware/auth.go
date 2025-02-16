package middleware

import (
	"location-share-backend/initializers"
	"location-share-backend/models"
	"location-share-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const SESSION_KEY = "session"

func checkTokenExistsAndValid(token string) (sessionId uuid.UUID, valid bool) {
	if token == "" {
		return uuid.Nil, false
	}

	// If the token starts with "Bearer ",
	// Cut the "Bearer " prefix
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Verify the token
	claims, err := utils.VerifyToken(token)

	if err != nil || claims == nil {
		return uuid.Nil, false
	}

	sessionUUID, err := uuid.Parse(claims["session_id"].(string))

	if err != nil {
		return uuid.Nil, false
	}

	return sessionUUID, true
}

func WriteAuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, valid := checkTokenExistsAndValid(ctx.GetHeader("Authorization"))

		if token == uuid.Nil || valid == false {
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Check if the token is valid
		db := initializers.DB

		var session models.Session
		var err = db.Where(&models.Session{
			ID:       token,
			ReadOnly: false,
		}).First(&session).Error

		if err != nil {
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set(SESSION_KEY, session)
		ctx.Next()
	}
}

func ReadAuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, valid := checkTokenExistsAndValid(ctx.GetHeader("Authorization"))

		if token == uuid.Nil || valid == false {
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Check if the token is valid
		db := initializers.DB

		var session models.Session
		var err = db.Where(&models.Session{
			ID: token,
		}).First(&session).Error

		if err != nil {
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set(SESSION_KEY, session)
		ctx.Next()
	}
}
