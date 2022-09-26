package server

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type Message struct {
	src     string
	dest    string
	content string
}

type Confirmed struct {
	message   Message
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

	// TODO:
	// Change implementation to write bites where the first n bytes are the address that the message came from
	fmt.Fprintf(c, message+"\n")
	c.Close()
}

func unicast_recieve(c net.Conn, client net.Conn) {

	// TODO:
	// Change implementation to Read Bytes where the first n bytes are the address that the message came from
	netData, err := bufio.NewReader(c).ReadString('\n')
	check(err)

	temp := strings.TrimSpace(string(netData)) + "\n"
	src := c.RemoteAddr().String()
	dest := c.LocalAddr().String()
	msg := Message{src, dest, temp}

	t := time.Now()
	myTime := t.Format(time.RFC3339)
	content := strings.Trim(msg.content, "\n")
	client.Write([]byte(myTime + "---" + "Received: \"" + content + "\" from " + msg.src + "\n"))

	c.Close()
}

func inbound(l net.Listener, client net.Conn) {
	for {
		c, err := l.Accept()
		check(err)
		go unicast_recieve(c, client)
	}
}

func outbound(delays [2]int, outgoing_messages chan Message, client net.Conn) {
	for {
		var msg Message = <-outgoing_messages
		delay := delays[0] + rand.Intn(delays[1]-delays[0])

		time.Sleep(time.Duration(delay) * time.Millisecond)
		unicast_send(msg.dest, msg.content)

		sent := Confirmed{msg, time.Now()}
		t := sent.timestamp
		myTime := t.Format(time.RFC3339)
		content := strings.Trim(sent.message.content, "\n")
		client.Write([]byte(myTime + "---" + "Sent: \"" + content + "\" to " + sent.message.dest + "\n"))
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
