package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/ABuarque/i2m/db"
	"github.com/labstack/echo"
)

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
