package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ABuarque/i2m/auth"
	"github.com/ABuarque/i2m/db"
	"github.com/ABuarque/i2m/twitter"
	"github.com/labstack/echo"
)

func main() {
	authSecret := os.Getenv("SECRET")
	if authSecret == "" {
		log.Fatal("missing SECRET environment variable")
	}
	authService := auth.New(authSecret)
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("missing MONGODB_URI environment variable")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("missing MONGO_NAME environment variable")
	}
	client, err := db.New(dbURL, dbName)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect with data base: %q", err))
	}
	log.Println("connected to data base!")
	apiKey := os.Getenv("TWITTER_API_KEY")
	if apiKey == "" {
		log.Fatal("missing TWITTER_API_KEY environment variable")
	}
	apiSecret := os.Getenv("TWITTER_API_SECRET")
	if apiSecret == "" {
		log.Fatal("missing TWITTER_API_SECRET environment variable")
	}
	accessToken := os.Getenv("TWITTER_ACESS_TOKEN")
	if accessToken == "" {
		log.Fatal("missing TWITTER_ACESS_TOKEN environment variable")
	}
	accessTokenSecret := os.Getenv("TWITTER_ACESS_TOKEN_SECRET")
	if accessToken == "" {
		log.Fatal("missing TWITTER_ACESS_TOKEN_SECRET environment variable")
	}
	twitterService := twitter.NewClient(apiKey, apiSecret, accessToken, accessTokenSecret)
	e := echo.New()
	e.Static("/static", "templates/assets")
	e.GET("/", homeHandler(client))
	e.GET("/login", loginHandler)
	e.POST("/login", loginAPIHandler(client, authService))
	e.GET("/dashboard", dashboardHandler(client, authService))
	e.POST("/new_post", createPostPage(authService))
	e.POST("/create_post", createPostHandler(client, authService, twitterService))
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing PORT environment variable")
	}
	log.Fatal(e.Start(":" + port))
}
