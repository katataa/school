package main

import (
	"flag"
	"log"
	"match-me/config"
	"match-me/models"
	"match-me/routes"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	devMode := flag.Bool("d", false, "Run server in developer mode (enables GraphQL Playground)")
	flag.Parse()

	config.ConnectDatabase()

	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	r := routes.SetupRouter(*devMode)

	log.Println("Server is running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
