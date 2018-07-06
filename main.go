package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Init router
	r := httprouter.New()

	// Add handler on /test
	r.GET("/test", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Simply write some test data for now
		fmt.Fprint(w, "Welcome!\n")
	})

	// Start the server
	http.ListenAndServe("localhost:8080", r)
}
