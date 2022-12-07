package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

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

type Driver struct {
	// DriverID  string `json:"DriverId"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	PhoneNo   string `json:"PhoneNo"`
	Email     string `json:"Email"`
	LicenseNo string `json:"LicenseNo"`
	Password  string `json:"Password"`
}

type passengerTrip struct {
	
	// PassengerID string `json:"PassengerId"`
	StartPostal string  `json:"StartPostal"`
	EndPostal   string  `json:"EndPostal"`
	StartTime   *string `json:"StartTime"`
	EndTime     *string `json:"EndTime"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/newPassenger", createPassengers).Methods("POST", "PATCH")
	router.HandleFunc("/api/v1/passengers/{id}", updatePassengers).Methods("GET", "PATCH")
	router.HandleFunc("/api/v1/drivers/{id}", driverInfo).Methods("POST", "PATCH")
	router.HandleFunc("/api/v1/trip/{passengerID}", createTrips).Methods("POST", "PATCH")
	router.HandleFunc("/api/v1/trips/{driverID}/{startEndTrip}", updateTrips).Methods("GET", "PATCH")

	//passenger retrieve all trips he has taken before in reverse chronological order
	// router.HandleFunc("/api/v1/trips/{passengerID}", retrieveTrips).Methods("GET")

	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))
}

func createPassengers(w http.ResponseWriter, r *http.Request) {

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

func createTrips(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Method == "POST" {
		if reqBody, err := ioutil.ReadAll(r.Body); err == nil {
			var trip passengerTrip
			if err := json.Unmarshal(reqBody, &trip); err == nil {
				if _, ok := isPassengerExist(params["passengerID"]); ok {
					createTripForPassenger(params["passengerID"], trip)
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
}

func updateTrips(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Method == "PATCH" {
		if reqBody, err := ioutil.ReadAll(r.Body); err == nil {
			var trip passengerTrip
			if err := json.Unmarshal(reqBody, &trip); err == nil {
				if _, ok := isDriverExist(params["driverID"]); ok {
					if _, ok := isTripExist(params["driverID"]); ok {
						if params["startEndTrip"] == "start" {
							startTrip(params["driverID"])
							w.WriteHeader(http.StatusCreated)
						} else if params["startEndTrip"] == "end" {
							endTrip(params["driverID"])
							w.WriteHeader(http.StatusCreated)
						} else {
							w.WriteHeader(http.StatusBadRequest)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
					}

				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	}
}

func driverInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Method == "PATCH" {
		if reqBody, err := ioutil.ReadAll(r.Body); err == nil {
			var driver map[string]interface{}
			if err := json.Unmarshal(reqBody, &driver); err == nil {
				if orig, ok := isDriverExist(params["id"]); ok {
					fmt.Println(driver)
					for key, value := range driver {
						switch key {
						case "FirstName":
							orig.FirstName = value.(string)
						case "LastName":
							orig.LastName = value.(string)
						case "PhoneNo":
							orig.PhoneNo = value.(string)
						case "Email":
							orig.Email = value.(string)
						case "LicenseNo":
							orig.LicenseNo = value.(string)
						case "Password":
							orig.Password = value.(string)
						}
					}
					updateDriver(params["id"], orig)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	} else if r.Method == "POST" {
		if resBody, err := ioutil.ReadAll(r.Body); err == nil {
			var driver Driver
			if err := json.Unmarshal(resBody, &driver); err == nil {
				if _, ok := isDriverExist(params["id"]); !ok {
					insertDriver(params["id"], driver)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusConflict)
					fmt.Println("Driver already exist")
				}
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	}
}

func updatePassengers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Method == "PATCH" {
		if reqBody, err := ioutil.ReadAll(r.Body); err == nil {
			var passenger map[string]interface{}
			if err := json.Unmarshal(reqBody, &passenger); err == nil {
				if orig, ok := isPassengerExist(params["id"]); ok {
					fmt.Println(passenger)
					for key, value := range passenger {
						switch key {
						case "FirstName":
							orig.FirstName = value.(string)
						case "LastName":
							orig.LastName = value.(string)
						case "PhoneNo":
							orig.PhoneNo = value.(string)
						case "Email":
							orig.Email = value.(string)
						case "Password":
							orig.Password = value.(string)
						}
					}
					updatePassenger(params["id"], orig)
					w.WriteHeader(http.StatusAccepted)
					fmt.Println("Updated")
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Println("Not Found")
				}
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Println("Bad Request")
			}
		}

	}
}

// func retrieveTrips(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	if r.Method == "GET" {

// }


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

func insertDriver(id string, driver Driver) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/trip_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Insert the driver
	_, err = db.Exec("INSERT INTO driver (DriverID, FirstName, LastName, PhoneNo, Email, LicenseNo, Password) VALUES (?,?,?,?,?,?,?)", id, driver.FirstName, driver.LastName, driver.PhoneNo, driver.Email, driver.LicenseNo, driver.Password)
	if err != nil {
		panic(err.Error())
	}
}

func updateDriver(id string, driver Driver) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/trip_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Update the driver
	_, err = db.Exec("UPDATE driver SET FirstName=?, LastName=?, PhoneNo=?, Email=?, LicenseNo=?, Password=? WHERE DriverID=?", driver.FirstName, driver.LastName, driver.PhoneNo, driver.Email, driver.LicenseNo, driver.Password, id)
	if err != nil {
		panic(err.Error())
	}
}

func updatePassenger(id string, passenger Passenger) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/trip_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Update the passenger
	pid, _ := strconv.Atoi(id)
	_, err = db.Exec("UPDATE passenger SET FirstName=?, LastName=?, PhoneNo=?, Email=?, Password=? WHERE PassengerID=?", passenger.FirstName, passenger.LastName, passenger.PhoneNo, passenger.Email, passenger.Password, pid)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
}

func findDriver() (string, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/trip_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Find the driverID that has not been assigned a trip (onTrip = 0)

	var driverID string
	if err := db.QueryRow("SELECT DriverID FROM driver WHERE onTrip = 0 LIMIT 1").Scan(&driverID); err == nil {
		_, err = db.Exec("UPDATE driver SET onTrip = 1 WHERE DriverID = ?", driverID)
		if err != nil {
			panic(err.Error())
		}
		return driverID, true
	} else {
		panic(err.Error())
	}
}

func updateDriverFound(id string) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/trip_db")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	pid, _ := strconv.Atoi(id)

	driverID, found := findDriver()
	if found {
		_, err := db.Exec("UPDATE trip SET DriverID = ? WHERE PassengerID = ?", driverID, pid)
		if err != nil {
			panic(err.Error())
		}
	}

}

func createTripForPassenger(id string, trip passengerTrip) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/trip_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	pid, _ := strconv.Atoi(id)

	_, err2 := db.Exec("INSERT INTO trip (PassengerID, StartPostal, EndPostal) VALUES (?,?,?)", pid, trip.StartPostal, trip.EndPostal)
	if err2 != nil {
		panic(err.Error())
	}

	_, err3 := db.Exec("UPDATE passenger SET onTrip = 1 WHERE PassengerID = ?", pid)
	if err3 != nil {
		panic(err.Error())
	}

	updateDriverFound(id)
}

func isDriverExist(id string) (Driver, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/trip_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	results := db.QueryRow("SELECT FirstName,LastName,PhoneNo,Email,LicenseNo,Password FROM driver WHERE DriverID = ?", id)
	var driver Driver
	err = results.Scan(&driver.FirstName, &driver.LastName, &driver.PhoneNo, &driver.Email, &driver.LicenseNo, &driver.Password)
	if err != nil {
		return Driver{}, false
	}
	return driver, true
}

func isPassengerExist(id string) (Passenger, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/trip_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	pId, _ := strconv.Atoi(id)
	results := db.QueryRow("SELECT FirstName,LastName,PhoneNo,Email,Password FROM passenger WHERE PassengerID = ?", pId)
	var passenger Passenger
	err = results.Scan(&passenger.FirstName, &passenger.LastName, &passenger.PhoneNo, &passenger.Email, &passenger.Password)
	if err != nil {
		fmt.Println(err)
		return Passenger{}, false
	}
	return passenger, true
}

func isTripExist(driverId string) (passengerTrip, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/trip_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	results := db.QueryRow("SELECT StartPostal,EndPostal FROM trip WHERE DriverID = ?", driverId)
	var trip passengerTrip
	err = results.Scan(&trip.StartPostal, &trip.EndPostal)
	if err != nil {
		return passengerTrip{}, false
	}
	return trip, true
}

func startTrip(tripID string) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/trip_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	tid, _ := strconv.Atoi(tripID)

	_, err = db.Exec("UPDATE trip SET StartTime = ? WHERE TripID = ?", time.Now(), tid)
	if err != nil {
		panic(err.Error())
	}
}

func endTrip(tripID string) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/trip_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	tid, _ := strconv.Atoi(tripID)

	_, err = db.Exec("UPDATE trip SET EndTime = ? WHERE TripID = ?", time.Now(), tid)
	if err != nil {
		panic(err.Error())
	}
	//Set onTrip to 0 for both driver and passenger
	_, err = db.Exec("UPDATE driver SET onTrip = 0 WHERE DriverID = (SELECT DriverID FROM trip WHERE TripID = ?)", tid)
	if err != nil {
		panic(err.Error())
	}
	_, err2 := db.Exec("UPDATE passenger SET onTrip = 0 WHERE PassengerID = (SELECT PassengerID FROM trip WHERE TripID = ?)", tid)
	if err2 != nil {
		panic(err.Error())
	}

}

// func getTrip(tripID string) (passengerTrip, bool) {