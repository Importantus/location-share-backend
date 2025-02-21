package fcm

import (
	"context"
	"fmt"
	"location-share-backend/initializers"
	"location-share-backend/models"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

// FCMMessage describes the structure of the message to be sent to FCM.
type FCMMessage struct {
	To       string                 `json:"to"`
	Priority string                 `json:"priority"`
	Data     map[string]interface{} `json:"data"`
}

var (
	firebaseApp     *firebase.App
	messagingClient *messaging.Client
)

// InitializeFirebase initializes the Firebase app and the messaging client.
// This function should be called once at application startup.
func InitializeFirebase() error {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		return fmt.Errorf("Error loading configuration: %v", err)
	}

	ctx := context.Background()
	opt := option.WithCredentialsFile(config.GoogleApplicationCredentials)

	firebaseApp, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("Error initializing Firebase app: %v", err)
	}

	messagingClient, err = firebaseApp.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving messaging client: %v", err)
	}

	return nil
}

// SendFCMNotification sends a push notification via FCM to the specified device.
func SendFCMNotification(token string) error {
	ctx := context.Background()

	message := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"action": "REQUEST_LOCATION_UPDATE",
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
	}

	response, err := messagingClient.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("Error sending message: %v", err)
	}
	log.Printf("Message sent successfully: %s", response)
	return nil
}

// RegisterFcmRequest expects the client to provide UserID and FCM token.
type RegisterFcmRequest struct {
	UserID   string `json:"userId"`
	FCMToken string `json:"fcmToken"`
}

// WakeUpRequest contains only the UserID.
type WakeUpRequest struct {
	UserID string `json:"userId"`
}

// WakeUpHandler searches for the active session based on the UserID and sends an FCM message.
func WakeUpHandler(c *gin.Context) {
	var req WakeUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Search for the active session.
	var session models.Session
	result := initializers.DB.Where("user_id = ? AND writing = ?", req.UserID, true).First(&session)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Active session not found"})
		return
	}

	// Send the push notification to the FCM device ID stored in the session.
	if err := SendFCMNotification(session.FCMDeviceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Wake-up notification sent"})
}
