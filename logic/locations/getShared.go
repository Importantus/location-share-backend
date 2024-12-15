package locations

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"
	"time"

	"github.com/google/uuid"
)

func GetSharedLocations(userId uuid.UUID) (locations []models.Location, error customerrors.APIError) {
	// Get all shares where the userId = shared_with and valid_until > now
	var shares []models.Share
	err := initializers.DB.Where("shared_with = ? AND (valid_until > ? OR valid_until IS NULL)", userId, time.Now()).Find(&shares).Error

	if err != nil {
		return nil, customerrors.ErrGettingLocation
	}

	// For each share, get the newest location
	for _, share := range shares {
		var location models.Location
		err := initializers.DB.Where("user_id = ?", share.SharedBy).Order("timestamp desc").First(&location).Error

		if err != nil {
			return nil, customerrors.ErrGettingLocation
		}

		locations = append(locations, location)
	}

	return locations, customerrors.Success
}
