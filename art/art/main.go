package main

import (
	"art/decode"
	"art/encode"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func main() {
	setupServer()
}

func setupServer() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			action := r.FormValue("action")
			inputText := r.FormValue("inputText")
			var result string
			var processErr error

			switch action {

			case "encode":
				result := encode.EncodeArt(inputText)
				// Replace [n] with <br> for HTML display
				result = strings.ReplaceAll(result, "[n]", "<br>")
				// Directly render the template with the result for the "encode" action
				tmpl.Execute(w, map[string]interface{}{"Result": result})
				return

			case "decode":
				result, processErr = decode.DecodeArt(inputText)
				if processErr != nil {
					if errors.Is(processErr, decode.ErrMalformedInput) {
						http.Error(w, "Malformed input", http.StatusBadRequest)
					} else {
						http.Error(w, "Processing error", http.StatusInternalServerError)
					}
					return
				}
				// For the "decode" action, render the template with the result as well
				tmpl.Execute(w, map[string]interface{}{"Result": result})
				return

			default:
				http.Error(w, "Invalid action", http.StatusBadRequest)
				return
			}
		} else {
			// For GET requests to "/", just serve the main page without any result
			tmpl.Execute(w, nil)
		}
	})

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
