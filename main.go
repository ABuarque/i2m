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

var client *db.Client

var a *auth.Auth

var t *twitter.Client

func main() {
	authSecret := os.Getenv("SECRET")
	if authSecret == "" {
		log.Fatal("missing SECRET environment variable")
	}
	a = auth.New(authSecret)
	dbURL := os.Getenv("MONGODB_URI")
	if dbURL == "" {
		log.Fatal("missing MONGODB_URI environment variable")
	}
	dbName := os.Getenv("MONGO_NAME")
	if dbName == "" {
		log.Fatal("missing MONGO_NAME environment variable")
	}
	c, err := db.New(dbURL, dbName)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect with data base: %q", err))
	}
	log.Println("connected to data base!")
	client = c
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
	t = twitter.NewClient(apiKey, apiSecret, accessToken, accessTokenSecret)
	e := echo.New()
	e.Static("/static", "templates/assets")
	e.GET("/", homeHandler)
	e.GET("/login", loginHandler)
	e.POST("/login", loginAPIHandler)
	e.GET("/dashboard", dashboardHandler)
	e.POST("/new_post", createPostPage)
	e.POST("/create_post", createPostHandler)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing PORT environment variable")
	}
	log.Fatal(e.Start(":" + port))
}
