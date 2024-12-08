package users

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"
	"location-share-backend/utils"

	"github.com/google/uuid"
)

func CreateUser(user models.UserCreate) (newUser models.User, error customerrors.APIError) {
	// Check if a user with the same username exists
	var existingUser models.User
	err := initializers.DB.Where(&models.User{
		Username: user.Username,
	}).First(&existingUser).Error

	if err == nil {
		return models.User{}, customerrors.ErrUserExists
	}

	// Check if a user with the same email exists
	err = initializers.DB.Where(&models.User{
		Email: user.Email,
	}).First(&existingUser).Error

	if err == nil {
		return models.User{}, customerrors.ErrEmailExists
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		return models.User{}, customerrors.ErrHashPassword
	}

	// Create the user
	newUser = models.User{
		ID:       uuid.New(),
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	err = initializers.DB.Create(&newUser).Error

	if err != nil {
		return models.User{}, customerrors.ErrCreateUser
	}

	return newUser, customerrors.Success
}
