package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"literary-lions/database"
)

var mu sync.Mutex

func AddToFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AddToFavoritesHandler called")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println("Error decoding request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Println("Book received:", book)

	userID, err := GetUserIDFromSession(database.DB, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	_, err = database.DB.Exec(`INSERT INTO favourite_books (user_id, title, author, cover, description, category, published_date) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID, book.Title, book.Author, book.Cover, book.Description, book.Category, book.PublishedDate)
	if err != nil {
		log.Println("Error inserting book into database:", err)
		http.Error(w, "Failed to add book to favorites", http.StatusInternalServerError)
		return
	}

	log.Println("Book added to favorites successfully")
	w.WriteHeader(http.StatusCreated)
}

func DisplayFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("DisplayFavoritesHandler called")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID, err := GetUserIDFromSession(database.DB, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := database.DB.Query(`SELECT title, author, cover, description, category, published_date FROM favourite_books WHERE user_id = ?`, userID)
	if err != nil {
		log.Println("Error fetching favorite books from database:", err)
		http.Error(w, "Failed to fetch favorite books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Title, &book.Author, &book.Cover, &book.Description, &book.Category, &book.PublishedDate)
		if err != nil {
			log.Println("Error scanning book:", err)
			http.Error(w, "Failed to scan book", http.StatusInternalServerError)
			return
		}
		log.Println("Fetched book:", book)
		books = append(books, book)
	}

	log.Println("Total books fetched:", len(books))

	data := TemplateData{
		Title: "Favourite Books",
		Data:  map[string]interface{}{"Books": books},
	}

	RenderTemplate(w, r, "favourites.html", data)
	log.Println("Template rendered successfully with books:", books)
}

func GetFavouriteBooksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFavouriteBooksHandler called")

	userID, err := GetUserIDFromSession(database.DB, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := database.DB.Query(`SELECT title FROM favourite_books WHERE user_id = ?`, userID)
	if err != nil {
		log.Println("Error fetching favorite books from database:", err)
		http.Error(w, "Failed to fetch favorite books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var favourites []string
	for rows.Next() {
		var title string
		err := rows.Scan(&title)
		if err != nil {
			log.Println("Error scanning book title:", err)
			http.Error(w, "Failed to scan book title", http.StatusInternalServerError)
			return
		}
		favourites = append(favourites, title)
	}

	log.Println("Favourite books fetched:", favourites)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(favourites)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func DeleteFromFavouritesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteFromFavouritesHandler called")
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var book struct {
		Title string `json:"title"`
	}
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println("Error decoding request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Println("Book to delete:", book)

	userID, err := GetUserIDFromSession(database.DB, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	_, err = database.DB.Exec(`DELETE FROM favourite_books WHERE user_id = ? AND title = ?`, userID, book.Title)
	if err != nil {
		log.Println("Error deleting book from database:", err)
		http.Error(w, "Failed to delete book from favorites", http.StatusInternalServerError)
		return
	}

	log.Println("Book deleted from favorites successfully")
	w.WriteHeader(http.StatusNoContent)
}
