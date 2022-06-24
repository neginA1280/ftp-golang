package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

////////////////////////////////////////////////////////////////

func handleClient(conn net.Conn, addr net.Addr) {

	fmt.Printf("\nAccepted connection from client : " + addr.String() + "\n")

	connected := true
	for connected {

		fmt.Println("\nWaiting for client [" + addr.String() + "] input ...")

		// get message from client
		message, err := bufio.NewReader(conn).ReadString('\n')
		errorCheck(err, "Couldn't read the new message from client")

		mes := strings.ReplaceAll(message, "\n", "")
		command, _, _ := strings.Cut(mes, " ")

		if command == "get" {
			fmt.Println("get")
		} else if command == "put" {
			fmt.Println("put")
		} else if command == "ls" {
			fmt.Println("ls")
		} else if command == "quit" {
			connected = false
			fmt.Printf("Client " + addr.String() + " left.\n")
		} else {
			fmt.Println("Invalid Command!")
		}

	}
	conn.Close()

}

/////////////////////////////////////////////////////////////////

func errorCheck(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		return
	}
}

/////////////////////////////////////////////////////////////////

func main() {
	fmt.Println("Start server...")

	// listen on port 8000
	server, err := net.Listen("tcp", ":8000")
	fmt.Println("Waiting for connections ...")

	errorCheck(err, "Couldn't listen on port 8000")

	defer server.Close()

	// run loop forever (or until ctrl-c)
	for {
		// accept connection
		conn, err := server.Accept()
		addr := conn.RemoteAddr()
		errorCheck(err, "Couldn't accept any new connections")

		// handle the connection (in seperate channel)
		go handleClient(conn, addr)
	}
}
