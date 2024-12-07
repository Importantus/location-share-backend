package sessions

import (
	"location-share-backend/initializers"
	"location-share-backend/models"
	"location-share-backend/utils"

	"github.com/google/uuid"
)

func CreateSession(sessionCreate *models.SessionCreate) (key string, err error) {
	// Check if the user exists
	var user models.User
	err = initializers.DB.Where(&models.User{
		Username: sessionCreate.Username,
	}).First(&user).Error

	if err != nil {
		return "", err
	}

	// Check if the password is correct
	if !utils.CheckPasswordHash(sessionCreate.Password, user.Password) {
		return "", err
	}

	// Create the session
	session := models.Session{
		ID:       uuid.New(),
		Name:     sessionCreate.Name,
		Writing:  sessionCreate.Writing,
		UserID:   user.ID,
		ReadOnly: sessionCreate.ReadOnly,
	}

	err = initializers.DB.Create(session).Error
	if err != nil {
		key, err = utils.CreateToken(session.ID.String())
	}

	return
}
