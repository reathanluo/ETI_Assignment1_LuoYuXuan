package main

import (
	//"encoding/json"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type User struct {
	UserID   int    `json:"UserID"`
	UserType string `json:"UserType"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type Passenger struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	PhoneNo   string `json:"PhoneNo"`
	Email     string `json:"Email"`
	Password  string `json:"Password"`
}
type Driver struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	PhoneNo   string `json:"PhoneNo"`
	Email     string `json:"Email"`
	IcNO      string `json:"IcNO"`
	LicenseNo string `json:"LicenseNo"`
	Password  string `json:"Password"`
}

type LoggedInUser struct {
	UserID    string `json:"UserID"`
	UserType  string `json:"UserType"`
	Email     string `json:"Email"`
	Password  string `json:"Password"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	PhoneNo   string `json:"PhoneNo"`
	IcNO      string `json:"IcNO"`
	LicenseNo string `json:"LicenseNo"`
}

var sqlDb = "user:password@tcp(127.0.0.1:3307)/trip_db"

func main() {
	router := mux.NewRouter()

	db, err := sql.Open("mysql", sqlDb)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	router.HandleFunc("/api/v1/passenger/{passengerId}", GetPassenger).Methods("GET", "POST", "OPTIONS")
	router.HandleFunc("/api/v1/driver/{driverId}", GetDriver).Methods("GET", "POST", "OPTIONS")
	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))
}

func GetPassenger(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", sqlDb)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	passengerId := vars["passengerId"]

	var loggedInUser LoggedInUser
	loggedInUser.UserID = passengerId
	loggedInUser.UserType = "passenger"

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else if r.Method == "GET" {

		results, err := db.Query("SELECT p.FirstName, p.LastName, p.PhoneNo, u.Email, u.Password from Passenger p inner join User u on p.PassengerID = u.UserID where p.PassengerId = ?", passengerId)
		if err != nil {
			panic(err.Error())
		}
		for results.Next() {
			err = results.Scan(&loggedInUser.FirstName, &loggedInUser.LastName, &loggedInUser.PhoneNo, &loggedInUser.Email, &loggedInUser.Password)
			if err != nil {
				panic(err.Error())
			}
		}
		LoggedInUser := LoggedInUser{
			UserID:    loggedInUser.UserID,
			UserType:  loggedInUser.UserType,
			Email:     loggedInUser.Email,
			Password:  loggedInUser.Password,
			FirstName: loggedInUser.FirstName,
			LastName:  loggedInUser.LastName,
			PhoneNo:   loggedInUser.PhoneNo,
			IcNO:      "",
			LicenseNo: "",
		}
		json.NewEncoder(w).Encode(LoggedInUser)
		return

	} else if r.Method == "POST" {
		// update passenger details
		reqBody, _ := ioutil.ReadAll(r.Body)
		var passenger Passenger
		json.Unmarshal(reqBody, &passenger)

		// update passenger details in database
		_, err := db.Exec("UPDATE Passenger SET FirstName = ?, LastName = ?, PhoneNo = ? WHERE PassengerId = ?", passenger.FirstName, passenger.LastName, passenger.PhoneNo, passengerId)
		if err != nil {
			panic(err.Error())
		}
		// update user details in database
		_, err2 := db.Exec("UPDATE User SET Email = ? WHERE UserID = ?", passenger.Email, passengerId)
		if err2 != nil {
			panic(err2.Error())
		}
		w.WriteHeader(http.StatusOK)

		//return the updated passenger details
		results, err := db.Query("SELECT p.FirstName, p.LastName, p.PhoneNo, u.Email, u.Password from Passenger p inner join User u on p.PassengerID = u.UserID where p.PassengerId = ?", passengerId)
		if err != nil {
			panic(err.Error())
		}
		for results.Next() {
			err = results.Scan(&loggedInUser.FirstName, &loggedInUser.LastName, &loggedInUser.PhoneNo, &loggedInUser.Email, &loggedInUser.Password)
			if err != nil {
				panic(err.Error())
			}
		}
		LoggedInUser := LoggedInUser{
			UserID:    loggedInUser.UserID,
			UserType:  loggedInUser.UserType,
			Email:     loggedInUser.Email,
			Password:  loggedInUser.Password,
			FirstName: loggedInUser.FirstName,
			LastName:  loggedInUser.LastName,
			PhoneNo:   loggedInUser.PhoneNo,
			IcNO:      "",
			LicenseNo: "",
		}
		json.NewEncoder(w).Encode(LoggedInUser)
		return
	}
}

