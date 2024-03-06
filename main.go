package main

import (
	"database/sql"
	"fmt"
	"log"
  "strconv"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
  var err error
	db, err = sql.Open("sqlite3", "/home/aday/databases/myfirstdb.db")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	defer db.Close()

	http.HandleFunc("/moisture/", handleGetById)
  http.HandleFunc("/moisture", handleGetAll)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGetAll(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, plant_name, moisture_level FROM moisture_readings")
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer rows.Close() 

	writeData(w, rows)
}

func handleGetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/moisture/"):])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT id, plant_name, moisture_level FROM moisture_readings WHERE id = ?", id)
	writeData(w, row) 
}

func writeData(w http.ResponseWriter, data interface{}) {
	switch rows := data.(type) {
	case *sql.Rows:
		for rows.Next() {
			var id int
			var plantName string
			var moistureLevel int
			rows.Scan(&id, &plantName, &moistureLevel)
			fmt.Fprintf(w, "ID: %d, Plant: %s, Moisture: %d\n", id, plantName, moistureLevel)
		}
	case *sql.Row:
		var id int
		var plantName string
		var moistureLevel int
		rows.Scan(&id, &plantName, &moistureLevel)
		fmt.Fprintf(w, "ID: %d, Plant: %s, Moisture: %d\n", id, plantName, moistureLevel)
	default:
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}
