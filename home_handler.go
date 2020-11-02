package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

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
