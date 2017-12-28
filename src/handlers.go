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
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	
	mux "github.com/julienschmidt/httprouter"
)

func checkErr(err error) {
   if err != nil {
	   panic(err)
   }
}

func sendResponse(w http.ResponseWriter, mau int) {
	var response Mau
	now:= time.Now()
	response.Mau = mau
	response.Timestamp = now.Unix()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func Logger(r *http.Request) {
	
	now:= time.Now()
	
	log.Printf(
		"%s\t%s\t%s\t%s",
		r.Method,
		r.RequestURI,
		//name,
		time.Since(now),
	)
}

func Index(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	
	Logger(r)

	fmt.Fprintf(w, "<h1>Hello, welcome to my blog2</h1>")
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	//fmt.Fprintf(w, "Hello, %s\n", p.ByName("anything"))
}

func RegisterMau(w http.ResponseWriter, r *http.Request, _ mux.Params) {

	var request Request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	checkErr(err)
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

	db, err := sql.Open("mysql", "test:test@tcp(golang_db:3306)/test")
	checkErr(err)

	// Prepare statement for selecting data
	var mau int
	err = db.QueryRow("SELECT sequence FROM active_users WHERE instance_id LIKE ?", request.InstanceID).Scan(&mau)
	if err == nil {
		sendResponse(w, mau)
	} else {
		tx, err := db.Begin()
		var row Request
		var newmau int
		err = tx.QueryRow("SELECT sequence, user_id, application_id, instance_id FROM active_users WHERE user_id = ? ORDER BY sequence DESC FOR UPDATE", request.UserID).Scan(&newmau, &row.UserID, &row.AppID, &row.InstanceID)
		if err != nil {
			newmau = 0
		}
		newmau += 1
		done := make(chan bool)
		insert := make(chan bool)
		go func(done chan bool) {
			sendResponse(w, newmau)
			done <- true
		}(done)
		go func(insert chan bool) {
			tx.Query("INSERT INTO active_users (instance_id, sequence, user_id, application_id) VALUES(?, ?, ?, ?)", request.InstanceID, newmau, request.UserID, request.AppID)
			tx.Commit()
			insert <- true
		}(insert)
		<-done
		<-insert
	}
	defer db.Close()
}
