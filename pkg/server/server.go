package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// type Network struct {
// 	id      string
// 	host    string
// 	port    string
// 	send    net.Conn
// 	receive net.Conn
// }

// var DNS map[string]Network

type Message struct {
	destination string
	message     string
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func parse_input(raw_input string) (string, string) {

	frags := strings.Split(raw_input, " ")
	dest := frags[1]
	message := strings.Join(frags[2:], " ")
	return dest, message
}

func unicast_send(destination string, message string) {

	c, err := net.Dial("tcp", destination)
	check(err)
	fmt.Fprintf(c, message+"\n")
}

// func unicast_recieve(source string, message string) {

// }

func Server(address string, addrMap map[string]string, delay [2]int) {

	// Start the server on the provided port
	port := ":" + address[10:]
	l, err := net.Listen("tcp", port)
	check(err)
	defer l.Close()

	// Open connection to the client for communication with client
	c, err := l.Accept()
	check(err)

	// Make channels to communicate with send and recieve goroutines
	// incoming := make(chan Message, 5)
	outgoing := make(chan Message, 5)

	// Start loop to read the buffer for messages from the client to send
	// This is the application layer
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		check(err)
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		// IMPLEMENT:
		// Parse client input to extract destination and message.
		// Add message to outgoing message buffer to be read from outgoing
		// transport GoRoutine.
		dst_id, msg := parse_input(netData)
		dest := addrMap[dst_id]
		var new_outgoing Message = Message{dest, msg}
		outgoing <- new_outgoing

		// IMPLEMENT
		// Check existing receive channels with other servers
		// If there are any messages, use c to send them to the client.
		// Implement through unicast_receive(src, msg)

		// Send message to the Client
		// IMPLEMENT
		// Need to update this block to send the message along
		// with the source of the message and the timestamp the
		// server received it

		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		c.Write([]byte(myTime))

	}

}
