package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type Network struct {
	id      string
	host    string
	port    string
	send    net.Conn
	receive net.Conn
}

var DNS map[string]Network

func check(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Server(address string, addrMap map[string]string, delay [2]int) {

	// Start the server on the provided port
	port := ":" + address[10:]
	l, err := net.Listen("tcp", port)
	check(err)
	defer l.Close()

	// Open connection to the client for communication with client
	c, err := l.Accept()
	check(err)

	// IMPLEMENT:
	// Loop over address map and attempt to establish connections with
	// all known servers.
	for id, address := range addrMap {
		// If a server is found, establish a connection, initialize a Network
		// variable, and add it to the array
		// Should look something like:
		var s net.Conn
		var r net.Conn
		n := Network{id, address[:10], address[10:], s, r}
		DNS[id] = n
	}
	fmt.Println(DNS)

	// Start loop to read the buffer for messages from the client to send
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		check(err)
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		// IMPLEMENT:
		// Check for any incoming connections from new servers trying to
		// establish a TCP connection. Not sure of the best implementation

		// IMPLEMENT:
		// Parse message from cleint (netData) and extract a destination.
		// Use unicast_send to send it to the correct server.
		// Implement through unicast_send(dest, msg)

		// IMPLEMENT
		// Check existing receive channels with other servers
		// If there are any messages, use c to send them to the client.
		// Implement through unicast_receive(src, msg)

		// Send message to the Client
		// IMPLEMENT
		// Need to update this block to send the message along
		// with the source of the message and the timestamp the
		// server received it
		fmt.Print("-> ", string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		c.Write([]byte(myTime))
	}

}
