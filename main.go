package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nazufel/raepublishing-website-api/controllers"
	mgo "gopkg.in/mgo.v2"
)

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
	r := httprouter.New()

	// Init UserController
	uc := controllers.NewUserController(getSession())
	// User Controllers
	r.GET("/users/:id", uc.GetUser)
	r.POST("/users/", uc.CreateUser)
	r.PUT("/users/", uc.CreateUser)
	r.DELETE("/users/:id", uc.DeleteUser)
	// Start the server
	http.ListenAndServe("localhost:3000", r)
}
