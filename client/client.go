package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var certificatePathPrefix = "/home/angelos/Desktop/Thesis_Stuff/certificates/out/"

func main() {
	help := flag.Bool("help", false, "Optional, prints usage info")
	// Hosts should become the container names in the future
	srvhost := flag.String("srvhost", "localhost", "The server's host name")
	port := flag.String("port", "443", "The https port, defaults to 443")
	caCertFile := flag.String("cacert", certificatePathPrefix+"ThesisCA.crt", "Required, the name of the CA that signed the server's certificate")
	clientCertFile := flag.String("clientcert", certificatePathPrefix+"client.crt", "Required, the name of the client's certificate file")
	clientKeyFile := flag.String("clientkey", certificatePathPrefix+"client.key", "Required, the file name of the clients's private key file")
	flag.Parse()

	usage := `usage:
	
client -clientcert <clientCertificateFile> -cacert <caFile> -clientkey <clientPrivateKeyFile> [-host <srvHostName> -help]
	
Options:
  -help       Optional, Prints this message
  -srvhost    Optional, the server's hostname, defaults to 'localhost'
  -clientcert Optional, the name the clients's certificate file
  -clientkey  Optional, the name the client's key certificate file
  -cacert     Required, the name of the CA that signed the server's certificate
 `

	if *help {
		fmt.Println(usage)
		return
	}
	if *caCertFile == "" {
		log.Fatalf("caCert is required but missing:\n%s", usage)
	}

	var cert tls.Certificate
	var err error
	if *clientCertFile != "" && *clientKeyFile != "" {
		cert, err = tls.LoadX509KeyPair(*clientCertFile, *clientKeyFile)
		if err != nil {
			log.Fatalf("Error creating x509 keypair from client cert file %s and client key file %s", *clientCertFile, *clientKeyFile)
		}
	}

	log.Printf("CAFile: %s", *caCertFile)
	caCert, err := ioutil.ReadFile(*caCertFile)
	if err != nil {
		log.Fatalf("Error opening cert file %s, Error: %s", *caCertFile, err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	}

	client := http.Client{Transport: t, Timeout: 15 * time.Second}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s:%s", *srvhost, *port), bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Fatalf("unable to create http request due to error %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		switch e := err.(type) {
		case *url.Error:
			log.Fatalf("url.Error received on http request: %s", e)
		default:
			log.Fatalf("Unexpected error received: %s", err)
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("unexpected error reading response body: %s", err)
	}

	fmt.Printf("\nResponse from server: \n\tHTTP status: %s\n\tBody: %s\n", resp.Status, body)
}
