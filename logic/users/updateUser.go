package users

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"
	"location-share-backend/utils"
)

func UpdateUser(userUpdate models.UserUpdate) (updatedUser models.User, error customerrors.APIError) {
	// Check if the user exists
	user, error := GetUser(userUpdate.ID.String())

	if error != customerrors.Success {
		return models.User{}, error
	}

	if userUpdate.Username != "" {
		// Check if a user with the same username exists
		var existingUser models.User
		err := initializers.DB.Where("id != ? AND username = ?", user.ID, userUpdate.Username).First(&existingUser).Error

		if err == nil {
			return models.User{}, customerrors.ErrUserExists
		}

		user.Username = userUpdate.Username
	}

	if userUpdate.Name != "" {
		user.Name = userUpdate.Name
	}

	if userUpdate.Email != "" {
		// Check if a user with the same email exists
		var existingUser models.User
		err := initializers.DB.Where("id != ? AND email = ?", user.ID, userUpdate.Email).First(&existingUser).Error

		if err == nil {
			return models.User{}, customerrors.ErrEmailExists
		}

		user.Email = userUpdate.Email
	}

	if userUpdate.Password != "" {
		hash, err := utils.HashPassword(userUpdate.Password)

		if err != nil {
			return models.User{}, customerrors.ErrHashPassword
		}

		user.Password = hash
	}

	err := initializers.DB.Save(&user).Error

	if err != nil {
		return models.User{}, customerrors.ErrUpdateUser
	}

	return user, customerrors.Success
}
