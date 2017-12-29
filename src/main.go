package main

import (
	"os"
    //"net"
    "fmt"
    "database/sql"
	"log"
	"net/http"
)

var db *sql.DB

func main() {

	router := NewRouter()

	port := os.Getenv("PORT")
    var err error
    if port == "" {
        log.Fatal("PORT environment variable was not set")
	}

    db, err = sql.Open("mysql", "test:test@tcp(go_db:3306)/test")
    db.SetMaxOpenConns(100)
    fmt.Printf("conecta")
    if err != nil {
        log.Panic(err)
    }

	log.Fatal(http.ListenAndServe(":" + port, router))
    defer db.Close()
}
