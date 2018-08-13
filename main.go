package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nazufel/raepublishing-website-api/controllers"
	mgo "gopkg.in/mgo.v2"
)

//comment change
//TODO: use env to inject db creds and location.
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")
	//defer s.Close()
	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	//defer s.Close()
	s.SetMode(mgo.Monotonic, true)
	return s
}

func main() {
	// Init router
	r := httprouter.New()

	//path extention to endpoints not at the root of the domain
	//+i.e.: http://raepublishing.com/api/v1/
	extension := "/api/v1"

	// Init UserController
	uc := controllers.NewUserController(getSession())
	// User Controllers

	// CREATE
	r.POST(extension+"/users/", uc.CreateUser)
	r.PUT(extension+"/users/", uc.CreateUser)
	// READ
	r.GET(extension+"/users/", uc.GetAllUsers)
	r.GET(extension+"/users/:id", uc.GetUsers)
	// UPDATE
	r.PATCH(extension+"/users/:id", uc.UpdateUser)
	r.PUT(extension+"/users/:id", uc.UpdateUser)
	r.POST(extension+"/users/:id", uc.UpdateUser)
	// DELETE
	r.DELETE(extension+"/users/:id", uc.DeleteUsers)
	// Start the server
	http.ListenAndServe("localhost:3000", r)
}
