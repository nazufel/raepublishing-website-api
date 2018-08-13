package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

/* Not sure if i'm going to even need this with oAuth. Kind of forgot i was going
+to do that. However, this is good practice for CRUD.
*/
//TODO: Add user profile pictures if oAuth doesn't pull them in from Facebook
//TODO: Add user Bio.

//Users struct to hold user data
type Users struct {
	// Users represents the structure of the resource, using bson to store in mongo
	// Undid "omitempty" since the mgo.Change's Change{} already impliments omitempty
	// ObjectId uses mongo's id service to assign a user id
	ID        bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	FirstName string        `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string        `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Username  string        `json:"username" bson:"username,omitempty"`
	Email     string        `json:"email,omitempty" bson:"email,omitempty"`
	Role      string        `json:"role,omitempty" bson:"role,omitempty"`
	Created   time.Time     `json:"created,omitempty" bson:"created,omitempty"`
	Updated   time.Time     `json:"updated,omitempty" bson:"updated,omitempty"`
	Bio       string        `json:"bio,omitempty" bson:"bio,omitempty"`
}

/*
Options for roles are as follows:
Admin - Admin of posts and users
Editor - Can write, edit, publish and remove articles for all users
Contributor - Can write, edit, publish, and remove their own articles
Writer - Can write and edit their own articles
Reader - Can comment on articles
*/
