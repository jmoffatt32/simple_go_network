package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Checks errors and kills process if error occurs
func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func FetchConfig() (map[string]string, [2]int) {
	// Safely open configuration file
	f, err := os.Open("pkg/config/.config")
	check(err)
	defer f.Close()

	addrMap := make(map[string]string) // Map a network id to a network address:port
	var delay [2]int                   // delay[0] = minDelay & delay[1] = maxDelay

	// Read the file
	lineNum := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// If reading the first line of the file...
		if lineNum == 0 {
			// ...use two integer values to initalize the delay array...
			var line = scanner.Text()
			delay[0], err = strconv.Atoi(line[:4])
			check(err)
			delay[1], err = strconv.Atoi(line[5:])
			check(err)
			// ...otherwise...
		} else {
			// ...read the line and map "id" --> "address:port"
			var line = string(scanner.Text())
			addrMap[string(line[0])] = line[2:11] + ":" + line[12:]
		}
		lineNum++
	}

	// Safely close the file
	err = scanner.Err()
	check(err)

	return addrMap, delay
}
