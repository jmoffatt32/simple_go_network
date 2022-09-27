package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Simple error check function to catch errors and stop the program.
// `check` will print the error message and then terminate the program.
func check(err error) {
	if err != nil {
		fmt.Print(err)
		return
	}
}

// Goroutine function to listen for "Sent" and "Received" notifications
// from the server. It takes the clients connection to the serer as an
// argument and reads for incoming notifications through the buffer.
func listening(c net.Conn) {
	for {
		// Reads response from server
		message, _ := bufio.NewReader(c).ReadString('\n')
		// Remove ">> " prompt characters from Stdout buffer, then...
		fmt.Fprint(os.Stdout, "\r \r")
		// ...print the message, finally...
		fmt.Print(message)
		// ...prompt for user input.
		fmt.Print(">> ")
	}
}

// Runs the command line client.
// `address` is the "host:port" address string of the server that the
// client will connect to.
// After connection is made, the client will start a goroutine to read
// messages from the server. Then, the client waits for user inputs which
// it sends to the server to process.
func MainClient(address string) {

	// Dial connection to the server
	c, err := net.Dial("tcp", address)
	check(err)

	// Start goroutine to read notifcations from the server
	go listening(c)

	// Prompt user for input
	fmt.Print(">> ")
	for {
		// Reads input from user, then...
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		check(err)
		// ...writes user's input to the connection with the server.
		c.Write([]byte(input + "\n"))
		fmt.Print(">> ")

		// Checks if user attempts to end the session
		if strings.TrimSpace(string(input)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}
