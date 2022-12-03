package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type Passenger struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	PhoneNo   string `json:"PhoneNo"`
	Email     string `json:"Email"`
	Password  string `json:"Password"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/passengers", passengers).Methods("POST", "PATCH")
	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))
}

func passengers(w http.ResponseWriter, r *http.Request) {

	// Get the request body
	if r.Method == "POST" {
		if reqBody, err := ioutil.ReadAll(r.Body); err == nil {
			var passenger Passenger
			if err := json.Unmarshal(reqBody, &passenger); err == nil {
				insertPassenger(passenger)
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func insertPassenger(passenger Passenger) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/trip_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Insert the passenger
	_, err = db.Exec("INSERT INTO passenger (FirstName, LastName, PhoneNo, Email, Password) VALUES (?,?,?,?,?)", passenger.FirstName, passenger.LastName, passenger.PhoneNo, passenger.Email, passenger.Password)
	if err != nil {
		panic(err.Error())
	}
}
