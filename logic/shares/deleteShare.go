package shares

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/logic/ws"
	"location-share-backend/models"

	"github.com/google/uuid"
)

func DeleteShare(shareId uuid.UUID, userId uuid.UUID) (error customerrors.APIError) {
	var share models.Share
	err := initializers.DB.Where("id = ? AND shared_by = ?", shareId, userId).First(&share).Error

	if err != nil {
		return customerrors.ErrDeleteShare
	}

	// Delete share
	err = initializers.DB.Where("id = ? AND shared_by = ?", shareId, userId).Delete(&models.Share{}).Error

	if err != nil {
		return customerrors.ErrDeleteShare
	}

	ws.BroadcastShareDeleted([]uuid.UUID{share.SharedWith, share.SharedBy}, shareId)

	return customerrors.Success
}
