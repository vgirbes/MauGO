package main

import (
	"encoding/json"
	"net/http"
	"io"
	"io/ioutil"
	"log"
	"time"
	"strconv"
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

func RegisterMau(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	var request Request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	checkErr(err)
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) 
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}


	var mau int

    if entry, errCache := cache.Get(request.Instance); errCache != nil {
    	err = db.QueryRow("SELECT sequence FROM activeUsers WHERE instance_id LIKE ?", request.Instance).Scan(&mau)

		if mau != 0 {
			var mauCache = strconv.Itoa(mau)
			cache.Set(request.Instance, []byte(mauCache))
			sendResponse(w, mau)
		} else {
			var row Request
			var newmau int

			tx, err := db.Begin()

			err = tx.QueryRow("SELECT sequence, user_id, application_id, instance_id FROM activeUsers WHERE user_id = ? ORDER BY sequence DESC FOR UPDATE", request.IdUser).Scan(&newmau, &row.IdUser, &row.IdApp, &row.Instance)

			if err != nil {
				newmau = 0
			}

			newmau += 1
			done := make(chan bool)
			insert := make(chan bool)

			go func(done chan bool) {
				var mauCache = strconv.Itoa(newmau)
				cache.Set(request.Instance, []byte(mauCache))
				sendResponse(w, newmau)
				done <- true
			}(done)

			go func(insert chan bool) {
				tx.Exec("INSERT INTO activeUsers (instance_id, sequence, user_id, application_id) VALUES(?, ?, ?, ?)", request.Instance, newmau, request.IdUser, request.IdApp)
				tx.Commit()
				insert <- true
			}(insert)

			<-done
		}
    } else {
    	aByteToInt, _ := strconv.Atoi(string(entry))
    	sendResponse(w, aByteToInt)
    }
}
