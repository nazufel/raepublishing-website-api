package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nazufel/raepublishing-website-api/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserController represents the controller for operating on the User resource
type UserController struct {
	// update the controller methods so they can access mongo
	session *mgo.Session
}

/*NewUserConroller represents the controller for updating new User Resources
feeding New UseConroller the *mgo.Session pointer allows it to use mongo */
func NewUserController(s *mgo.Session) *UserController {
	// init mongo
	return &UserController{s}
}

// GET - GetUser retrieves an individual user resource
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Get user id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}
	// Grab id
	oid := bson.ObjectIdHex(id)

	// Init user
	u := models.User{}

	// Get users
	// DB("go_test_tutorial") == is the Database to use. C("users") == the collection.
	if err := uc.session.DB("go_test_tutorial").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}
	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //200
	fmt.Fprintf(w, "%s", uj)
}

// POST - Controller for creating a new user
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub a user to be populated from the body
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.Id = bson.NewObjectId()

	// Write the user to mongo
	uc.session.DB("go_rest_tutorial").C("users").Insert(u)

	// Marshal provided interface into JSON structure
	//TODO: understand Marshal
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //201
	fmt.Fprintf(w, "%s", uj)
}

// RemoveUser removes an existing user resource DELETE
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: will write logic for this when DB is implimented. Just posting status code for now
	w.WriteHeader(http.StatusOK) //200
}
