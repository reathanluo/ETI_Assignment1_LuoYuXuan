package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	// "io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Trip struct {
	TripID string `json:"TripID" `

	DriverID      string `json:"DriverID"`
	DriverName    string `json:"DriverName"`
	PassengerID   string `json:"PassengerID"`
	PassengerName string `json:"PassengerName"`
	LicenseNo     string `json:"LicenseNo"`
	StartTime     string `json:"StartTime"`
	EndTime       string `json:"EndTime"`
	StartPostal   string `json:"StartPostal"`
	EndPostal     string `json:"EndPostal"`
}

var sqlDb = "user:password@tcp(127.0.0.1:3307)/trip_db"

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/trips/alltrips/{userID}", GetTrips).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/trips/create/{userID}", CreateTrip).Methods("GET","POST", "OPTIONS")
	router.HandleFunc("/api/v1/trips/currenttrip/{driverID}", GetCurrentTrip).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/trips/currenttripns/{driverID}", GetCurrentTripNotStart).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/trips/start/{driverID}", startTrip).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/trips/end/{driverID}", endTrip).Methods("POST", "OPTIONS")
	fmt.Println("Listening at port 5003")
	log.Fatal(http.ListenAndServe(":5003", router))
}

func GetTrips(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		var trips []Trip
		params := mux.Vars(r)
		userID := params["userID"]
		db, err := sql.Open("mysql", sqlDb)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		result, err := db.Query("SELECT t.TripID, t.PassengerID, t.DriverID, CONCAT_WS('',d.FirstName,' ',d.LastName) AS Name,d.LicenseNo,t.StartTime, t.EndTime, t.StartPostal, t.EndPostal FROM trip t inner join Driver d on t.DriverID = d.DriverID where t.PassengerID = ? AND EndTime IS NOT NULL ORDER BY StartTime DESC", userID)
		if err != nil {
			panic(err.Error())
		}
		defer result.Close()
		for result.Next() {
			var trip Trip
			err := result.Scan(&trip.TripID, &trip.PassengerID, &trip.DriverID, &trip.DriverName, &trip.LicenseNo, &trip.StartTime, &trip.EndTime, &trip.StartPostal, &trip.EndPostal)
			if err != nil {
				panic(err.Error())
			}
			trips = append(trips, trip)
		}
		json.NewEncoder(w).Encode(trips)
	}

}

func CreateTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	} else if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		userID := params["userID"]
		db, err := sql.Open("mysql", sqlDb)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		var currentTrip Trip
		result, err6 := db.Query("SELECT t.TripID, t.PassengerID, t.DriverID, CONCAT_WS('',d.FirstName,' ',d.LastName) AS Name,d.LicenseNo, t.StartPostal, t.EndPostal FROM trip t inner join Driver d on t.DriverID = d.DriverID where t.PassengerID = ? AND t.EndTime IS NULL", userID)
		if err6 != nil {
			panic(err.Error())
		}
		defer result.Close()
		for result.Next() {
			err := result.Scan(&currentTrip.TripID, &currentTrip.PassengerID, &currentTrip.DriverID, &currentTrip.DriverName, &currentTrip.LicenseNo, &currentTrip.StartPostal, &currentTrip.EndPostal)
			if err != nil {
				panic(err.Error())
			}
		}
		json.NewEncoder(w).Encode(currentTrip)
	}else {
		w.Header().Set("Content-Type", "application/json")
		var trip Trip
		parms := mux.Vars(r)
		trip.PassengerID = parms["userID"]
		_ = json.NewDecoder(r.Body).Decode(&trip)
		db, err := sql.Open("mysql", sqlDb)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		_, err2 := db.Exec("INSERT INTO trip (PassengerID, StartPostal, EndPostal) VALUES (?, ?, ?)", trip.PassengerID, trip.StartPostal, trip.EndPostal)
		if err2 != nil {
			panic(err.Error())
		}

		_, err3 := db.Exec("UPDATE passenger SET onTrip = 1 WHERE PassengerID = ?", trip.PassengerID)
		if err3 != nil {
			panic(err.Error())
		}

		driverID, found := findDriver()
		if found {
			_, err4 := db.Exec("UPDATE trip SET DriverID = ? WHERE PassengerID = ? AND EndTime IS NULL", driverID, trip.PassengerID)
			if err4 != nil {
				panic(err.Error())
			}
			_, err5 := db.Exec("UPDATE driver SET onTrip = 1 WHERE DriverID = ?", driverID)
			if err5 != nil {
				panic(err.Error())
			}
		}
		var currentTrip Trip
		result, err6 := db.Query("SELECT t.TripID, t.PassengerID, t.DriverID, CONCAT_WS('',d.FirstName,' ',d.LastName) AS Name,d.LicenseNo, t.StartPostal, t.EndPostal FROM trip t inner join Driver d on t.DriverID = d.DriverID where t.PassengerID = ? AND t.EndTime IS NULL", trip.PassengerID)

		if err6 != nil {
			panic(err.Error())
		}
		defer result.Close()
		for result.Next() {
			err := result.Scan(&currentTrip.TripID, &currentTrip.PassengerID, &currentTrip.DriverID, &currentTrip.DriverName, &currentTrip.LicenseNo, &currentTrip.StartPostal, &currentTrip.EndPostal)
			if err != nil {
				panic(err.Error())
			}
		}
		json.NewEncoder(w).Encode(currentTrip)
		
	}

}

