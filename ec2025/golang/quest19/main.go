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

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q19_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q19_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q19_p3.txt`, "the input for part 3")

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

type wallInfo struct {
	x, y, height int
}

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	walls := make([]wallInfo, 0)
	for scanner.Scan() {
		numStrs := strings.Split(scanner.Text(), ",")
		x, _ := strconv.Atoi(numStrs[0])
		y, _ := strconv.Atoi(numStrs[1])
		height, _ := strconv.Atoi(numStrs[2])
		walls = append(walls, wallInfo{x, y, height})
	}
	// height off floor -> negative number of up flaps
	myLocs := make(map[int]int)
	myLocs[0] = 0
	currX := 0
	for _, wall := range walls {
		newLocs := make(map[int]int)
		xDist := wall.x - currX
		for duckY, duckUpFlapsNeg := range myLocs {
			for tgt := wall.y; tgt < wall.y+wall.height; tgt++ {
				yDist := duckY - tgt
				if (xDist%2 == 0) != (yDist%2 == 0) {
					continue
				}
				// upFlaps is zero when yDist == xDist
				// upFlaps is 1 when yDist == xDist - 2
				upFlaps := (xDist - yDist) / 2
				if upFlaps < 0 || upFlaps > xDist {
					continue
				}
				newLocs[tgt] = min(newLocs[tgt], duckUpFlapsNeg-upFlaps)
			}
		}
		currX = wall.x
		myLocs = newLocs
		// fmt.Println(myLocs)
	}
	minUpFlaps := 999999999999
	for _, upFlapsNeg := range myLocs {
		minUpFlaps = min(minUpFlaps, -upFlapsNeg)
	}
	return minUpFlaps
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	walls1 := make([]wallInfo, 0)
	for scanner.Scan() {
		numStrs := strings.Split(scanner.Text(), ",")
		x, _ := strconv.Atoi(numStrs[0])
		y, _ := strconv.Atoi(numStrs[1])
		height, _ := strconv.Atoi(numStrs[2])
		walls1 = append(walls1, wallInfo{x, y, height})
	}
	walls2 := make([][]wallInfo, 0)
	currX := 0
	for idx := 0; idx < len(walls1); idx++ {
		idxhigh := idx + 1
		for ; idxhigh < len(walls1); idxhigh++ {
			if walls1[idxhigh].x != walls1[idx].x {
				break
			}
		}
		walls2 = append(walls2, walls1[idx:idxhigh])
		idx = idxhigh - 1
	}
	// height off floor -> negative number of up flaps
	myLocs := make(map[int]int)
	myLocs[0] = 0
	currX = 0
	for _, walls := range walls2 {
		newLocs := make(map[int]int)
		xDist := walls[0].x - currX
		for _, wall := range walls {
			for duckY, duckUpFlapsNeg := range myLocs {
				for tgt := wall.y; tgt < wall.y+wall.height; tgt++ {
					yDist := duckY - tgt
					if (xDist%2 == 0) != (yDist%2 == 0) {
						continue
					}
					// upFlaps is zero when yDist == xDist
					// upFlaps is 1 when yDist == xDist - 2
					upFlaps := (xDist - yDist) / 2
					if upFlaps < 0 || upFlaps > xDist {
						continue
					}
					newLocs[tgt] = min(newLocs[tgt], duckUpFlapsNeg-upFlaps)
				}
			}
		}
		currX = walls[0].x
		myLocs = newLocs
		// fmt.Println(myLocs)
	}
	minUpFlaps := 999999999999
	for _, upFlapsNeg := range myLocs {
		minUpFlaps = min(minUpFlaps, -upFlapsNeg)
	}
	return minUpFlaps
}

func doProblem3(data []byte) any {
	return doProblem2(data)
}
