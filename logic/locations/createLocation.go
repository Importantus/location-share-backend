package locations

import (
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/models"
	"location-share-backend/utils"

	"github.com/google/uuid"
)

func CreateLocation(sessionId uuid.UUID, locationCreate models.LocationCreate) (location models.Location, error customerrors.APIError) {
	utils.CopyStruct(&locationCreate, &location)

	location.ID = uuid.New()
	location.SessionID = sessionId

	err := initializers.DB.Create(&location)

	if err.Error != nil {
		return location, customerrors.ErrCreateLocation
	}

	return location, customerrors.Success
}
