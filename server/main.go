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
)

func main() {

	port_number := handeCommandLineParameters()
	port_string := fmt.Sprintf(":%d", port_number)

	listener, err := net.Listen("tcp", port_string)
	if err != nil {
		log.Fatal("Error listening:", err)
	}

	log.Printf("Opened tcp listener at port: %d", port_number)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

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

func handeCommandLineParameters() int {

	log.Println("number of parameters passed: ", len(os.Args))

	for i, v := range os.Args {
		log.Println("param: ", "[", i, "]: ", v, "\n\r")
	}

	if len(os.Args) > 2 {
		log.Println("Error: Too many parameters, the server only needs the port number as parameter")
		os.Exit(1)
	}

	port_number, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Println("Error: the parameter must be a port number")
		os.Exit(2)
	}
	fmt.Println("port number: ", port_number)

	return port_number
}
