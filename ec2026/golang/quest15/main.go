package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"flag"
	"fmt"
	"os"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q15_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q15_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q15_p3.txt`, "the input for part 3")

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

func addpt(a, b point) point {
	return point{a.x + b.x, a.y + b.y}
}

func leftpt(a point) point {
	return point{-a.y, a.x}
}

func rightpt(a point) point {
	return point{a.y, -a.x}
}

type cartState struct {
	position, direction point
	sulfurs             int
}

// only for p2
type cartOrientation struct {
	position, direction point
}

type qType struct {
	dist  int
	state cartState

	// this field is only used in the first part of p3; mostly it's nil
	prevMidlineUp *qType
}

type P1Queue []*qType

// implementation of heap.Interface, mostly from the golang docs
func (pq P1Queue) Len() int { return len(pq) }

func (pq P1Queue) Less(i, j int) bool {
	// We want Pop to give us the lowest dist
	return pq[i].dist < pq[j].dist
}

func (pq P1Queue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *P1Queue) Push(x any) {
	item := x.(*qType)
	*pq = append(*pq, item)
}

func (pq *P1Queue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // for gc
	*pq = old[0 : n-1]
	return item
}

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	width := 0
	height := 0
	grid := make(map[point]byte)
	sSpot := point{}
	sulfurName := make(map[point]int)
	for scanner.Scan() {
		width = max(width, len(scanner.Bytes()))
		for y, ch := range scanner.Bytes() {
			here := point{height, y}
			grid[here] = ch
			if ch == 'S' {
				sSpot = here
			}
			if ch == '@' {
				sulfurName[here] = 1 << len(sulfurName)
			}
		}
		height++
	}
	initState := cartState{position: sSpot, direction: point{0, 1}, sulfurs: 0}
	cache := make(map[cartState]bool)
	neighbors := func(where cartState) []cartState {
		retval := make([]cartState, 0, 3)
		fwdSpot := cartState{position: addpt(where.position, where.direction), direction: where.direction, sulfurs: where.sulfurs}
		leftSpot := cartState{position: addpt(where.position, leftpt(where.direction)), direction: leftpt(where.direction), sulfurs: where.sulfurs}
		rightSpot := cartState{position: addpt(where.position, rightpt(where.direction)), direction: rightpt(where.direction), sulfurs: where.sulfurs}
		for _, st := range []*cartState{&fwdSpot, &leftSpot, &rightSpot} {
			switch grid[st.position] {
			case '.', 'E':
				retval = append(retval, *st)
			case '@':
				st.sulfurs |= sulfurName[st.position]
				retval = append(retval, *st)
			}
		}
		return retval
	}
	var myQueue P1Queue
	myQueue = append(myQueue, &qType{dist: 0, state: initState})
	heap.Init(&myQueue)
	for myQueue.Len() > 0 {
		me := heap.Pop(&myQueue).(*qType)
		if cache[me.state] {
			continue
		}
		cache[me.state] = true
		if grid[me.state.position] == 'E' && me.state.sulfurs+1 == 1<<len(sulfurName) {
			return me.dist
		}
		nbrs := neighbors(me.state)
		for _, nbr := range nbrs {
			heap.Push(&myQueue, &qType{dist: me.dist + 1, state: nbr})
		}
	}

	return -1
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	width := 0
	height := 0
	grid := make(map[point]byte)
	sSpot := point{}
	sulfurName := make(map[point]int)
	for scanner.Scan() {
		width = max(width, len(scanner.Bytes()))
		for y, ch := range scanner.Bytes() {
			here := point{height, y}
			grid[here] = ch
			if ch == 'S' {
				sSpot = here
			}
			if ch == '@' {
				sulfurName[here] = 1 << len(sulfurName)
			}
		}
		height++
	}
	initState := cartState{position: sSpot, direction: point{0, 1}, sulfurs: 0}
	neighbors := func(where cartState) []cartState {
		retval := make([]cartState, 0, 3)
		fwdSpot := cartState{position: addpt(where.position, where.direction), direction: where.direction, sulfurs: where.sulfurs}
		leftSpot := cartState{position: addpt(where.position, leftpt(where.direction)), direction: leftpt(where.direction), sulfurs: where.sulfurs}
		rightSpot := cartState{position: addpt(where.position, rightpt(where.direction)), direction: rightpt(where.direction), sulfurs: where.sulfurs}
		for _, st := range []*cartState{&fwdSpot, &leftSpot, &rightSpot} {
			switch grid[st.position] {
			case '.', 'E':
				retval = append(retval, *st)
			case '@':
				st.sulfurs++
				retval = append(retval, *st)
			}
		}
		return retval
	}
	var myQueue P1Queue
	myQueue = append(myQueue, &qType{dist: 0, state: initState})
	heap.Init(&myQueue)
	// note that the cache key type doesn't include sulfurs this time
	cache := make(map[cartOrientation]int)
	for myQueue.Len() > 0 {
		me := heap.Pop(&myQueue).(*qType)
		meOrientation := cartOrientation{direction: me.state.direction, position: me.state.position}
		if cVal, ok := cache[meOrientation]; ok && cVal >= me.state.sulfurs {
			// Our cache this time also throws stuff away if we've already
			// been here facing the same direction with more sulfurs,
			// not just when the sulfur count was exactly equal
			continue
		}
		cache[meOrientation] = me.state.sulfurs
		if grid[me.state.position] == 'E' && me.state.sulfurs >= 2026 {
			return me.dist
		}
		nbrs := neighbors(me.state)
		for _, nbr := range nbrs {
			heap.Push(&myQueue, &qType{dist: me.dist + 1, state: nbr})
		}
	}

	return -1
}

func doProblem3(data []byte) any {
	sulfurGoal := 202620262026
	scanner := bufio.NewScanner(bytes.NewReader(data))
	width := 0
	height := 0
	grid := make(map[point]byte)
	sSpot := point{}
	sulfurName := make(map[point]int)
	for scanner.Scan() {
		width = max(width, len(scanner.Bytes()))
		for y, ch := range scanner.Bytes() {
			here := point{height, y}
			grid[here] = ch
			if ch == 'S' {
				sSpot = here
			}
			if ch == '@' {
				sulfurName[here] = 1 << len(sulfurName)
			}
		}
		height++
	}
	midLineRecordSpots := make(map[point]bool)
	for spot := range grid {
		if spot.x == sSpot.x {
			switch grid[point{spot.x - 1, spot.y}] {
			case '.', '@':
				midLineRecordSpots[spot] = true
			}
		}
	}
	midLineRecord := make(map[point]*qType)
	initState := cartState{position: sSpot, direction: point{0, 1}, sulfurs: 0}
	neighbors := func(where cartState) []cartState {
		retval := make([]cartState, 0, 3)
		fwdSpot := cartState{position: addpt(where.position, where.direction), direction: where.direction, sulfurs: where.sulfurs}
		leftSpot := cartState{position: addpt(where.position, leftpt(where.direction)), direction: leftpt(where.direction), sulfurs: where.sulfurs}
		rightSpot := cartState{position: addpt(where.position, rightpt(where.direction)), direction: rightpt(where.direction), sulfurs: where.sulfurs}
		for _, st := range []*cartState{&fwdSpot, &leftSpot, &rightSpot} {
			switch grid[st.position] {
			case '.':
				retval = append(retval, *st)
			case 'E':
				if st.sulfurs == sulfurGoal {
					retval = append(retval, *st)
				}
			case '@':
				st.sulfurs++
				if st.sulfurs <= sulfurGoal {
					retval = append(retval, *st)
				}
			}
		}
		return retval
	}
	{
		var myQueue P1Queue
		myQueue = append(myQueue, &qType{dist: 0, state: initState})
		heap.Init(&myQueue)
		cache0 := make(map[cartState]bool)
		cache1 := make(map[cartState]bool)
		for myQueue.Len() > 0 && len(midLineRecord) < len(midLineRecordSpots) {
			me := heap.Pop(&myQueue).(*qType)
			if me.prevMidlineUp == nil {
				if cache0[me.state] {
					continue
				}
				cache0[me.state] = true
			} else {
				if cache1[me.state] {
					continue
				}
				cache1[me.state] = true
			}
			nbrs := neighbors(me.state)
			pMidUp := me.prevMidlineUp
			if me.state.direction == (point{-1, 0}) && midLineRecordSpots[me.state.position] {
				if pMidUp != nil && me.state.position == pMidUp.state.position && midLineRecord[me.state.position] == nil {
					midLineRecord[me.state.position] = me
					continue
				}
				pMidUp = me
			}
			if me.state.direction.x == 0 && me.state.position.x == sSpot.x {
				pMidUp = nil
			}
			for _, nbr := range nbrs {
				heap.Push(&myQueue, &qType{dist: me.dist + 1, state: nbr, prevMidlineUp: pMidUp})
			}
		}
	}
	{
		var myQueue P1Queue
		largestSulfurStep := 0
		for _, loopedState := range midLineRecord {
			sulfurStep := loopedState.state.sulfurs - loopedState.prevMidlineUp.state.sulfurs
			largestSulfurStep = max(largestSulfurStep, sulfurStep)
		}
		for _, loopedState := range midLineRecord {
			sulfurStep := loopedState.state.sulfurs - loopedState.prevMidlineUp.state.sulfurs
			distStep := loopedState.dist - loopedState.prevMidlineUp.dist
			nSteps := (sulfurGoal - 2*largestSulfurStep - loopedState.state.sulfurs) / sulfurStep
			endSulfurs := loopedState.state.sulfurs + nSteps*sulfurStep
			endDist := loopedState.dist + nSteps*distStep
			myQueue = append(myQueue,
				&qType{
					dist: endDist,
					state: cartState{
						sulfurs:   endSulfurs,
						position:  loopedState.state.position,
						direction: loopedState.state.direction},
					prevMidlineUp: loopedState,
				})
			// fmt.Println("dist: ", endDist, "starting from", loopedState.state.position, "with sulfur step", sulfurStep)
		}
		heap.Init(&myQueue)
		cache := make(map[cartState]bool)
		for myQueue.Len() > 0 {
			me := heap.Pop(&myQueue).(*qType)
			if cache[me.state] {
				continue
			}
			cache[me.state] = true
			if grid[me.state.position] == 'E' {
				return me.dist
			}
			//			fmt.Println(me.dist, "starting from", me.prevMidlineUp.state.position)
			nbrs := neighbors(me.state)
			for _, nbr := range nbrs {
				heap.Push(&myQueue, &qType{dist: me.dist + 1, state: nbr, prevMidlineUp: me.prevMidlineUp})
			}
		}
	}
	return -1
}
