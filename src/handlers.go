package main

import (
	"encoding/json"
	"fmt"
	"net/http"
//	"html"
	"io"
	"io/ioutil"
	//"strconv"
	"log"
	"time"
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	
	mux "github.com/julienschmidt/httprouter"
)

func Logger(r *http.Request) {
	
	start:= time.Now()
	
	log.Printf(
		"%s\t%s\t%s\t%s",
		r.Method,
		r.RequestURI,
		//name,
		time.Since(start),
	)
}

func Index(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	
	Logger(r)

	fmt.Fprintf(w, "<h1>Hello, welcome to my blog2</h1>")
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	//fmt.Fprintf(w, "Hello, %s\n", p.ByName("anything"))
}

func MauCreate(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	
	Logger(r)

	var request Request
	var response Mau
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// Save JSON to Todo struct
	if err := json.Unmarshal(body, &request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	start:= time.Now()
	db, err := sql.Open("mysql", "test:test@/test")

	if err != nil {
	    panic(err.Error())
	}

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO maus (InstanceID, AppID, UserID, CreationDate) VALUES( ?, ?, ?, ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtIns.Close()

	// Prepare statement for inserting data
	stmtOut, err := db.Prepare("INSERT INTO maus (InstanceID, AppID, UserID, CreationDate) VALUES( ?, ?, ?, ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtOut.Close()

	defer db.Close()
	
	log.Printf(
		"%s\t%s\t%s\t%s",
		r.Method,
		r.RequestURI,
		//name,
		start.Unix(),
		request.AppID,
	)

	//t := RepoCreateTodo(todo)
	response.Mau = 23
	response.Timestamp = start.Unix()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
