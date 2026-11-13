package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q14_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q14_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q14_p3.txt`, "the input for part 3")

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

type point struct {
	x, y int
}

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	width := 0
	height := 0
	grid := make(map[point]byte)
	for scanner.Scan() {
		width = max(width, len(scanner.Bytes()))
		for y, ch := range scanner.Bytes() {
			grid[point{height, y}] = ch
		}
		height++
	}
	ans := 0
	// horiz segments
	for rowbelow := range height + 1 {
		for col := range width {
			if grid[point{rowbelow, col}] != grid[point{rowbelow - 1, col}] {
				ans++
			}
		}
	}
	// vert segments
	for row := range height {
		for colright := range width + 1 {
			if grid[point{row, colright}] != grid[point{row, colright - 1}] {
				ans++
			}
		}
	}
	// handCount, handDir, handPos, handMoves := parsePortal(scanner)
	// return workPortal(handCount, handDir, handPos, handMoves)
	return ans
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	width := 0
	height := 0
	pieces := make(map[byte]bool)
	grid := make(map[point]byte)
	graph := make(map[byte]map[byte]int)
	pieces[0] = true
	graph[0] = make(map[byte]int)
	for scanner.Scan() {
		width = max(width, len(scanner.Bytes()))
		for y, ch := range scanner.Bytes() {
			grid[point{height, y}] = ch
			pieces[ch] = true
			if _, ok := graph[ch]; !ok {
				graph[ch] = make(map[byte]int)
			}
		}
		height++
	}
	for rowbelow := range height + 1 {
		for col := range width {
			if grid[point{rowbelow, col}] != grid[point{rowbelow - 1, col}] {
				graph[grid[point{rowbelow, col}]][grid[point{rowbelow - 1, col}]]++
				graph[grid[point{rowbelow - 1, col}]][grid[point{rowbelow, col}]]++
			}
		}
	}
	// vert segments
	for row := range height {
		for colright := range width + 1 {
			if grid[point{row, colright}] != grid[point{row, colright - 1}] {
				graph[grid[point{row, colright}]][grid[point{row, colright - 1}]]++
				graph[grid[point{row, colright - 1}]][grid[point{row, colright}]]++
			}
		}
	}
	ans := 0
	treePieces := make(map[byte]bool)
	treePieces[0] = true
	for len(treePieces) < len(pieces) {
		bestTarget := byte(0)
		bestDist := math.MaxInt
		for existing := range treePieces {
			existingOut := graph[existing]
			for potTgt, dist := range existingOut {
				if dist < bestDist && !treePieces[potTgt] {
					bestTarget = potTgt
					bestDist = dist
				}
			}
		}
		treePieces[bestTarget] = true
		ans += bestDist
	}
	return ans
}
func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	width := 0
	height := 0
	pieces := make(map[byte]bool)
	grid := make(map[point]byte)
	graph := make(map[byte]map[byte]int)
	pieces[0] = true
	graph[0] = make(map[byte]int)
	for scanner.Scan() {
		width = max(width, len(scanner.Bytes()))
		for y, ch := range scanner.Bytes() {
			grid[point{height, y}] = ch
			pieces[ch] = true
			if _, ok := graph[ch]; !ok {
				graph[ch] = make(map[byte]int)
			}
		}
		height++
	}
	for rowbelow := range height + 1 {
		for col := range width {
			if grid[point{rowbelow, col}] != grid[point{rowbelow - 1, col}] {
				graph[grid[point{rowbelow, col}]][grid[point{rowbelow - 1, col}]]++
				graph[grid[point{rowbelow - 1, col}]][grid[point{rowbelow, col}]]++
			}
		}
	}
	// vert segments
	for row := range height {
		for colright := range width + 1 {
			if grid[point{row, colright}] != grid[point{row, colright - 1}] {
				graph[grid[point{row, colright}]][grid[point{row, colright - 1}]]++
				graph[grid[point{row, colright - 1}]][grid[point{row, colright}]]++
			}
		}
	}
	var piecesList []byte
	for piece := range pieces {
		piecesList = append(piecesList, piece)
	}
	bestAns := math.MaxInt
	for idx1 := range len(piecesList) - 1 {
		for idx2 := idx1 + 1; idx2 < len(piecesList); idx2++ {
			chosenA := piecesList[idx1]
			chosenB := piecesList[idx2]
			if chosenA == 0 || chosenB == 0 {
				continue
			}
			treePiecesTotal := make(map[byte]bool)
			treePiecesA := make(map[byte]bool)
			treePiecesB := make(map[byte]bool)
			treePiecesA[chosenA] = true
			treePiecesB[chosenB] = true
			treePiecesTotal[chosenA] = true
			treePiecesTotal[chosenB] = true
			for connected := range graph[chosenA] {
				treePiecesA[connected] = true
				treePiecesTotal[connected] = true
			}
			for connected := range graph[chosenB] {
				treePiecesB[connected] = true
				treePiecesTotal[connected] = true
			}
			ans := 0
			for len(treePiecesTotal) < len(pieces) {
				bestTarget := byte(0)
				bestSource := byte(0)
				bestDist := math.MaxInt
				for source := range treePiecesTotal {
					sourceOut := graph[source]
					for potTgt, dist := range sourceOut {
						if dist < bestDist && !treePiecesTotal[potTgt] {
							bestTarget = potTgt
							bestDist = dist
							bestSource = source
						}
					}
				}
				treePiecesTotal[bestTarget] = true
				if treePiecesA[bestSource] {
					treePiecesA[bestTarget] = true
				}
				if treePiecesB[bestSource] {
					treePiecesB[bestTarget] = true
				}
				ans += bestDist
			}
			// Now, are treeA and treeB connected?
			isConnected := false
			for piece := range treePiecesA {
				if treePiecesB[piece] {
					isConnected = true
					break
				}
			}
			// Damn. we gotta go connect the trees
			if !isConnected {
				bestDist := math.MaxInt
				for source := range treePiecesA {
					for target, dist := range graph[source] {
						if treePiecesB[target] && dist < bestDist {
							bestDist = dist
						}
					}
				}
				ans += bestDist
			}
			if ans < bestAns {
				// fmt.Println("New better answer found; chosen ", string([]byte{chosenA, ' ', chosenB}), "for ans", ans)
				bestAns = ans
			}
		}
	}
	return bestAns
}
