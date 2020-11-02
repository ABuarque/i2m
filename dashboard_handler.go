package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/ABuarque/i2m/auth"
	"github.com/ABuarque/i2m/db"
	"github.com/labstack/echo"
)

func dashboardHandler(client *db.Client, authService *auth.Auth) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := c.QueryParam("authorization")
		if authorization == "" {
			return c.JSON(http.StatusForbidden, "Acesso negado!")
		}
		ok, err := authService.IsValid(authorization)
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
}
