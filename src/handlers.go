package main

import (
	"encoding/json"
	"net/http"
	"io"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	
	mux "github.com/julienschmidt/httprouter"
)

func checkErr(err error) {
   if err != nil {
	   //panic(err)
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
		//panic(err)
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

func RegisterMau(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	var request Request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	checkErr(err)
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	//db, err := sql.Open("mysql", "test:test@tcp(go_db:3306)/test")
	//checkErr(err)

	var mau int
	err = db.QueryRow("SELECT sequence FROM active_users WHERE instance_id LIKE ?", request.InstanceID).Scan(&mau)

	if mau != 0 {
		//fmt.Printf("Old mau", mau)
		sendResponse(w, mau)
		//defer db.Close()
	} else {
		//fmt.Printf("SELECT sequence FROM active_users WHERE instance_id LIKE %s", request.InstanceID)
		fmt.Printf("New mau %d", mau)
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
			tx.Exec("INSERT INTO active_users (instance_id, sequence, user_id, application_id) VALUES(?, ?, ?, ?)", request.InstanceID, newmau, request.UserID, request.AppID)

			//fmt.Printf("INSERT INTO active_users (instance_id, sequence, user_id, application_id) VALUES('%s', %d, %d, %d)", request.InstanceID, newmau, request.UserID, request.AppID)
			tx.Commit()
			//defer db.Close()
			insert <- true
		}(insert)
		<-done
	}

	//defer db.Close()
}
