package utils

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
)

func GetTLSConfig(host, caCertFile string, certOpt tls.ClientAuthType) *tls.Config {
	var caCert []byte
	var err error
	var caCertPool *x509.CertPool
	if certOpt > tls.RequestClientCert {
		caCert, err = ioutil.ReadFile(caCertFile)
		if err != nil {
			log.Fatal("Error opening cert file", caCertFile, ", error ", err)
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	}

	return &tls.Config{
		ServerName: host,
		// ClientAuth: tls.NoClientCert,				// Client certificate will not be requested and it is not required
		// ClientAuth: tls.RequestClientCert,			// Client certificate will be requested, but it is not required
		// ClientAuth: tls.RequireAnyClientCert,		// Client certificate is required, but any client certificate is acceptable
		// ClientAuth: tls.VerifyClientCertIfGiven,		// Client certificate will be requested and if present must be in the server's Certificate Pool
		// ClientAuth: tls.RequireAndVerifyClientCert,	// Client certificate will be required and must be present in the server's Certificate Pool
		ClientAuth: certOpt,
		ClientCAs:  caCertPool,
		MinVersion: tls.VersionTLS12, // TLS versions below 1.2 are considered insecure - see https://www.rfc-editor.org/rfc/rfc7525.txt for details
	}
}
func InitSecretKey(filePath string) string {
	// Read the secret key for the sessions. Keys are periodically rotated on the server with a cronjob

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	sessionSecretKey := ExtractSessionSecretKey(string(content))
	if sessionSecretKey == "" {
		fmt.Println("SESSION_SECRET_KEY not found in the file")
		panic("SESSION_SECRET_KEY environment variable not set")
	}

	return sessionSecretKey
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

func AuthenticateUser(username, password string) bool {

	for _, user := range Users {
		if user.Username == username && user.Password == password {
			return true
		}
	}
	return false
}

func FindUser(users []User, targetUsername string) *User {
	for i := range users {
		if users[i].Username == targetUsername {
			return &users[i]
		}
	}
	return nil
}

func CreateUserObjectFromDB(username string) {

	maxRetries := 10
	retryInterval := time.Second * 2

	if err := waitForDBReady(maxRetries, retryInterval); err != nil {
		log.Fatal(err)
	}

	// connectionString := "user=angelos password=example host=user_auth_db dbname=auth_db sslmode=disable port=5432" // For containers

	//Connect to the database
	connectionString := "user=angelos password=example host=localhost dbname=auth_db sslmode=disable port=5432"

	var db *sql.DB
	var err error

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Perform database operations on tables
	// Should probably create views for users login info and authorization info seperately
	rows, err := db.Query("SELECT * FROM Users WHERE username=$1;", username)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Parse the appropriate fields in database and create a User object
	for rows.Next() {
		new_user := new(User)
		err := rows.Scan(&new_user.ID, &new_user.Username, &new_user.Password, (*pq.StringArray)(&new_user.Groups), &new_user.TrustLevel)
		if err != nil {
			log.Fatal(err)
		}
		Users = append(Users, *new_user)
	}
}

func waitForDBReady(maxRetries int, retryInterval time.Duration) error {

	connectionString := "user=angelos password=example host=device_container dbname=ligma_db sslmode=disable port=5432"
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			fmt.Printf("Attempt %d: Database connection failed - %s\n", i+1, err)
			time.Sleep(retryInterval)
		}
	}
	if db == nil {
		return fmt.Errorf("max retries reached database is still not ready")
	}

	return nil
}
