package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func manufacturersHandler(w http.ResponseWriter, r *http.Request) {
	var manufacturers []Manufacturer
	err := fetchAPI("manufacturers", &manufacturers)
	if err != nil {
		http.Error(w, "Failed to fetch manufacturers data", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(manufacturers)
	if err != nil {
		http.Error(w, "Failed to marshal manufacturers data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	var categories []Category
	err := fetchAPI("categories", &categories)
	if err != nil {
		http.Error(w, "Failed to fetch categories data", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, "Failed to marshal categories data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func carModelsHandler(w http.ResponseWriter, r *http.Request) {
	var carModels []CarModel
	err := fetchAPI("carModels", &carModels)
	if err != nil {
		http.Error(w, "Failed to fetch car models data", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(carModels)
	if err != nil {
		http.Error(w, "Failed to marshal car models data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getManufacturerByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid manufacturer ID", http.StatusBadRequest)
		return
	}

	var manufacturer Manufacturer
	err = fetchAPI(fmt.Sprintf("manufacturers/%d", id), &manufacturer)
	if err != nil {
		http.Error(w, "Failed to fetch manufacturer data", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(manufacturer)
	if err != nil {
		http.Error(w, "Failed to marshal manufacturer data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func searchCarModels(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	searchName := query.Get("name")
	searchManufacturer := query.Get("manufacturer")
	searchCategory := query.Get("category")

	var results []CarModel
	err := fetchAPI("carModels", &results)
	if err != nil {
		http.Error(w, "Failed to fetch car models data", http.StatusInternalServerError)
		return
	}

	var filteredResults []CarModel
	for _, model := range results {
		if (searchName == "" || strings.Contains(strings.ToLower(model.Name), strings.ToLower(searchName))) &&
			(searchManufacturer == "" || strconv.Itoa(model.ManufacturerID) == searchManufacturer) &&
			(searchCategory == "" || strconv.Itoa(model.CategoryID) == searchCategory) {
			filteredResults = append(filteredResults, model)
		}
	}

	jsonResponse, err := json.Marshal(filteredResults)
	if err != nil {
		http.Error(w, "Failed to marshal search results", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func carModelDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid car model ID", http.StatusBadRequest)
		return
	}

	var foundModel *CarModel
	for _, model := range data.CarModels {
		if model.ID == id {
			foundModel = &model
			break
		}
	}

	if foundModel == nil {
		http.NotFound(w, r)
		return
	}

	response := map[string]interface{}{
		"carModel": foundModel,
	}

	manufacturer := getManufacturerByID(foundModel.ManufacturerID)
	if manufacturer != nil {
		response["manufacturer"] = manufacturer
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal car model details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func compareCarModelsHandler(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["ids"]
	var results []CarModel

	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid car model ID", http.StatusBadRequest)
			return
		}

		var model CarModel
		err = fetchAPI(fmt.Sprintf("carModels/%d", id), &model)
		if err != nil {
			http.Error(w, "Failed to fetch car model data", http.StatusInternalServerError)
			return
		}

		results = append(results, model)
	}

	if len(results) == 0 {
		http.NotFound(w, r)
		return
	}

	jsonResponse, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Failed to marshal comparison data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func personalizedRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	recommendations := getRecentViewedCars()
	log.Println("Generated recommendations:", recommendations)

	if len(recommendations) == 0 {
		log.Println("No recommendations available")
		http.Error(w, "No recommendations available", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(recommendations)
	if err != nil {
		log.Println("Failed to encode recommendations:", err)
		http.Error(w, "Failed to encode recommendations", http.StatusInternalServerError)
	}
}

func getRecentViewedCars() []CarModel {
	var recentCars []CarModel
	for _, carID := range data.RecommendedIDs {
		for _, car := range data.CarModels {
			if car.ID == carID {
				recentCars = append(recentCars, car)
				break
			}
		}
	}
	log.Println("Recent viewed cars:", recentCars)
	return recentCars
}

func trackInteractionHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	carModelIDStr := r.FormValue("car_model_id")
	carModelID, err := strconv.Atoi(carModelIDStr)
	if err != nil {
		http.Error(w, "Invalid car model ID", http.StatusBadRequest)
		return
	}

	trackUserInteraction(carModelID)
	log.Println("Tracked interaction:", carModelID)
	log.Println("Current user interactions (RecommendedIDs):", data.RecommendedIDs)

	http.Redirect(w, r, "/recommendations.html", http.StatusSeeOther)
}

func trackUserInteraction(carModelID int) {
	for _, id := range data.RecommendedIDs {
		if id == carModelID {
			return
		}
	}
	data.RecommendedIDs = append(data.RecommendedIDs, carModelID)
	if len(data.RecommendedIDs) > 3 {
		data.RecommendedIDs = data.RecommendedIDs[1:]
	}
	log.Println("Updated RecommendedIDs:", data.RecommendedIDs)
}

func getManufacturerName(id int) string {
	var manufacturer Manufacturer
	err := fetchAPI(fmt.Sprintf("manufacturers/%d", id), &manufacturer)
	if err != nil {
		return ""
	}
	return manufacturer.Name
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	results := searchDatabase(query)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func searchDatabase(query string) []CarModel {
	var results []CarModel
	lowerQuery := strings.ToLower(query)
	var carModels []CarModel
	err := fetchAPI("carModels", &carModels)
	if err != nil {
		log.Printf("Error fetching car models: %v", err)
		return results
	}

	for _, car := range carModels {
		if strings.Contains(strings.ToLower(car.Name), lowerQuery) ||
			strings.Contains(strings.ToLower(getManufacturerName(car.ManufacturerID)), lowerQuery) ||
			strings.Contains(fmt.Sprint(car.Year), lowerQuery) {
			results = append(results, car)
		}
	}
	return results
}

func ErrorHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Internal Server Error: %v", err)
				http.Error(w, "Internal Server Error. Please try again later.", http.StatusInternalServerError)
			}
		}()
		f(w, r)
	}
}
