package main

import "fmt"

func main() {
	//	Create a server and keep it alive
	//	Expect connection from PEP with a username
	//	Retreive password from user and hash it
	//	Connect to user and group database and retreive the stored hash value
	//	Verify or discard the connection and send reply to PEP
	fmt.Println("Hi from authentication!")
}
