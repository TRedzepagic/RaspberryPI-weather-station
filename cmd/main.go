package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// SensorValues structure that contains all the individual measured values
type SensorValues struct {
	ID          string `json:"id"`
	Temperature string `json:"temperature"`
	Pressure    string `json:"pressure"`
	Altitude    string `json:"altitude"`
	Time        string `json:"time"`
}
type dbHandler struct {
	db *sql.DB
}

type dbconfig struct {
	DbDriver string `json:"dbdriver"`
	DbUser   string `json:"dbuser"`
	DbPass   string `json:"dbpass"`
	DbName   string `json:"dbname"`
}

func dbConn() *dbHandler {
	data, err := ioutil.ReadFile("./dbconf.json")

	if err != nil {
		fmt.Println("error opening configuration", err.Error())
	}

	var databaseConfig dbconfig
	err = json.Unmarshal(data, &databaseConfig)
	if err != nil {
		fmt.Println("unmarshalling error: ", err.Error())
	}

	dbDriver := databaseConfig.DbDriver
	dbUser := databaseConfig.DbUser
	dbPass := databaseConfig.DbPass
	dbName := databaseConfig.DbName

	// Check
	fmt.Println(databaseConfig.DbDriver)
	fmt.Println(databaseConfig.DbUser)
	fmt.Println(databaseConfig.DbPass)
	fmt.Println(databaseConfig.DbName)

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		fmt.Println(err)
	}

	var DataBase dbHandler
	DataBase.db = db
	return &DataBase
}

func getTime() string {
	time := fmt.Sprint(time.Now().Format("15:04:05"))
	return time
}

func (dbHandler *dbHandler) getReadings(w http.ResponseWriter, req *http.Request) {
	var sensVal SensorValues
	var readingSlice []SensorValues

	rows, err := dbHandler.db.Query("SELECT id, Temperature,Pressure,Altitude,Time FROM READINGS ORDER BY Time DESC LIMIT 10")
	defer rows.Close()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		err := rows.Scan(&sensVal.ID, &sensVal.Temperature, &sensVal.Pressure, &sensVal.Altitude, &sensVal.Time)
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusInternalServerError)
			return
		}
		readingSlice = append(readingSlice, sensVal)
	}
	bytes, _ := json.MarshalIndent(readingSlice, "", " ")
	fmt.Fprintf(w, string(bytes))
}

func (dbHandler *dbHandler) postReading(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	var sensVal SensorValues
	err = json.Unmarshal(body, &sensVal)

	if sensVal.Temperature == "" || sensVal.Altitude == "" || sensVal.Pressure == "" {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusBadRequest)
		fmt.Println("Bad request, empty fields")
		return
	}

	fmt.Print(sensVal)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sensVal.Time = getTime()
	stmt, err := dbHandler.db.Prepare("INSERT INTO READINGS(Temperature, Pressure, Altitude, Time) VALUES(?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(sensVal.Temperature, sensVal.Pressure, sensVal.Altitude, sensVal.Time)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}
}
func (dbHandler *dbHandler) deleteReading(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if params["id"] == "" {
		fmt.Println("No ID")
		http.Error(w, "No ID", http.StatusNoContent)
		return
	}

	stmt, err := dbHandler.db.Prepare(("DELETE FROM READINGS WHERE id=?"))
	defer stmt.Close()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	result, err := stmt.Exec(params["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}
}

func (dbHandler *dbHandler) updateReading(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if params["id"] == "" {
		log.Println("No ID")
		http.Error(w, "No ID", http.StatusNoContent)
		return
	}
	fmt.Println("ID PARAM: " + params["id"])
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	var sensVal SensorValues
	err = json.Unmarshal(body, &sensVal)
	if sensVal.Temperature == "" || sensVal.Altitude == "" || sensVal.Pressure == "" {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusBadRequest)
		fmt.Println("Bad request, empty fields")
		return
	}
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	sensVal.Time = getTime()
	result, err := dbHandler.db.Exec("UPDATE READINGS SET Temperature = ?, Pressure = ?, Altitude = ?, Time = ? where id = ?", sensVal.Temperature, sensVal.Pressure, sensVal.Altitude, sensVal.Time, params["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}

}

// AccessControl middleware function inserts access control parameters
func AccessControl(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "3.5") //firefox
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if req.Method == "OPTIONS" {
			return
		}
		handler.ServeHTTP(w, req)
	}

}

func main() {
	DBconn := dbConn()
	defer DBconn.db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/getReadings", AccessControl(DBconn.getReadings))
	router.HandleFunc("/postReading", AccessControl(DBconn.postReading))
	router.HandleFunc("/deleteReading/{id}", AccessControl(DBconn.deleteReading))
	router.HandleFunc("/updateReading/{id}", AccessControl(DBconn.updateReading))
	log.Println("Server started...")
	http.ListenAndServe(":8090", router)

}
