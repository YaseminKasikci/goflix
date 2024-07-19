package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

func(s *server) HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Welcome to Goflix")
	}
}

// Handles token creation by validating user credentials and returning a JWT if valid.
func(s *server) HandleTokenCreate() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type response struct {
		Token string `jason:"token"`
	}

	type responseError struct {
		Error string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request){

		//parsing login body
		req := request{}
		err := s.decode(w, r, &req)
		if err != nil {
			msg := fmt.Sprintf("Cannot parse login body. err=%v, err", err)
			log.Println(msg)
			s.respond(w, r, responseError{
				Error: msg,
			}, http.StatusBadRequest)
			return 
			}
			//check credentials
			found, err := s.store.FindUSer(req.Username, req.Password)
			if err != nil {
				msg := fmt.Sprintf("Cannot find user. err=%v", err)
				s.respond(w, r, responseError{
					Error: msg,
				}, http.StatusInternalServerError)
				return
			}
		if !found {
			s.respond(w, r, responseError{
				Error: "Invalid credentials",
			}, http.StatusUnauthorized)
			return
		}
		//generate JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": req.Username,
			"exp": time.Now().Add(time.Hour + time.Duration(1)).Unix(),
			"iat": time.Now().Unix(),
			//issue at time
		}) 
			// Sign token
			tokenStr, err := token.SignedString([]byte(JWT_APP_KEY))
			if err != nil {
				msg := fmt.Sprintf("Cannot generate JWT. err=%v", err)
				s.respond(w, r, responseError{
					Error: msg,
				}, http.StatusInternalServerError)
				return
			}
			//return token in response
			s.respond(w, r, response{
				Token: tokenStr,
			}, http.StatusOK)
	}
}
