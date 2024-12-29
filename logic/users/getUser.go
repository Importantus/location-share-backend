package users

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"

	"github.com/google/uuid"
)

func GetUser(userId string) (user models.User, error customerrors.APIError) {
	id, err := uuid.Parse(userId)

	if err != nil {
		return models.User{}, customerrors.ErrInvalidUUID
	}

	// Check if the user exists
	err = initializers.DB.Where(&models.User{
		ID: id,
	}).First(&user).Error

	if err != nil {
		return models.User{}, customerrors.ErrUserNotFound
	}

	return user, customerrors.Success
}
