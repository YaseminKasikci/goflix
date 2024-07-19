package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Store interface {
	Open() error // ouvrir connect vers sql
	Close() error 

	GetMovies() ([]*Movie, error)
	GetMovieById(id int64) (*Movie, error)
	CreateMovie(m *Movie) error 
	FindUSer(username string, password string) (bool, error)
}


type dbStore struct {
	db *sqlx.DB
}

// schema defines the SQL statements to create the movie and user tables
var schema =`
CREATE TABLE IF NOT EXISTS movie
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT,
	release_date TEXT,
	duration INTERGER,
	trailer_url TEXT
);


CREATE TABLE IF NOT EXISTS user
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user TEXT,
	password TEXT
);
`
//open the connection to the database
func (store *dbStore) Open() error {	
	db, err := sqlx.Connect("sqlite3", "goflix.db")
	if err != nil {
		return err
	}
	log.Println("Connected to DB")
	db.MustExec(schema) //mustexec fais planté le programme si ca ne fonctionne pas 

	store.db = db
	return nil
}

// close terminates the connection to the database
func (store *dbStore) Close() error {
	return store.db.Close()
}

// GetMovies retrieves all movies from the database and returns them as a slice of Movie pointers.
func (store *dbStore) GetMovies() ([]*Movie, error) {
	var movies []*Movie
	err := store.db.Select(&movies, "SELECT * FROM movie") // select remplis le tableau &movies
	if err != nil {
		return movies, err
	}
	return movies, nil
}

// GetMovieById retrieves a specific movie by its ID from the database.
func (store *dbStore) GetMovieById(id int64) (*Movie, error) {
	var movie = &Movie{}
	err := store.db.Get(movie, "SELECT * FROM movie WHERE id=$1", id)
	if err != nil { 
		return movie, nil
	}
	return movie, err
}

// CreateMovie inserts a new movie into the database and updates the movie's ID with the newly assigned ID.
func (store *dbStore) CreateMovie(m *Movie) error {
	res, err := store.db.Exec("INSERT INTO movie (title, release_date, duration, trailer_url) VALUES (?, ?, ?, ?)",m.Title, m.ReleaseDate, m.Duration, m.TrailerURL)
	if err != nil {
		return err
	}
	m.ID, err = res.LastInsertId()
	return err 
}

// FindUSer checks if the provided username and password match a user in the database.
func (store *dbStore) FindUSer(username string, password string) (bool, error) {
	var count int 
	err := store.db.Get(&count,  "SELECT COUNT(id) FROM user WHERE user=$1 AND password=$2", username, password)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}