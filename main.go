package main

import (
	"fmt"
	"os"
	"tcp-network/config"
)

func main() {

	// Take ID input from command line and initialize node address and delay
	id := os.Args[1]
	addrMap, delay := config.FetchConfig()
	address := addrMap[id]

	fmt.Println(address)
	fmt.Println(delay)

	// Implement packages to start a client and a package to start a server
}
