package main

import (
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

////////////////////////////////////////////////////////////////////

func ls(conn net.Conn) {

	f, err := os.Open("./server/")
	if err != nil {
		// TODO
		fmt.Println("Didn't find any files in this directory.")
		return
	}

	files, err := f.Readdir(0)
	if err != nil {
		// TODO
		fmt.Println(err)
		return
	}

	list := ""
	// converts file format to string format
	for _, v := range files {
		list += string(v.Name()) + "\n"
	}

	// send data to client
	conn.Write([]byte(list))

	fmt.Println("ls command output sent successfully!!!")

}

////////////////////////////////////////////////////////////////

func get(conn net.Conn, fileName string) {

	//Create a new file by file name
	file, err := os.Create("server/" + fileName)
	errorCheck(err, "Couldn't create a new file.")
	defer file.Close()

	// set a deadline for reading side of the connection
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)

		fmt.Println("received chunks : [" + string(buf[:n]) + "]")
		//Write received data to local file
		file.Write(buf[:n])

		if err != nil {
			fmt.Printf("receive file complete. \n")
			isSendingFile = false
			return
		}

	}

}

////////////////////////////////////////////////////////////////

func put(conn net.Conn, fileName string) {

	// receiving client commands
	file, err := os.Open("server/" + fileName)
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
				errorCheck(err, "Couldn't send file from server to client")
			}
			return
		}

		//sending chunk bytes to client
		_, err = conn.Write(buf[:n])
	}

}

////////////////////////////////////////////////////////////////

func handleClient(conn net.Conn, addr net.Addr) {

	fmt.Printf("\nAccepted connection from client : " + addr.String() + "\n")
	if isSendingFile == true {
		fmt.Println("getting file from client")
		get(conn, fileName)
		return
	}
	connected := true
	for connected {

		fmt.Println("\nWaiting for client [" + addr.String() + "] input ...")

		// get message from client
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		errorCheck(err, "Couldn't read the new message from client")

		mes := strings.ReplaceAll(string(buf[:n]), "\n", "")
		command, fileName, _ := strings.Cut(mes, " ")

		if command == "get" {
			put(conn, fileName)
		} else if command == "put" {
			isSendingFile = true
			fmt.Println("start receiving file")
			get(conn, fileName)
		} else if command == "ls" {
			ls(conn)
		} else if command == "quit" {
			connected = false
			fmt.Printf("Client " + addr.String() + " left.\n")
		} else {
			fmt.Println("Invalid Command!")
		}
		// break
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
