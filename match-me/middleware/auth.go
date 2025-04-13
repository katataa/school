package middleware

import (
	"log"
	"match-me/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Handle WebSocket requests
		if isWebSocketRequest(c) {
			tokenString = c.Query("token") // Extract token from query parameter
			if tokenString == "" {
				log.Println("Token is missing in WebSocket query")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
		} else {
			// Handle standard HTTP requests
			tokenString = c.GetHeader("Authorization")
			if tokenString == "" {
				log.Println("Authorization header is missing")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
				c.Abort()
				return
			}

			// Validate "Bearer" format
			if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
				log.Println("Invalid token format:", tokenString)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
				c.Abort()
				return
			}

			tokenString = tokenString[7:] // Remove "Bearer " prefix
		}

		// Parse and validate the token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			log.Println("Token parsing error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract user ID from claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			log.Println("Invalid token payload")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		log.Printf("Authenticated user: %d\n", uint(userID))
		c.Set("user_id", uint(userID)) // Set user ID in the request context
		c.Next()
	}
}

// Helper function to determine if the request is a WebSocket request
func isWebSocketRequest(c *gin.Context) bool {
	return strings.HasPrefix(c.Request.URL.Path, "/ws/")
}
