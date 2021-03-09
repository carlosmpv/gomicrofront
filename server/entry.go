package main

import (
	"net/http"
)

// Router define api routes
func Router(mux *http.ServeMux) {

}

// Entrypoint function to handle every request
// Returns a status code and error (nil if none)
// Used intercept any request
func Entrypoint(w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}
