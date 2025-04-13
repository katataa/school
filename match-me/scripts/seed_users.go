package main

import (
	"fmt"
	"log"
	"match-me/config"
	"match-me/models"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
)

var locations = []struct {
	Name      string
	Latitude  float64
	Longitude float64
}{
	{"Tallinn", 59.4370, 24.7536},
	{"Tartu", 58.3780, 26.7325},
	{"Narva", 59.3797, 28.1791},
	{"Pärnu", 58.3859, 24.4971},
	{"Kohtla-Järve", 59.3986, 27.2731},
	{"Viljandi", 58.3639, 25.5900},
	{"Rakvere", 59.3464, 26.3553},
	{"Maardu", 59.4767, 24.9774},
	{"Kuressaare", 58.2481, 22.5039},
	{"Võru", 57.8339, 27.0170},
	{"Jõhvi", 59.3592, 27.4133},
	{"Haapsalu", 58.9431, 23.5400},
	{"Paide", 58.8856, 25.5576},
	{"Keila", 59.3031, 24.4153},
	{"Rapla", 58.9956, 24.7972},
}

var genders = []string{"Male", "Female", "Non-binary", "Other", "Prefer not to say"}
var lookingForOptions = []string{"Male", "Female", "Non-binary", "Any", "Other"}

var interestsCategories = map[string][]string{
	"Sports":   {"Football", "Basketball", "Tennis", "Swimming", "Hiking"},
	"Food":     {"Chinese", "Italian", "Mexican", "Indian", "Japanese", "Eating out", "Cooking at home"},
	"Culture":  {"Art", "History", "Museums", "Theater", "Traveling", "Languages", "Religion"},
	"Games":    {"Counter-Strike", "Minecraft", "Sims 4", "League of Legends", "Valorant", "Genshin Impact"},
	"MoviesTV": {"Action", "Comedy", "Drama", "Sci-fi", "Romance", "Documentaries", "Anime", "Horror"},
}

var aboutMeTexts = []string{
	"I love exploring new places and meeting new people!",
	"Big fan of movies and video games.",
	"Passionate about fitness and outdoor adventures.",
	"Foodie at heart, always trying new recipes.",
	"Music is my life! Let's talk about our favorite artists.",
	"Fluent in memes and sarcasm.",
	"Traveling is my therapy. Have been to 10+ countries!",
	"Let’s grab a coffee and talk about books!",
	"Looking for someone who shares my love for dogs!",
	"Aspiring astronaut, just waiting for my spaceship.",
}

func SeedUsers() {
	config.ConnectDatabase()

	if err := config.DB.Exec("DELETE FROM users").Error; err != nil {
		log.Fatalf("Failed to clear users: %v", err)
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		randomLocation := locations[rand.Intn(len(locations))]
		randomGender := genders[rand.Intn(len(genders))]
		randomLookingFor := lookingForOptions[rand.Intn(len(lookingForOptions))]
		randomAge := rand.Intn(40) + 18
		randomProfilePic := fmt.Sprintf("https://i.pravatar.cc/150?img=%d", rand.Intn(70)+1)
		randomAboutMe := aboutMeTexts[rand.Intn(len(aboutMeTexts))]

		userInterests := ""
		categoryKeys := []string{"Sports", "Food", "Culture", "Games", "MoviesTV"}
		selectedCategories := rand.Perm(len(categoryKeys))[:rand.Intn(3)+1]

		for _, idx := range selectedCategories {
			category := categoryKeys[idx]
			interestList := interestsCategories[category]
			userInterests += interestList[rand.Intn(len(interestList))] + ", "
		}
		userInterests = userInterests[:len(userInterests)-2]

		user := models.User{
			Name:           faker.Name(),
			Email:          faker.Email(),
			Password:       "password123",
			Age:            randomAge,
			Gender:         randomGender,
			LookingFor:     randomLookingFor,
			Interests:      userInterests,
			Location:       randomLocation.Name,
			Latitude:       randomLocation.Latitude,
			Longitude:      randomLocation.Longitude,
			Info:           randomAboutMe,
			ProfilePicture: randomProfilePic,
		}

		if err := config.DB.Create(&user).Error; err != nil {
			log.Fatalf("Failed to insert user: %v", err)
		}
	}

	fmt.Println("✅ 100 fake users added with random profile pictures, locations, hobbies, and infos!")
}
