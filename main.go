package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nazufel/raepublishing-website-api/controllers"
	mgo "gopkg.in/mgo.v2"
)

//TODO: use env to inject db creds and location.
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}

func main() {
	// Init router
	/* TODO: Evaluate whether to continue using the julienschmidt router or Gorilla Mux.
	+ Gorilla Mux does allow multiple methods per route. It could be nice to have POST and PUT
	+ for create and update users*/
	r := httprouter.New()

	// Init UserController
	uc := controllers.NewUserController(getSession())
	// User Controllers
	r.GET("/users/:id", uc.GetUser)
	r.GET("/users/", uc.GetAllUsers)
	r.POST("/users/", uc.CreateUser)
	r.PUT("/users/", uc.CreateUser)
	r.DELETE("/users/:id", uc.DeleteUser)
	// Start the server
	http.ListenAndServe("localhost:3000", r)
}
