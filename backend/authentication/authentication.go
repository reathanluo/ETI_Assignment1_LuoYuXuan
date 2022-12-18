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

type Credential struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type LoggedInUser struct {
	UserID    int    `json:"UserID"`
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

	router.HandleFunc("/api/v1/signup/{userType}", SignUpUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/login", Login).Methods("POST", "OPTIONS")

	fmt.Println("Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", router))
}

func SignUpUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", sqlDb)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	params := mux.Vars(r)
	userType := params["userType"]
	req, _ := ioutil.ReadAll(r.Body)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else if r.Method == "POST" {
		if userType == "passenger" {
			var p Passenger
			json.Unmarshal(req, &p)

			isExist := emailExist(db, p.Email)
			if isExist {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode("Email exists! Please use other email.")
				return
			} else if p.Email == "" || p.FirstName == "" || p.LastName == "" || p.PhoneNo == "" || p.Password == "" {
				w.WriteHeader(http.StatusNotAcceptable)
				json.NewEncoder(w).Encode("Empty input!")
				return
			} else {
				rowsinserted := insertIntoUserAndPassengerDB(userType, p)
				if rowsinserted == 1 {
					w.WriteHeader(http.StatusAccepted)
					json.NewEncoder(w).Encode("Passenger created")
					return
				} else {
					w.WriteHeader(http.StatusExpectationFailed)
					json.NewEncoder(w).Encode("Failed create new passenger")
					return
				}
			}
		} else if userType == "driver" {
			var d Driver
			json.Unmarshal(req, &d)

			isExist := emailExist(db, d.Email)
			if isExist {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode("Email exists! Please use other email.")
				return
			} else {
				rowsinserted := insertIntoUserAndDriverDB(userType, d)
				if rowsinserted == 1 {
					w.WriteHeader(http.StatusAccepted)
					json.NewEncoder(w).Encode("Driver created")
					return
				} else {
					w.WriteHeader(http.StatusExpectationFailed)
					json.NewEncoder(w).Encode("Failed create new driver")
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credential
	var user User

	db, err := sql.Open("mysql", sqlDb)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		exists := emailExist(db, credentials.Email)
		if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Email Not Found!")
			return
		}

		checked := checkPassword(db, credentials.Email, credentials.Password)
		if !checked {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Wrong password!")
			return
		}

		user, err = retrieveUserInfo(db, credentials.Email)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			panic(err.Error())
		}

		if user.UserType == "passenger" {
			passenger, err := retrievePassengerInfo(db, credentials.Email)
			if err != nil {
				panic(err.Error())
			}
			LoggedInUser := LoggedInUser{
				UserID:    user.UserID,
				UserType:  user.UserType,
				Email:     passenger.Email,
				Password:  passenger.Password,
				FirstName: passenger.FirstName,
				LastName:  passenger.LastName,
				PhoneNo:   passenger.PhoneNo,
				IcNO:      "",
				LicenseNo: "",
			}
			json.NewEncoder(w).Encode(LoggedInUser)
			return

		} else if user.UserType == "driver" {
			driver, err := retrieveDriverInfo(db, credentials.Email)
			if err != nil {
				panic(err.Error())
			}
			LoggedInUser := LoggedInUser{
				UserID:    user.UserID,
				UserType:  user.UserType,
				Email:     driver.Email,
				Password:  driver.Password,
				FirstName: driver.FirstName,
				LastName:  driver.LastName,
				PhoneNo:   driver.PhoneNo,
				IcNO:      driver.IcNO,
				LicenseNo: driver.LicenseNo,
			}
			json.NewEncoder(w).Encode(LoggedInUser)
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func insertIntoUserAndPassengerDB(userType string, p Passenger) int {
	db, err := sql.Open("mysql", sqlDb)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insertUser, err := db.Exec("INSERT INTO User (UserType, Email, Password) VALUES (?,?,?)", userType, p.Email, p.Password)
	if err != nil {
		panic(err.Error())
	}
	userID, err2 := insertUser.LastInsertId()
	if err2 != nil {
		panic(err.Error())
	}

	insertPassenger, err := db.Exec("INSERT INTO Passenger (PassengerID,FirstName,LastName,PhoneNo) VALUES (?,?,?,?)", userID, p.FirstName, p.LastName, p.PhoneNo)
	if err != nil {
		panic(err.Error())
	}

	insertRows, err := insertPassenger.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return int(insertRows)
}

func insertIntoUserAndDriverDB(userType string, d Driver) int {
	db, err := sql.Open("mysql", sqlDb)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insertUser, err := db.Exec("INSERT INTO User (UserType, Email, Password) VALUES (?,?,?)", userType, d.Email, d.Password)
	if err != nil {
		panic(err.Error())
	}
	userID, err2 := insertUser.LastInsertId()
	if err2 != nil {
		panic(err.Error())
	}

	insertPassenger, err := db.Exec("INSERT INTO Driver VALUES (?,?,?,?,?,?,?)", userID, d.FirstName, d.LastName, d.PhoneNo, d.IcNO, d.LicenseNo, 0)
	if err != nil {
		panic(err.Error())
	}

	insertRows, err := insertPassenger.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return int(insertRows)
}

func emailExist(db *sql.DB, email string) bool {
	result,err := db.Query("SELECT UserID from User WHERE Email = ?", email)
	if err != nil {
		panic(err.Error())
	}
	if result.Next() {
		return true
	}
	return false
}

func checkPassword(db *sql.DB, email string, password string) bool {
	results, err := db.Query("SELECT UserID from User WHERE Email = ? and password = ?", email, password)
	if err != nil {
		panic(err.Error())
	}
	if results.Next() {
		return true
	}
	return false
}

func retrieveUserInfo(db *sql.DB, email string) (User, error) {
	results, err := db.Query("SELECT * from User WHERE Email = ?", email)
	if err != nil {
		panic(err.Error())
	}
	var user User
	for results.Next() {
		err = results.Scan(&user.UserID, &user.UserType, &user.Email, &user.Password)
		if err != nil {
			panic(err.Error())
		}
	}
	return user, err
}

func retrievePassengerInfo(db *sql.DB, email string) (Passenger, error) {
	results, err := db.Query("SELECT p.FirstName, p.LastName, p.PhoneNo, u.Email, u.Password FROM User u INNER JOIN Passenger p ON u.UserID = p.PassengerID WHERE Email = ?", email)

	if err != nil {
		panic(err.Error())
	}
	var p Passenger
	for results.Next() {
		err = results.Scan(&p.FirstName, &p.LastName, &p.PhoneNo, &p.Email, &p.Password)
		if err != nil {
			panic(err.Error())
		}
	}
	return p, err
}

func retrieveDriverInfo(db *sql.DB, email string) (Driver, error) {
	results, err := db.Query("SELECT d.FirstName, d.LastName, d.PhoneNo, u.Email, d.IcNO,d.LicenseNo, u.Password FROM User u INNER JOIN Driver d ON u.UserID = d.DriverID WHERE Email = ?", email)

	if err != nil {
		panic(err.Error())
	}
	var d Driver
	for results.Next() {
		err = results.Scan(&d.FirstName, &d.LastName, &d.PhoneNo, &d.Email, &d.IcNO, &d.LicenseNo, &d.Password)
		if err != nil {
			panic(err.Error())
		}
	}
	return d, err
}
