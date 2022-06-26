package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////

var isSendingFile bool = false
var fileName string

////////////////////////////////////////////////////////////////

func ls(conn net.Conn) {

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	errorCheck(err, "Couldn't read the list from server")

	fmt.Println(string(buf[:n]))
}

////////////////////////////////////////////////////////////////

func put(conn net.Conn, fileName string) {

	// receiving client commands
	file, err := os.Open("client/" + fileName)
	errorCheck(err, "Couldn't open file "+fileName)
	defer file.Close()

	buf := make([]byte, 1024)
	for {

		//Read chunk bytes from file
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("sending file completed.")
			} else {
				errorCheck(err, "Couldn't send file from client to server")
			}
			return
		}

		//sending chunk bytes to client
		_, err = conn.Write(buf[:n])
	}

}

////////////////////////////////////////////////////////////////

func get(conn net.Conn, fileName string) {

	//Create a new file by file name
	file, err := os.Create("client/" + fileName)
	errorCheck(err, "Couldn't create a new file.")
	defer file.Close()

	// set a deadline for reading side of the connection
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)

		//Write received data to local file
		file.Write(buf[:n])

		if err != nil {
			fmt.Printf("receive file complete. \n")
			isSendingFile = false
			return
		}

	}

}

/////////////////////////////////////////////////////////////////

func handleCommand(conn net.Conn) {

	// if send/receive data channel is active don't receive the client commands.
	// instead, go to send/receive file channel
	if isSendingFile == true {
		get(conn, fileName)
		return
	}

	// an infinit loop for receiving client commands until quit
	for {

		fmt.Print("ftp> ")
		reader := bufio.NewReader(os.Stdin)

		message, err := reader.ReadString('\n')
		errorCheck(err, "Couldn't read the message from input")

		// send command to server
		conn.Write([]byte(message))

		mes := strings.ReplaceAll(message, "\n", "")
		command, fileName, _ := strings.Cut(mes, " ")

		if command == "get" {
			isSendingFile = true
			get(conn, fileName)
		} else if command == "put" {
			put(conn, fileName)
		} else if command == "ls" {
			ls(conn)
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
