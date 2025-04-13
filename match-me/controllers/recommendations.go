package controllers

import (
	"log"
	"match-me/config"
	"match-me/models"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Recommendation struct {
	User  models.User
	Score int
}

var recentRecommendations = make(map[uint][]uint)

func GetRecommendations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	locationFilter := c.Query("location")
	ageFilter := c.Query("age")
	hobbiesFilter := c.Query("hobbies")
	genderFilter := c.Query("gender")
	mode := c.Query("mode")

	var currentUser models.User
	if err := config.DB.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Name == "" || currentUser.Info == "" || currentUser.Interests == "" || currentUser.Location == "" || currentUser.Age <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Complete your profile before getting recommendations."})
		return
	}

	var declinedUsers []models.DeclinedUser
	if err := config.DB.Where("user_id = ? OR declined_user_id = ?", userID, userID).Find(&declinedUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching declined users"})
		return
	}

	declinedUserIDs := make([]uint, len(declinedUsers))
	for i, d := range declinedUsers {
		if d.UserID == userID {
			declinedUserIDs[i] = d.DeclinedUserID
		} else {
			declinedUserIDs[i] = d.UserID
		}
	}

	var connections []models.Connection
	if err := config.DB.Where("sender_id = ? OR receiver_id = ?", userID, userID).Find(&connections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching connected users"})
		return
	}

	connectedUserIDs := make([]uint, 0)
	for _, conn := range connections {
		if conn.SenderID == userID {
			connectedUserIDs = append(connectedUserIDs, conn.ReceiverID)
		} else {
			connectedUserIDs = append(connectedUserIDs, conn.SenderID)
		}
	}

	query := config.DB.Where("id != ?", currentUser.ID)
	if len(declinedUserIDs) > 0 {
		query = query.Where("id NOT IN ?", declinedUserIDs)
	}
	if len(connectedUserIDs) > 0 {
		query = query.Where("id NOT IN ?", connectedUserIDs)
	}

	if mode == "all" || mode == "location" {
		if locationFilter != "" {
			query = query.Where("LOWER(location) ILIKE ?", "%"+strings.ToLower(locationFilter)+"%")
		}
	}

	if ageFilter != "" && (mode == "all" || mode == "age") {
		requestedAge, err := strconv.Atoi(ageFilter)
		if err == nil {
			query = query.Where("age BETWEEN ? AND ?", requestedAge-3, requestedAge+3)
		}
	}
	if genderFilter != "" && (mode == "all" || mode == "gender") {
		query = query.Where("LOWER(gender) = ?", strings.ToLower(genderFilter))
	}
	if hobbiesFilter != "" && (mode == "all" || mode == "hobbies") {
		query = query.Where("LOWER(interests) ILIKE ?", "%"+strings.ToLower(hobbiesFilter)+"%")
	}

	var potentialMatches []models.User
	if err := query.Find(&potentialMatches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching potential matches"})
		return
	}

	recommendations := []Recommendation{}
	for _, match := range potentialMatches {
		if currentUser.Latitude != 0 && currentUser.Longitude != 0 &&
			match.Latitude != 0 && match.Longitude != 0 {

			dist := distanceKm(currentUser.Latitude, currentUser.Longitude, match.Latitude, match.Longitude)
			log.Printf("Distance between %s and %s: %.2f km\n", currentUser.Name, match.Name, dist)

			if dist > currentUser.PreferredRadius {
				continue
			}
		} else {
			if mode == "all" || mode == "location" {
				if locationFilter == "" {
					if !strings.Contains(strings.ToLower(match.Location), strings.ToLower(currentUser.Location)) {
						log.Printf("Skipping %s, location %s does not match %s", match.Name, match.Location, currentUser.Location)
						continue
					}
				}
			}
		}

		score := calculateMatchScore(currentUser, match)
		recommendations = append(recommendations, Recommendation{User: match, Score: score})
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	resp := []gin.H{}
	recommendedIDs := make([]uint, 0, len(recommendations))

	for _, rec := range recommendations {
		recommendedIDs = append(recommendedIDs, rec.User.ID)
		resp = append(resp, gin.H{
			"id":    rec.User.ID,
			"score": rec.Score,
		})
	}

	recentRecommendations[currentUser.ID] = recommendedIDs

	c.JSON(http.StatusOK, gin.H{"recommendations": resp})
}

