package controllers

// comment
// more comments
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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

// ############# //
// CRUD HANDLERS //
// ############# //

// CREATE: POST, PUT

//CreateUser Controller for creating a new user
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//TODO: verify username and email are unique
	// shorten db string
	col := uc.session.DB("go_rest_tutorial").C("users")
	// Stub a user to be populated from the body
	us := models.Users{}

	// Convert Hex Object Id into string for redrect
	//sid := bson.ObjectId(us.ID).String()

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&us)

	// Get ObjectID from DB and assign to user
	us.ID = bson.NewObjectId()

	//Set user created time
	us.Created = time.Now().Local()

	//Set user updated time
	us.Updated = time.Now().Local()

	// Write the user to mongo
	col.Insert(us)

	// Write content-type and return statuscode and original payload
	w.Header().Set("Content-Type", "application/json")
	//http.Redirect(w, r, path+sid, http.StatusMovedPermanently) //302
	w.WriteHeader(http.StatusCreated) //201
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

	/* not a good implimentation. need to find another way to return all users.
	I think returning the slice []models.Users is messing with the browser.
	Using the browser returns an empty array.
	*/
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

// UPDATE - PATCH, POST, PUT methods

//UpdateUser Controller for updating user document
func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// shorten db string
	col := uc.session.DB("go_rest_tutorial").C("users")
	// Init user model - probably could be pulled out as a global package variable.
	us := models.Users{}

	// Setup Decoder and decode HTTP request body into User struct
	json.NewDecoder(r.Body).Decode(&us)

	// Grab id from the router
	id := p.ByName("id")

	// Verify id is ObjectId hex representation, otherwise return status not found
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// Convert id to ObjectId for update query
	oid := bson.ObjectIdHex(id)
	log.Print(oid)

	//Set user updated time
	us.Updated = time.Now().Local()

	// Write updates to DB
	col.Update(bson.M{"_id": oid}, bson.M{"$set": bson.M{"firstname": us.FirstName, "lastname": us.LastName, "username": us.Username, "email": us.Email, "role": us.Role, "updated": us.Updated, "bio": us.Bio}})

	uj, _ := json.Marshal(us)

	// Write content-type and return statuscode and original payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) //202
	fmt.Fprintf(w, "%s", uj)
}

// DELETE - DELETE method

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
