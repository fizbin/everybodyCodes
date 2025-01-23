package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"flag"
	"fmt"
	"os"
)

type point struct{ x, y int }

type heapItem struct {
	totGuess  int
	dist      int
	spot      point
	curheight int
	parent    *heapItem
}

type MyHeap []*heapItem

func (pq MyHeap) Len() int { return len(pq) }

func (pq MyHeap) Less(i, j int) bool {
	return pq[i].totGuess < pq[j].totGuess
}

func (pq MyHeap) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *MyHeap) Push(x any) {
	*pq = append(*pq, x.(*heapItem))
}

func (pq *MyHeap) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func neighbors(current heapItem, gridMap map[point]int, endLoc point) []heapItem {
	retval := make([]heapItem, 0, 6)
	manhattanToEnd := abs(current.spot.x-endLoc.x) + abs(current.spot.y-endLoc.y)
	for _, step := range []int{9, 1} {
		newHeight := (current.curheight + step) % 10
		newGuess := current.dist + 1 + manhattanToEnd + min(newHeight, 10-newHeight)
		retval = append(retval, heapItem{totGuess: newGuess, dist: current.dist + 1, spot: current.spot, curheight: newHeight, parent: &current})
	}
	sx := current.spot.x
	sy := current.spot.y
	for _, loc := range []point{{sx + 1, sy}, {sx - 1, sy}, {sx, sy + 1}, {sx, sy - 1}} {
		if nbrh, ok := gridMap[loc]; ok && nbrh == current.curheight {
			manhattanToEnd := abs(loc.x-endLoc.x) + abs(loc.y-endLoc.y)
			retval = append(retval, heapItem{totGuess: current.dist + 1 + manhattanToEnd + min(nbrh, 10-nbrh), dist: current.dist + 1, spot: loc, curheight: nbrh, parent: &current})
		}
	}
	return retval
}

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q13_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 1:", doProblem(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 2 || argsWithoutProg[1] == "" {
			infile = "../input/everybody_codes_e2024_q13_p2.txt"
		} else {
			infile = argsWithoutProg[1]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 2:", doProblem(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 3 || argsWithoutProg[2] == "" {
			infile = "../input/everybody_codes_e2024_q13_p3.txt"
		} else {
			infile = argsWithoutProg[2]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 3:", doProblem(data))
	}
}

func doProblem(data []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	gridMap := make(map[point]int)
	startLocs := make([]point, 0)
	var endLoc point
	rowIdx := 0
	for scanner.Scan() {
		for colIdx, ch := range scanner.Text() {
			switch ch {
			case 'S':
				sc := point{rowIdx, colIdx}
				gridMap[sc] = 0
				startLocs = append(startLocs, sc)
			case 'E':
				endLoc = point{rowIdx, colIdx}
				gridMap[endLoc] = 0
			case '#', ' ':
				gridMap[point{rowIdx, colIdx}] = 99
			default:
				gridMap[point{rowIdx, colIdx}] = int(ch - '0')
			}
		}
		rowIdx += 1
	}

	type stRec struct {
		where  point
		height int
	}
	beenThere := make(map[stRec]bool)
	pq := make(MyHeap, 0)
	for _, startLoc := range startLocs {
		pq = append(pq, &heapItem{totGuess: 0, dist: 0, spot: startLoc, curheight: 0})
	}
	heap.Init(&pq)
	for len(pq) > 0 {
		current := heap.Pop(&pq).(*heapItem)
		if beenThere[stRec{current.spot, current.curheight}] {
			continue
		}
		beenThere[stRec{current.spot, current.curheight}] = true
		if current.spot.x == endLoc.x && current.spot.y == endLoc.y {
			return current.dist
		}
		nbrs := neighbors(*current, gridMap, endLoc)
		for _, nbr := range nbrs {
			heap.Push(&pq, &nbr)
		}
	}
	panic("Couldn't find way out!")
}
