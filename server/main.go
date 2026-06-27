package main

/* references:
- go by example: https://gobyexample.com/
- tcp server: https://gobyexample.com/tcp-server
- command line arguments: https://gobyexample.com/command-line-arguments

*/

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type client struct {
	name       string
	connection *net.Conn
}

const max_clients int = 50

var clients_array [max_clients]client
var clients_mutex sync.Mutex

func main() {

	port_number := handeCommandLineParameters(os.Args)
	port_string := fmt.Sprintf(":%d", port_number)

	for _, v := range clients_array {
		v.name = ""
		v.connection = nil
	}

	listener, err := net.Listen("tcp", port_string)
	if err != nil {
		log.Fatal("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()

	hostname, _ := os.Hostname()
	log.Printf("Server ready! host: %s port: %d\n\r", hostname, port_number)

	for { // for-ever

		// Accept an incomming request from a client
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}

		// Allocate a new entry in the array of clients for the new client
		// create a new thread to handle the client
		// 'clients' may be modified by other threads, so it has to be protected with a semaphore

		// Lock so only one goroutine at a time can access clients array
		clients_mutex.Lock()
		for _, v := range clients_array {
			if v.connection != nil { // find and empty slot
				// allocate the new client there
				v.connection = &conn
				log.Println("accepted new connection: ", conn.RemoteAddr().String())

				// Handle the connection in a new goroutine.
				// The loop then returns to accepting, so that
				// multiple connections may be served concurrently.
				go handleClient(conn)
			}
		}
		clients_mutex.Unlock()
	}
}

func handleClient(conn net.Conn) {

	defer conn.Close()

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Read error: %v", err)
		return
	}

	ackMsg := strings.ToUpper(strings.TrimSpace(message))
	response := fmt.Sprintf("ACK: %s\n", ackMsg)
	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Printf("Server write error: %v", err)
	}
}

func handeCommandLineParameters(parameters []string) int {

	log.Println("number of parameters passed: ", len(parameters))

	for i, v := range parameters {
		log.Println("param: ", "[", i, "]: ", v, "\n\r")
	}

	if len(parameters) > 2 {
		log.Println("Error: Too many parameters, the server only needs the port number as parameter")
		os.Exit(2)
	}

	port_number, err := strconv.Atoi(parameters[1])
	if err != nil {
		log.Println("Error: the parameter must be a port number")
		os.Exit(3)
	}
	fmt.Println("port number: ", port_number)

	return port_number
}

func readFirstWord(message string) string {
	// reference: https://pkg.go.dev/strings#SplitN
	firstWord := strings.SplitN(message, " ", 2)
	return firstWord[0]
}
