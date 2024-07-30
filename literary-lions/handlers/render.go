package handlers

import (
	"log"
	"net/http"
)

type TemplateData struct {
	Title          string
	LoggedIn       bool
	Data           interface{}
	Username       string
	Email          string
	FirstName      string
	LastName       string
	Age            string
	Gender         string
	PopularBooks   []PopularBook
	Categories     []Category
	ProfilePicture string
	Posts          []Post
}

type PopularBook struct {
	Title  string
	Author string
	Genre  string
	Cover  string
}

type Category struct {
	Name string
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	session, _ := r.Cookie("session_token")
	loggedIn := session != nil

	templateData := TemplateData{
		Title:    "Literary Lions Forum",
		LoggedIn: loggedIn,
	}

	if td, ok := data.(TemplateData); ok {
		templateData = td
		templateData.LoggedIn = loggedIn
	} else {
		templateData.Data = data
	}
	if templateData.Data == nil {
		templateData.Data = make(map[string]interface{})
	}
	templateData.LoggedIn = loggedIn // Ensure LoggedIn status is set

	err := templates.ExecuteTemplate(w, tmpl, templateData)
	if err != nil {
		log.Printf("Error rendering template %s: %v", tmpl, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
