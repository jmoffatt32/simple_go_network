package server

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

// Message struct holds message objects. It stores the source address, and destination
// addresses of the message, both as strings. Addtionally, it stores a string with the
// content of the message.
type Message struct {
	src     string
	dest    string
	content string
}

// Simple error check function to catch errors and stop the program.
// `check` will print the error message and then terminate the program.
func check(err error) {
	if err != nil {
		fmt.Print(err)
		return
	}
}

// Helper function to parse the users input. The function expects a string
// with 2 or more spaces and returns two string values, the network ID of the
// destination process and the content of the message.
func parse_input(input string) (string, string) {

	frags := strings.Split(input, " ")
	dest := frags[1]                        // Extract network ID
	message := strings.Join(frags[2:], " ") // Extract message content
	return dest, message
}

// Sends message a message to the destination address. `destination` is the
// "host:port" of the destination process and `message` is the content of the message
// to be send.
func unicast_send(destination string, message string) {

	// Dial a connection to the destination process
	c, err := net.Dial("tcp", destination)
	check(err)
	// Write the message via the TCP connection
	c.Write([]byte(message + "\n"))
	c.Close()
}

// Recieves an icoming message from a foriegn process. `c` is a network connection that
// has come in from another process and `client` is the network connection the server has
// to the main client (the command line). The message is read and the content of the message
// is delivered to the main client.
func unicast_recieve(c net.Conn, client net.Conn) {

	// Read incoming message from the network connection
	netData, err := bufio.NewReader(c).ReadString('\n')
	check(err)

	// Trim the message and create the messge struct variable to organize the data.
	temp := strings.TrimSpace(string(netData))
	content := strings.Trim(temp, "\n")
	src := c.RemoteAddr().String()
	dest := c.LocalAddr().String()
	msg := Message{src, dest, content}
	// Technically, the use of a message struct in this case is redundant, as we could simply
	// reference `src`, `dest`, and `msg` directly for the remainder of the fucntion, but it
	// is useful for organization to use the struct.

	// Get the timestamp of the message reception
	t := time.Now()
	myTime := t.Format(time.RFC3339)

	// Write the reception message to the client
	client.Write([]byte("->: " + myTime + "---" + "Received: \"" + content + "\" from " + msg.src + "\n"))
	c.Close()
}

// Goroutine function to listen for incoming messages.
// It takes the server's listener to listen for incoming connections and the
// server's connection to the command line client to write the "Received" notifcations
// back to the command line.
// The for loop is blocked at the first line until an icoming connection is heard.
// The connection is then passed to `unicast_receive` to handle delivery of the
// message and the loop returns to the first line to listen for new connections.
func inbound(l net.Listener, client net.Conn) {
	for {
		// Wait for and accept incoming connections
		c, err := l.Accept()
		check(err)

		// Handle new connection and deliver the message
		go unicast_recieve(c, client)
	}
}

// Goroutine function to handle sending ooutgoing messages.
// It takes the array holding the min/max delay values, a Message channel which holds
// queue of messages to be sent, and the servers connection to the command line client
// to write the "Sent" notifications back to the command line.
// The for loop is blocked at the first line until a message is waiting in the
// `outgoing_messages` channel. Once their is a message in the queue, it generates a
// random time delay and starts a goroutine to handle the delay and to send the message via
// `unicast_receive`. Simultaneously, it writes back the "Sent" notifcation to the client.
func outbound(delays [2]int, outgoing_messages chan Message, client net.Conn) {
	for {
		// Wait for and load outbound messages from the queue
		var msg Message = <-outgoing_messages

		// Generate random time delay
		delay := delays[0] + rand.Intn(delays[1]-delays[0])
		// Goroutine to handle...
		go func() {
			// ...simulating the network delay, then...
			time.Sleep(time.Millisecond * time.Duration(delay))
			// ... delivering the message.
			unicast_send(msg.dest, msg.content)
		}()

		// Get the timestamp that the server processed the message
		t := time.Now()
		myTime := t.Format(time.RFC3339)

		// Write the "Sent" notification back to the client
		content := strings.Trim(msg.content, "\n")
		client.Write([]byte("->: " + myTime + "---" + "Sent: \"" + content + "\" to " + msg.dest + "\n"))
	}
}

// Run the messaging server.
// `address` is the "host:port" string of the network to host the server on.
// `addrMap` maps network IDs to "host:port" format address of other possible
// networks to send messages to. `delay` holds the min/max values of the simulated
// network delay.
// The server intializes on the provided port and starts Goroutines to handle accepting
// incoming messages and sending outgoing messages. Then the for loop waits for input
// from the client to send messages to other servers.
func Server(address string, addrMap map[string]string, delay [2]int) {

	// Start the server on the provided port
	port := ":" + address[10:]
	l, err := net.Listen("tcp", port)
	check(err)
	defer l.Close()

	// Open connection to the client for communication with client
	c, err := l.Accept()
	check(err)

	// Make a channel for the server to fill with a queue of outbound messages
	// to be sent.
	outgoing := make(chan Message, 5)

	// Start a goroutine to handle incoming connections and deliver messages and
	// another goroutine to handle sending outgoing messages from the "outgoing" channel.
	go inbound(l, c)
	go outbound(delay, outgoing, c)

	// Read the buffer from the client for users input..
	for {
		// ...read the input from the client until a "\n" character...
		netData, err := bufio.NewReader(c).ReadString('\n')
		check(err)

		// ... if necessary, stop the server...
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		// ...else, parse the users input and extract the
		// destination and content for the message.
		id, msg := parse_input(netData)
		src := c.LocalAddr().String()
		dest := addrMap[id]

		// Place this message into the outgoing channel for it to be
		// sent by the goroutine running "outbound".
		outgoing <- Message{src, dest, msg}

	}

}
