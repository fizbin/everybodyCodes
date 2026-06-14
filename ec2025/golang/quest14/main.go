package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"slices"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q14_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q14_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q14_p3.txt`, "the input for part 3")

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
	return doPartsOneTwo(data, 10)
}

func doPartsOneTwo(data []byte, nRounds int) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}
	totalLive := 0
	for range nRounds {
		newGrid := make([][]byte, 0)
		for rowIdx, row := range grid {
			newRow := make([]byte, 0, len(row))
			for colIdx, ch := range row {
				neighbors := 0
				if ch == '#' {
					neighbors++
				}
				for _, rowdiff := range []int{-1, 1} {
					for _, coldiff := range []int{-1, 1} {
						nrow := rowdiff + rowIdx
						ncol := coldiff + colIdx
						if (0 <= nrow) && (nrow < len(grid)) && (0 <= ncol) && (ncol < len(row)) && (grid[nrow][ncol] == '#') {
							neighbors++
						}
					}
				}
				if neighbors%2 == 0 {
					newRow = append(newRow, '#')
					totalLive++
				} else {
					newRow = append(newRow, '.')
				}
			}
			newGrid = append(newGrid, newRow)
		}
		grid = newGrid
	}
	return totalLive
}

func doProblem2(data []byte) any {
	return doPartsOneTwo(data, 2025)
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	midgrid := make([][]byte, 0)
	for scanner.Scan() {
		midgrid = append(midgrid, []byte(scanner.Text()))
	}
	nRounds := 1000000000
	grid := make([][]byte, 0, 34)
	for range 34 {
		grid = append(grid, []byte(".................................."))
	}
	ans := 0
	referenceRound := 0
	referenceAns := 0
	initialWait := 5
	var referenceGrid [][]byte
	for round := 0; round <= nRounds; round++ {
		newGrid := make([][]byte, 0)
		for rowIdx, row := range grid {
			newRow := make([]byte, 0, len(row))
			for colIdx, ch := range row {
				neighbors := 0
				if ch == '#' {
					neighbors++
				}
				for _, rowdiff := range []int{-1, 1} {
					for _, coldiff := range []int{-1, 1} {
						nrow := rowdiff + rowIdx
						ncol := coldiff + colIdx
						if (0 <= nrow) && (nrow < len(grid)) && (0 <= ncol) && (ncol < len(row)) && (grid[nrow][ncol] == '#') {
							neighbors++
						}
					}
				}
				if neighbors%2 == 0 {
					newRow = append(newRow, '#')
					// totalLive++
				} else {
					newRow = append(newRow, '.')
				}
			}
			newGrid = append(newGrid, newRow)
		}
		grid = newGrid
		centerMatches := true
		for centRow := (34 - len(midgrid)) / 2; centRow < (34+len(midgrid))/2; centRow++ {
			for centCol := (34 - len(midgrid[0])) / 2; centCol < (34+len(midgrid[0]))/2; centCol++ {
				if grid[centRow][centCol] != midgrid[centRow-(34-len(midgrid))/2][centCol-(34-len(midgrid[0]))/2] {
					centerMatches = false
				}
			}
		}
		if centerMatches {
			livingCount := 0
			for _, row := range grid {
				for _, ch := range row {
					if ch == '#' {
						livingCount++
					}
				}
			}
			// fmt.Println("Round", round+1, "living", livingCount)
			ans += livingCount
			if referenceGrid == nil {
				initialWait--
				if initialWait == 0 {
					referenceAns = ans
					referenceRound = round
					referenceGrid = make([][]byte, 0, len(grid))
					for _, row := range grid {
						referenceGrid = append(referenceGrid, slices.Clone(row))
					}
					// fmt.Println("Reference made at round", round)
				}
			} else {
				matchesReference := true
				for rowIdx, row := range grid {
					if matchesReference {
						for colIdx, ch := range row {
							if ch != referenceGrid[rowIdx][colIdx] {
								matchesReference = false
								break
							}
						}
					}
				}
				if matchesReference {
					ansIncrement := ans - referenceAns
					roundIncrement := round - referenceRound
					cyclesToAdd := (nRounds - 1 - round) / roundIncrement
					// fmt.Print("Jumping from round ", round)
					round += cyclesToAdd * roundIncrement
					ans += cyclesToAdd * ansIncrement
					// fmt.Println(" to round", round)
					referenceGrid = nil
				}
			}
			// for _, row := range grid {
			// 	fmt.Println(string(row))
			// }
		}
	}
	return ans
}
