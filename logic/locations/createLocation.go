package locations

import (
	"fmt"
	"location-share-backend/customerrors"
	"location-share-backend/initializers"
	"location-share-backend/logic/shares"
	"location-share-backend/logic/ws"
	"location-share-backend/models"
	"location-share-backend/utils"

	"github.com/google/uuid"
)

func CreateLocation(session models.Session, locationCreate models.LocationCreate) (location models.Location, error customerrors.APIError) {
	utils.CopyStruct(&locationCreate, &location)

	location.ID = uuid.New()
	location.SessionID = session.ID

	err := initializers.DB.Create(&location)

	if err.Error != nil {
		return location, customerrors.ErrCreateLocation
	}

	// Get shares of the user
	shares, error := shares.GetShares(session.UserID)
	if error != customerrors.Success {
		fmt.Println("Error getting shares for user:", session.UserID.String())
		return location, customerrors.ErrCreateLocation
	}

	// Get user IDs from shares
	userIDs := make([]uuid.UUID, 0)
	for _, share := range shares {
		if share.SharedBy == session.UserID {
			userIDs = append(userIDs, share.SharedWith)
		}
	}

	// Broadcast the new location to all connected clients
	ws.BroadcastLocation(userIDs, location)

	return location, customerrors.Success
}
