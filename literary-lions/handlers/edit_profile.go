package handlers

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func EditProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("EditProfileHandler called")

		userID, err := GetUserIDFromSession(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Method == "POST" {
			firstName := r.FormValue("first_name")
			lastName := r.FormValue("last_name")
			username := r.FormValue("username")
			age, err := strconv.Atoi(r.FormValue("age"))
			if err != nil {
				http.Error(w, "Invalid age", http.StatusBadRequest)
				return
			}
			gender := r.FormValue("gender")

			file, _, err := r.FormFile("cover")
			var profilePicturePath string
			if err == nil {
				defer file.Close()
				profilePicturePath = "./static/uploads/" + username + ".jpg"

				if err := os.MkdirAll(filepath.Dir(profilePicturePath), os.ModePerm); err != nil {
					log.Println("Error creating directory:", err)
					http.Error(w, "Server error, unable to create directory for profile picture.", http.StatusInternalServerError)
					return
				}

				f, err := os.Create(profilePicturePath)
				if err != nil {
					log.Println("Error creating file:", err)
					http.Error(w, "Server error, unable to save profile picture.", http.StatusInternalServerError)
					return
				}
				defer f.Close()

				_, err = io.Copy(f, file)
				if err != nil {
					log.Println("Error saving profile picture:", err)
					http.Error(w, "Server error, unable to save profile picture.", http.StatusInternalServerError)
					return
				}
			} else {
				err = db.QueryRow("SELECT profile_picture FROM users WHERE id = ?", userID).Scan(&profilePicturePath)
				if err != nil {
					log.Println("Error retrieving existing profile picture:", err)
					http.Error(w, "Server error, unable to retrieve existing profile picture.", http.StatusInternalServerError)
					return
				}
			}

			stmt, err := db.Prepare("UPDATE users SET first_name = ?, last_name = ?, username = ?, age = ?, gender = ?, profile_picture = ? WHERE id = ?")
			if err != nil {
				log.Println("Error preparing statement:", err)
				http.Error(w, "Server error, unable to update your profile.", http.StatusInternalServerError)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(firstName, lastName, username, age, gender, profilePicturePath, userID)
			if err != nil {
				log.Println("Error executing statement:", err)
				http.Error(w, "Server error, unable to update your profile.", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}

		var username, email, firstName, lastName, gender, profilePicture sql.NullString
		var age sql.NullInt64
		err = db.QueryRow("SELECT username, email, first_name, last_name, gender, age, profile_picture FROM users WHERE id = ?", userID).Scan(&username, &email, &firstName, &lastName, &gender, &age, &profilePicture)
		if err != nil {
			log.Println("Error retrieving user information from users table:", err)
			http.Error(w, "Server error, unable to retrieve profile.", http.StatusInternalServerError)
			return
		}

		log.Println("User information retrieved successfully")

		ageStr := "Not specified"
		if age.Valid {
			ageStr = strconv.Itoa(int(age.Int64))
		}

		data := TemplateData{
			Username:       username.String,
			Email:          email.String,
			FirstName:      firstName.String,
			LastName:       lastName.String,
			Gender:         gender.String,
			Age:            ageStr,
			ProfilePicture: profilePicture.String + "?t=" + time.Now().Format("20060102150405"), // Add timestamp to URL
		}

		log.Println("Rendering edit_profile.html template")
		RenderTemplate(w, r, "edit_profile.html", data)
	}
}
