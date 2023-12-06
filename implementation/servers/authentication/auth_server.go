package main

import (
	"authentication/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var store *sessions.CookieStore

// The form values are only available to the login handler
// Verify access to the session from all handlers with proper session name
// Find a way to share the user who wants to log in amongst handlers (Probably via a shared user variable in types.go?)

func initContainer() {
	// Userlist and secret key initialization

	// time.Sleep(5000 * time.Millisecond)
	utils.CreateUserObjectFromDB()

	filePath := "/etc/profile.d/session_secret.sh"
	sessionSecretKey := utils.InitSecretKey(filePath)
	store = sessions.NewCookieStore([]byte(sessionSecretKey))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login handler hello")

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Password authentication
		if !utils.AuthenticateUser(username, password) {
			fmt.Println("Password authentication failed, redirecting to /error")
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}

		sessionName := utils.GenerateUniqueSessionName(username)
		session, _ := store.Get(r, sessionName)

		// Check fingerprint after successful password authentication
		fingerprintAuthenticator := &utils.MockFingerprintAuthenticator{}
		if fingerprintAuthenticator.CheckFingerprint(username) {
			// MFA completed, redirect to welcome page
			fmt.Println("MFA completed, redirect to /welcome")
			session.Values["authenticated"] = true
			session.Save(r, w)
			http.Redirect(w, r, "/welcome", http.StatusSeeOther)
			return
		} else {
			fmt.Println("Fingerprint check failed, redirect to /error")
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
	}

	// Display the login form.
	http.ServeFile(w, r, "index.html")
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome handler hello")
	http.ServeFile(w, r, "templates/welcome.html")
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	// Handle error logic and log the failed attempt
	fmt.Println("Error handler hello")
	http.Error(w, "Authentication error", http.StatusUnauthorized)
}

func main() {

	initContainer()

	// Create a router and handle all routes
	r := mux.NewRouter()

	r.HandleFunc("/login", loginHandler).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/welcome", welcomeHandler)
	r.HandleFunc("/error", errorHandler)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))
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
