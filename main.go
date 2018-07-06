package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nazufel/raepublishing-website-api/controllers"
)

func main() {
	// Init router
	r := httprouter.New()

	// Init UserController
	uc := controllers.NewUserController()
	// Add handler on /test
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user/", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	// Start the server
	http.ListenAndServe("localhost:3000", r)
}
