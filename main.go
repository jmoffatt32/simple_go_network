package main

import (
	"fmt"
	"os"
	"strconv"
	"tcp-network/pkg/config"
)

func inputError() {
	fmt.Println("Please provide a network ID: {1, 2, 3, 4}")
}

func main() {

	// Take ID input from command line and initialize node address and delay
	if len(os.Args) < 2 {
		inputError()
		return
	}

	// Set the network ID and check for errors taking the command line input
	id := os.Args[1]
	val, err := strconv.Atoi(os.Args[1])
	if err != nil || val < 1 || val > 4 {
		inputError()
		return
	}

	// Initialize configuration variables
	addrMap, delay := config.FetchConfig()
	address := addrMap[id]
	fmt.Println(address)
	fmt.Println(delay)

	// Implement package to start a client and a package to start a server...
}
