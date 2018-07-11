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

//NewUserController represents the controller for updating new User Resources
func NewUserController(s *mgo.Session) *UserController {
	// init mongo
	return &UserController{s}
}

// CRUD HANDLERS //

// CREATE: POST, PUT

//CreateUser Controller for creating a new user
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//TODO: verify username and email are unique
	// Stub a user to be populated from the body
	u := models.Users{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.ID = bson.NewObjectId()

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

// READ: GET //

//GetUsers retrieves an individual user resource
func (uc UserController) GetUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId hex representation, otherwise return status not found
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// ObjectIdHex returns an ObjectId from the provided hex representation.
	oid := bson.ObjectIdHex(id)

	// composite literal
	u := models.Users{}

	// Fetch user
	if err := uc.session.DB("go_rest_tutorial").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", uj)
}

//TODO: Decide if returning an array of users instead of each individual users is okay.

//GetAllUsers returns all users in the collection
func (uc UserController) GetAllUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var users []models.Users

	// Fetch users
	err := uc.session.DB("go_rest_tutorial").C("users").Find(bson.M{}).All(&users)
	if err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	uj, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Users{})
	fmt.Fprintf(w, "%s\n", uj)
}

//UpdateUsers controller to update user document fields
func (uc UserController) UpdateUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//decoder := json.NewDecoder(r.Body)
	s := uc.session.DB("go_rest_tutorial").C("users")

	// get the user id from the httprouter parameter
	id := p.ByName("id")

	// verify id is Objectid
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// get ObjectId
	oid := bson.ObjectIdHex(id)

	//MongoDB query, build the changes
	changeFirstName := mgo.Change{
		/*Breakthrough: feeding the Query.Apply method u.FirstName gives it the
		+empty values in the models.Users{} struct. Need to parse the json http
		+payload for specific fields, ie: {"firstname": "Bobby"} and update only
		+those fields.
		*/
		Update:    bson.M{"$set": bson.M{"firstname": "Bobby"}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	// store updated document in result variable
	var result bson.M

	// apply the changes to the document(s)
	_, err := s.Find(bson.M{"_id": oid}).Apply(changeFirstName, &result)

	if err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// DeleteUsers removes an existing user resource DELETE
func (uc UserController) DeleteUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// get the user id from the httprouter parameter
	id := p.ByName("id")

	// verify id is Objectid
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// get ObjectId
	oid := bson.ObjectIdHex(id)

	// delete user
	err := uc.session.DB("go_rest_tutorial").C("users").RemoveId(oid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
