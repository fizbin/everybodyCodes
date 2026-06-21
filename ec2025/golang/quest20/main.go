package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q20_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q20_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q20_p3.txt`, "the input for part 3")

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
	grid := make([][]byte, 0)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}
	tot := 0

	pch := byte('.')
	for rowIdx, row := range grid {
		for colIdx, ch := range row {
			if ch == 'T' {
				if pch == 'T' {
					tot++
				}
				if (rowIdx > 0) && (colIdx%2 == rowIdx%2) && (grid[rowIdx-1][colIdx] == 'T') {
					tot++
				}
			}
			pch = ch
		}
	}

	return tot
}

type coord struct{ x, y int }

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	sSpot := coord{-1, -1}
	rowIdx := 0
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
		sIdx := strings.Index(scanner.Text(), "S")
		if sIdx >= 0 {
			sSpot = coord{rowIdx, sIdx}
		}
		rowIdx++
	}
	seen := make(map[coord]int)
	seen[sSpot] = 0
	recent := []coord{sSpot}
	for {
		newRecent := make([]coord, 0)
		for _, r := range recent {
			newRs := make([]coord, 0, 3)
			if r.y > 0 {
				newRs = append(newRs, coord{r.x, r.y - 1})
			}
			if r.y < len(grid[r.x])-1 {
				newRs = append(newRs, coord{r.x, r.y + 1})
			}
			if r.x%2 == r.y%2 {
				if r.x > 0 {
					newRs = append(newRs, coord{r.x - 1, r.y})
				}
			} else {
				if r.x < len(grid)-1 {
					newRs = append(newRs, coord{r.x + 1, r.y})
				}
			}
			for _, newR := range newRs {
				if _, already := seen[newR]; already {
					continue
				}
				if grid[newR.x][newR.y] == 'T' {
					newRecent = append(newRecent, newR)
					seen[newR] = seen[r] + 1
				}
				if grid[newR.x][newR.y] == 'E' {
					return seen[r] + 1
				}
			}
		}
		recent = newRecent
	}
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	sSpot := coord{-1, -1}
	rowIdx := 0
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
		sIdx := strings.Index(scanner.Text(), "S")
		if sIdx >= 0 {
			sSpot = coord{rowIdx, sIdx}
		}
		rowIdx++
	}
	seen := make(map[coord]int)
	seen[sSpot] = 0
	recent := []coord{sSpot}
	for {
		newRecent := make([]coord, 0)
		for _, rold := range recent {
			r := transformCCW(rold, rowIdx)
			newRs := make([]coord, 0, 4)
			newRs = append(newRs, r)
			if r.y > 0 {
				newRs = append(newRs, coord{r.x, r.y - 1})
			}
			if r.y < len(grid[r.x])-1 {
				newRs = append(newRs, coord{r.x, r.y + 1})
			}
			if r.x%2 == r.y%2 {
				if r.x > 0 {
					newRs = append(newRs, coord{r.x - 1, r.y})
				}
			} else {
				if r.x < len(grid)-1 {
					newRs = append(newRs, coord{r.x + 1, r.y})
				}
			}
			for _, newR := range newRs {
				if _, already := seen[newR]; already {
					continue
				}
				if grid[newR.x][newR.y] == 'T' {
					newRecent = append(newRecent, newR)
					seen[newR] = seen[rold] + 1
				}
				if grid[newR.x][newR.y] == 'E' {
					return seen[rold] + 1
				}
			}
		}
		recent = newRecent
	}
}

func transformCCW(duck coord, rowCount int) coord {
	// a duck at (rowCount-1, rowCount-1) goes to (0, 2*rowCount-2)
	// a duck at (0, 2*rowCount-2) goes to (0, 0)
	// (0, 0) goes to (rowCount-1, rowCount-1)
	// a row at idx N has edge pieces at column N and at colum 2*rowCount-2-N
	// (N, N) goes to (rowCount-1-N, rowCount-1+N)
	// (N, 2*rowCount-2-N) goes to (0, 2*N)

	working := coord{duck.x, duck.x}
	workingTx := coord{rowCount - 1 - duck.x, rowCount - 1 + duck.x}

	for working != duck {
		// move working one to the right, move Tx one to the up-left
		working.y++
		if workingTx.x%2 == workingTx.y%2 {
			workingTx.x--
		} else {
			workingTx.y--
		}
	}
	return workingTx
}
