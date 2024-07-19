package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type jsonMovie struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	ReleaseDate string `json:"release_date"` //date de publication
	Duration int `json:"duration"`
	TrailerURL string `json:"trailer_url"` 
}

// Handles requests to list all movies, returning them as a JSON response.
func (s *server) handleMovieList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){

		// Retrieve the list of movies from the store
		movies, err := s.store.GetMovies()
		if err != nil {
			log.Printf("Cannot load movies. err=%v\n", err)

			s.respond(w,r,nil, http.StatusInternalServerError)
			return
		}

		// Convert the list of movies to a list of jsonMovie
		var resp = make([]jsonMovie, len(movies))
		for i, m := range movies {
			resp[i] = mapMovieToJson(m)
		}
		// Send the list of movies as a JSON response
		s.respond(w, r, resp, http.StatusOK)
	}
}

// Handles requests for movie details by ID, returning the movie as a JSON response.
func (s *server) handleMovieDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the movie ID from the URL parameters
		vars := mux.Vars(r) 
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.Printf("Cannir parse id to int. err=%v", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		// Retrieve the movie details by ID
		m, err := s.store.GetMovieById(id)
		if err != nil {
			log.Printf("Cannot load movie err=%v", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}
		// Convert the movie details to JSON format
		var resp = mapMovieToJson(m)
		s.respond(w, r, resp, http.StatusOK)
	}
}

// Handles movie creation requests by parsing the request body, storing the movie in the database, and returning the movie details as a JSON response.
func (s *server) handleMovieCreate() http.HandlerFunc {
	type request struct {
		Title string `json:"title"`
		ReleaseDate string `json:"release_date"` //date de publication
		Duration int `json:"duration"`
		TrailerURL string `json:"trailer_url"` 
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body to extract movie details
		req := request{}
		err := s.decode(w, r, &req)
		if err != nil {
			log.Printf("Cannot parse movie body. err=%v", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}
		// create a Movie instance
		m := &Movie{
			ID:0,
			Title: req.Title,
			ReleaseDate: req.ReleaseDate,
			Duration: req.Duration,
			TrailerURL: req.TrailerURL,

		}

		// store the movie in database
		err = s.store.CreateMovie(m)
		if err != nil {
			log.Printf("Cannot create movie in db err=%v", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		var resp = mapMovieToJson(m)
		s.respond(w, r, resp, http.StatusOK)
	}
}

// Converts a Movie struct to a jsonMovie struct for JSON serialization.
func mapMovieToJson(m *Movie) jsonMovie {
	return jsonMovie {
		ID: m.ID,
		Title: m.Title,
		ReleaseDate: m.ReleaseDate,
		Duration: m.Duration,
		TrailerURL: m.TrailerURL,
	}
}