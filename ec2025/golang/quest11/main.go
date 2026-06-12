package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q11_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q11_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q11_p3.txt`, "the input for part 3")

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
	duckData := make([]int, 0)
	for scanner.Scan() {
		val, _ := strconv.Atoi(scanner.Text())
		duckData = append(duckData, val)
	}
	nextShouldHaveFewer := true
	for round := 1; round <= 10; round++ {
		foundMove := false
		for idx := 0; idx < len(duckData)-1; idx++ {
			if nextShouldHaveFewer {
				if duckData[idx] > duckData[idx+1] {
					foundMove = true
					duckData[idx]--
					duckData[idx+1]++
				}
			} else {
				if duckData[idx] < duckData[idx+1] {
					foundMove = true
					duckData[idx+1]--
					duckData[idx]++
				}
			}
		}
		if nextShouldHaveFewer && !foundMove {
			round--
			nextShouldHaveFewer = false
		}
	}
	chkSum := 0
	for idx := 1; idx <= len(duckData); idx++ {
		chkSum += idx * duckData[idx-1]
	}
	return chkSum
}

// func duckRowDebug(round any, duckData []int) {
// 	duckTot := 0
// 	for _, ducks := range duckData {
// 		duckTot += ducks
// 	}
// 	duckAvg := duckTot / len(duckData)
// 	bigRows := 0
// 	smallRows := 0
// 	duckExcess := 0
// 	for _, ducks := range duckData {
// 		if ducks > duckAvg {
// 			bigRows++
// 			duckExcess += ducks - duckAvg
// 		}
// 		if ducks < duckAvg {
// 			smallRows++
// 		}
// 	}
// 	fmt.Println(round, "Total", duckTot, "rows", len(duckData), "excess", duckExcess, "bigrows", bigRows, "smallrows", smallRows)
// }

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	duckData := make([]int, 0)
	for scanner.Scan() {
		val, _ := strconv.Atoi(scanner.Text())
		duckData = append(duckData, val)
	}
	nextShouldHaveFewer := true
	round := uint64(1)
	for ; ; round++ {
		foundMove := false
		for idx := 0; idx < len(duckData)-1; idx++ {
			if nextShouldHaveFewer {
				if duckData[idx] > duckData[idx+1] {
					foundMove = true
					duckData[idx]--
					duckData[idx+1]++
				}
			} else {
				if duckData[idx] < duckData[idx+1] {
					foundMove = true
					duckData[idx+1]--
					duckData[idx]++
				}
			}
		}
		if !foundMove {
			if nextShouldHaveFewer {
				round--
				nextShouldHaveFewer = false
			} else {
				return (round - 1)
			}
		}
		if !nextShouldHaveFewer {
			// calculate excess, add it to round:
			duckTot := 0
			for _, ducks := range duckData {
				duckTot += ducks
			}
			duckAvg := duckTot / len(duckData)
			duckExcess := uint64(0)
			for _, ducks := range duckData {
				if ducks > duckAvg {
					duckExcess += uint64(ducks - duckAvg)
				}
			}
			return round + duckExcess
		}
	}
}

func doProblem3(data []byte) any {
	return doProblem2(data)
}
