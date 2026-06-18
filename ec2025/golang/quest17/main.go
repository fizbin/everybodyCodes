package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"flag"
	"fmt"
	"os"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q17_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q17_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q17_p3.txt`, "the input for part 3")

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
	grid := make([][]byte, 0)
	atLoc := coord{0, 0}
	rows := 0
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
		atCol := strings.Index(scanner.Text(), "@")
		if atCol >= 0 {
			atLoc = coord{rows, atCol}
		}
		rows++
	}
	sum := 0
	for rowIdx, row := range grid {
		for colIdx, ch := range row {
			rowdiff := rowIdx - atLoc.x
			coldiff := colIdx - atLoc.y
			if rowdiff*rowdiff+coldiff*coldiff <= 100 {
				if ch != '@' {
					sum += int(ch - '0')
				}
			}
		}
	}
	return sum
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	atLoc := coord{0, 0}
	rows := 0
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
		atCol := strings.Index(scanner.Text(), "@")
		if atCol >= 0 {
			atLoc = coord{rows, atCol}
		}
		rows++
	}
	destroyed := make(map[coord]bool)
	maxDestructionSum := 0
	maxDestructionRadius := 0
	shouldExitRadiusLoop := false
	for radius := 1; !shouldExitRadiusLoop; radius++ {
		sum := 0
		for rowIdx, row := range grid {
			for colIdx, ch := range row {
				myloc := coord{rowIdx, colIdx}
				if destroyed[myloc] {
					continue
				}
				rowdiff := rowIdx - atLoc.x
				coldiff := colIdx - atLoc.y
				if rowdiff*rowdiff+coldiff*coldiff <= radius*radius {
					destroyed[myloc] = true
					if ch != '@' {
						sum += int(ch - '0')
					}
					if (colIdx == len(row)-1) || (colIdx == 0) || (rowIdx == rows-1) || (rowIdx == 0) {
						shouldExitRadiusLoop = true
					}
				}
			}
		}
		// fmt.Println("p2 debug", sum, radius, sum*radius)
		if sum > maxDestructionSum {
			maxDestructionSum = sum
			maxDestructionRadius = radius
		}
	}
	return maxDestructionSum * maxDestructionRadius
}

type searchLoc struct {
	where     coord
	waypoints int
	time      int
	previous  *searchLoc
}

type searchHeap []searchLoc

func (sh searchHeap) Len() int {
	return len(sh)
}

func (sh searchHeap) Less(i, j int) bool {
	return sh[i].time < sh[j].time
}

func (sh searchHeap) Swap(i, j int) {
	sh[i], sh[j] = sh[j], sh[i]
}

func (sh *searchHeap) Push(x any) {
	*sh = append(*sh, x.(searchLoc))
}

func (sh *searchHeap) Pop() any {
	n := len(*sh)
	retval := (*sh)[n-1]
	*sh = (*sh)[0 : n-1]
	return retval
}

func sqdist(a, b coord) int {
	xdiff := a.x - b.x
	ydiff := a.y - b.y
	return xdiff*xdiff + ydiff*ydiff
}

// waypoints: Just S initially = 0; To the left of @ = 1; exactly below @ = 2; to the right of @ = 3; back at S = 4
func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make(map[coord]int)
	atLoc := coord{0, 0}
	sLoc := coord{0, 0}
	cols := 0
	rows := 0
	for scanner.Scan() {
		cols = max(cols, len(scanner.Text()))
		atCol := strings.Index(scanner.Text(), "@")
		if atCol >= 0 {
			atLoc = coord{rows, atCol}
		}
		sCol := strings.Index(scanner.Text(), "S")
		if sCol >= 0 {
			sLoc = coord{rows, sCol}
		}
		for colIdx, ch := range []byte(scanner.Text()) {
			if '0' <= ch && ch <= '9' {
				grid[coord{rows, colIdx}] = int(ch - '0')
			}
		}
		rows++
	}
	for radius := 2; ; radius++ {
		seen := make(map[searchLoc]bool)
		myHeap := searchHeap{}
		heap.Init(&myHeap)
		heap.Push(&myHeap, searchLoc{sLoc, 0, 0, nil})
		for len(myHeap) > 0 {
			current := heap.Pop(&myHeap).(searchLoc)
			if current.time >= 30*(radius+1) {
				continue
			}
			if sqdist(current.where, atLoc) <= radius*radius {
				continue
			}
			if seen[searchLoc{current.where, current.waypoints, 0, nil}] {
				continue
			}
			if current.where == sLoc {
				if current.waypoints >= 3 {
					if false {
						fmt.Println("Found at radius", radius, "time", current.time, "where", current.where)
						debugPrintAll(grid, rows, cols, &current)
					}
					return current.time * radius
				}
			} else {
				if _, ok := grid[current.where]; !ok {
					continue
				}
			}
			// fmt.Println("Examining", current)
			seen[searchLoc{current.where, current.waypoints, 0, nil}] = true
			for _, jump := range []coord{{-1, 0}, {0, -1}, {1, 0}, {0, 1}} {
				nbr := coord{current.where.x + jump.x, current.where.y + jump.y}
				nbrWaypoints := current.waypoints
				switch current.waypoints {
				case 0:
					if nbr.x == atLoc.x && nbr.y < atLoc.y {
						nbrWaypoints++
					}
				case 1:
					if nbr.x > atLoc.x && nbr.y == atLoc.y {
						nbrWaypoints++
					}
				case 2:
					if nbr.x == atLoc.x && nbr.y > atLoc.y {
						nbrWaypoints++
					}
				case 3:
					if nbr.x > atLoc.x && nbr.y == atLoc.y {
						nbrWaypoints--
					}
				}
				heap.Push(&myHeap, searchLoc{nbr, nbrWaypoints, current.time + grid[nbr], &current})
			}
		}
		if radius > rows {
			return -1
		}
		// fmt.Println("Radius", radius, "failed")
	}
}

func debugPrintAll(grid map[coord]int, rows, cols int, loc *searchLoc) {
	visited := make(map[coord]bool, 0)
	for s := loc; s != nil; s = s.previous {
		visited[s.where] = true
	}
	for rowIdx := range rows {
		for colIdx := range cols {
			where := coord{rowIdx, colIdx}
			if visited[where] {
				fmt.Printf("\x1b[7m")
			}
			if val, ok := grid[where]; ok {
				fmt.Printf("%d", val)
			} else {
				fmt.Printf(".")
			}
			if visited[where] {
				fmt.Printf("\x1b[0m")
			}
		}
		fmt.Println()
	}
}
