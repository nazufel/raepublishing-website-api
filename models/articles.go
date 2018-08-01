package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Note: I have not implimented the controller for this yet. Just mocking what the data strcuture could look like.

//Articles struct to hold article data
type Articles struct {
	// Articles represents the structure of the resource, using bson to store in mongo
	// ObjectId uses mongo's id service to assign a user id
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Title     string        `json:"tile" bson:"title"`
	Subtitle  string        `json:"subtitle" bson:"subtitle"`
	Author    *Users        `json:"author" bson:"author"`
	Editor    *Users        `json:"editor" bson:"editor"`
	Published time.Time     `json:"published" bson:"published"`
	Updated   time.Time     `json:"updated" bson:"updated"`
	Body      string        `json:"body" bson:"body"`
	// TODO: be able add multiple categories to one post
	//Categories []bytes       `json:"categories" bson:"categories"`
	// TODO: be able to add multiple comments from people who are logged in from oAuth.
	// TODO: impliment oAuth, store users in DB, and be able reference user documents from article/comment documents,
}
