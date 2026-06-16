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
	// startTime := time.Now()
	maxFirstBoom := 0
	maxFirstCoord := coord{-1, -1}
	knownExplosions := make(map[coord]map[coord]bool)
	for focusVal := byte('0'); focusVal <= '9'; focusVal++ {
		for me, megrid := range grid {
			if (focusVal != megrid) || (knownExplosions[me] != nil) {
				continue
			}
			exploded := explode(grid, me, knownExplosions)
			if len(exploded) > maxFirstBoom {
				maxFirstCoord = me
				maxFirstBoom = len(exploded)
			}
		}
	}
	maxFirstMap := knownExplosions[maxFirstCoord]
	// fmt.Println("So far1", time.Since(startTime))
	maxSecondCoord := coord{-1, -1}
	maxSecondBoom := 0
	for r := range row {
		for c := range col {
			me := coord{r, c}
			if knownExplosions[maxSecondCoord][me] || maxFirstMap[me] {
				continue
			}
			newBoom := 0
			for c := range knownExplosions[me] {
				if !maxFirstMap[c] {
					newBoom++
				}
			}
			if newBoom > maxSecondBoom {
				maxSecondBoom = newBoom
				maxSecondCoord = me
			}
		}
	}
	// fmt.Println("So far2", time.Since(startTime))
	maxThirdCoord := coord{-1, -1}
	maxThirdBoom := 0
	for r := range row {
		for c := range col {
			me := coord{r, c}
			if knownExplosions[maxThirdCoord][me] || knownExplosions[maxSecondCoord][me] || maxFirstMap[me] {
				continue
			}
			newBoom := 0
			for c := range knownExplosions[me] {
				if !(maxFirstMap[c] || knownExplosions[maxSecondCoord][c]) {
					newBoom++
				}
			}
			if newBoom > maxThirdBoom {
				maxThirdBoom = newBoom
				maxThirdCoord = me
			}
		}
	}
	// fmt.Println("So far3", time.Since(startTime))
	return (maxFirstBoom + maxSecondBoom + maxThirdBoom)
}

func explode(grid map[coord]byte, start coord, knownExplosions map[coord]map[coord]bool) map[coord]bool {
	exploded := make(map[coord]bool)
	recent := []coord{start}
	startval := grid[start]
	sameAsStart := make([]coord, 0)
	exploded[start] = true
	for len(recent) > 0 {
		newrecent := make([]coord, 0, len(recent))
		for _, boom := range recent {
			for _, boomdir := range []coord{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
				nxt := coord{boom.x + boomdir.x, boom.y + boomdir.y}
				if !exploded[nxt] {
					nxtval := grid[nxt]
					if '0' <= nxtval && nxtval <= grid[boom] {
						if prevKnown, ok := knownExplosions[nxt]; ok {
							for k := range prevKnown {
								exploded[k] = true
							}
						} else {
							newrecent = append(newrecent, nxt)
							exploded[nxt] = true
							if nxtval == startval {
								sameAsStart = append(sameAsStart, nxt)
							}
						}
					}
				}
			}
		}
		recent = newrecent
	}
	knownExplosions[start] = exploded
	for _, c := range sameAsStart {
		knownExplosions[c] = exploded
	}
	return exploded
}
