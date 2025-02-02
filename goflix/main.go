package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Goflix")
		
		if err := run(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n",err)
			os.Exit(1)
		}
	
}

func run() error  {
		//Server creation 
		srv := newServer()
		srv.store = &dbStore{}
		err := srv.store.Open()
		if err != nil {
			return err
		}
		
		http.HandleFunc("/", srv.serveHTTP)
		log.Printf("Serving HTTP on port 8080")
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			return err
		}

		defer srv.store.Close()
		return nil
}