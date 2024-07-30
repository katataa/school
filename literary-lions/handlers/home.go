package handlers

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	popularBooks := []PopularBook{
		{"Dune", "Frank Herbert", "Science Fiction", "/static/img/dune.jpg"},
		{"To Kill a Mockingbird", "Harper Lee", "Classic Literature", "/static/img/mockingbird.jpg"},
		{"1984", "George Orwell", "Dystopian", "/static/img/1984.jpg"},
		{"Pride and Prejudice", "Jane Austen", "Romance", "/static/img/prideandprejudice.jpg"},
		{"The Great Gatsby", "F. Scott Fitzgerald", "Classic Literature", "/static/img/gatsby.jpg"},
		{"The Catcher in the Rye", "J.D. Salinger", "Classic Literature", "/static/img/catherrye.jpg"},
		{"Harry Potter and the Sorcerer's Stone", "J.K. Rowling", "Fantasy", "/static/img/harrypotter.jpg"},
		{"The Hobbit", "J.R.R. Tolkien", "Fantasy", "/static/img/thehobbit.jpg"},
		{"The Handmaid's Tale", "Margaret Atwood", "Dystopian", "/static/img/handmaid.jpg"},
		{"The Da Vinci Code", "Dan Brown", "Thriller", "/static/img/davinki.jpg"},
	}

	categories := []Category{
		{"Romance"},
		{"Science Fiction"},
		{"Fantasy"},
		{"Mystery"},
		{"Horror"},
		{"Non-Fiction"},
	}

	cookie, err := r.Cookie("session_token")
	loggedIn := err == nil && cookie != nil

	if !loggedIn {
		RenderTemplate(w, r, "index.html", TemplateData{})
		return
	}

	data := TemplateData{
		Title:        "Literary Lions Forum",
		LoggedIn:     loggedIn,
		PopularBooks: popularBooks,
		Categories:   categories,
	}

	RenderTemplate(w, r, "home.html", data)
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			email := r.FormValue("email")
			password := r.FormValue("password")
			firstName := r.FormValue("first_name")
			lastName := r.FormValue("last_name")
			age := r.FormValue("age")
			gender := r.FormValue("gender")

			// check for existing username or email
			var existingID int
			err := db.QueryRow("SELECT id FROM users WHERE username = ? OR email = ?", username, email).Scan(&existingID)
			if err == nil {
				// user with this username or email already exists
				RenderTemplate(w, r, "register.html", TemplateData{
					Data: map[string]string{
						"ErrorMessage": "Username or email already exists. Please choose a different one.",
					},
				})
				return
			} else if err != sql.ErrNoRows {
				http.Error(w, "Server error, unable to create your account.", http.StatusInternalServerError)
				return
			}

			file, _, err := r.FormFile("cover")
			var profilePicturePath string
			if err == nil {
				defer file.Close()
				profilePicturePath = "./static/uploads/" + username + ".jpg"
				f, err := os.Create(profilePicturePath)
				if err != nil {
					http.Error(w, "Server error, unable to save profile picture.", http.StatusInternalServerError)
					return
				}
				defer f.Close()
				_, err = io.Copy(f, file)
				if err != nil {
					http.Error(w, "Server error, unable to save profile picture.", http.StatusInternalServerError)
					return
				}
			} else {
				profilePicturePath = "/static/uploads/default-pic.jpg" // default profile picture
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Server error, unable to create your account.", http.StatusInternalServerError)
				return
			}

			stmt, err := db.Prepare("INSERT INTO users(username, email, password, first_name, last_name, age, gender, profile_picture, registration_date) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				http.Error(w, "Server error, unable to create your account.", http.StatusInternalServerError)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(username, email, hashedPassword, firstName, lastName, age, gender, profilePicturePath, time.Now().Format(time.RFC3339))
			if err != nil {
				http.Error(w, "Server error, unable to create your account.", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		RenderTemplate(w, r, "register.html", TemplateData{})
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			email := r.FormValue("email")
			password := r.FormValue("password")

			log.Println("Attempting login - Email:", email, "Password:", password)

			var storedPassword string
			var userID int
			err := db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &storedPassword)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Println("No user found with email:", email)
					http.Redirect(w, r, "/login_email", http.StatusSeeOther)
				} else {
					log.Println("Error retrieving user:", err)
					http.Error(w, "Server error, unable to login.", http.StatusInternalServerError)
				}
				return
			}

			log.Println("User found - ID:", userID, "Stored Password:", storedPassword)

			err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
			if err != nil {
				log.Println("Password mismatch for user ID:", userID)
				http.Redirect(w, r, "/login_password", http.StatusSeeOther)
				return
			}

			sessionToken := uuid.New().String()
			_, err = db.Exec("INSERT INTO sessions (user_id, token, created_at) VALUES (?, ?, ?)", userID, sessionToken, time.Now().Format(time.RFC3339))
			if err != nil {
				log.Println("Error creating session:", err)
				http.Error(w, "Server error, unable to create session.", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:  "session_token",
				Value: sessionToken,
				Path:  "/",
			})

			log.Println("User logged in successfully:", email)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		RenderTemplate(w, r, "login.html", nil)
	}
}

func ProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromSession(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
			ProfilePicture: profilePicture.String,
		}

		RenderTemplate(w, r, "profile.html", data)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
