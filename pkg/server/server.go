package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Server(address string) {

	// Start the server on the provided port
	port := ":" + address[10:]
	l, err := net.Listen("tcp", port)
	check(err)
	defer l.Close()

	// Open connection to the client for communication with client
	c, err := l.Accept()
	check(err)

	// Start loop to read the buffer for messages from the client to send
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		check(err)
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		// IMPLEMENT:
		// Parse message from cleint (netData) and extract a destination.
		// Use unicast_send to send it to the correct server

		// IMPLEMENT
		// Check existing reveive channels with other servers
		// If there are any messages, use c to send them to the client

		// Send message to the Client
		// This message is the timestamp at which the server processed the message
		fmt.Print("-> ", string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		c.Write([]byte(myTime))
	}
}
