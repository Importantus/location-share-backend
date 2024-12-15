package shares

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"
	"time"

	"github.com/google/uuid"
)

func CreateShare(from uuid.UUID, with uuid.UUID, validUntil *time.Time) (share models.Share, error customerrors.APIError) {
	// Check if a share from from to with already exists where valid_until is in the future or nil
	var existingShare models.Share

	initializers.DB.Where("shared_by = ? AND shared_with = ? AND (valid_until > ? OR valid_until IS NULL)", from, with, time.Now()).First(&existingShare)

	if existingShare.ID != uuid.Nil {
		return existingShare, customerrors.ErrShareExists
	}

	// Create share
	share = models.Share{
		ID:         uuid.New(),
		ValidUntil: validUntil,
		SharedWith: with,
		SharedBy:   from,
	}

	initializers.DB.Create(&share)

	return share, customerrors.Success
}
