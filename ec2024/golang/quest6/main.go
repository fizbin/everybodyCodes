package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type queueElem struct {
	node  string
	depth int
	pnode *queueElem
}

func findPowerful(network map[string][]string, abbrev bool) string {
	finals := make(map[int][]queueElem)
	var q []queueElem = []queueElem{{"RR", 1, nil}}
	for len(q) > 0 {
		me := q[0]
		q = q[1:]
		if val, ok := network[me.node]; ok {
			for _, next := range val {
				if next != "ANT" && next != "BUG" {
					q = append(q, queueElem{next, me.depth + 1, &me})
				}
			}
		} else {
			if me.node == "@" {
				finals[me.depth] = append(finals[me.depth], me)
			}
		}
	}
	var topNode *queueElem
	for _, val := range finals {
		if len(val) == 1 {
			topNode = &(val[0])
		}
	}
	if topNode == nil {
		log.Fatal("Couldn't find a unique top node")
	}
	var names []string
	for topNode != nil {
		if abbrev {
			names = append(names, topNode.node[0:1])
		} else {
			names = append(names, topNode.node)
		}
		topNode = topNode.pnode
	}
	slices.Reverse(names)
	return strings.Join(names, "")
}

func allPaths(start string, network map[string][]string, abbrev bool) []string {
	var retval []string
	if val, ok := network[start]; ok {
		for _, next := range val {
			for _, pth := range allPaths(next, network, abbrev) {
				if abbrev {
					retval = append(retval, start[0:1]+pth)
				} else {
					retval = append(retval, start+pth)
				}
			}
		}

		return retval
	} else {
		return []string{start}
	}
}

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q06_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		network := make(map[string][]string)
		for scanner.Scan() {
			splitLine := strings.Split(scanner.Text(), ":")
			splitDests := strings.Split(splitLine[1], ",")
			network[splitLine[0]] = splitDests
		}
		lenToPath := make(map[int][]string)
		for _, pth := range allPaths("RR", network, false) {
			plen := len(pth)
			lenToPath[plen] = append(lenToPath[plen], pth)
		}
		var ans string
		for _, val := range lenToPath {
			if len(val) == 1 {
				ans = val[0]
			}
		}

		fmt.Println("Part 1:", ans)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q06_p2.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		network := make(map[string][]string)
		for scanner.Scan() {
			splitLine := strings.Split(scanner.Text(), ":")
			splitDests := strings.Split(splitLine[1], ",")
			network[splitLine[0]] = splitDests
		}
		lenToPath := make(map[int][]string)
		for _, pth := range allPaths("RR", network, true) {
			plen := len(pth)
			lenToPath[plen] = append(lenToPath[plen], pth)
		}
		var ans string
		for _, val := range lenToPath {
			if len(val) == 1 {
				ans = val[0]
			}
		}

		fmt.Println("Part 2:", ans)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q06_p3.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		network := make(map[string][]string)
		for scanner.Scan() {
			splitLine := strings.Split(scanner.Text(), ":")
			splitDests := strings.Split(splitLine[1], ",")
			network[splitLine[0]] = splitDests
		}
		fmt.Println("Part 3:", findPowerful(network, true))
	}
}
