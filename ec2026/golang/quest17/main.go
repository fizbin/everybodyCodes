package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q17_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q17_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q17_p3.txt`, "the input for part 3")

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

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	width := 0
	for scanner.Scan() {
		grid = append(grid, bytes.Clone(scanner.Bytes()))
		width = max(width, len(scanner.Bytes()))
	}
	height := len(grid)
	neighbors := func(pt point) []point {
		nb1 := point{pt.x - 1, pt.y}
		nb2 := point{pt.x + 1, pt.y}
		nb3 := point{pt.x, pt.y - 1}
		nb4 := point{pt.x, pt.y + 1}
		retval := make([]point, 0)
		for _, nb := range []point{nb1, nb2, nb3, nb4} {
			if nb.x >= 0 && nb.x < height && nb.y >= 0 && nb.y < width {
				if grid[nb.x][nb.y] == '#' {
					retval = append(retval, nb)
				}
			}
		}
		return retval
	}
	dist := make(map[point]int)
	dist[point{0, 0}] = 0
	queue := []point{{0, 0}}
	for len(queue) > 0 {
		me := queue[0]
		queue = queue[1:]
		mydist := dist[me]
		if me == (point{height - 1, width - 1}) {
			return mydist
		}
		nbs := neighbors(me)
		for _, nb := range nbs {
			if other, ok := dist[nb]; ok && other <= mydist+1 {
				continue
			}
			dist[nb] = mydist + 1
			queue = append(queue, nb)
		}
	}
	return -1
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	width := 0
	for scanner.Scan() {
		grid = append(grid, bytes.Clone(scanner.Bytes()))
		width = max(width, len(scanner.Bytes()))
	}
	// height := len(grid)
	crustmap := make(map[point]bool)
	for xidx, row := range grid {
		for yidx, ch := range row {
			if ch == '#' {
				crustmap[point{xidx, yidx}] = true
			}
		}
	}
	type elfLoc struct {
		crustCenter, offset point
	}
	neighbors := func(w elfLoc) []elfLoc {
		potential := []elfLoc{{addpt(w.crustCenter, leftpt(w.offset)), w.offset},
			{w.crustCenter, leftpt(w.offset)}}
		if crustmap[addpt(w.crustCenter, leftpt(w.offset))] {
			potential = append(potential, elfLoc{addpt(addpt(w.crustCenter, w.offset), leftpt(w.offset)), rightpt(w.offset)})
		}
		retval := make([]elfLoc, 0)
		for _, pot := range potential {
			if crustmap[pot.crustCenter] {
				if !crustmap[addpt(pot.crustCenter, pot.offset)] {
					retval = append(retval, pot)
				}
			}
		}
		return retval
	}
	where := elfLoc{point{0, 0}, point{0, -1}}
	dist := 1
	for where != (elfLoc{point{0, 0}, point{-1, 0}}) {
		dist++
		nbs := neighbors(where)
		if len(nbs) != 1 {
			fmt.Println(where, nbs)
			panic("Bad neighbors!")
		}
		where = nbs[0]
	}
	return dist
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	width := 0
	for scanner.Scan() {
		grid = append(grid, bytes.Clone(scanner.Bytes()))
		width = max(width, len(scanner.Bytes()))
	}
	height := len(grid)
	crustmap := make(map[point]bool)
	for xidx, row := range grid {
		for yidx, ch := range row {
			if ch == '#' {
				crustmap[point{xidx, yidx}] = true
			}
		}
	}
	type elfLoc struct {
		crustCenter, offset point
	}
	neighbors := func(w elfLoc) []elfLoc {
		// remember, keep crust to the elf's left
		potential := []elfLoc{
			{addpt(w.crustCenter, leftpt(w.offset)), w.offset}, // straight ahead
			{w.crustCenter, leftpt(w.offset)},                  // turn left
		}
		if crustmap[addpt(w.crustCenter, leftpt(w.offset))] {
			// turn right
			// make sure you don't try to turn right at a spot where two crust pieces meet diagonally
			potential = append(potential, elfLoc{addpt(addpt(w.crustCenter, w.offset), leftpt(w.offset)), rightpt(w.offset)})
		}
		retval := make([]elfLoc, 0)
		for _, pot := range potential {
			if crustmap[pot.crustCenter] {
				if !crustmap[addpt(pot.crustCenter, pot.offset)] {
					retval = append(retval, pot)
				}
			}
		}
		return retval
	}
	where := elfLoc{point{0, 0}, point{0, -1}}
	dist1 := 1
	for where != (elfLoc{point{height - 1, width - 1}, point{1, 0}}) {
		dist1++
		nbs := neighbors(where)
		if len(nbs) != 1 {
			fmt.Println(where, nbs)
			panic("Bad neighbors!")
		}
		where = nbs[0]
	}
	where = elfLoc{point{height - 1, width - 1}, point{0, 1}}
	dist2 := 1
	for where != (elfLoc{point{0, 0}, point{-1, 0}}) {
		dist2++
		nbs := neighbors(where)
		if len(nbs) != 1 {
			fmt.Println(where, nbs)
			panic("Bad neighbors!")
		}
		where = nbs[0]
	}
	return min(dist1, dist2)
}
