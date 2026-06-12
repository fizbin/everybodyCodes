package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q12_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q12_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q12_p3.txt`, "the input for part 3")

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

type coord struct {
	x, y int
}

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make(map[coord]byte)
	row := 0
	for scanner.Scan() {
		for col, val := range scanner.Bytes() {
			if '0' <= val && val <= '9' {
				grid[coord{row, col}] = val
			}
		}
		row++
	}
	exploded := make(map[coord]bool)
	exploded[coord{0, 0}] = true
	recent := []coord{{0, 0}}
	for len(recent) > 0 {
		newrecent := make([]coord, 0, len(recent))
		for _, boom := range recent {
			for _, boomdir := range []coord{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
				nxt := coord{boom.x + boomdir.x, boom.y + boomdir.y}
				if !exploded[nxt] {
					nxtval := grid[nxt]
					if '0' <= nxtval && nxtval <= grid[boom] {
						newrecent = append(newrecent, nxt)
						exploded[nxt] = true
					}
				}
			}
		}
		recent = newrecent
	}
	return len(exploded)
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make(map[coord]byte)
	row := 0
	col := 0
	for scanner.Scan() {
		for c, val := range scanner.Bytes() {
			if '0' <= val && val <= '9' {
				grid[coord{row, c}] = val
				col = max(c+1, col)
			}
		}
		row++
	}
	exploded := make(map[coord]bool)
	exploded[coord{0, 0}] = true
	exploded[coord{row - 1, col - 1}] = true
	recent := []coord{{0, 0}, {row - 1, col - 1}}
	for len(recent) > 0 {
		newrecent := make([]coord, 0, len(recent))
		for _, boom := range recent {
			for _, boomdir := range []coord{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
				nxt := coord{boom.x + boomdir.x, boom.y + boomdir.y}
				if !exploded[nxt] {
					nxtval := grid[nxt]
					if '0' <= nxtval && nxtval <= grid[boom] {
						newrecent = append(newrecent, nxt)
						exploded[nxt] = true
					}
				}
			}
		}
		recent = newrecent
	}
	return len(exploded)
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make(map[coord]byte)
	row := 0
	col := 0
	for scanner.Scan() {
		for c, val := range scanner.Bytes() {
			if '0' <= val && val <= '9' {
				grid[coord{row, c}] = val
				col = max(c+1, col)
			}
		}
		row++
	}
	maxFirstCoord := coord{-1, -1}
	var maxFirstMap map[coord]bool
	for r := range row {
		for c := range col {
			me := coord{r, c}
			if maxFirstMap[me] {
				continue
			}
			exploded := explode(grid, []coord{me}, nil)
			if maxFirstMap == nil || len(exploded) > len(maxFirstMap) {
				maxFirstCoord = me
				maxFirstMap = exploded
			}
		}
	}
	maxSecondCoord := coord{-1, -1}
	var maxSecondMap map[coord]bool
	for r := range row {
		for c := range col {
			me := coord{r, c}
			if maxFirstMap[me] || maxSecondMap[me] {
				continue
			}
			exploded := explode(grid, []coord{me}, maxFirstMap)
			if maxSecondMap == nil || len(exploded) > len(maxSecondMap) {
				maxSecondCoord = me
				maxSecondMap = exploded
			}
		}
	}
	maxThirdCoord := coord{-1, -1}
	var maxThirdMap map[coord]bool
	for r := range row {
		for c := range col {
			me := coord{r, c}
			if maxSecondMap[me] || maxThirdMap[me] {
				continue
			}
			exploded := explode(grid, []coord{me}, maxSecondMap)
			if maxThirdMap == nil || len(exploded) > len(maxThirdMap) {
				maxThirdCoord = me
				maxThirdMap = exploded
			}
		}
	}
	return len(explode(grid, []coord{maxFirstCoord, maxSecondCoord, maxThirdCoord}, nil))
}

func explode(grid map[coord]byte, start []coord, alreadyExploded map[coord]bool) map[coord]bool {
	exploded := make(map[coord]bool)
	recent := make([]coord, 0, len(start))
	for _, c := range start {
		exploded[c] = true
		recent = append(recent, c)
	}
	for k, v := range alreadyExploded {
		if v {
			exploded[k] = true
		}
	}
	for len(recent) > 0 {
		newrecent := make([]coord, 0, len(recent))
		for _, boom := range recent {
			for _, boomdir := range []coord{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
				nxt := coord{boom.x + boomdir.x, boom.y + boomdir.y}
				if !exploded[nxt] {
					nxtval := grid[nxt]
					if '0' <= nxtval && nxtval <= grid[boom] {
						newrecent = append(newrecent, nxt)
						exploded[nxt] = true
					}
				}
			}
		}
		recent = newrecent
	}
	return exploded
}
