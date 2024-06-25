package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Specifications struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsepower"`
	Transmission string `json:"transmission"`
	Drivetrain   string `json:"drivetrain"`
}

type CarModel struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	ManufacturerID int            `json:"manufacturerId"`
	CategoryID     int            `json:"categoryId"`
	Year           int            `json:"year"`
	Specifications Specifications `json:"specifications"`
	Image          string         `json:"image"`
}

type Manufacturer struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Data struct {
	Manufacturers  []Manufacturer `json:"manufacturers"`
	Categories     []Category     `json:"categories"`
	CarModels      []CarModel     `json:"carModels"`
	RecommendedIDs []int          `json:"recommendedIDs"`
}

// global variable to hold the data
var data Data

const apiBaseURL = "http://localhost:8080/api"

func fetchAPI(endpoint string, target interface{}) error {
	url := fmt.Sprintf("%s/%s", apiBaseURL, endpoint)
	log.Printf("Fetching API URL: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch data from API: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func loadData() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		err := fetchAPI("manufacturers", &data.Manufacturers)
		if err != nil {
			log.Fatalf("Error fetching manufacturers data: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := fetchAPI("categories", &data.Categories)
		if err != nil {
			log.Fatalf("Error fetching categories data: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := fetchAPI("carModels", &data.CarModels)
		if err != nil {
			log.Fatalf("Error fetching car models data: %v", err)
		}
	}()

	wg.Wait()
}

func getManufacturerByID(id int) *Manufacturer {
	for _, manufacturer := range data.Manufacturers {
		if manufacturer.ID == id {
			return &manufacturer
		}
	}
	return nil
}