func calculateMatchScore(user1, user2 models.User) int {
	score := 0

	// Interest Score
	interestScore := calculateInterestScore(user1.Interests, user2.Interests)
	log.Printf("Interest Score: %d", interestScore)
	score += interestScore

	// Location Match Score
	if user1.Location == user2.Location {
		log.Printf("Location Match Score: %d", 30)
		score += 30
	}

	// Age Score
	ageDifference := abs(user1.Age - user2.Age)
	var ageScore int
	if ageDifference <= 2 {
		ageScore = 20
	} else if ageDifference <= 10 {
		ageScore = 20 - (ageDifference * 2)
	} else {
		ageScore = 0
	}
	log.Printf("Age Score: %d", ageScore)
	score += ageScore

	// Profile Completeness Score
	if user2.ProfilePicture != "" && user2.Interests != "" {
		log.Printf("Profile Completeness Score: %d", 10)
		score += 10
	}

	// Gender Preference Matching
	genderMatchScore := calculateGenderPreferenceMatch(user1, user2)
	log.Printf("Gender Preference Match Score: %d", genderMatchScore)
	score += genderMatchScore

	log.Printf("Total Score: %d", score)
	return score
}

func calculateGenderPreferenceMatch(user1, user2 models.User) int {
	score := 0

	// Normalize gender and lookingFor to lowercase for case-insensitive comparison
	user1LookingFor := strings.ToLower(user1.LookingFor)
	user2LookingFor := strings.ToLower(user2.LookingFor)
	user1Gender := strings.ToLower(user1.Gender)
	user2Gender := strings.ToLower(user2.Gender)

	// User1's LookingFor matches User2's Gender
	if matchesPreference(user1LookingFor, user2Gender) {
		score += 20
		log.Printf("User1's lookingFor (%s) matches User2's gender (%s). +20 points", user1.LookingFor, user2.Gender)
	}

	// User2's LookingFor matches User1's Gender
	if matchesPreference(user2LookingFor, user1Gender) {
		score += 20
		log.Printf("User2's lookingFor (%s) matches User1's gender (%s). +20 points", user2.LookingFor, user1.Gender)
	}

	return score
}

func matchesPreference(lookingFor, gender string) bool {
	if lookingFor == "any" {
		return true
	}
	return lookingFor == gender
}

func distanceKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0
	dLat := deg2rad(lat2 - lat1)
	dLon := deg2rad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(lat1))*math.Cos(deg2rad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	dist := earthRadius * c
	return dist
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180
}

func calculateInterestScore(interests1, interests2 string) int {
	interests1List := parseInterests(interests1)
	interests2List := parseInterests(interests2)
	common := intersect(interests1List, interests2List)
	log.Printf("Interests1: %v, Interests2: %v, Common: %v", interests1List, interests2List, common)
	return len(common) * 10 // Assign 10 points per common interest
}

func parseInterests(interests string) []string {
	if interests == "" {
		return []string{}
	}
	interestList := strings.Split(interests, ",")
	parsed := []string{}
	for _, interest := range interestList {
		trimmed := strings.TrimSpace(strings.ToLower(interest))
		if trimmed != "" {
			parsed = append(parsed, trimmed)
		}
	}
	return parsed
}

func intersect(slice1, slice2 []string) []string {
	set := make(map[string]bool)
	for _, v := range slice1 {
		set[v] = true
	}
	var intersection []string
	for _, v := range slice2 {
		if set[v] {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func DeclineRecommendation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var input struct {
		RequestID uint `json:"request_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	declinedUser := models.DeclinedUser{
		UserID:         userID.(uint),
		DeclinedUserID: input.RequestID,
	}
	if err := config.DB.Create(&declinedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decline user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User declined successfully"})
}
