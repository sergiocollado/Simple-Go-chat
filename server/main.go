package main

/* references:
- tutorials: https://go.dev/doc/tutorial/
- go by example: https://gobyexample.com/
- tcp server: https://gobyexample.com/tcp-server
- command line arguments: https://gobyexample.com/command-line-arguments
- https://www.freecodecamp.org/news/how-to-build-a-production-grade-distributed-chatroom-in-go-full-handbook/
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

	for i := range clients_array {
		clients_array[i].name = ""
		clients_array[i].connection = nil
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
			log.Println("Error accepting connection:", err)
			continue
		}

		//log.Println("accepting connection from: ", conn.RemoteAddr().String())

		// Allocate a new entry in the array of clients for the new client
		// create a new thread(goroutine) to handle the client
		// 'clients' may be modified by other threads(goroutines), so it has to be protected with a semaphore/mutex
		// This design actually is lacking against the C10K problem (https://www.youtube.com/watch?v=L0jMBrCEQNQ)

		// Lock so only one goroutine at a time can access clients array
		clients_mutex.Lock()
		//log.Println("inside mutex")
		for _, v := range clients_array {

			//log.Printf("client %d, name: %s, connection: %p", i, v.name, v.connection)
			if v.connection == nil { // find and empty slot
				// allocate the new client there
				v.connection = &conn
				log.Println("accepted new connection: ", conn.RemoteAddr().String())

				// Handle the connection in a new goroutine.
				// The loop then returns to accepting, so that
				// multiple connections may be served concurrently.
				go handleClient(conn)

				//log.Println("break")
				break
			}
		}
		clients_mutex.Unlock()
		//log.Println("just got outside of the mutex")
	}
}

func handleClient(conn net.Conn) {

	//log.Println("handling connection from: ", conn.RemoteAddr().String())

	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		log.Printf("Read message: %s", message)

		//ackMsg := strings.ToUpper(strings.TrimSpace(message))
		//response := fmt.Sprintf("ACK: %s\n", ackMsg)
		//_, err = conn.Write([]byte(response))
		//if err != nil {
		//	log.Printf("Server write error: %v", err)
		//}
	}
}

func handeCommandLineParameters(parameters []string) int {

	// reference: command line arguments: https://gobyexample.com/command-line-arguments

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
	firstWord := strings.SplitN(strings.TrimSpace(message), " ", 2)
	return firstWord[0]
}
