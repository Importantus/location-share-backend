package users

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
)

func DeleteUser(userId string) (error customerrors.APIError) {
	// Check if the user exists
	user, error := GetUser(userId)

	if error != customerrors.Success {
		return error
	}

	// Delete the user
	err := initializers.DB.Delete(&user).Error

	if err != nil {
		return customerrors.ErrDeleteUser
	}

	return customerrors.Success
}
