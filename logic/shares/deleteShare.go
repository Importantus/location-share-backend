package shares

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"

	"github.com/google/uuid"
)

func DeleteShare(shareId uuid.UUID, userId uuid.UUID) (error customerrors.APIError) {
	// Delete share
	err := initializers.DB.Where("id = ? AND shared_by = ?", shareId, userId).Delete(&models.Share{}).Error

	if err != nil {
		return customerrors.ErrDeleteShare
	}

	return customerrors.Success
}
