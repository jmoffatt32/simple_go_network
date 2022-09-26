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
	address string
	content string
}

type Confirmed struct {
	message   Message
	direction string
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
	c.Close()
}

func unicast_recieve(c net.Conn, client net.Conn) {

	netData, err := bufio.NewReader(c).ReadString('\n')
	check(err)

	temp := strings.TrimSpace(string(netData)) + "\n"
	src := c.RemoteAddr().String()
	msg := Message{src, temp}

	t := time.Now()
	myTime := t.Format(time.RFC3339)
	content := strings.Trim(msg.content, "\n")
	client.Write([]byte(myTime + "---" + "Recieved: \"" + content + "\" to " + msg.address + "\n"))

	c.Close()
}

func incoming_routine(l net.Listener, client net.Conn) {
	for {
		c, err := l.Accept()
		check(err)
		go unicast_recieve(c, client)
	}
}

func outgoing_routine(delays [2]int, outgoing_messages chan Message, client net.Conn) {
	for {
		var msg Message = <-outgoing_messages
		delay := delays[0] + rand.Intn(delays[1]-delays[0])

		time.Sleep(time.Duration(delay) * time.Millisecond)
		unicast_send(msg.address, msg.content)

		sent := Confirmed{msg, "OUT", time.Now()}
		t := sent.timestamp
		myTime := t.Format(time.RFC3339)
		content := strings.Trim(sent.message.content, "\n")
		client.Write([]byte(myTime + "---" + "Sent: \"" + content + "\" to " + sent.message.address + "\n"))
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

	outgoing := make(chan Message, 5)
	go incoming_routine(l, c)
	go outgoing_routine(delay, outgoing, c)
	for {
		// Read the input from the client...
		netData, err := bufio.NewReader(c).ReadString('\n')
		check(err)
		// ... if necessary, stop the server.
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		dest, msg := parse_input(netData)
		address := addrMap[dest]
		outgoing <- Message{address, msg}
	}

}
