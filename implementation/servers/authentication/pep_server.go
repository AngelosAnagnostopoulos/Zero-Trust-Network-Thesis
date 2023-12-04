package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	utils "authentication/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

// TODO: Create the encryption keys for the CookieStore
var (
	store = sessions.NewCookieStore([]byte(secretKey))
)

// User represents a user with authentication credentials.
var users = []User{}

type User struct {
	Username string
	Password string
}
type FormData struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type FingerprintAuthenticator interface {
	utils.CheckFingerprint(username string) bool
}

type MockFingerprintAuthenticator struct{}

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

func CreateUserObjectFromDB() {

	maxRetries := 10
	retryInterval := time.Second * 2

	if err := waitForDBReady(maxRetries, retryInterval); err != nil {
		log.Fatal(err)
	}

	//Connect to the database
	connectionString := "user=angelos password=example host=test_database dbname=test_db sslmode=disable port=5432"

	var db *sql.DB
	var err error

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Perform database operations on tables
	rows, err := db.Query("SELECT username,passwd_hash FROM Users;")
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	// Parse the appropriate fields in database and create a User object
	for rows.Next() {
		new_user := new(User)
		err := rows.Scan(&new_user.Username, &new_user.Password)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, *new_user)
	}
}

func AuthenticateUser(username, password string) (*User, error) {

	for _, user := range users {
		if user.Username == username && user.Password == password {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("Authentication failed")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is already authenticated
	session, _ := store.Get(r, "session-name")
	// fmt.Println(session.Values["authenticated"])
	if authenticated, ok := session.Values["authenticated"].(bool); ok && authenticated {
		// If already authenticated, redirect to the welcome page
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
		fmt.Println("Ne")
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := AuthenticateUser(username, password)
		if err != nil {
			// http.Error(w, "Authentication failed", http.StatusUnauthorized)
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		// Check fingerprint after successful authentication
		fingerprintAuthenticator := &MockFingerprintAuthenticator{}
		if fingerprintAuthenticator.CheckFingerprint(username) {
			session.Values["authenticated"] = true
			session.Save(r, w)

			fmt.Fprintf(w, "Authentication successful. Welcome, %s! Fingerprint check passed.", user.Username)
			return
		} else {
			// Fingerprint check failed
			http.Error(w, "Fingerprint check failed", http.StatusUnauthorized)
			return
		}
	}

	// Display the login form.
	http.ServeFile(w, r, "static/index.html")
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome! You are already authenticated.")
}

func ExtractSessionSecretKey(content string) string {
	// Find the line that contains "SESSION_SECRET_KEY"
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "export SESSION_SECRET_KEY=") {
			return strings.TrimPrefix(line, "export SESSION_SECRET_KEY=")
		}
	}
	return ""
}

func main() {

	// Initialize the users list. The search done on the auth function is O(n) and can probably be optimized w/ a hash table
	createUserObjectFromDB()

	// Read the secret key for the sessions. Keys are periodically rotated on the server with a cronjob
	filePath := "/etc/profile.d/session_secret.sh"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	sessionSecretKey := extractSessionSecretKey(string(content))
	if sessionSecretKey == "" {
		fmt.Println("SESSION_SECRET_KEY not found in the file")
		panic("SESSION_SECRET_KEY environment variable not set")
	}

	fmt.Printf("SESSION_SECRET_KEY: %s\n", sessionSecretKey)

	// Create a router and handle all routes
	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/welcome", welcomeHandler).Methods(http.MethodGet)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))

	http.Handle("/", r)

	// Create a channel to synchronize server startup
	serverReady := make(chan struct{})

	// Run the server in a goroutine
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
		close(serverReady)
	}()

	// Wait for the server to be ready
	<-serverReady

}
