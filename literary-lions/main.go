package main

import (
	"literary-lions/database"
	"literary-lions/handlers"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := database.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}
	defer database.Close()

	handlers.InitializeTemplates()

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/register", handlers.RegisterHandler(database.DB))
	http.HandleFunc("/login", handlers.LoginHandler(database.DB))
	http.HandleFunc("/profile", handlers.ProfileHandler(database.DB))
	http.HandleFunc("/categories", handlers.CategoriesHandler)
	http.HandleFunc("/search-books", handlers.SearchBooksHandler)
	http.HandleFunc("/discussion", handlers.DisplayPostsHandler(database.DB))
	http.HandleFunc("/create_post", handlers.CreatePostHandler(database.DB))
	http.HandleFunc("/create_comment", handlers.CreateCommentHandler(database.DB))
	http.HandleFunc("/like_dislike", handlers.LikeDislikeHandler(database.DB))
	http.HandleFunc("/edit_profile", handlers.EditProfileHandler(database.DB))
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/add-to-favorites", handlers.AddToFavoritesHandler)
	http.HandleFunc("/favourites", handlers.DisplayFavoritesHandler)
	http.HandleFunc("/get-favourite-books", handlers.GetFavouriteBooksHandler)
	http.HandleFunc("/delete-from-favourites", handlers.DeleteFromFavouritesHandler)
	http.HandleFunc("/my_posts", handlers.MyPostsHandler(database.DB))
	http.HandleFunc("/liked_posts", handlers.LikedPostsHandler(database.DB))
	http.HandleFunc("/user_profile", handlers.UserProfileHandler(database.DB)) // Ensure this is only registered once

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/login_prompt", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Rendering login_prompt.html")
		handlers.RenderTemplate(w, r, "login_prompt.html", nil)
	})
	http.HandleFunc("/login_email", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Rendering login_email.html")
		handlers.RenderTemplate(w, r, "login_email.html", nil)
	})
	http.HandleFunc("/login_password", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Rendering login_password.html")
		handlers.RenderTemplate(w, r, "login_password.html", nil)
	})

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
