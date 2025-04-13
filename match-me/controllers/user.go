package controllers

import (
	"io"
	"log"
	"match-me/config"
	"match-me/models"
	"match-me/utils"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(c *gin.Context) {
	requestedID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	currentUserID := c.GetUint("user_id")

	var user models.User
	if err := config.DB.First(&user, requestedID).Error; err != nil {
		// Return 404 if not found
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	authorized := checkUserAuthorization(currentUserID, uint(requestedID))
	if !authorized {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              user.ID,
		"name":            user.Name,
		"profile_picture": user.ProfilePicture,
	})
}

func GetUserProfile(c *gin.Context) {
	requestedID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	currentUserID := c.GetUint("user_id")

	var user models.User
	if err := config.DB.First(&user, requestedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	if !checkUserAuthorization(currentUserID, uint(requestedID)) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   user.ID,
		"info": user.Info,
	})
}

func GetUserBio(c *gin.Context) {
	requestedID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	currentUserID := c.GetUint("user_id")

	var user models.User
	if err := config.DB.First(&user, requestedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	if !checkUserAuthorization(currentUserID, uint(requestedID)) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"name":      user.Name,
		"age":       user.Age,
		"location":  user.Location,
		"interests": user.Interests,
		"info":      user.Info,
		"gender":    user.Gender,
	})
}

func checkUserAuthorization(currentUserID, requestedUserID uint) bool {
	if currentUserID == requestedUserID {
		return true
	}

	var conn models.Connection
	err := config.DB.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			currentUserID, requestedUserID, requestedUserID, currentUserID).
		First(&conn).Error

	if err == nil {
		if conn.Status == "pending" || conn.Status == "accepted" {
			return true
		}
	}
	if recIDs, ok := recentRecommendations[currentUserID]; ok {
		for _, rid := range recIDs {
			if rid == requestedUserID {
				return true
			}
		}
	}

	return false
}

func GetMe(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              user.ID,
		"name":            user.Name,
		"email":           user.Email,
		"profile_picture": user.ProfilePicture,
	})
}

func GetMeProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   user.ID,
		"info": user.Info,
	})
}

func GetMeBio(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"name":      user.Name,
		"age":       user.Age,
		"location":  user.Location,
		"interests": user.Interests,
		"info":      user.Info,
		"gender":    user.Gender,
	})
}

func RegisterUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User registered successfully!",
		"token":   token,
	})

}

func LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		log.Println("Unauthorized: user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Println("Fetching profile for user_id:", userID)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		log.Println("User not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	log.Println("Fetched user profile:", user)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"id":               user.ID,
			"name":             user.Name,
			"email":            user.Email,
			"info":             user.Info,
			"interests":        user.Interests,
			"location":         user.Location,
			"age":              user.Age,
			"profilePicture":   user.ProfilePicture,
			"latitude":         user.Latitude,
			"longitude":        user.Longitude,
			"preferred_radius": user.PreferredRadius,
			"gender":           user.Gender,
			"lookingFor":       user.LookingFor,
		},
	})
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	info := c.PostForm("info")
	interests := c.PostForm("interests")
	location := c.PostForm("location")
	name := c.PostForm("name")
	ageStr := c.PostForm("age")
	latStr := c.PostForm("latitude")
	lonStr := c.PostForm("longitude")
	prefRadiusStr := c.PostForm("preferredRadius")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age format"})
		return
	}

	latitude, _ := strconv.ParseFloat(latStr, 64)
	longitude, _ := strconv.ParseFloat(lonStr, 64)
	preferredRadius, _ := strconv.ParseFloat(prefRadiusStr, 64)
	gender := c.PostForm("gender")
	lookingFor := c.PostForm("lookingFor")

	file, fileHeader, err := c.Request.FormFile("profilePicture")
	var profilePicturePath string
	if err == nil {
		profilePicturePath = "uploads/" + fileHeader.Filename
		out, err := os.Create(profilePicturePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		defer out.Close()
		io.Copy(out, file)
	} else if err != http.ErrMissingFile {
		log.Println("Error uploading file:", err)
	}

	if name == "" || info == "" || interests == "" || location == "" || gender == "" || lookingFor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Name = name
	user.Info = info
	user.Interests = strings.TrimSpace(interests)
	user.Location = location
	user.Age = age
	user.Latitude = latitude
	user.Longitude = longitude
	user.PreferredRadius = preferredRadius

	if profilePicturePath != "" {
		user.ProfilePicture = profilePicturePath
	} else if user.ProfilePicture == "" {
		user.ProfilePicture = "uploads/default-profile.png"
	}
	user.Gender = gender
	user.LookingFor = lookingFor

	if err := config.DB.Save(&user).Error; err != nil {
		log.Println("Error saving user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Profile updated successfully", "data": user})
}

func RemoveProfilePicture(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.ProfilePicture = "uploads/default-profile.png"

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove profile picture"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile picture removed successfully",
		"data": gin.H{
			"profilePicture": user.ProfilePicture,
		},
	})
}
