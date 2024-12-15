package shares

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"

	"github.com/google/uuid"
)

func GetShares(userId uuid.UUID) (shares []models.Share, error customerrors.APIError) {
	// Get all shares that are shared with or by the user
	err := initializers.DB.Where("shared_by = ? OR shared_with = ?", userId, userId).Find(&shares).Error

	if err != nil {
		return nil, customerrors.ErrGettingShare
	}

	return shares, customerrors.Success
}
