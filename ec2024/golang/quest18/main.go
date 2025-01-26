package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
)

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q18_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 1:", doProblem1(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 2 || argsWithoutProg[1] == "" {
			infile = "../input/everybody_codes_e2024_q18_p2.txt"
		} else {
			infile = argsWithoutProg[1]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 2:", doProblem1(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 3 || argsWithoutProg[2] == "" {
			infile = "../input/everybody_codes_e2024_q18_p3.txt"
		} else {
			infile = argsWithoutProg[2]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 3:", doProblem3(data))
	}
}

type point struct{ x, y int }

func neighbors(grid [][]byte, height int, width int, current point) []point {
	sx := current.x
	sy := current.y
	retval := make([]point, 0, 4)
	for _, candidate := range []point{{sx - 1, sy}, {sx, sy - 1}, {sx, sy + 1}, {sx + 1, sy}} {
		if candidate.x >= 0 && candidate.x < height && candidate.y >= 0 && candidate.y < width {
			if grid[candidate.x][candidate.y] != '#' {
				retval = append(retval, candidate)
			}
		}
	}
	return retval
}

func doProblem1(data []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	palmTrees := make(map[point]bool)
	startLocs := make([]point, 0)
	for scanner.Scan() {
		row := scanner.Bytes()
		for colIdx, ch := range row {
			if ch == 'P' {
				palmTrees[point{len(grid), colIdx}] = true
			} else if ch == '.' && (colIdx == 0 || colIdx == len(row)-1) {
				startLocs = append(startLocs, point{len(grid), colIdx})
			}
		}
		grid = append(grid, slices.Clone(row))
	}
	lastPalmTreeDist := 0
	beenThere := make(map[point]int)
	q := slices.Clone(startLocs)
	for len(palmTrees) > 0 && len(q) > 0 {
		// fmt.Println(len(q))
		current := q[0]
		q = q[1:]
		mydist := beenThere[current]
		if palmTrees[current] {
			lastPalmTreeDist = mydist
			delete(palmTrees, current)
		}
		for _, nbr := range neighbors(grid, len(grid), len(grid[0]), current) {
			if _, ok := beenThere[nbr]; ok {
				continue
			}
			beenThere[nbr] = mydist + 1
			q = append(q, nbr)
		}
	}
	if len(palmTrees) > 0 {
		for rowIdx := range len(grid) {
			for colIdx := range len(grid[0]) {
				ch := grid[rowIdx][colIdx]
				if _, ok := beenThere[point{rowIdx, colIdx}]; ok {
					if ch == 'P' {
						ch = 'W'
					} else {
						ch = '~'
					}
				}
				fmt.Print(string([]byte{ch}))
			}
			fmt.Println()
		}
	}
	return strconv.Itoa(lastPalmTreeDist)
}

func doProblem3(data []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	palmTrees := make(map[point]bool)
	for scanner.Scan() {
		row := scanner.Bytes()
		for colIdx, ch := range row {
			if ch == 'P' {
				palmTrees[point{len(grid), colIdx}] = true
			}
		}
		grid = append(grid, slices.Clone(row))
	}
	type queueType struct {
		spot    point
		palmIdx int
		dist    int
	}
	q := make([]queueType, 0)
	palmCount := 0
	for tree := range palmTrees {
		q = append(q, queueType{tree, palmCount, 0})
		palmCount++
	}
	beenThere := make([][][]bool, 0)
	vistedCount := make([][]int, 0)
	distSum := make([][]int, 0)
	for range len(grid) {
		btRow := make([][]bool, 0)
		vcRow := make([]int, len(grid[0]))
		dsRow := make([]int, len(grid[0]))
		for range len(grid[0]) {
			btRow = append(btRow, make([]bool, palmCount))
		}
		beenThere = append(beenThere, btRow)
		vistedCount = append(vistedCount, vcRow)
		distSum = append(distSum, dsRow)
	}
	goodspotMinSum := math.MaxInt
	for len(q) > 0 {
		current := q[0]
		q = q[1:]
		if beenThere[current.spot.x][current.spot.y][current.palmIdx] {
			continue
		}
		beenThere[current.spot.x][current.spot.y][current.palmIdx] = true
		vistedCount[current.spot.x][current.spot.y]++
		distSum[current.spot.x][current.spot.y] += current.dist
		if vistedCount[current.spot.x][current.spot.y] == palmCount {
			if grid[current.spot.x][current.spot.y] == '.' {
				goodspotMinSum = min(goodspotMinSum, distSum[current.spot.x][current.spot.y])
			}
		}
		for _, nbr := range neighbors(grid, len(grid), len(grid[0]), current.spot) {
			q = append(q, queueType{nbr, current.palmIdx, current.dist + 1})
		}
	}
	return goodspotMinSum
}
