package locations

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"
	"time"

	"github.com/google/uuid"
)

func DeleteLocations(userId uuid.UUID, from time.Time, to time.Time) (error customerrors.APIError) {
	// Delete all locations that are shared with or by the user
	err := initializers.DB.Where("user_id = ? AND created_at BETWEEN ? AND ?", userId, from, to).Delete(&models.Location{}).Error

	if err != nil {
		return customerrors.ErrDeleteLocation
	}

	return customerrors.Success
}
