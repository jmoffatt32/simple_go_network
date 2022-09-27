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
	f, err := os.Open("pkg/config/.config")
	check(err)
	defer f.Close()

	var delay [2]int
	IDList := make([]string, 0)
	IPList := make([]string, 0)
	PortList := make([]string, 0)
	addrMap := make(map[string]string)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		x := strings.Split(string(scanner.Text()), " ")
		for i, v := range x {
			if i%3 == 0 {
				IDList = append(IDList, v)
			}
			if i%3 == 1 {
				IPList = append(IPList, v)
			}
			if i%3 == 2 {
				PortList = append(PortList, v)
			}
		}
	}

	delay[0], err = strconv.Atoi(IDList[0])
	delay[1], err = strconv.Atoi(IPList[0])
	check(err)
	IDList = IDList[1:]
	IPList = IPList[1:]

	for i := 0; i < len(IDList); i++ {
		addrMap[IDList[i]] = IPList[i] + ":" + PortList[i]
	}

	return addrMap, delay
}
