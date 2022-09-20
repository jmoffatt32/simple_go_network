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

	addrMap := make(map[string]string)
	var delays [2]int
	scanner := bufio.NewScanner(f)
	lineNum := 0

	// Read the file
	for scanner.Scan() {
		// If reading the first line of the file...
		if lineNum == 0 {
			// ...use two integer values to initalize the delays array...
			var line = scanner.Text()
			delays[0], err = strconv.Atoi(line[:4])
			check(err)
			delays[1], err = strconv.Atoi(line[5:])
			check(err)
			// ...otherwise...
		} else {
			// ...read the line and initialize a field in the port
			var line = string(scanner.Text())
			addrMap[string(line[0])] = line[2:11] + ":" + line[12:]
		}
		lineNum++
	}

	err = scanner.Err()
	check(err)

	return addrMap, delays
}
