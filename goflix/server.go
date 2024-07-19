package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const JWT_APP_KEY = "training.go"

// server represents the HTTP server with a router and a data store.
type server struct {
	router *mux.Router
	store Store
}

// Creates a new server instance with a router and initializes routes.
func newServer() *server {
	s := &server{
		router: mux.NewRouter(),
	}
	s.routes()
	return s
}

// Handles HTTP requests by applying the logRequestMiddleware to the router's ServeHTTP method
func (s *server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	logRequestMiddleware(s.router.ServeHTTP).ServeHTTP(w,r)
}

// Sends a JSON response with the specified status code; logs an error if JSON encoding fails.
func (s *server) respond(w http.ResponseWriter, _ *http.Request, data interface{}, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
// si je n'ai pas de donné
	if data == nil {
		return
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Cannot format json. err=%v\n", err)
	}
}

// Decodes the JSON request body into the provided struct
func (s *server) decode(w http.ResponseWriter, r *http.Request, v interface {}) error {
	return json.NewDecoder(r.Body).Decode(v)
}