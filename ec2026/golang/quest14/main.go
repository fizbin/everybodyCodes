package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"slices"
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

	// Prim's algorithm was fast enough for part 2, but for part 3 we're
	// redoing stuff enough times that I want to a faster approach, so try
	// Kruskal's instead (but with a simple node->component map; we aren't
	// going to go implement a full union-find for this)
	type edge struct {
		length         int
		source, target byte
	}

	edges := make([]edge, 0, len(graph)/2)
	for source, outgoing := range graph {
		for target, dist := range outgoing {
			if source < target {
				edges = append(edges, edge{dist, source, target})
			}
		}
	}
	slices.SortFunc(edges, func(a, b edge) int {
		if a.length != b.length {
			return a.length - b.length
		}
		if a.source != b.source {
			return (int(a.source) - int(b.source))
		}
		return (int(a.target) - int(b.target))
	})

	var piecesList []byte
	for piece := range pieces {
		piecesList = append(piecesList, piece)
	}
	slices.Sort(piecesList)
	bestAns := math.MaxInt
	// skipping index 0 which will == 0
	for idx1 := 1; idx1 < len(piecesList); idx1++ {
		for idx2 := idx1 + 1; idx2 < len(piecesList); idx2++ {
			chosenA := piecesList[idx1]
			chosenB := piecesList[idx2]
			classToPieces := make(map[byte][]byte)
			pieceToClass := make(map[byte]byte)
			for _, piece := range piecesList {
				pieceToClass[piece] = piece
			}
			{
				aClass := pieceToClass[chosenA]
				for connected := range graph[chosenA] {
					aClass = min(aClass, pieceToClass[connected])
				}
				pieceToClass[chosenA] = aClass
				for connected := range graph[chosenA] {
					pieceToClass[connected] = aClass
				}
				bWorldTouchesAWorld := false
				bClass := pieceToClass[chosenB]
				// I know, 4 is overkill most of the time, but not always
				for connected := range graph[chosenB] {
					if pieceToClass[connected] == aClass {
						bWorldTouchesAWorld = true
					}
					bClass = min(bClass, pieceToClass[connected])
				}
				pieceToClass[chosenB] = bClass
				for connected := range graph[chosenB] {
					pieceToClass[connected] = bClass
				}
				if bWorldTouchesAWorld && bClass < aClass {
					pieceToClass[chosenA] = bClass
					for connected := range graph[chosenA] {
						pieceToClass[connected] = bClass
					}
				}
			}
			for piece, class := range pieceToClass {
				classToPieces[class] = append(classToPieces[class], piece)
			}

			ans := 0
			for edgesIdx := 0; len(classToPieces) > 1; edgesIdx++ {
				myEdge := &edges[edgesIdx]
				if pieceToClass[myEdge.source] != pieceToClass[myEdge.target] {
					newClass := min(pieceToClass[myEdge.source], pieceToClass[myEdge.target])
					oldClass := max(pieceToClass[myEdge.source], pieceToClass[myEdge.target])
					for _, piece := range classToPieces[oldClass] {
						pieceToClass[piece] = newClass
					}
					classToPieces[newClass] = append(classToPieces[newClass], classToPieces[oldClass]...)
					delete(classToPieces, oldClass)
					ans += myEdge.length
				}
			}
			if ans < bestAns {
				// fmt.Println("New better answer found; chosen ", string([]byte{chosenA, ' ', chosenB}), "for ans", ans)
				bestAns = ans
			}
		}
	}
	return bestAns
}
