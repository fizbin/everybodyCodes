package main

import (
	"bufio"
	"bytes"
	"cmp"
	"container/heap"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q15_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q15_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q15_p3.txt`, "the input for part 3")

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

type coord struct{ x, y int }

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	dataline := scanner.Text()
	wallDirs := strings.Split(dataline, ",")
	parser := regexp.MustCompile(`([LR])(\d+)`)
	grid := make(map[coord]byte)
	spot := coord{0, 0}
	heading := coord{-1, 0}
	for _, wallThing := range wallDirs {
		parts := parser.FindStringSubmatch(wallThing)
		switch parts[1][0] {
		case 'L':
			heading = coord{-heading.y, heading.x}
		case 'R':
			heading = coord{heading.y, -heading.x}
		default:
			panic("Bad direction")
		}
		dist, _ := strconv.Atoi(parts[2])
		for range dist {
			spot = coord{spot.x + heading.x, spot.y + heading.y}
			grid[spot] = '#'
		}
	}
	// spot is now End location
	grid[spot] = 'E'
	dists := make(map[coord]int)
	dists[coord{0, 0}] = 0
	recent := []coord{{0, 0}}
	for /* ever */ {
		newRecent := make([]coord, 0)
		for _, b := range recent {
			thisDist := dists[b]
			for _, jump := range []coord{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
				c := coord{b.x + jump.x, b.y + jump.y}
				if grid[c] == 'E' {
					return thisDist + 1
				}
				if grid[c] == '#' {
					continue
				}
				if _, already := dists[c]; already {
					continue
				}
				dists[c] = thisDist + 1
				newRecent = append(newRecent, c)
			}
		}
		recent = newRecent
	}
}

func doProblem2(data []byte) any {
	return doProblem1(data)
}

type heapItem struct {
	where          coord
	distSoFar      int
	guessRemaining int
}

type searchHeap []heapItem

func (sh searchHeap) Len() int {
	return len(sh)
}

func (sh searchHeap) Less(i, j int) bool {
	return ((sh[i].distSoFar + sh[i].guessRemaining) < (sh[j].distSoFar + sh[j].guessRemaining))
}

func (sh searchHeap) Swap(i, j int) {
	sh[i], sh[j] = sh[j], sh[i]
}

func (sh *searchHeap) Push(x any) {
	*sh = append(*sh, x.(heapItem))
}

func (sh *searchHeap) Pop() any {
	n := len(*sh)
	retval := (*sh)[n-1]
	*sh = (*sh)[0 : n-1]
	return retval
}

// allVals must be sorted!
func goUp[A cmp.Ordered](val A, allVals []A) A {
	// fmt.Print("goUp ", val, " ")
	where, isreal := slices.BinarySearch(allVals, val)
	if where >= len(allVals) {
		// fmt.Println("->", val)
		return val
	}
	if isreal {
		if where >= len(allVals)-1 {
			// fmt.Println("->", val)
			return val
		}
		// fmt.Println("->", allVals[where+1])
		return allVals[where+1]
	}
	// fmt.Println("->", allVals[where])
	return allVals[where]
}

// allVals must be sorted!
func goDown[A cmp.Ordered](val A, allVals []A) A {
	// fmt.Print("goDown ", val, " ")
	where, _ := slices.BinarySearch(allVals, val)
	if where <= 0 {
		// fmt.Println("->", val)
		return val
	}
	// fmt.Println("->", allVals[where-1])
	return allVals[where-1]
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	dataline := scanner.Text()
	wallDirs := strings.Split(dataline, ",")
	parser := regexp.MustCompile(`([LR])(\d+)`)
	spot := coord{0, 0}
	heading := coord{-1, 0}
	realCoordsMap := make(map[int]bool)
	realCoordsMap[-1] = true
	realCoordsMap[0] = true
	realCoordsMap[1] = true
	for _, wallThing := range wallDirs {
		parts := parser.FindStringSubmatch(wallThing)
		switch parts[1][0] {
		case 'L':
			heading = coord{-heading.y, heading.x}
		case 'R':
			heading = coord{heading.y, -heading.x}
		default:
			panic("Bad direction")
		}
		dist, _ := strconv.Atoi(parts[2])
		spot = coord{spot.x + dist*heading.x, spot.y + dist*heading.y}
		realCoordsMap[spot.x] = true
		realCoordsMap[spot.x+1] = true
		realCoordsMap[spot.x-1] = true
		realCoordsMap[spot.y] = true
		realCoordsMap[spot.y+1] = true
		realCoordsMap[spot.y-1] = true
	}
	realCoords := make([]int, 0, len(realCoordsMap))
	for v := range realCoordsMap {
		realCoords = append(realCoords, v)
	}
	slices.Sort(realCoords)
	grid := make(map[coord]byte)
	spot = coord{0, 0}
	heading = coord{-1, 0}
	for _, wallThing := range wallDirs {
		parts := parser.FindStringSubmatch(wallThing)
		switch parts[1][0] {
		case 'L':
			heading = coord{-heading.y, heading.x}
		case 'R':
			heading = coord{heading.y, -heading.x}
		default:
			panic("Bad direction")
		}
		dist, _ := strconv.Atoi(parts[2])
		tspot := coord{spot.x + dist*heading.x, spot.y + dist*heading.y}
		for spot != tspot {
			switch heading {
			case coord{1, 0}:
				spot.x = goUp(spot.x, realCoords)
			case coord{0, 1}:
				spot.y = goUp(spot.y, realCoords)
			case coord{-1, 0}:
				spot.x = goDown(spot.x, realCoords)
			case coord{0, -1}:
				spot.y = goDown(spot.y, realCoords)
			}
			grid[spot] = '#'
		}
	}
	// spot is now End location
	grid[spot] = 'E'
	seen := make(map[coord]bool)
	sheap := searchHeap{{coord{0, 0}, 0, 0}}
	heap.Init(&sheap)
	// fmt.Println(realCoords)
	for len(sheap) > 0 {
		working := heap.Pop(&sheap).(heapItem)
		here := working.where
		if seen[here] {
			continue
		}
		// fmt.Println("at", here)
		seen[here] = true
		for _, jump := range []coord{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
			c := here
			jumpSize := 0
			switch jump {
			case coord{1, 0}:
				jumpSize = c.x
				c.x = goUp(c.x, realCoords)
				jumpSize = c.x - jumpSize
			case coord{0, 1}:
				jumpSize = c.y
				c.y = goUp(c.y, realCoords)
				jumpSize = c.y - jumpSize
			case coord{-1, 0}:
				jumpSize = c.x
				c.x = goDown(c.x, realCoords)
				jumpSize -= c.x
			case coord{0, -1}:
				jumpSize = c.y
				c.y = goDown(c.y, realCoords)
				jumpSize -= c.y
			}
			if grid[c] == 'E' {
				return working.distSoFar + jumpSize
			}
			if grid[c] == '#' {
				continue
			}
			estimateX := c.x - spot.x
			estimateY := c.y - spot.y
			if estimateX < 0 {
				estimateX = -estimateX
			}
			if estimateY < 0 {
				estimateY = -estimateY
			}
			heap.Push(&sheap, heapItem{c, working.distSoFar + jumpSize, estimateX + estimateY})
		}
	}
	return 0
}
