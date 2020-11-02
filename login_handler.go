package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ABuarque/i2m/auth"
	"github.com/ABuarque/i2m/db"
	"github.com/ABuarque/i2m/encryption"
	"github.com/labstack/echo"
)

func loginHandler(c echo.Context) error {
	template := template.Must(template.ParseFiles("templates/login.html"))
	var html bytes.Buffer
	err := template.Execute(&html, nil)
	if err != nil {
		return c.HTML(http.StatusOK, "<h1>Error</h1>")
	}
	return c.HTML(http.StatusOK, string(html.Bytes()))
}

func loginAPIHandler(client *db.Client, authService *auth.Auth) echo.HandlerFunc {
	return func(c echo.Context) error {
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
		authorization, err := authService.GetToken(&auth.TokenClaims{ID: u.ID, Email: u.Email})
		if err != nil {
			return c.HTML(http.StatusOK, "<h1>Error</h1>")
		}
		return c.Redirect(http.StatusFound, fmt.Sprintf("/dashboard?authorization=%s", authorization))
	}
}
