package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

const (
	timeout = 15 // in seconds
)

//Client manages all iteractions with mongodb
type Client struct {
	client *mongo.Client
	dbName string
}

//New returns an db connection instance that can be used for CRUD opetations
func New(dbURL, dbName string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB at link [%s], error %v", dbURL, err)
	}
	return &Client{
		client: client,
		dbName: dbName,
	}, nil
}

// FindByEmail finds an user with email
func (db *Client) FindByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	var user User
	filter := bson.M{"email": email}
	if err := db.client.Database(db.dbName).Collection(usersCollection).FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to fetch user from DB, error %v", err)
	}
	return &user, nil
}

// SaveUser saves a new user
func (db *Client) SaveUser(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	if _, err := db.client.Database(db.dbName).Collection(usersCollection).InsertOne(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to persist user data into db, error %v", err)
	}
	return user, nil
}

// SavePost saves a new post
func (db *Client) SavePost(post *Post) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	if _, err := db.client.Database(db.dbName).Collection(postsCollection).InsertOne(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to persist post data into db, error %v", err)
	}
	return post, nil
}

// GetPosts retrieve all posts
func (db *Client) GetPosts() ([]Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	var posts []Post
	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", -1}})
	cursor, err := db.client.Database(db.dbName).Collection(postsCollection).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts from db, error %v", err)
	}
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("failed to get posts from db, error %v", err)
	}
	return posts, nil
}
