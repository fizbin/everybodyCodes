package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q13_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q13_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q13_p3.txt`, "the input for part 3")

const doingPart = 3

func main() {
	flag.Parse()
	if doingPart >= 1 {
		infile := *input1
		data, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		fmt.Println("p1:", doProblem1(data))
	}
	if doingPart >= 2 {
		infile := *input2
		data, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		fmt.Println("p2:", doProblem2(data))
	}
	if doingPart >= 3 {
		infile := *input3
		data, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		fmt.Println("p3:", doProblem3(data))
	}
}

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	nums := make([]int, 0)
	for scanner.Scan() {
		val, _ := strconv.Atoi(scanner.Text())
		nums = append(nums, val)
	}
	// initially, wheel has "1" at position 0, and has len(nums)+1 spots
	// so number at position 0 in nums ends up at position 1 on wheel
	// number at position 1 in nums ends up at position -1 on wheel
	// .................. 2 ........................... 2
	// .................. 3 ........................... -2
	// .................. 4 ........................... 3
	//                    5                             -3

	rotSpot := 2025 % (len(nums) + 1)
	if rotSpot == 0 {
		return 1
	}
	negRotSpot := len(nums) + 1 - rotSpot
	if rotSpot <= negRotSpot {
		return nums[(rotSpot-1)*2]
	}
	return nums[(negRotSpot-1)*2+1]
}

type numRange struct {
	von, zu int
}

func doPartsTwoThree(data []byte, rotAmount int) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	parser := regexp.MustCompile(`(\d+)-(\d+)`)
	isClockwise := true
	cwRanges := make([]numRange, 0)
	ccwRanges := make([]numRange, 0)
	dialPop := 1
	for scanner.Scan() {
		parts := parser.FindStringSubmatch(scanner.Text())
		lft, _ := strconv.Atoi(parts[1])
		rgt, _ := strconv.Atoi(parts[2])
		dialPop += rgt - lft + 1
		if isClockwise {
			cwRanges = append(cwRanges, numRange{lft, rgt})
		} else {
			ccwRanges = append(ccwRanges, numRange{lft, rgt})
		}
		isClockwise = !isClockwise
	}
	slices.Reverse(ccwRanges)
	rotSpot := rotAmount % dialPop
	// fmt.Println("RotSpot start", rotSpot)
	if rotSpot == 0 {
		return 1
	}
	rotSpot--
	for _, numPair := range cwRanges {
		rangeSize := numPair.zu - numPair.von + 1
		if rotSpot < rangeSize {
			return numPair.von + rotSpot
		}
		rotSpot -= rangeSize
	}
	// fmt.Println("RotSpot after cw", rotSpot)
	for _, numPair := range ccwRanges {
		rangeSize := numPair.zu - numPair.von + 1
		if rotSpot < rangeSize {
			// fmt.Println("returning from range", numPair, "and rotSpot", rotSpot)
			return numPair.zu - rotSpot
		}
		rotSpot -= rangeSize
	}
	return -1 // error value
}

func doProblem2(data []byte) any {
	return doPartsTwoThree(data, 20252025)
}

func doProblem3(data []byte) any {
	return doPartsTwoThree(data, 202520252025)
}
