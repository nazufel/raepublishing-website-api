package models

import "gopkg.in/mgo.v2/bson"

/* Not sure if i'm going to even need this with oAuth. Kind of forgot i was going
+to do that. However, this is good practice for CRUD.
*/
//TODO: Encrypt user passwords
//TODO: Add user profile pictures if oAuth doesn't pull them in from Facebook
//TODO: Add user Bio.

//Users struct to hold user data
type Users struct {
	// Users represents the structure of the resource, using bson to store in mongo
	// ObjectId uses mongo's id service to assign a user id
	ID        bson.ObjectId `json:"id" bson:"_id"`
	FirstName string        `json:"firstname" bson:"firstname"`
	LastName  string        `json:"lastname" bson:"lastname"`
	Username  string        `json:"username" bson:"username"`
	Email     string        `json:"email" bson:"email"`
	Role      string        `json:"role" bson:"role"`
	Bio       string        `json:"bio" bson:"bio"`
}

/*
Options for roles are as follows:
Admin - Admin of posts and users
Editor - Can write, edit, publish and remove articles for all users
Contributor - Can write, edit, publish, and remove their own articles
Writer - Can write and edit their own articles
Reader - Can comment on articles
*/
