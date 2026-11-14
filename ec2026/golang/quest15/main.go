package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q15_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q15_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q15_p3.txt`, "the input for part 3")

const doingPart = 3

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
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

// I call the coords x and y, but they're really row and column - "x" is
// vertical and increases as you go down, "y" is horizontal
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

type qType struct {
	dist  int
	state cartState

	// this field is only used in p3; mostly it's nil
	loopStart *qType
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
		fwdSpot := cartState{
			position:  addpt(where.position, where.direction),
			direction: where.direction}
		leftSpot := cartState{
			position:  addpt(where.position, leftpt(where.direction)),
			direction: leftpt(where.direction)}
		rightSpot := cartState{
			position:  addpt(where.position, rightpt(where.direction)),
			direction: rightpt(where.direction)}
		for _, st := range []*cartState{&fwdSpot, &leftSpot, &rightSpot} {
			st.sulfurs = where.sulfurs
			switch grid[st.position] {
			case '.':
				retval = append(retval, *st)
			case 'E':
				if st.sulfurs+1 == 1<<len(sulfurName) {
					retval = append(retval, *st)
				}
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
		if grid[me.state.position] == 'E' {
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
	eSpot := point{}
	for scanner.Scan() {
		width = max(width, len(scanner.Bytes()))
		for y, ch := range scanner.Bytes() {
			here := point{height, y}
			grid[here] = ch
			switch ch {
			case 'S':
				sSpot = here
			case 'E':
				eSpot = here
			}
		}
		height++
	}
	initState := cartState{position: sSpot, direction: point{0, 1}, sulfurs: 0}
	neighbors := func(where *cartState) []cartState {
		retval := make([]cartState, 0, 3)
		fwdSpot := cartState{
			position:  addpt(where.position, where.direction),
			direction: where.direction}
		leftSpot := cartState{
			position:  addpt(where.position, leftpt(where.direction)),
			direction: leftpt(where.direction)}
		rightSpot := cartState{
			position:  addpt(where.position, rightpt(where.direction)),
			direction: rightpt(where.direction)}
		for _, st := range []*cartState{&fwdSpot, &leftSpot, &rightSpot} {
			st.sulfurs = where.sulfurs
			switch grid[st.position] {
			case '.':
				retval = append(retval, *st)
			case 'E':
				if st.sulfurs >= 2026 {
					retval = append(retval, *st)
				}
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

	// the cache key type doesn't include sulfurs this time; this time the
	// cache is a map to number of sulfurs, so we can throw away paths if
	// we already handled at least that many sulfurs at that spot+dir
	type cartOrientation struct {
		position, direction point
	}
	cache := make(map[cartOrientation]int)
	for myQueue.Len() > 0 {
		me := heap.Pop(&myQueue).(*qType)
		meOrientation := cartOrientation{
			direction: me.state.direction,
			position:  me.state.position}
		if cVal, ok := cache[meOrientation]; ok && cVal >= me.state.sulfurs {
			continue
		}
		cache[meOrientation] = me.state.sulfurs
		if me.state.position == eSpot {
			return me.dist
		}
		nbrs := neighbors(&me.state)
		ndist := me.dist + 1
		// optimization loop; the idea is that once you're in one of the
		// half-loops, you'll have to finish the loop because you can't
		// turn, so might as well just follow all those steps now.
		for len(nbrs) == 1 && nbrs[0].position != eSpot {
			nbrs = neighbors(&nbrs[0])
			ndist++
		}
		for _, nbr := range nbrs {
			heap.Push(&myQueue, &qType{dist: ndist, state: nbr})
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
	var eSpot, sSpot point
	for scanner.Scan() {
		width = max(width, len(scanner.Bytes()))
		for y, ch := range scanner.Bytes() {
			here := point{height, y}
			grid[here] = ch
			switch ch {
			case 'S':
				sSpot = here
			case 'E':
				eSpot = here
			}
		}
		height++
	}

	// Strategy:
	// - Do a Dijkstra run through the grid to find all the loops - a loop
	//   starts when you step off the midline and is complete when
	//   you get back to the starting spot.
	// - "fast forward" all those loops to just a bit before we have exactly
	//   enough sulfur
	// - Do a second Dijkstra run with all those fast-forwarded loops as
	//   startpoints

	loopRecord := make(map[point]*qType) // loop start position -> looped queue entry

	initState := cartState{position: sSpot, direction: point{0, 1}, sulfurs: 0}
	neighbors := func(where *cartState) []cartState {
		retval := make([]cartState, 0, 3)
		fwdSpot := cartState{
			position:  addpt(where.position, where.direction),
			direction: where.direction,
			sulfurs:   where.sulfurs}
		leftSpot := cartState{
			position:  addpt(where.position, leftpt(where.direction)),
			direction: leftpt(where.direction),
			sulfurs:   where.sulfurs}
		rightSpot := cartState{
			position:  addpt(where.position, rightpt(where.direction)),
			direction: rightpt(where.direction),
			sulfurs:   where.sulfurs}
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
		type cacheType struct {
			state      cartState
			loopStartP point
		}
		cache := make(map[cacheType]bool)
		for myQueue.Len() > 0 {
			me := heap.Pop(&myQueue).(*qType)
			cacheVal := cacheType{state: me.state}
			if me.loopStart != nil {
				loopStartP := me.loopStart.state.position
				// bail if someone got our loop already
				if loopRecord[loopStartP] != nil {
					continue
				}
				cacheVal.loopStartP = loopStartP
			}
			if cache[cacheVal] {
				continue
			}
			cache[cacheVal] = true
			loopStart := me.loopStart
			// blank the "loop start" pointer on midline travel
			if me.state.direction.x == 0 && me.state.position.x == sSpot.x {
				loopStart = nil
			}
			if me.state.direction.y == 0 && me.state.position.x == sSpot.x+me.state.direction.x {
				// vertical travel and we just stepped off the midline

				if loopStart == nil { // new loop
					if loopRecord[me.state.position] != nil {
						// we already have this loop
						continue
					}
					loopStart = me
				} else if me.loopStart.state.position == me.state.position {
					// we finished a loop
					loopRecord[me.state.position] = me
					continue
				}
			}
			nbrs := neighbors(&me.state)
			ndist := me.dist + 1
			for len(nbrs) == 1 && nbrs[0].position != eSpot {
				nbrs = neighbors(&nbrs[0])
				ndist++
			}
			for _, nbr := range nbrs {
				heap.Push(&myQueue, &qType{dist: ndist, state: nbr, loopStart: loopStart})
			}
		}
	}
	{
		var myQueue P1Queue
		for _, loopedState := range loopRecord {
			sulfurStep := loopedState.state.sulfurs - loopedState.loopStart.state.sulfurs
			distStep := loopedState.dist - loopedState.loopStart.dist
			nSteps := (sulfurGoal - 2*sulfurStep - loopedState.state.sulfurs) / sulfurStep
			endSulfurs := loopedState.state.sulfurs + nSteps*sulfurStep
			endDist := loopedState.dist + nSteps*distStep
			myQueue = append(myQueue,
				&qType{
					dist: endDist,
					state: cartState{
						sulfurs:   endSulfurs,
						position:  loopedState.state.position,
						direction: loopedState.state.direction},
					loopStart: loopedState,
				})
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
				// fmt.Println(me.dist, "starting from", me.prevMidlineUp.state.position)
				return me.dist
			}
			ndist := me.dist + 1
			nbrs := neighbors(&me.state)
			for len(nbrs) == 1 && nbrs[0].position != eSpot {
				nbrs = neighbors(&nbrs[0])
				ndist++
			}
			for _, nbr := range nbrs {
				heap.Push(&myQueue, &qType{dist: ndist, state: nbr, loopStart: me.loopStart})
			}
		}
	}
	return -1
}
