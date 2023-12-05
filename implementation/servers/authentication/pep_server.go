package main

import (
	"authentication/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var store *sessions.CookieStore

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		_, err := utils.AuthenticateUser(username, password)
		if err != nil {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		sessionName := utils.GenerateUniqueSessionName(username)
		session, _ := store.Get(r, sessionName)

		// Check fingerprint after successful password authentication
		fingerprintAuthenticator := &utils.MockFingerprintAuthenticator{}
		if fingerprintAuthenticator.CheckFingerprint(username) {
			session.Values["authenticated"] = true
			session.Save(r, w)
			// MFA completed, redirect to welcome page
			http.Redirect(w, r, "/welcome", http.StatusSeeOther)
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
	session, _ := store.Get(r, "session-name")
	if authenticated, ok := session.Values["authenticated"].(bool); ok && authenticated {
		// Serve the welcome.html file
		http.ServeFile(w, r, "/welcome.html")
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the login.html file on the root route
	http.ServeFile(w, r, "/login.html")
}

func main() {

	time.Sleep(5000 * time.Millisecond)
	// Initialize the users list. The search done on the auth function is O(n) and can probably be optimized w/ a hash table
	utils.CreateUserObjectFromDB()

	// Read the secret key for the sessions. Keys are periodically rotated on the server with a cronjob
	filePath := "/etc/profile.d/session_secret.sh"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	sessionSecretKey := utils.ExtractSessionSecretKey(string(content))
	if sessionSecretKey == "" {
		fmt.Println("SESSION_SECRET_KEY not found in the file")
		panic("SESSION_SECRET_KEY environment variable not set")
	}
	store = sessions.NewCookieStore([]byte(sessionSecretKey))

	// Create a router and handle all routes
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))

	r.HandleFunc("/login", loginHandler).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/welcome", welcomeHandler).Methods(http.MethodGet)
	r.HandleFunc("/", rootHandler)

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
