package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Simple error check function to catch errors and stop the program.
// `check` will print the error message and then terminate the program.
func check(err error) {
	if err != nil {
		fmt.Print(err)
		return
	}
}

// Configure network settings from "pkg/config/.config" file.
// Opens file and reads scans line by line. To configure delay
// settings and all possible addresses. Returns a map of network IDs
// to network addresses and an array of min/max delay values.
func FetchConfig() (map[string]string, [2]int) {

	// Safely open the ".config" file
	f, err := os.Open("pkg/config/.config")
	check(err)
	defer f.Close()

	// Initialize variables. delay [2]int will be exported. IDList, IPList, and PortList will be collated and exported in addrMap
	var delay [2]int
	IDList := make([]string, 0)
	IPList := make([]string, 0)
	PortList := make([]string, 0)
	addrMap := make(map[string]string)

	scanner := bufio.NewScanner(f)
	// Begin scanning the file...
	for scanner.Scan() {
		x := strings.Split(string(scanner.Text()), " ")
		for i, v := range x {
			if i%3 == 0 { // Add everything from first column to IDList
				IDList = append(IDList, v)
			}
			if i%3 == 1 { // Add everything from second column to IPList
				IPList = append(IPList, v)
			}
			if i%3 == 2 { // Add everything from third column to PortList
				PortList = append(PortList, v)
			}
		}
	}

	// First items in IDList and IPList are the min and max delay
	// We take these values, convert to int, and assign to delay[0] for minimum delay, and delay[1] for maximum delay
	delay[0], err = strconv.Atoi(IDList[0])
	check(err)
	delay[1], err = strconv.Atoi(IPList[0])
	check(err)
	// Remove the min and max delay values from the ID and IP list
	IDList = IDList[1:]
	IPList = IPList[1:]

	// Assign the values from IDList, IPList and PortList to addrMap
	for i := 0; i < len(IDList); i++ {
		addrMap[IDList[i]] = IPList[i] + ":" + PortList[i]
	}

	return addrMap, delay
}
