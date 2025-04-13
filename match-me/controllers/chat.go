package controllers

import (
	"log"
	"match-me/config"
	"match-me/models"
	"match-me/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wsConnections = make(map[uint]*websocket.Conn)
var wsChatListConnections = make(map[uint]*websocket.Conn)
var activeChats = make(map[uint]uint)

func WebSocketChatHandler(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Println("Token is missing from query string")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		log.Printf("Invalid token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		log.Println("Invalid token payload")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDFloat)

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	wsConnections[uint(userID)] = conn
	log.Printf("User %d connected to DM socket, total DM connections: %d\n", uint(userID), len(wsConnections))

	defer func() {
		conn.Close()
		delete(wsConnections, userID)
		log.Printf("User %d disconnected from DM socket, total DM connections: %d\n", uint(userID), len(wsConnections))
	}()

	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			log.Println("WebSocket connection closed normally")
			break
		} else if err != nil {
			log.Println("Error reading WebSocket message:", err)
			break
		}

		var chat models.Chat
		dbErr := config.DB.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
			uint(userID), message.ReceiverID, message.ReceiverID, uint(userID)).First(&chat).Error
		if dbErr != nil {
			chat = models.Chat{
				User1ID: uint(userID),
				User2ID: message.ReceiverID,
			}
			if createErr := config.DB.Create(&chat).Error; createErr != nil {
				log.Println("Error creating new chat:", createErr)
				continue
			}
		}

		message.ChatID = chat.ID
		message.SenderID = uint(userID)
		message.Timestamp = time.Now()

		if dbErr := config.DB.Create(&message).Error; dbErr != nil {
			log.Println("Error saving message:", dbErr)
			continue
		}

		if receiverConn, ok := wsConnections[message.ReceiverID]; ok {
			if sendErr := receiverConn.WriteJSON(message); sendErr != nil {
				log.Println("Error sending WebSocket message to DM socket:", sendErr)
			}
		} else {
			log.Printf("User %d is offline (DM socket), message saved in DB", message.ReceiverID)
		}

		if chatListConn, ok := wsChatListConnections[message.ReceiverID]; ok {
			outMsg := map[string]interface{}{
				"type":        "new_message",
				"chat_id":     message.ChatID,
				"sender_id":   message.SenderID,
				"receiver_id": message.ReceiverID,
				"content":     message.Content,
				"timestamp":   message.Timestamp.Format(time.RFC3339),
				"unread_count": func(receiverID uint, chatID uint) int {
					// check if the receiver is currently viewing this exact chat
					if activeChatID, ok := activeChats[message.ReceiverID]; ok && activeChatID == message.ChatID {
						// mark ALL unread messages in that chat for this user as read
						config.DB.Model(&models.Message{}).
							Where("chat_id = ? AND receiver_id = ? AND is_read = ?", chatID, receiverID, false).
							Update("is_read", true)
					}

					var count int64
					config.DB.Model(&models.Message{}).
						Where("chat_id = ? AND receiver_id = ? AND is_read = ?", chatID, receiverID, false).
						Count(&count)
					return int(count)
				}(message.ReceiverID, message.ChatID),
			}
			if sendErr := chatListConn.WriteJSON(outMsg); sendErr != nil {
				log.Println("Error sending to chat list ws (receiver):", sendErr)
			}
		}

		if chatListConn, ok := wsChatListConnections[message.SenderID]; ok {
			outMsg := map[string]interface{}{
				"type":        "new_message",
				"chat_id":     message.ChatID,
				"sender_id":   message.SenderID,
				"receiver_id": message.ReceiverID,
				"content":     message.Content,
				"timestamp":   message.Timestamp.Format(time.RFC3339),
			}
			if sendErr := chatListConn.WriteJSON(outMsg); sendErr != nil {
				log.Println("Error sending to chat list ws (sender):", sendErr)
			}
		}
	}
}

