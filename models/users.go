package models

import "gopkg.in/mgo.v2/bson"

//Users struct to hold user data
type Users struct {
	// User represents the structure of the resource, using bson to store in mongo
	// ObjectId uses mongo's id service to assign a user id
	ID     bson.ObjectId `json:"id" bson:"_id"`
	Name   string        `json:"name" bson:"name"`
	Gender string        `json:"gender" bson:"gender"`
	Age    int           `json:"age" bson:"age"`
}
