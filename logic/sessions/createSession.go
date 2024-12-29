package sessions

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"
	"location-share-backend/utils"

	"github.com/google/uuid"
)

func CreateSession(sessionCreate *models.SessionCreate) (key string, id uuid.UUID, err customerrors.APIError) {
	// Check if the user exists
	var user models.User
	dbErr := initializers.DB.Where(&models.User{
		Username: sessionCreate.Username,
	}).First(&user).Error

	if dbErr != nil {
		return "", uuid.UUID{}, customerrors.ErrUserNotFound
	}

	// Check if the password is correct
	if !utils.CheckPasswordHash(sessionCreate.Password, user.Password) {
		return "", uuid.UUID{}, customerrors.ErrInvalidPassword
	}

	// Create the session
	session := models.Session{
		ID:       uuid.New(),
		Name:     sessionCreate.Name,
		Writing:  sessionCreate.Writing,
		UserID:   user.ID,
		ReadOnly: sessionCreate.ReadOnly,
	}

	createErr := initializers.DB.Create(&session).Error
	if createErr != nil {
		return "", uuid.UUID{}, customerrors.ErrSessionCreationFailed
	}

	key, tokenErr := utils.CreateToken(session.ID.String())

	if tokenErr != nil {
		return "", uuid.UUID{}, customerrors.ErrTokenCreationFailed
	}

	return key, session.ID, customerrors.Success
}
