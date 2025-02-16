package routers

import (
	"location-share-backend/logic/ws"
	"location-share-backend/middleware"
	"location-share-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Websocket(router *gin.RouterGroup, hub *ws.Hub) {
	router.GET("", middleware.ReadAuthRequired(), ws.WSHandler(func(c *gin.Context) (uuid.UUID, bool) {
		// Hier wird der authentifizierte User aus dem Context gelesen.
		// Deine Middleware hat z. B. unter middleware.SESSION_KEY die Session gespeichert.
		sess, exists := c.Get(middleware.SESSION_KEY)
		if !exists {
			return uuid.Nil, false
		}
		session, ok := sess.(models.Session)
		if !ok {
			return uuid.Nil, false
		}
		return session.UserID, true
	}))
}
