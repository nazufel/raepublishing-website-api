package controllers

import (
	"encoding/json"
	"fmt"
	"log"
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

//UpdateUser function to update a user fields
func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var (
		fn = "firstname"
	)
	//shorten session string
	s := uc.session.DB("go_rest_tutorial").C("users")
	// init User model
	var us models.Users

	// get the user id from the httprouter parameter and ensure it's a bson object from the DB
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}
	// set var for ObjectId
	oid := bson.ObjectIdHex(id)

	//read the request message and parse the fields
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&us)
	if err != nil {
		log.Fatal(err)
	}
	//MongoDB query, build the changes
	change := mgo.Change{
		// Now to need to loop through users scruct
		Update:    bson.M{"$set": bson.M{fn: us.FirstName}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	// store updated document in result variable
	var result bson.M

	// apply the changes to the document(s)
	_, err = s.Find(bson.M{"_id": oid}).Apply(change, &result)

	if err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Users{})
}

/*
//UpdateUsersFirstname controller to update user document fields
func (uc UserController) UpdateUsersFirstname(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//read the request message and parse the fields
	decoder := json.NewDecoder(r.Body)
	var us models.Users
	err := decoder.Decode(&us)
	if err != nil {
		log.Fatal(err)
	}

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
	change := mgo.Change{
		// Now to need to loop through users scruct
		Update:    bson.M{"$set": bson.M{"firstname": us.FirstName}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	// store updated document in result variable
	var result bson.M

	// apply the changes to the document(s)
	_, err = s.Find(bson.M{"_id": oid}).Apply(change, &result)

	if err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Users{})
}

//UpdateUsersLastname updates the user's lastname field
func (uc UserController) UpdateUsersLastname(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//read the request message and parse the fields
	decoder := json.NewDecoder(r.Body)
	var us models.Users
	err := decoder.Decode(&us)
	if err != nil {
		log.Fatal(err)
	}

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
	change := mgo.Change{
		// Now to need to loop through users scruct
		Update:    bson.M{"$set": bson.M{"lastname": us.LastName}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	// store updated document in result variable
	var result bson.M

	// apply the changes to the document(s)
	_, err = s.Find(bson.M{"_id": oid}).Apply(change, &result)

	if err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Users{})
}

//r.PATCH("/users/:id", uc.UpdateUsersUsername)
//UpdateUsersUsername updates the user's username field
func (uc UserController) UpdateUsersUsername(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//read the request message and parse the fields
	decoder := json.NewDecoder(r.Body)
	var us models.Users
	err := decoder.Decode(&us)
	if err != nil {
		log.Fatal(err)
	}

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
	change := mgo.Change{
		// Now to need to loop through users scruct
		Update:    bson.M{"$set": bson.M{"username": us.Username}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	// store updated document in result variable
	var result bson.M

	// apply the changes to the document(s)
	_, err = s.Find(bson.M{"_id": oid}).Apply(change, &result)

	if err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Users{})
}

//r.PATCH("/users/:id", uc.UpdateUsersEmail)
//UpdateUsersLastname updates the user's lastname field
func (uc UserController) UpdateUsersEmail(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//read the request message and parse the fields
	decoder := json.NewDecoder(r.Body)
	var us models.Users
	err := decoder.Decode(&us)
	if err != nil {
		log.Fatal(err)
	}

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
	change := mgo.Change{
		// Now to need to loop through users scruct
		Update:    bson.M{"$set": bson.M{"email": us.Email}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	// store updated document in result variable
	var result bson.M

	// apply the changes to the document(s)
	_, err = s.Find(bson.M{"_id": oid}).Apply(change, &result)

	if err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Users{})
}

//r.PATCH("/users/:id", uc.UpdateUsersUsername)
//UpdateUsersUsername updates the user's username field
func (uc UserController) UpdateUsersRole(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//read the request message and parse the fields
	decoder := json.NewDecoder(r.Body)
	var us models.Users
	err := decoder.Decode(&us)
	if err != nil {
		log.Fatal(err)
	}

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
	change := mgo.Change{
		// Now to need to loop through users scruct
		Update:    bson.M{"$set": bson.M{"role": us.Role}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	// store updated document in result variable
	var result bson.M

	// apply the changes to the document(s)
	_, err = s.Find(bson.M{"_id": oid}).Apply(change, &result)

	if err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Users{})
}

//r.PATCH("/users/:id", uc.UpdateUsersRole)
//UpdateUsersLastname updates the user's lastname field
func (uc UserController) UpdateUsersBio(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//read the request message and parse the fields
	decoder := json.NewDecoder(r.Body)
	var us models.Users
	err := decoder.Decode(&us)
	if err != nil {
		log.Fatal(err)
	}

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
	change := mgo.Change{
		// Now to need to loop through users scruct
		Update:    bson.M{"$set": bson.M{"bio": us.Bio}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	// store updated document in result variable
	var result bson.M

	// apply the changes to the document(s)
	_, err = s.Find(bson.M{"_id": oid}).Apply(change, &result)

	if err != nil {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Users{})
}

*/

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
