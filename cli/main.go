package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ABuarque/i2m/db"
	"github.com/ABuarque/i2m/encryption"
)

func main() {
	name := flag.String("name", "", "user name")
	email := flag.String("email", "", "user email")
	password := flag.String("password", "", "user password")
	dbURL := flag.String("dbURL", "", "URL connection to DB")
	dbName := flag.String("dbName", "", "DB name")
	flag.Parse()
	if *name == "" {
		log.Fatal("inform user name")
	}
	if *email == "" {
		log.Fatal("inform user email")
	}
	if *password == "" {
		log.Fatal("inform user password")
	}
	if *dbName == "" {
		log.Fatal("inform DB name")
	}
	if *dbURL == "" {
		log.Fatal("inform db URL")
	}
	dbService, err := db.New(*dbURL, *dbName)
	if err != nil {
		log.Fatalf("failed to connect to DB, error %v", err)
	}
	if err := addNewUser(*name, *email, *password, dbService); err != nil {
		log.Fatalf("failed to add new user, %v\n", err)
	}
}

func addNewUser(name, email, password string, dbService *db.Client) error {
	encryptedPassword, err := encryption.Encrypt(password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password, error %v", err)
	}
	user := db.User{
		Name:     name,
		Email:    email,
		Password: encryptedPassword,
		CreatedAt: time.Now(),
	}
	if _, err := dbService.SaveUser(&user); err != nil {
		return fmt.Errorf("failed to save new user on DB, error %v", err)
	}
	log.Printf("added user %s\n", user.Name)
	return nil
}
