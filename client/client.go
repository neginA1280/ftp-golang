package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

/////////////////////////////////////////////////////////////////

func handleCommand(conn net.Conn) {
	// an infinit loop for receiving client commands until quit
	for {
		// what to send?
		fmt.Print("ftp> ")
		reader := bufio.NewReader(os.Stdin)

		message, err := reader.ReadString('\n')
		errorCheck(err, "Couldn't read the message from input")

		// send command to server
		fmt.Fprintf(conn, message)

		mes := strings.ReplaceAll(message, "\n", "")
		command, _, _ := strings.Cut(mes, " ")

		if command == "get" {
			fmt.Println("get command")
		} else if command == "put" {
			fmt.Println("put command")
		} else if command == "ls" {
			fmt.Println("ls command")
		} else if command == "quit" {
			conn.Close()
			break
		} else {
			fmt.Println("Invalid Command!")
		}
	}
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

	// connect to server
	conn, err := net.Dial("tcp", "localhost:8000")
	errorCheck(err, "Couldn't connect to the server")

	handleCommand(conn)

}
