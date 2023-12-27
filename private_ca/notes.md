# Create an automated workflow for the private_ca:

Clone the certificate tool "certstrap" and build the project

```
$ git clone https://github.com/square/certstrap
$ cd certstrap
$ go build
$ chmod +x certstrap && export PATH=$PATH:$(pwd)
```

Create a directory for the certificates (All subsequent commands are run from this directory)

```
$ mkdir certificates
```

Create a CA for our project

```
$ certstrap init --common-name "ThesisCA"
```

### Generate encryption keys and certificates
```
$ certstrap request-cert --domain  "pep_server"
$ certstrap request-cert --domain  "client"
```

### Sign the certificates of the server and client wtih the CA
```
$ certstrap sign pep_server --CA ThesisCA
$ certstrap sign client --CA ThesisCA
```

> The certificates are IP bound / Device bound, we create a simple mapping and deliver them to the appropriate containers: device name --> certificate

### Send the files to the appropriate machine/container 
(We just move them to the corresponding directories for testing)

> Send the CA certificates to the server and client respectively
Send the client certificate and key to the client
Send the server certificate and key to the server

