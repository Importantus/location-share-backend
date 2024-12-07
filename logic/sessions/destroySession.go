package sessions

import (
	"location-share-backend/initializers"
	"location-share-backend/models"
)

func DestroySession(session models.Session) (err error) {
	err = initializers.DB.Delete(&session).Error
	return
}
