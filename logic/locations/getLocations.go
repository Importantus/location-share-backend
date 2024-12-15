package locations

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"
	"time"

	"github.com/google/uuid"
)

func GetLocations(userId uuid.UUID, from time.Time, to time.Time) (locations []models.Location, error customerrors.APIError) {
	// Get all locations that are shared with or by the user
	err := initializers.DB.Where("user_id = ? AND created_at BETWEEN ? AND ?", userId, from, to).Find(&locations).Error

	if err != nil {
		return nil, customerrors.ErrGettingLocation
	}

	return locations, customerrors.Success
}