func findDriver() (string, bool) {
	db, err := sql.Open("mysql", sqlDb)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var driverID string
	if err := db.QueryRow("SELECT DriverID FROM driver WHERE onTrip = 0 LIMIT 1").Scan(&driverID); err == nil {
		if err != nil {
			panic(err.Error())
		}
		return driverID, true
	} else {
		panic(err.Error())
	}

}

func startTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		var trip Trip
		parms := mux.Vars(r)
		trip.DriverID = parms["driverID"]
		_ = json.NewDecoder(r.Body).Decode(&trip)
		db, err := sql.Open("mysql", sqlDb)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		_, err2 := db.Exec("UPDATE trip SET StartTime = ? WHERE DriverID = ? AND EndTime IS NULL", time.Now(), trip.DriverID)
		if err2 != nil {
			panic(err.Error())
		}
	}
}

func endTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		var trip Trip
		parms := mux.Vars(r)
		trip.DriverID = parms["driverID"]
		_ = json.NewDecoder(r.Body).Decode(&trip)
		db, err := sql.Open("mysql", sqlDb)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		_, err4 := db.Exec("UPDATE passenger SET onTrip = 0 WHERE PassengerID = (SELECT PassengerID FROM trip WHERE DriverID = ? AND EndTime IS NULL)", trip.DriverID)
		if err4 != nil {
			panic(err.Error())
		}

		_, err2 := db.Exec("UPDATE trip SET EndTime = ? WHERE DriverID = ? AND EndTime IS NULL", time.Now(), trip.DriverID)
		if err2 != nil {
			panic(err.Error())
		}
		_, err3 := db.Exec("UPDATE driver SET onTrip = 0 WHERE DriverID = ?", trip.DriverID)
		if err3 != nil {
			panic(err.Error())
		}

	}
}

func GetCurrentTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		parms := mux.Vars(r)
		db, err := sql.Open("mysql", sqlDb)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		var currentTrip Trip

		result, err2 := db.Query("SELECT t.TripID, t.PassengerID, t.DriverID, CONCAT_WS('',p.FirstName,' ',p.LastName) AS Name, t.StartTime, t.StartPostal, t.EndPostal FROM trip t inner join Passenger p on t.PassengerID = p.PassengerID where t.DriverID = ? AND t.EndTime IS NULL", parms["driverID"])
		if err2 != nil {
			panic(err.Error())
		}
		defer result.Close()
		for result.Next() {
			err := result.Scan(&currentTrip.TripID, &currentTrip.PassengerID, &currentTrip.DriverID, &currentTrip.PassengerName, &currentTrip.StartTime, &currentTrip.StartPostal, &currentTrip.EndPostal)
			if err != nil {
				panic(err.Error())
			}
		}
		json.NewEncoder(w).Encode(currentTrip)
	}
}

func GetCurrentTripNotStart(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		parms := mux.Vars(r)
		db, err := sql.Open("mysql", sqlDb)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		var currentTrip Trip

		result, err2 := db.Query("SELECT t.TripID, t.PassengerID, t.DriverID, CONCAT_WS('',p.FirstName,' ',p.LastName) AS Name, t.StartPostal, t.EndPostal FROM trip t inner join Passenger p on t.PassengerID = p.PassengerID where t.DriverID = ? AND t.EndTime IS NULL", parms["driverID"])
		if err2 != nil {
			panic(err.Error())
		}
		defer result.Close()
		for result.Next() {
			err := result.Scan(&currentTrip.TripID, &currentTrip.PassengerID, &currentTrip.DriverID, &currentTrip.PassengerName, &currentTrip.StartPostal, &currentTrip.EndPostal)
			if err != nil {
				panic(err.Error())
			}
		}
		json.NewEncoder(w).Encode(currentTrip)
	}
}
