package controllers

import (
	"log"
	"match-me/config"
	"match-me/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendConnectionRequest sends a connection request to another user
func SendConnectionRequest(c *gin.Context) {
	senderID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		ReceiverID uint `json:"receiver_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	connection := models.Connection{
		SenderID:   senderID.(uint),
		ReceiverID: input.ReceiverID,
		Status:     "pending",
	}

	if err := config.DB.Create(&connection).Error; err != nil {
		log.Printf("Error creating connection request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send connection request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection request sent"})
}

// GetConnectionRequests fetches all pending connection requests for the authenticated user
func GetConnectionRequests(c *gin.Context) {
	receiverID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var connections []models.Connection
	if err := config.DB.Preload("Sender").Where("receiver_id = ? AND status = ?", receiverID, "pending").Find(&connections).Error; err != nil {
		log.Printf("Error fetching connection requests: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch connection requests"})
		return
	}

	var response []gin.H
	for _, connection := range connections {
		response = append(response, gin.H{
			"id":     connection.ID,
			"sender": connection.Sender,
			"status": connection.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{"requests": response})
}

// AcceptConnectionRequest accepts a connection request
func AcceptConnectionRequest(c *gin.Context) {
	var input struct {
		RequestID uint `json:"request_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var connection models.Connection
	if err := config.DB.First(&connection, input.RequestID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Connection request not found"})
		return
	}

	connection.Status = "accepted"
	if err := config.DB.Save(&connection).Error; err != nil {
		log.Printf("Error accepting connection request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept connection request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection request accepted"})
}

func DeclineConnectionRequest(c *gin.Context) {
	var input struct {
		RequestID uint `json:"request_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var connection models.Connection
	if err := config.DB.First(&connection, input.RequestID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Connection request not found"})
		return
	}

	// Add an entry to DeclinedUser for both users
	declinedByReceiver := models.DeclinedUser{
		UserID:         connection.ReceiverID,
		DeclinedUserID: connection.SenderID,
	}
	declinedBySender := models.DeclinedUser{
		UserID:         connection.SenderID,
		DeclinedUserID: connection.ReceiverID,
	}
	if err := config.DB.Create(&declinedByReceiver).Error; err != nil {
		log.Printf("Error saving declined user (receiver perspective): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decline connection request"})
		return
	}
	if err := config.DB.Create(&declinedBySender).Error; err != nil {
		log.Printf("Error saving declined user (sender perspective): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decline connection request"})
		return
	}

	// Delete the connection request
	if err := config.DB.Delete(&connection).Error; err != nil {
		log.Printf("Error deleting connection request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decline connection request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection request declined"})
}

// GetConnections fetches all accepted connections for the authenticated user
func GetConnections(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var connections []models.Connection
	if err := config.DB.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? OR receiver_id = ?) AND status = ?", userID, userID, "accepted").
		Find(&connections).Error; err != nil {
		log.Printf("Error fetching connections: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch connections"})
		return
	}

	var connectedIDs []uint
	for _, connection := range connections {
		if connection.SenderID == userID {
			connectedIDs = append(connectedIDs, connection.ReceiverID)
		} else {
			connectedIDs = append(connectedIDs, connection.SenderID)
		}
	}

	// The spec says just return a list of { "id": ... }
	result := []gin.H{}
	for _, cid := range connectedIDs {
		result = append(result, gin.H{"id": cid})
	}

	c.JSON(http.StatusOK, gin.H{"connections": result})
}

func DisconnectUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Remove the connection
	if err := config.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, input.UserID, input.UserID, userID).
		Delete(&models.Connection{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disconnect user"})
		return
	}

	// Delete the chat and associated messages
	var chat models.Chat
	if err := config.DB.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		userID, input.UserID, input.UserID, userID).First(&chat).Error; err == nil {
		// Delete messages associated with the chat
		if err := config.DB.Where("chat_id = ?", chat.ID).Delete(&models.Message{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete messages"})
			return
		}
		// Delete the chat itself
		if err := config.DB.Delete(&chat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User disconnected and chat removed successfully"})
}
