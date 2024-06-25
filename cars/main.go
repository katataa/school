package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	loadData()
	setupStaticFileServing()
	setupRouteHandlers()

	port := ":8081"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func setupStaticFileServing() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func setupRouteHandlers() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			log.Println("Serving index.html")
			http.ServeFile(w, r, "./static/index.html")
		} else if r.URL.Path == "/details.html" {
			log.Println("Serving details.html")
			http.ServeFile(w, r, "./static/details.html")
		} else if strings.HasPrefix(r.URL.Path, "/static/") || strings.HasPrefix(r.URL.Path, "/img/") || strings.HasPrefix(r.URL.Path, "/muudpildid/") {
			log.Println("Serving static/image file:", r.URL.Path)
			http.StripPrefix("/", http.FileServer(http.Dir("."))).ServeHTTP(w, r)
		} else {
			log.Println("File not found:", r.URL.Path)
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/carModels", carModelsHandler)
	http.HandleFunc("/carModelDetail", carModelDetailHandler)
	http.HandleFunc("/compareCarModels", compareCarModelsHandler)
	http.HandleFunc("/searchCarModels", searchCarModels)
	http.HandleFunc("/search", searchHandler)

	http.HandleFunc("/likeCar", likeCarHandler)
	http.HandleFunc("/likedCars", likedCarsHandler)
	http.HandleFunc("/track-interaction", trackInteractionHandler)
	http.HandleFunc("/recommendations", personalizedRecommendationsHandler)

	http.HandleFunc("/manufacturers", manufacturersHandler)
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/manufacturer", getManufacturerByIDHandler)

	http.HandleFunc("/recommendations.html", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving recommendations.html")
		http.ServeFile(w, r, "./static/recommendations.html")
	})
}
