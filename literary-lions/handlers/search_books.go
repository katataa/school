package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"literary-lions/database"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Book struct {
	Title         string `json:"title"`
	Author        string `json:"author"`
	Cover         string `json:"cover"`
	Description   string `json:"description"`
	Category      string `json:"category"`
	PublishedDate string `json:"publishedDate"`
	IsFavorite    bool   `json:"isFavorite"`
}

func SearchBooksHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	authorFilter := r.URL.Query().Get("author")
	categoryFilter := r.URL.Query().Get("category")
	publishedDateFilter := r.URL.Query().Get("publishedDate")

	if query == "" && authorFilter == "" && categoryFilter == "" && publishedDateFilter == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	googleBooksAPI := "https://www.googleapis.com/books/v1/volumes"
	searchURL := fmt.Sprintf("%s?q=%s", googleBooksAPI, url.QueryEscape(query))

	if query != "" {
		searchURL += fmt.Sprintf("+intitle:%s", url.QueryEscape(query))
	}
	if authorFilter != "" {
		searchURL += fmt.Sprintf("+inauthor:%s", url.QueryEscape(authorFilter))
	}
	if categoryFilter != "" {
		searchURL += fmt.Sprintf("+subject:%s", url.QueryEscape(categoryFilter))
	}

	resp, err := http.Get(searchURL)
	if err != nil {
		http.Error(w, "Error fetching data from Google Books API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Error response from Google Books API: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		http.Error(w, "Error decoding response from Google Books API", http.StatusInternalServerError)
		return
	}

	items, ok := result["items"].([]interface{})
	if !ok {
		http.Error(w, "No books found", http.StatusNotFound)
		return
	}

	query = strings.ToLower(query)
	authorFilter = strings.ToLower(authorFilter)
	categoryFilter = strings.ToLower(categoryFilter)
	publishedDateFilter = strings.ToLower(publishedDateFilter)

	var books []Book
	var userID int
	userLoggedIn := false

	userID, err = GetUserIDFromSession(database.DB, r)
	if err == nil {
		userLoggedIn = true
	}

	for _, item := range items {
		volumeInfo := item.(map[string]interface{})["volumeInfo"].(map[string]interface{})
		title := getString(volumeInfo["title"])
		author := getFirstAuthor(volumeInfo["authors"])
		category := getFirstCategory(volumeInfo["categories"])
		publishedDate := getString(volumeInfo["publishedDate"])
		book := Book{
			Title:         title,
			Author:        author,
			Cover:         getThumbnail(volumeInfo["imageLinks"]),
			Description:   getString(volumeInfo["description"]),
			Category:      category,
			PublishedDate: publishedDate,
		}

		if userLoggedIn {
			err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM favourite_books WHERE user_id = ? AND title = ? AND author = ? AND published_date = ?)", userID, title, author, publishedDate).Scan(&book.IsFavorite)
			if err != nil && err != sql.ErrNoRows {
				log.Println("Error checking favorite status:", err)
				http.Error(w, "Error checking favorite status", http.StatusInternalServerError)
				return
			}
		}

		// ensure the book matches all provided filters (case-insensitive and partial matching)
		if (query == "" || strings.Contains(strings.ToLower(book.Title), query)) &&
			(authorFilter == "" || strings.Contains(strings.ToLower(book.Author), authorFilter)) &&
			(categoryFilter == "" || strings.Contains(strings.ToLower(book.Category), categoryFilter)) &&
			(publishedDateFilter == "" || strings.HasPrefix(strings.ToLower(book.PublishedDate), publishedDateFilter)) {
			books = append(books, book)
		}
	}

	if len(books) == 0 {
		http.Error(w, "No books found matching the criteria", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"books": books,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func getFirstAuthor(authors interface{}) string {
	if authors == nil {
		return "Unknown"
	}
	authorList, ok := authors.([]interface{})
	if !ok || len(authorList) == 0 {
		return "Unknown"
	}
	return authorList[0].(string)
}

func getFirstCategory(categories interface{}) string {
	if categories == nil {
		return "Unknown"
	}
	categoryList, ok := categories.([]interface{})
	if !ok || len(categoryList) == 0 {
		return "Unknown"
	}
	return categoryList[0].(string)
}

func getThumbnail(imageLinks interface{}) string {
	if imageLinks == nil {
		return ""
	}
	imageMap, ok := imageLinks.(map[string]interface{})
	if !ok {
		return ""
	}
	if thumbnail, ok := imageMap["thumbnail"].(string); ok {
		return thumbnail
	}
	return ""
}

func getString(value interface{}) string {
	if value == nil {
		return ""
	}
	str, ok := value.(string)
	if !ok {
		return ""
	}
	return str
}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "categories.html", nil)
}
