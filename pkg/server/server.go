package server

import (
	"bufio"
	"fmt"
	"math/rand"
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

type OutgoingMessage struct {
	destination string
	content     string
}

type Confirmed struct {
	direction string
	address   string
	content   string
	status    bool
	timestamp time.Time
}

func check(err error) {
	if err != nil {
		fmt.Print(err)
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

func unicast_recieve(source string, message string) {

}

func outgoing_routine(delays [2]int, outgoing_messages chan OutgoingMessage, status chan Confirmed) {

	for {
		var msg OutgoingMessage = <-outgoing_messages
		delay := delays[0] + rand.Intn(delays[1]-delays[0])

		time.Sleep(time.Duration(delay) * time.Millisecond)
		unicast_send(msg.destination, msg.content)

		trimmed := strings.Trim(msg.content, "\n")
		verified := Confirmed{"OUT", msg.destination, trimmed, true, time.Now()}
		status <- verified
	}
}

func incoming_routine(listener net.Listener, incoming chan Confirmed) {

	for {
		c, err := listener.Accept()
		check(err)

		source := c.RemoteAddr().String()
		message, err := bufio.NewReader(c).ReadString('\n')
		check(err)
		unicast_recieve(source, message)
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

	// Make channels to communicate with send and recieve goroutines
	confirmed := make(chan Confirmed, 5)
	outgoing := make(chan OutgoingMessage, 5)

	// Start a go routine to handle sending outgoing messages. This accepts a channel to read
	// messages to be sent as well as a channel to tell the server when a message has been sent.
	go outgoing_routine(delay, outgoing, confirmed)

	// NEED TO FINISH
	// Start a go routine to handle incoming connections. This routine should implement unicast_receive
	// to deleiver the message recieved from the source process
	go incoming_routine(l, confirmed)

	// Create a buffered channel to hold all the outputs to be sent
	// read to the client at the end of each iteration.
	outputs := make(chan Confirmed, 5)
	// Start loop to read the buffer for messages from the client to send
	// This is the application layer
	for {
		// Read the input from the client...
		netData, err := bufio.NewReader(c).ReadString('\n')
		check(err)
		// ... if necessary, stop the server.
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		// Parse users input to extract destination and content of the message to be sent...
		dst_id, content := parse_input(netData)
		// ... convert network ID to "host:port"...
		dest := addrMap[dst_id]
		// ... create a new Message and send it to the "outgoing_routine" via
		// the outgoing channel.
		var new_msg OutgoingMessage = OutgoingMessage{dest, content}
		outgoing <- new_msg

		// ... accept Confirmed variable from the "outgoing_routine" once the message
		// has been sent.
		var verified Confirmed = <-confirmed
		outputs <- verified

		// IMPLEMENT:
		// Read channel from "incoming_routine" to check for any incoming messages.

		// Loop over outputs and empty the channel by writing each output to the client.
		for element := range outputs {
			t := element.timestamp
			myTime := t.Format(time.RFC3339)
			c.Write([]byte(myTime + "---" + "Sent: \"" + element.content + "\" to " + element.address + "\n"))
		}

	}

}
