package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q03_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q03_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q03_p3.txt`, "the input for part 3")

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
		sizeStr := scanner.Text()
		if scanner.Err() != nil {
			panic("Scanner error")
		}
		sizeStrs := strings.Split(sizeStr, ",")
		sizes := make([]int, 0, len(sizeStrs))
		for _, sStr := range sizeStrs {
			val, err := strconv.Atoi(sStr)
			if err != nil {
				panic("Bad size " + sStr)
			}
			sizes = append(sizes, val)
		}
		slices.Sort(sizes)
		totSz := 0
		currentSz := sizes[len(sizes)-1] + 1
		for idx := len(sizes); idx > 0; idx-- {
			val := sizes[idx-1]
			if val < currentSz {
				totSz += val
				currentSz = val
			}
		}
		fmt.Println("p1", totSz)
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
		sizeStr := scanner.Text()
		if scanner.Err() != nil {
			panic("Scanner error")
		}
		sizeStrs := strings.Split(sizeStr, ",")
		sizes := make([]int, 0, len(sizeStrs))
		for _, sStr := range sizeStrs {
			val, err := strconv.Atoi(sStr)
			if err != nil {
				panic("Bad size " + sStr)
			}
			sizes = append(sizes, val)
		}
		slices.Sort(sizes)
		totSz := 0
		usedVals := 0
		currentSz := sizes[0] - 1
		for _, val := range sizes {
			if val > currentSz {
				totSz += val
				currentSz = val
				usedVals++
				if usedVals >= 20 {
					break
				}
			}
		}
		fmt.Println("p2", totSz)
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
		sizeStr := scanner.Text()
		if scanner.Err() != nil {
			panic("Scanner error")
		}
		sizeStrs := strings.Split(sizeStr, ",")
		sizes := make([]int, 0, len(sizeStrs))
		for _, sStr := range sizeStrs {
			val, err := strconv.Atoi(sStr)
			if err != nil {
				panic("Bad size " + sStr)
			}
			sizes = append(sizes, val)
		}
		maxoccur := 0
		mymap := make(map[int]int)
		for _, val := range sizes {
			mymap[val] += 1
			maxoccur = max(maxoccur, mymap[val])
		}
		fmt.Println("p3", maxoccur)
	}
}
