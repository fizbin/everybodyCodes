package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q01_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q01_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q01_p3.txt`, "the input for part 3")

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

func nextPerfectSquare(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	sqrtLow, sqrtHigh := 1, n
	for sqrtHigh-sqrtLow > 1 {
		mid := (sqrtLow + sqrtHigh) / 2
		midSq := mid * mid
		if midSq == n {
			return n
		} else if midSq < n {
			sqrtLow = mid
		} else {
			sqrtHigh = mid
		}
	}
	return (sqrtHigh * sqrtHigh)
}

func doProblem1(data []byte) any {
	ans := 0
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	specStr := scanner.Text()
	if scanner.Err() != nil {
		panic("Scanner error")
	}
	specBits := strings.Split(specStr, ",")
	for _, numStr := range specBits {
		myInt, _ := strconv.Atoi(numStr)
		nextSq := nextPerfectSquare(myInt)
		ans += nextSq - myInt
	}
	return ans
}
func doProblem2(data []byte) any {
	ans := 0
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	specStr := scanner.Text()
	if scanner.Err() != nil {
		panic("Scanner error")
	}
	specBits := strings.Split(specStr, ",")
	for _, numStr := range specBits {
		myInt, _ := strconv.Atoi(numStr)
		nextSq := max(myInt, 8)
		for nextSq%4 != 0 {
			nextSq++
		}
		ans += nextSq - myInt
	}
	return ans
}
func doProblem3(data []byte) any {
	ans := 0
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	specStr := scanner.Text()
	if scanner.Err() != nil {
		panic("Scanner error")
	}
	specBits := strings.Split(specStr, ",")
	for _, numStr := range specBits {
		myInt, _ := strconv.Atoi(numStr)
		nextSq := max(myInt, 8)
		nextSq = nextPerfectSquare(nextSq + 1)
		if nextSq%2 == 0 {
			nextSq = nextPerfectSquare(nextSq + 1)
		}
		nextSq -= 1
		ans += nextSq - myInt
	}
	return ans
}
