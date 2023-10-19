package utilities

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

type FormData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func waitForDBReady(maxRetries int, retryInterval time.Duration) error {

	time.Sleep(time.Second * 2)

	connectionString := "user=angelos password=example host=device_container dbname=ligma_db sslmode=disable port=5432"
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			fmt.Printf("Attempt %d: Database connection failed - %s\n", i+1, err)
			time.Sleep(retryInterval)
		}
		// fmt.Println("Database is ready!")
		return nil
	}

	if db == nil {
		return fmt.Errorf("max retries reached database is still not ready")
	}

	return nil
}

func AuthenticateUserMain(w http.ResponseWriter, r *http.Request) bool {

	formData := _authenticateUserParseJSON(w, r)
	rows := _authenticateUserConToDB()
	check := _authenticateUserChecks(formData, rows)

	return check
}

func _authenticateUserParseJSON(w http.ResponseWriter, r *http.Request) FormData {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
	}

	name := r.FormValue("name")
	password := r.FormValue("password")

	data := FormData{
		Name:     name,
		Password: password,
	}

	// Send JSON response
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(jsonData)

	return data
}

func _authenticateUserConToDB() *sql.Rows {
	// fmt.Println("Connecting to sql database")

	maxRetries := 10
	retryInterval := time.Second * 2

	// fmt.Println("Waiting for the database to be ready...")
	if err := waitForDBReady(maxRetries, retryInterval); err != nil {
		log.Fatal(err)
	}

	// Construct the database connection string
	connectionString := "user=angelos password=example host=device_container dbname=ligma_db sslmode=disable port=5432"

	var db *sql.DB
	var err error

	// fmt.Println("Connecting to the database...")
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Perform database operations
	rows, err := db.Query("SELECT username FROM ligma")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Connected to database!")

	return rows
}

func _authenticateUserChecks(data FormData, rows *sql.Rows) bool {
	defer rows.Close()
	if data.Name == "" {
		return false
	}

	// Process the query results
	for rows.Next() {
		var db_username string
		err := rows.Scan(&db_username)
		if err != nil {
			log.Fatal(err)
			return false
		}
		if data.Name == db_username {
			return true
		}
	}
	return false
}
