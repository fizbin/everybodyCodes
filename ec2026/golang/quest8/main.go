package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"slices"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q08_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q08_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q08_p3.txt`, "the input for part 3")

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
	r, c int
}

func addCoords(a, b coord) coord {
	return coord{a.r + b.r, a.c + b.c}
}

func doProblem1(data []byte) any {
	grid := make([][]byte, 0)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		grid = append(grid, slices.Clone(scanner.Bytes()))
	}

	ans := 0
	shaded := make(map[coord]int)
	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[row]); col++ {
			spotChar := grid[row][col]
			if spotChar >= '0' && spotChar <= '9' {
				here := coord{row, col}
				spotVal := int(spotChar - '0')
				if spotVal > shaded[here] {
					ans++
					for shadeDist := 1; shadeDist < spotVal; shadeDist++ {
						there := coord{row + shadeDist, col + shadeDist}
						shaded[there] = max(shaded[there], spotVal-shadeDist)
					}
				}
			}
		}
	}
	return ans
}

func doProblem2(data []byte) any {
	grid := make([][]byte, 0)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		grid = append(grid, slices.Clone(scanner.Bytes()))
	}

	ans := 0
	for _, shadeDir := range []coord{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, {0, 2}, {0, -2}} {
		colStart := 0
		colEnd := len(grid[0])
		colDir := 1
		if shadeDir.c < 0 {
			colStart, colEnd = colEnd, colStart
			colStart -= 1
			colEnd -= 1
			colDir = -1
		}
		shaded := make(map[coord]int)
		for col := colStart; col != colEnd; col += colDir {
			for row := 0; row < len(grid); row++ {
				spotChar := grid[row][col]
				if spotChar >= '0' && spotChar <= '9' {
					here := coord{row, col}
					spotVal := int(spotChar - '0')
					if spotVal > shaded[here] {
						ans++
						there := here
						for shadeDist := 1; shadeDist < spotVal; shadeDist++ {
							there = addCoords(there, shadeDir)
							shaded[there] = max(shaded[there], spotVal-shadeDist)
						}
					}
				}
			}
		}
	}
	return ans
}

func doProblem3(data []byte) any {
	grid := make([][]byte, 0)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		grid = append(grid, slices.Clone(scanner.Bytes()))
	}

	shadeDirs := []coord{{1, 1}, {1, -1}, {0, -2}, {-1, -1}, {-1, 1}, {0, 2}}
	ans := 0
	warriors := 0
	grow1 := make(map[coord]bool)
	grow2 := make(map[coord]bool)
	grow3 := make(map[coord]bool)
	for round := 0; warriors < 100000; round++ {
		// fmt.Println(round, warriors, ans, "growing:", len(grow1), len(grow2), len(grow3))
		// if round > 30 {
		// 	panic("STOP")
		// }
		grow0 := make(map[coord]bool)
		shadeDir := shadeDirs[round%6]
		colStart := 0
		colEnd := len(grid[0])
		colDir := 1
		if shadeDir.c < 0 {
			colStart, colEnd = colEnd, colStart
			colStart -= 1
			colEnd -= 1
			colDir = -1
		}
		shaded := make(map[coord]int)
		hitSpots := make(map[coord]bool)
		for col := colStart; col != colEnd; col += colDir {
			for row := 0; row < len(grid); row++ {
				// fmt.Println(row, col, colStart, colEnd, colDir)
				spotChar := grid[row][col]
				if spotChar >= '0' && spotChar <= '9' {
					here := coord{row, col}
					spotVal := int(spotChar - '0')
					if spotVal > shaded[here] {
						ans++
						hitSpots[here] = true
						there := here
						for shadeDist := 1; shadeDist < spotVal; shadeDist++ {
							there = addCoords(there, shadeDir)
							shaded[there] = max(shaded[there], spotVal-shadeDist)
						}
					}
				}
			}
		}
		for spot := range hitSpots {
			spotChar := grid[spot.r][spot.c]
			if spotChar == '9' {
				grid[spot.r][spot.c] = 'X'
				grow0[spot] = true
				warriors++
			} else if (spotChar <= '8') && (spotChar >= '0') {
				grid[spot.r][spot.c]++
			}
		}
		for spot := range grow3 {
			grid[spot.r][spot.c] = '1'
		}
		grow3 = grow2
		grow2 = grow1
		grow1 = grow0
	}
	return ans
}