func WebSocketChatListHandler(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Println("Token is missing from query string")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		log.Printf("Invalid token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDFloat)

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	wsChatListConnections[userID] = conn
	log.Printf("Chat list WebSocket connected for user %d", userID)

	defer func() {
		conn.Close()
		delete(wsChatListConnections, userID)
		delete(activeChats, userID)
		log.Printf("Chat list WebSocket disconnected for user %d", userID)
	}()

	for {
		var incoming map[string]interface{}
		if err := conn.ReadJSON(&incoming); err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("Chat list WebSocket closed for user %d", userID)
				break
			}
			log.Println("Error reading from chat list WebSocket:", err)
			break
		}

		if incoming["type"] == "active_chat" {
			if otherUserIDFloat, ok := incoming["chat_id"].(float64); ok {
				otherUserID := uint(otherUserIDFloat)

				var chat models.Chat
				dbErr := config.DB.
					Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
						userID, otherUserID, otherUserID, userID).
					First(&chat).Error

				if dbErr != nil {
					log.Printf("Chat row not found for %d & %d: %v", userID, otherUserID, dbErr)
					continue
				}

				activeChats[userID] = chat.ID
				log.Printf("User %d is now viewing chat ID: %d (for user pair: %d, %d)",
					userID, chat.ID, userID, otherUserID)
			}
		}

		if incoming["type"] == "inactive_chat" {
			delete(activeChats, userID)
			log.Printf("User %d is no longer viewing any chat", userID)
		}
	}
}

func GetChatMessages(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	otherUserIDStr := c.Param("userId")
	otherUserID, err := strconv.ParseUint(otherUserIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var connection models.Connection
	err = config.DB.Where("(sender_id = ? AND receiver_id = ? OR sender_id = ? AND receiver_id = ?) AND status = ?",
		userID, uint(otherUserID), uint(otherUserID), userID, "accepted").First(&connection).Error
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not connected with this user."})
		return
	}

	var chat models.Chat
	err = config.DB.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		userID, uint(otherUserID), uint(otherUserID), userID).First(&chat).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	cursor := c.Query("cursor")
	limit := 20
	var messages []models.Message
	query := config.DB.Where("chat_id = ?", chat.ID).Order("timestamp DESC").Limit(limit)

	if cursor != "" {
		cursor = strings.Replace(cursor, " ", "+", -1)
		parsedTime, err := time.Parse(time.RFC3339, cursor)
		if err != nil {
			log.Printf("Invalid cursor format: %s, error: %v", cursor, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor format"})
			return
		}
		query = query.Where("timestamp < ?", parsedTime)
	}

	if err := query.Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	reverseMessages := make([]models.Message, len(messages))
	for i, msg := range messages {
		reverseMessages[len(messages)-1-i] = msg
	}

	nextCursor := ""
	if len(messages) > 0 {
		nextCursor = messages[len(messages)-1].Timestamp.Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, gin.H{"messages": reverseMessages, "nextCursor": nextCursor})
}

func SendMessage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		ReceiverID uint   `json:"receiver_id" binding:"required"`
		Content    string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var chat models.Chat
	if err := config.DB.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		userID, input.ReceiverID, input.ReceiverID, userID).First(&chat).Error; err != nil {
		chat = models.Chat{User1ID: userID.(uint), User2ID: input.ReceiverID}
		if createErr := config.DB.Create(&chat).Error; createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
			return
		}
	}

	message := models.Message{
		ChatID:    chat.ID,
		SenderID:  userID.(uint),
		Content:   input.Content,
		Timestamp: time.Now(),
	}

	if err := config.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func GetChats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var chats []models.Chat
	if err := config.DB.Preload("User1").Preload("User2").
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Find(&chats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}

	response := []gin.H{}
	for _, chat := range chats {
		otherUser := chat.User1
		if otherUser.ID == userID.(uint) {
			otherUser = chat.User2
		}

		var latestMessage models.Message
		if err := config.DB.Where("chat_id = ?", chat.ID).Order("timestamp DESC").First(&latestMessage).Error; err != nil {
			latestMessage = models.Message{}
		}

		response = append(response, gin.H{
			"id":                       chat.ID,
			"user_id":                  otherUser.ID,
			"name":                     otherUser.Name,
			"profile_picture":          otherUser.ProfilePicture,
			"latest_message":           latestMessage.Content,
			"latest_message_timestamp": latestMessage.Timestamp.Format("2006-01-02T15:04:05Z"),
			"unread_count": func() int {
				var count int64
				config.DB.Model(&models.Message{}).
					Where("chat_id = ? AND receiver_id = ? AND is_read = ?", chat.ID, userID, false).
					Count(&count)
				return int(count)
			}(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"chats": response})
}

func MarkMessagesAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	chatID := c.Param("chatId")

	if err := config.DB.Model(&models.Message{}).
		Where("chat_id = ? AND receiver_id = ?", chatID, userID).
		Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark messages as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Messages marked as read"})
}
