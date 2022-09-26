package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func listening(c net.Conn) {
	for {
		// Reads response from server
		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print(message)
	}
}

func MainClient(address string) {

	c, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	go listening(c)
	for {
		// Reads input from user and sends it to server
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")

		// Checks if user attempts to end the session
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}
