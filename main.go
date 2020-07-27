package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ABuarque/i2m/auth"
	"github.com/ABuarque/i2m/db"
	"github.com/ABuarque/i2m/encryption"
	"github.com/ABuarque/i2m/twitter"
	"github.com/labstack/echo"
)

var client *db.Client

var a *auth.Auth

var t *twitter.Client

type shortPost struct {
	ID          string
	Title       string
	Date        string
	Description string
	Link        string
}

func homeHandler(c echo.Context) error {
	p, err := client.GetPosts()
	if err != nil {
		log.Println(fmt.Sprintf("failed to retrieve files from db with error %q", err))
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	var shortPosts []shortPost
	for _, post := range p {
		shortPosts = append(shortPosts, shortPost{
			Title:       post.Title,
			Date:        post.Date,
			Description: post.Info,
			Link:        post.Link,
		})
	}
	sps := struct {
		Posts []shortPost
	}{
		shortPosts,
	}
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	var html bytes.Buffer
	err = tmpl.Execute(&html, sps)
	if err != nil {
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	return c.HTML(http.StatusOK, string(html.Bytes()))
}

func loginHandler(c echo.Context) error {
	template := template.Must(template.ParseFiles("templates/login.html"))
	var html bytes.Buffer
	err := template.Execute(&html, nil)
	if err != nil {
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	return c.HTML(http.StatusOK, string(html.Bytes()))
}

func loginAPIHandler(c echo.Context) error {
	r := c.Request()
	email := r.FormValue("email")
	password := r.FormValue("password")
	u, err := client.FindByEmail(email)
	if err != nil {
		log.Println(fmt.Sprintf("not found email %s, got error %q", email, err))
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	match, err := encryption.Check(password, u.Password)
	if err != nil || !match {
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	authorization, err := a.GetToken(&auth.TokenClaims{ID: u.ID, Email: u.Email})
	if err != nil {
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	return c.Redirect(http.StatusFound, fmt.Sprintf("/dashboard?authorization=%s", authorization))
}

func dashboardHandler(c echo.Context) error {
	authorization := c.QueryParam("authorization")
	if authorization == "" {
		return c.JSON(http.StatusForbidden, "Acesso negado!")
	}
	ok, err := a.IsValid(authorization)
	if !ok || err != nil {
		return c.JSON(http.StatusForbidden, "Acesso negado!")
	}
	p, err := client.GetPosts()
	if err != nil {
		log.Println("failed to get posts from db, got ", err)
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	var shortPosts []shortPost
	for _, post := range p {
		shortPosts = append(shortPosts, shortPost{
			Title: post.Title,
			Date:  post.Date,
			ID:    post.ID,
		})
	}
	template := template.Must(template.ParseFiles("templates/dashboard.html"))
	var html bytes.Buffer
	data := struct {
		Authorization string
		Posts         []shortPost
	}{
		authorization,
		shortPosts,
	}
	err = template.Execute(&html, data)
	if err != nil {
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	return c.HTML(http.StatusOK, string(html.Bytes()))
}

func createPostPage(c echo.Context) error {
	authorization := c.QueryParam("authorization")
	if authorization == "" {
		return c.JSON(http.StatusForbidden, "Acesso negado!")
	}
	ok, err := a.IsValid(authorization)
	if !ok || err != nil {
		return c.JSON(http.StatusForbidden, "Acesso negado!")
	}
	template := template.Must(template.ParseFiles("templates/createPost.html"))
	var html bytes.Buffer
	data := struct {
		Authorization string
	}{
		authorization,
	}
	err = template.Execute(&html, data)
	if err != nil {
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	return c.HTML(http.StatusOK, string(html.Bytes()))
}

func createPostHandler(c echo.Context) error {
	authorization := c.QueryParam("authorization")
	if authorization == "" {
		return c.JSON(http.StatusForbidden, "Acesso negado!")
	}
	ok, err := a.IsValid(authorization)
	if !ok || err != nil {
		return c.JSON(http.StatusForbidden, "Acesso negado!")
	}
	r := c.Request()
	title := r.FormValue("title")
	info := r.FormValue("info")
	link := r.FormValue("link")
	post := db.Post{
		Title:     title,
		Info:      info,
		Link:      link,
		Date:      getDate(),
		CreatedAt: time.Now(),
	}
	_, err = client.SavePost(&post)
	if err != nil {
		log.Println(fmt.Sprintf("failed to save post on db, got %q", err))
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	tweet := fmt.Sprintf("checkout my new post: %s", link)
	err = t.Post(tweet)
	if err != nil {
		log.Println(fmt.Sprintf("failed to make tweet, got error %q", err))
	}
	log.Println(fmt.Sprintf("new tweet made: %s ", tweet))
	return c.Redirect(http.StatusFound, fmt.Sprintf("/dashboard?authorization=%s", authorization))
}

func getDate() string {
	year, month, _ := time.Now().Date()
	return fmt.Sprintf("%s, %d", month.String(), year)
}

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