func GetDriver(w http.ResponseWriter, r *http.Request) {


	db, err := sql.Open("mysql", sqlDb)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	driverId := vars["driverId"]

	var loggedInUser LoggedInUser
	loggedInUser.UserID = driverId
	loggedInUser.UserType = "driver"

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else if r.Method == "GET" {

		results, err := db.Query("SELECT d.FirstName, d.LastName, d.PhoneNo, d.IcNo, d.LicenseNo, u.Email, u.Password from Driver d inner join User u on d.DriverID = u.UserID where d.DriverId = ?", driverId)
		if err != nil {
			panic(err.Error())
		}
		for results.Next() {
			err = results.Scan(&loggedInUser.FirstName, &loggedInUser.LastName, &loggedInUser.PhoneNo, &loggedInUser.IcNO, &loggedInUser.LicenseNo, &loggedInUser.Email, &loggedInUser.Password)
			if err != nil {
				panic(err.Error())
			}
		}
		LoggedInUser := LoggedInUser{
			UserID:    loggedInUser.UserID,
			UserType:  loggedInUser.UserType,
			Email:     loggedInUser.Email,
			Password:  loggedInUser.Password,
			FirstName: loggedInUser.FirstName,
			LastName:  loggedInUser.LastName,
			PhoneNo:   loggedInUser.PhoneNo,
			IcNO:      loggedInUser.IcNO,
			LicenseNo: loggedInUser.LicenseNo,
		}
		json.NewEncoder(w).Encode(LoggedInUser)
		return

	} else if r.Method == "POST" {
		// update driver details
		reqBody, _ := ioutil.ReadAll(r.Body)
		var driver Driver
		json.Unmarshal(reqBody, &driver)

		// update driver details in database
		_, err := db.Exec("UPDATE Driver SET FirstName = ?, LastName = ?, PhoneNo = ?, LicenseNo = ? WHERE DriverId = ?", driver.FirstName, driver.LastName, driver.PhoneNo, driver.LicenseNo, driverId)
		if err != nil {
			panic(err.Error())
		}
		// update user details in database

		_, err2 := db.Exec("UPDATE User SET Email = ? WHERE UserID = ?", driver.Email, driverId)
		if err2 != nil {
			panic(err2.Error())
		}
		w.WriteHeader(http.StatusOK) //200

		//return the updated driver details
		results, err := db.Query("SELECT d.FirstName, d.LastName, d.PhoneNo, d.IcNo, d.LicenseNo, u.Email, u.Password from Driver d inner join User u on d.DriverID = u.UserID where d.DriverId = ?", driverId)
		if err != nil {
			panic(err.Error())
		}
		for results.Next() {
			err = results.Scan(&loggedInUser.FirstName, &loggedInUser.LastName, &loggedInUser.PhoneNo, &loggedInUser.IcNO, &loggedInUser.LicenseNo, &loggedInUser.Email, &loggedInUser.Password)
			if err != nil {
				panic(err.Error())
			}
		}
		LoggedInUser := LoggedInUser{
			UserID:    loggedInUser.UserID,
			UserType:  loggedInUser.UserType,
			Email:     loggedInUser.Email,
			Password:  loggedInUser.Password,
			FirstName: loggedInUser.FirstName,
			LastName:  loggedInUser.LastName,
			PhoneNo:   loggedInUser.PhoneNo,
			IcNO:      loggedInUser.IcNO,
			LicenseNo: loggedInUser.LicenseNo,
		}
		json.NewEncoder(w).Encode(LoggedInUser)
		return
	}
}
