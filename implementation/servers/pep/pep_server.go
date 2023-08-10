package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	w.Write([]byte("This is my website!\n"))
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	w.Write([]byte("Hello, HTTP!\n"))
}

func waitForDBReady(maxRetries int, retryInterval time.Duration) error {

	time.Sleep(time.Second * 5)

	connectionString := "user=angelos password=example host=device_container dbname=ligma_db sslmode=disable port=5432"
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connectionString)
		if err == nil {
			fmt.Println("Database is ready!")
			return nil
		}

		fmt.Printf("Attempt %d: Database connection failed - %s\n", i+1, err)
		time.Sleep(retryInterval)
	}

	if db == nil {
		return fmt.Errorf("max retries reached database is still not ready")
	}

	return nil
}

func conToSql() {
	maxRetries := 10
	retryInterval := time.Second * 5

	fmt.Println("Waiting for the database to be ready...")
	if err := waitForDBReady(maxRetries, retryInterval); err != nil {
		log.Fatal(err)
		return
	}

	// Construct the database connection string
	connectionString := "user=angelos password=example host=device_container dbname=ligma_db sslmode=disable port=5432"

	var db *sql.DB
	var err error

	fmt.Println("Connecting to the database...")
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	// Perform database operations
	rows, err := db.Query("SELECT ligma_id, username FROM ligma")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	fmt.Println("Connected to database!")

	// Process the query results
	for rows.Next() {
		var ligmaID int
		var username string
		err := rows.Scan(&ligmaID, &username)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ligmaID, username)
	}
}

func main() {

	conToSql()

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)

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

//	Create a server and keep alive
//	Wait for user connection with a device certificate
//	Connect to the device certificate storage
//	Verify certificate validity
//	On success, redirect user to authorization server
//	Wait for authorization token response (or error)
//	Send user information to PE and await response
//	On success, connect to the requested resource
//	Fetch the resource and return it to the user
//	Close the connection
