package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	LikedCars []int
}

var users = make(map[string]*User)
var usersMutex sync.Mutex

func likeCarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	carModelIDStr := r.URL.Query().Get("car_model_id")
	carModelID, err := strconv.Atoi(carModelIDStr)
	if err != nil {
		http.Error(w, "Invalid car model ID", http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	user, exists := users[userID]
	if !exists {
		user = &User{}
		users[userID] = user
	}

	log.Printf("Received like/unlike request for user ID: %s, car ID: %d", userID, carModelID)

	isLiked := false
	for _, id := range user.LikedCars {
		if id == carModelID {
			isLiked = true
			break
		}
	}

	if isLiked {
		// unlike the car
		for i, id := range user.LikedCars {
			if id == carModelID {
				user.LikedCars = append(user.LikedCars[:i], user.LikedCars[i+1:]...)
				break
			}
		}
		log.Printf("Unliked car ID: %d for user ID: %s", carModelID, userID)
	} else {
		// like the car
		user.LikedCars = append(user.LikedCars, carModelID)
		log.Printf("Liked car ID: %d for user ID: %s", carModelID, userID)
	}

	w.WriteHeader(http.StatusOK)
}

func likedCarsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("User-ID")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	user, exists := users[userID]
	usersMutex.Unlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var likedCarModels []CarModel
	for _, carID := range user.LikedCars {
		var car CarModel
		err := fetchAPI(fmt.Sprintf("carModels/%d", carID), &car)
		if err == nil {
			likedCarModels = append(likedCarModels, car)
		}
	}

	jsonResponse, err := json.Marshal(likedCarModels)
	if err != nil {
		http.Error(w, "Failed to marshal liked cars data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
