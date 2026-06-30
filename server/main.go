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
const code_version string = "0.1"

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
		for index, v := range clients_array {

			//log.Printf("client %d, name: %s, connection: %p", i, v.name, v.connection)
			if v.connection == nil { // find and empty slot
				// allocate the new client there
				v.connection = &conn
				log.Println("accepted new connection: ", conn.RemoteAddr().String())

				// Handle the connection in a new goroutine.
				// The loop then returns to accepting, so that
				// multiple connections may be served concurrently.
				go handleClient(index)

				//log.Println("break")
				break
			}
		}
		clients_mutex.Unlock()
		//log.Println("just got outside of the mutex")
	}
}

func handleClient(index int) {

	//log.Println("handling connection from: ", conn.RemoteAddr().String())

	connection := *clients_array[index].connection

	defer connection.Close()

	reader := bufio.NewReader(connection)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		if checkLEAVE(message) {
			HandleLEAVE(index)
			return // this will kill the goroutine
		} else if !(ifIsCommandExecute(index, message)) {
			// commands have been executied
		} else {
			// HandleBroadcast(index)
			log.Printf("Read message: %s", message)
		}
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

func trimAllSpace(s string) string {
	// reference: https://www.danielmorell.com/blog/how-to-trim-whitespace-from-a-string-in-go#trim_duplicate
	return strings.Join(strings.Fields(s), " ")
}

func checkJOIN(message string) string {
	words := strings.SplitN(trimAllSpace(message), " ", 3)
	if words[0] == "JOIN" {
		return words[1]
	} else {
		return ""
	}
}

func checkLEAVE(message string) bool {
	return "LEAVE" == readFirstWord(message)
}

func checkWHO(message string) bool {
	return "WHO" == readFirstWord(message)
}

func checkHELP(message string) bool {
	return "HELP" == readFirstWord(message)
}

func checkVERSION(message string) bool {
	return "VERSION" == readFirstWord(message)
}

func ifIsCommandExecute(index int, message string) bool {
	name := checkJOIN(message)
	if name != "" {
		HandleJOIN(index, name)
		return true
	}
	if checkWHO(message) {
		HandleWHO(index)
		return true
	}
	if checkHELP(message) {
		HandleHELP(index)
		return true
	}
	if checkVERSION(message) {
		HandleVERSION(index)
		return true
	}
	return false
}

/* HandleWHO: send out client names in response to the WHO command
 */
func HandleWHO(index int) {

	message := ""
	name := clients_array[index].name
	conn := *clients_array[index].connection

	for _, v := range clients_array {
		if v.name != "" && conn != nil {
			message += fmt.Sprintf("\n%s", name)
		}
	}

	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Printf("Error when sending the WHO message\n")
	}
}

/* HandleJOIN: Add client name in position index in the
*  array of clients
 */
func HandleJOIN(index int, nameToJoin string) {

	message := ""
	conn := *clients_array[index].connection
	name := clients_array[index].name

	if index > max_clients || index < 0 {
		log.Printf("HandleJOIN, incorrect index value: %d", index)
	}

	// Make sure that client did not already join
	if name != "" {
		message = fmt.Sprintf("Already joined as %s", clients_array[index].name)
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Printf("handleJOIN: server write error: %v", err)
		}
		return
	}

	// we have to block here, in case several clients want to subscribe at the same time
	clients_mutex.Lock()
	clients_array[index].name = nameToJoin
	clients_mutex.Unlock() // release the lock

	// let the other users know about the new user
	// HandleBroadcast(index, message)

}

/* HandleLEAVE: Remove client form position index in array clients
 * and close the connection to the client
 */
func HandleLEAVE(index int) {

	message := ""
	conn := *clients_array[index].connection
	name := clients_array[index].name

	if conn == nil && name == "" {
		return
	}

	// we have to block here, in case several clients want to subscribe at the same time
	clients_mutex.Lock()
	// conn.Close() // this is done in the go routine defer: defer connection.Close()
	clients_array[index].name = ""
	clients_array[index].connection = nil
	clients_mutex.Unlock() // release the lock

	message = fmt.Sprintf("%s just leaved the chat room\n", name)
	log.Printf("%s", message)
	// HandleBroadcast(index, message)

	// kill the go routine
	// references: https://appliedgo.net/spotlight/how-goroutines-want-to-exit/#:~:text=To%20exit%20a%20goroutine%20using,signal%20the%20goroutine%20to%20exit.
}

func HandleVERSION(index int) {

	conn := *clients_array[index].connection
	version := fmt.Sprintf("version: %s", code_version)

	_, err := conn.Write([]byte(version))
	if err != nil {
		log.Printf("handleVERSION: server write error: %v", err)
	}
}

func HandleHELP(index int) {

	conn := *clients_array[index].connection

	message := "You can use the commands:\n"
	message += " - JOIN <name> : join to the chat with alias <name>\n"
	message += " - WHO         : enumerate the chat participands\n"
	message += " - LEAVE       : leave the chat\n"
	message += " - VERSION     : get the version of the server\n"
	message += " - HELP        : list the possible commands\n\n"

	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Printf("handleVERSION: server write error: %v", err)
	}
}
