package main

import (
	"fmt"
	"os"
	"strconv"
	"tcp-network/pkg/client"
	"tcp-network/pkg/config"
	"tcp-network/pkg/server"
	"time"
)

// Helper function to promt the user to provide a network ID as
// a command line argument.
func input_error() {
	fmt.Println("Please provide a network ID: {1, 2, 3, 4}")
}

func main() {

	// Take ID input from command line and initialize node address and delay
	if len(os.Args) < 2 {
		input_error()
		return
	}

	// Set the network ID and check for errors taking the command line input
	id := os.Args[1]
	val, err := strconv.Atoi(os.Args[1])
	if err != nil || val < 1 || val > 4 {
		input_error()
		return
	}

	// Initialize configuration variables
	addrMap, delay := config.FetchConfig()
	address := addrMap[id]
	fmt.Println(address)
	fmt.Println(delay)

	// Launch server to run for this process
	go server.Server(addrMap[id], addrMap, delay)
	time.Sleep(time.Millisecond * 15) // Needed so client doesn't try to reach server before server has started
	client.MainClient(address)
}
