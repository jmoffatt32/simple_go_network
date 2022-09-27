package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func FetchConfig() (map[string]string, [2]int) {
	// Open file safely
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
