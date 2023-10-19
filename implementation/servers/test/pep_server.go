package main

import (
	"fmt"
	"net/http"
	"os"

	utils "test/utils"

	_ "github.com/lib/pq"
)

func authenticateUser(w http.ResponseWriter, r *http.Request) {

	if utils.AuthenticateUserMain(w, r) {
		// Handle successful authentication response here
		// Send auth token to user
		http.Redirect(w, r, "/success", http.StatusSeeOther)
	} else {
		// Handle unsuccessful authentication response here
		// Probably add a counter and raise an alarm or log the attempt or sth
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func getSuccess(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/success.html")
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
