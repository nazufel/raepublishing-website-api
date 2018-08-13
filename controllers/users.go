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

//Path for url redirects
var path = "http://localhost:3000/api/v1"

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

//TODO: set up goroutines in the UpdateUser method that use the below methods to update different fields of the user resource.
//+This will expose only one endpoint for updating the entire structure. Something like:
/*
```
if r.Body.Firstname != nil {
	go UpdateUserFirstname(r)
}
if r.Body.Lastname != nil {
	go UpdateUserLastname(r)
}
...
```
*/
//CreateUser Controller for creating a new user
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
	//col.Update(oid, us)
	col.Update(bson.M{"_id": oid}, bson.M{"$set": bson.M{"firstname": us.FirstName, "lastname": us.LastName, "username": us.Username}})

	uj, _ := json.Marshal(us)

	// Write content-type and return statuscode and original payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) //202
	fmt.Fprintf(w, "%s", uj)
}

//UpdateUsersFirstname controller to update user document fields
func (uc UserController) UpdateUsersFirstname(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//read the request message and parse the fields
	decoder := json.NewDecoder(r.Body)
	var us models.Users
	err := decoder.Decode(&us)
	if err != nil {
		log.Fatal(err)
	}

	col := uc.session.DB("go_rest_tutorial").C("users")

	// get the user id from the httprouter parameter
	id := p.ByName("id")

	// verify id is Objectid
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// get ObjectId
	oid := bson.ObjectIdHex(id)
	// set ObjectId as a string for url redirect
	//sid := bson.ObjectId(id).String()

	//MongoDB query, build the changes
	change := mgo.Change{
		// Now to need to loop through users scruct
		Update:    bson.M{"$set": bson.M{"firstname": us.FirstName, "lastname": us.LastName}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	// store updated document in result variable
	//	var result bson.M

	// apply the changes to the document(s)
	_, err = col.Find(bson.M{"_id": oid}).Apply(change, &us)
	if err != nil {
		fmt.Println(err)
		//w.WriteHeader(http.StatusNotFound) //404
		return
	}

	//Write the updated time
	updatedTime := mgo.Change{
		// Now to need to loop through users scruct
		Update:    bson.M{"$set": bson.M{"updated": time.Now()}},
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}
	// store updated document in result variable
	var updatedResult bson.M

	// apply the changes to the document
	_, err = col.Find(bson.M{"_id": us.ID}).Apply(updatedTime, &updatedResult)
	if err != nil {
		fmt.Println(err)
		//w.WriteHeader(http.StatusNotFound) //404
		return
	}

	uj, _ := json.Marshal(us)

	// Write content-type and return statuscode and original payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) //202
	fmt.Fprintf(w, "%s", uj)
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

//UpdateUsersEmail updates the user's lastname field
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

//UpdateUsersRole updates the user's username field
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

//UpdateUsersBio updates the user's lastname field
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
