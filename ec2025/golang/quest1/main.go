package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q01_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q01_p2.txt`, "the input for part 1")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q01_p3.txt`, "the input for part 1")

func main() {
	flag.Parse()
	{
		infile := *input1
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		nameStr := scanner.Text()
		names := strings.Split(nameStr, ",")
		scanner.Scan()
		scanner.Scan()
		instrStr := scanner.Text()
		instructions := strings.Split(instrStr, ",")
		nameidx := 0
		for _, instruction := range instructions {
			lrIndicator := instruction[0]
			distanceUint64, err := strconv.ParseUint(instruction[1:], 10, 32)
			if err != nil {
				panic("Bad instruction " + instruction)
			}
			distance := (int)(distanceUint64)
			if lrIndicator == 'L' {
				nameidx = max(0, nameidx-distance)
			} else {
				nameidx = min(len(names)-1, nameidx+distance)
			}
		}
		fmt.Println("p1:", names[nameidx])
	}
	{
		infile := *input2
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		nameStr := scanner.Text()
		names := strings.Split(nameStr, ",")
		scanner.Scan()
		scanner.Scan()
		instrStr := scanner.Text()
		instructions := strings.Split(instrStr, ",")
		nameidx := 0
		for _, instruction := range instructions {
			lrIndicator := instruction[0]
			distanceUint64, err := strconv.ParseUint(instruction[1:], 10, 32)
			if err != nil {
				panic("Bad instruction " + instruction)
			}
			distance := (int)(distanceUint64)
			if lrIndicator == 'L' {
				nameidx = ((nameidx-distance)%len(names) + len(names)) % len(names)
			} else {
				nameidx = (nameidx + distance) % len(names)
			}
		}
		fmt.Println("p2:", names[nameidx])
	}
	{
		infile := *input3
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		nameStr := scanner.Text()
		names := strings.Split(nameStr, ",")
		scanner.Scan()
		scanner.Scan()
		instrStr := scanner.Text()
		instructions := strings.Split(instrStr, ",")
		for _, instruction := range instructions {
			lrIndicator := instruction[0]
			distanceUint64, err := strconv.ParseUint(instruction[1:], 10, 32)
			if err != nil {
				panic("Bad instruction " + instruction)
			}
			var nameidx int
			distance := (int)(distanceUint64)
			if lrIndicator == 'L' {
				nameidx = len(names) - (distance % len(names))
			} else {
				nameidx = distance % len(names)
			}
			tmp := names[0]
			names[0] = names[nameidx]
			names[nameidx] = tmp
		}
		fmt.Println("p3:", names[0])
	}
}
