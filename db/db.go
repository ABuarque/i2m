package db

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	usersCollection = "users"
	postsCollection = "posts"
)

// User is the struct for users collections
type User struct {
	ID                       string    `json:"id" bson:"_id,omitempty"`
	Name                     string    `json:"name" bson:"name"`
	Email                    string    `json:"email" bson:"email"`
	Password                 string    `json:"password" bson:"password"`
	CreatedAt                time.Time `json:"createdAt" bson:"createdAt"`
	PasswordRedefinitionCode string    `json:"passwordRedefinitionCode," bson:"passwordRedefinitionCode,omitempty"`
}

// Post is the struct of posts
type Post struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Title     string    `json:"tile" bson:"title"`
	Date      string    `json:"date" bson:"date"`
	Info      string    `json:"info" bson:"info"`
	Link      string    `json:"link" bson:"link"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

//Client manages all iteractions with mongodb
type Client struct {
	client *mgo.Database
	dbName string
}

//New returns an db connection instance that can be used for CRUD opetations
func New(url, database string) (*Client, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	c := session.DB(database)
	c.C(usersCollection).EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true})
	return &Client{
		client: session.DB(database),
		dbName: database,
	}, nil
}

// FindByEmail finds an user with email
func (db *Client) FindByEmail(email string) (*User, error) {
	var profile User
	err := db.client.C(usersCollection).Find(bson.M{"email": email}).One(&profile)
	if err != nil {
		return nil, fmt.Errorf("email %s not found on database, got error %q", email, err)
	}
	return &profile, nil
}

// SaveUser saves a new user
func (db *Client) SaveUser(user *User) (*User, error) {
	return user, db.client.C(usersCollection).Insert(user)
}

// SavePost saves a new post
func (db *Client) SavePost(post *Post) (*Post, error) {
	return post, db.client.C(postsCollection).Insert(post)
}

// GetPosts retrieve all posts
func (db *Client) GetPosts() ([]Post, error) {
	var posts []Post
	sortBy := []string{"-createdAt"}
	err := db.client.C(postsCollection).Find(bson.M{}).Sort(strings.Join(sortBy, ",")).All(&posts)
	return posts, err
}
