package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	Username       string
	Email          string
	FirstName      string
	LastName       string
	Age            sql.NullInt64
	Gender         string
	ProfilePicture string
}

func GetUserIDFromSession(db *sql.DB, r *http.Request) (int, error) {
	session, err := r.Cookie("session_token")
	if err != nil {
		return 0, err
	}

	var userID int
	err = db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", session.Value).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func UserProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "Username not specified", http.StatusBadRequest)
			return
		}

		var user User
		err := db.QueryRow("SELECT username, email, first_name, last_name, age, gender, profile_picture FROM users WHERE username = ?", username).Scan(&user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Age, &user.Gender, &user.ProfilePicture)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				log.Println("Error retrieving user profile:", err)
				http.Error(w, "Server error, unable to retrieve user profile", http.StatusInternalServerError)
			}
			return
		}

		log.Printf("Retrieved user profile: %+v", user)

		ageStr := "Not specified"
		if user.Age.Valid {
			ageStr = strconv.Itoa(int(user.Age.Int64))
		}

		data := TemplateData{
			Title:          user.Username + "'s Profile",
			Username:       user.Username,
			Email:          user.Email,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Age:            ageStr,
			Gender:         user.Gender,
			ProfilePicture: user.ProfilePicture,
			LoggedIn:       false, // Update this based on session
		}

		RenderTemplate(w, r, "user_profile.html", data)
	}
}
