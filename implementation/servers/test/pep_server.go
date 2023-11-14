package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var (
	store = sessions.NewCookieStore([]byte("your-secret-key"))
)

// User represents a user with authentication credentials.
type User struct {
	Username string
	Password string
}

func authenticateUser(w http.ResponseWriter, r *http.Request) {
	if AuthenticateUserMain(w, r) {
		// Handle successful authentication response here
		// Send auth token to user
		http.Redirect(w, r, "/success", http.StatusSeeOther)
	} else {
		fmt.Println("Fail")
		// Handle unsuccessful authentication response here
		// Probably add a counter and raise an alarm or log the attempt or sth
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func getSuccess(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/success.html")
}

type FormData struct {
	Name     string `json:"username"`
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
	check := _authenticateUserConToDB(formData)

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

	username := r.FormValue("username")
	password := r.FormValue("password")

	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
	}

	data := FormData{
		Name:     username,
		Password: password,
	}

	return data
}

func _authenticateUserConToDB(data FormData) bool {

	maxRetries := 10
	retryInterval := time.Second * 2

	if err := waitForDBReady(maxRetries, retryInterval); err != nil {
		log.Fatal(err)
	}

	connectionString := "user=angelos password=example host=test_database dbname=test_db sslmode=disable port=5432"

	var db *sql.DB
	var err error

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Perform database operations
	rows, err := db.Query("SELECT username,passwd_hash FROM Users;")
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	if data.Name == "" {
		return false
	}

	// Process the query results
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.Username, &user.Password)
		if err != nil {
			log.Fatal(err)
			return false
		}
		if data.Name == user.Username && data.Password == user.Password {
			return true
		}
	}
	return false

}

func main() {

	// Serve an HTML page
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// Route handlers
	http.HandleFunc("/authenticate", authenticateUser)
	http.HandleFunc("/success", getSuccess)

	// Create a channel to synchronize server startup
	serverReady := make(chan struct{})

	// Run the server in a goroutine
	go func() {
		err := http.ListenAndServe(":4444", nil)
		if err != nil {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
		close(serverReady)
	}()

	// Wait for the server to be ready
	<-serverReady

}
