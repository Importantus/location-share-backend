package users

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"
)

/**
* Returns a list of users. If a userId is provided, it will return a list with only that user.
* If no userId is provided, it will return a list with all users.
 */
func ListUsers(userId string) (users []models.User, error customerrors.APIError) {
	if userId != "" {
		user, err := GetUser(userId)
		if err == customerrors.Success {
			users = append(users, user)
		}
	}

	if len(users) == 0 {
		// Get all users
		err := initializers.DB.Find(&users).Error
		if err != nil {
			return []models.User{}, customerrors.ErrListUsers
		}
	}

	return users, customerrors.Success
}
