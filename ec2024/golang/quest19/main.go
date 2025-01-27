package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
)

var showFull = flag.Bool("verbose", false, "Add to see the ASCII art")

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q19_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 1:", doProblem1(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 2 || argsWithoutProg[1] == "" {
			infile = "../input/everybody_codes_e2024_q19_p2.txt"
		} else {
			infile = argsWithoutProg[1]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 2:", doProblem2(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 3 || argsWithoutProg[2] == "" {
			infile = "../input/everybody_codes_e2024_q19_p3.txt"
		} else {
			infile = argsWithoutProg[2]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 3:", doProblem3(data))
	}
}

// affects grid "in place", dir is 1 for R/clockwise, -1 for left/CCW
func rotate(grid [][]byte, key string) {
	// 0 1 2
	// 7   3
	// 6 5 4
	rowOffsets := []int{-1, -1, -1, 0, 1, 1, 1, 0}
	colOffsets := []int{-1, 0, 1, 1, 1, 0, -1, -1}
	keyIdx := 0
	var dir int
	for rowIdx := 1; rowIdx < len(grid)-1; rowIdx++ {
		for colIdx := 1; colIdx < len(grid[0])-1; colIdx++ {
			dirByte := key[keyIdx]
			keyIdx++
			keyIdx %= len(key)
			switch dirByte {
			case 'L':
				dir = -1
			case 'R':
				dir = 1
			default:
				log.Fatal("Bad key character", dirByte)
			}
			oldStuff := make([]byte, 8)
			for idx := range rowOffsets {
				oldStuff[idx] = grid[rowIdx+rowOffsets[idx]][colIdx+colOffsets[idx]]
			}
			for idx := range rowOffsets {
				grid[rowIdx+rowOffsets[idx]][colIdx+colOffsets[idx]] = oldStuff[(idx-dir+8)%8]
			}
		}
	}
}

func findMessage(grid [][]byte) string {
	for _, row := range grid {
		c1Idx := slices.Index(row, '>')
		if c1Idx == -1 {
			continue
		}
		c2Idx := slices.Index(row, '<')
		if c2Idx < c1Idx {
			return "***BAD FORMAT***"
		}
		return string(row[c1Idx+1 : c2Idx])
	}
	return "***NOT FOUND***"
}

func doProblem1(data []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	key := scanner.Text()
	scanner.Scan() // blank line
	grid := make([][]byte, 0)
	for scanner.Scan() {
		grid = append(grid, slices.Clone(scanner.Bytes()))
	}
	rotate(grid, key)
	if *showFull {
		for _, row := range grid {
			fmt.Println(string(row))
		}
	}
	return findMessage(grid)
}

func doProblem2(data []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	key := scanner.Text()
	scanner.Scan() // blank line
	grid := make([][]byte, 0)
	for scanner.Scan() {
		grid = append(grid, slices.Clone(scanner.Bytes()))
	}
	for range 100 {
		rotate(grid, key)
	}
	if *showFull {
		for _, row := range grid {
			fmt.Println(string(row))
		}
	}
	return findMessage(grid)
}

type point struct{ x, y int }

func doProblem3(data []byte) string {
	const bigRotCount = 1048576000
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	key := scanner.Text()
	scanner.Scan() // blank line
	grid := make([][]byte, 0)
	for scanner.Scan() {
		grid = append(grid, slices.Clone(scanner.Bytes()))
	}
	// fmt.Println("DBG key size:", len(key))
	// fmt.Println("DBG grid dim:", len(grid), "rows of", len(grid[0]), "cols")

	blankGrid := make([][]byte, len(grid))
	for rowIdx, row := range grid {
		blankGrid[rowIdx] = make([]byte, len(row))
	}
	transform := make(map[point]point)
	for rowIdx := range grid {
		for colIdx := range len(grid[0]) {
			startPoint := point{rowIdx, colIdx}
			blankGrid[rowIdx][colIdx] = 1
			rotate(blankGrid, key)
			var endPoint point
			var foundIt bool
			for endRowIdx, row := range blankGrid {
				cidx := slices.Index(row, 1)
				if cidx >= 0 {
					endPoint = point{endRowIdx, cidx}
					row[cidx] = 0
					foundIt = true
					break
				}
			}
			if !foundIt {
				panic("Rotated into nothing????!?")
			}
			transform[startPoint] = endPoint
		}
	}
	// fmt.Println("DBG transform len", len(transform))

	cycles := make([][]point, 0)
	for len(transform) > 0 {
		var cPoint point
		for spot := range transform {
			cPoint = spot
			break
		}
		cPoint0 := cPoint
		cycle := make([]point, 0)
		for {
			if cPointN, ok := transform[cPoint]; ok {
				cycle = append(cycle, cPoint)
				delete(transform, cPoint)
				cPoint = cPointN
			} else {
				if cPoint != cPoint0 {
					log.Fatal("Started at ", cPoint0, "but ended at", cPoint)
				}
				break
			}
		}
		// fmt.Println("DBG cycle len", len(cycle), "remaining in transform", len(transform))
		slices.Reverse(cycle)
		cycles = append(cycles, cycle)
	}

	for _, cycle := range cycles {
		stash := make([]byte, len(cycle))
		for idx, pt := range cycle {
			stash[idx] = grid[pt.x][pt.y]
		}
		for idx, pt := range cycle {
			srcIdx := (idx + bigRotCount) % len(cycle)
			grid[pt.x][pt.y] = stash[srcIdx]
		}
	}

	if *showFull {
		for _, row := range grid {
			fmt.Println(string(row))
		}
	}
	return findMessage(grid)
}
