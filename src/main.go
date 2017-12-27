package main

import (
	"os"
	"log"
	"net/http"
)

func main() {

	router := NewRouter()

	port := os.Getenv("PORT")
    if port == "" {
        log.Fatal("PORT environment variable was not set")
	}

	log.Fatal(http.ListenAndServe(":" + port, router))
}
