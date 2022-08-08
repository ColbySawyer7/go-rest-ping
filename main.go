package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

// location represents the logical destination of ping operation
type Location struct {
	ID               string  `json:"id"`
	Destination      string  `json:"destination`
	Name             string  `json:"name"`
	Success_Rate     float64 `json:"success_rate"`
	Polling_Interval int     `json:"polling_interval"`
}

//Establish Connection to Database
func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	rows, err := db.Query("SELECT * FROM location")
	if err != nil {
		log.Fatal(err)
	}

	var locations []Location

	for rows.Next() {
		var location Location
		rows.Scan(&location.ID, &location.Destination, &location.Name, &location.Success_Rate, &location.Polling_Interval)
		locations = append(locations, location)
	}

	locationBytes, _ := json.MarshalIndent(locations, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(locationBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	var l Location
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO location (id, destination, name, success_rate, polling_interval) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, l.Destination, l.Name, l.Success_Rate, l.Polling_Interval)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func runPingTest(polling_interval int) float64 {
	var success_rate = 0.0
	return success_rate
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/locations/", GETHandler).Methods("GET")
	router.HandleFunc("/locations/", POSTHandler).Methods("POST")

	fmt.Println("Server at 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
