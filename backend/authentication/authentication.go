package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	IcNO      string `json:"IcNumber"`
	LicenseNo string `json:"LicenseNo"`
	Password  string `json:"Password"`
}

var sqlDb = "user:password@tcp(127.0.0.1:3307)/trip_db"

func main() {
	router := mux.NewRouter()

	db, err := sql.Open("mysql", sqlDb)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	router.HandleFunc("/api/auth/signup/{userType}", SignUpUser).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/auth/logout", Logout)
}

func emailExist(db *sql.DB, email string) bool {
	results, err := db.Query("SELECT UserID from User WHERE Email = ?", email)
	if err != nil {
		panic(err.Error())
	}
	var id int
	for results.Next() {
		err = results.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		return true

	}
	return false
}

func checkPassword(db *sql.DB, email string, password string) bool {
	results, err := db.Query("SELECT UserID from User WHERE Email = ? and password = ?", email, password)
	if err != nil {
		panic(err.Error())
	}
	var id int
	for results.Next() {
		err = results.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
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
		err = results.Scan(&d.FirstName, &d.LastName, &d.PhoneNo, &d.Email, d.IcNO, d.LicenseNo, &d.Password)
		if err != nil {
			panic(err.Error())
		}
	}
	return d, err
}

func SignUpUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", sqlDb)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	params := mux.Vars(r)
	userType := params["userType"]
	reqBody, _ := ioutil.ReadAll(r.Body)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else if r.Method == "POST" {
		if userType == "passenger" {
			var p Passenger
			json.Unmarshal(reqBody, &p)

			isExist := emailExist(db, p.Email)
			if isExist {
				w.WriteHeader(http.StatusConflict) //409
				json.NewEncoder(w).Encode("Email exists! Please use other email.")
				return
			} else if p.Email == "" || p.FirstName == "" || p.LastName == "" || p.PhoneNo == "" || p.Password == "" {
				w.WriteHeader(http.StatusNotAcceptable) //406
				json.NewEncoder(w).Encode("Found empty content")
				return
			} else {
				insertIntoUserAndPassengerDB(userType, p)

				userQueryStatement := fmt.Sprintf(`
		INSERT INTO User(user_type, email_address, password)
		VALUES ('%s', '%s', '%s')`, user_type, passenger.EmailAddress, passenger.Password)
				result, err := db.Exec(userQueryStatement)
				if err != nil {
					panic(err.Error())
					return
				}
				// Get the auto-gen ID
				id, err := result.LastInsertId()
				if err != nil {
					panic(err.Error())
					return
				}

				// Upon getting the ID, now insert into Passenger table
				queryStatement := fmt.Sprintf(`
		INSERT INTO Passenger
		VALUES (%d, '%s', '%s', '%s')`,
					id,
					passenger.FirstName,
					passenger.LastName,
					passenger.MobileNumber)

				result2, err := db.Exec(queryStatement)
				if err != nil {
					panic(err.Error())
				}
				rows_affected, err := result2.RowsAffected()
				if err != nil {
					panic(err.Error())
				}
				if rows_affected == 1 {
					w.WriteHeader(http.StatusAccepted) //202
					json.NewEncoder(w).Encode("Passenger Insert Successfully")
					return
				} else {
					w.WriteHeader(http.StatusExpectationFailed) //417
					json.NewEncoder(w).Encode("Error with inserting")
					return
				}
			}

		} else if user_type == "rider" {
			var rider Rider
			json.Unmarshal(reqBody, &rider)

			// check if email exists in the User table
			isExist := checkEmailIsExist(db, rider.EmailAddress)
			if isExist {
				w.WriteHeader(http.StatusConflict) //409
				json.NewEncoder(w).Encode("Duplicated account: " + rider.EmailAddress)
				return
			} else if rider.EmailAddress == "" || rider.FirstName == "" || rider.LastName == "" || rider.MobileNumber == "" || rider.Password == "" || rider.IcNumber == "" || rider.CarLicNumber == "" {
				w.WriteHeader(http.StatusNotAcceptable) //406
				json.NewEncoder(w).Encode("Found empty content")
				return
			} else {
				// Insert into User table
				userQueryStatement := fmt.Sprintf(`
		INSERT INTO User(user_type, email_address, password)
		VALUES ('%s', '%s', '%s')`, user_type, rider.EmailAddress, rider.Password)
				result, err := db.Exec(userQueryStatement)
				if err != nil {
					panic(err.Error())
					return
				}
				id, err := result.LastInsertId()
				if err != nil {
					panic(err.Error())
					return
				}

				// Upon getting the ID, now insert into Rider table
				queryStatement := fmt.Sprintf(`
		INSERT INTO Rider
		VALUES (%d, '%s', '%s', '%s', '%s', '%s')`,
					id,
					rider.FirstName,
					rider.LastName,
					rider.MobileNumber,
					rider.IcNumber,
					rider.CarLicNumber)

				results, err := db.Exec(queryStatement)
				if err != nil {
					panic(err.Error())
				}
				rows_affected, err := results.RowsAffected()
				if err != nil {
					panic(err.Error())
				}
				if rows_affected == 1 {
					w.WriteHeader(http.StatusAccepted) //202
					json.NewEncoder(w).Encode("Rider Insert Successfully")
					return
				} else {
					w.WriteHeader(http.StatusExpectationFailed) //417
					json.NewEncoder(w).Encode("Error with inserting")
					return
				}
			}

		} else { // If the query param is not "Passenger" nor "Rider"
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}

	} else { // Other req methods
		w.WriteHeader(http.StatusNotFound) //404
		return
	}
}

func insertIntoUserAndPassengerDB(userType string, p Passenger) {
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

	insertPassenger, err := db.Exec("INSERT INTO Passenger")

}
