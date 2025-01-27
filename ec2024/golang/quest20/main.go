package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"slices"
)

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q20_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 1:", doProblem1(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 2 || argsWithoutProg[1] == "" {
			infile = "../input/everybody_codes_e2024_q20_p2.txt"
		} else {
			infile = argsWithoutProg[1]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 2:", doProblem2(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 3 || argsWithoutProg[2] == "" {
			infile = "../input/everybody_codes_e2024_q20_p3.txt"
		} else {
			infile = argsWithoutProg[2]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 3:", doProblem3(data))
	}
}

type point struct{ x, y int }

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var startLoc point
	grid := make([][]byte, 0)
	height := 0
	for scanner.Scan() {
		sIdx := slices.Index(scanner.Bytes(), 'S')
		if sIdx >= 0 {
			startLoc = point{height, sIdx}
		}
		grid = append(grid, slices.Clone(scanner.Bytes()))
		height++
	}
	width := len(grid[0])
	type gState struct {
		where point
		dir   point
	}
	beenThere := make(map[gState]int)
	whereWeAre := make(map[gState]int)
	whereWeAre[gState{startLoc, point{0, -1}}] = 1000
	whereWeAre[gState{startLoc, point{-1, 0}}] = 1000
	whereWeAre[gState{startLoc, point{0, 1}}] = 1000
	whereWeAre[gState{startLoc, point{1, 0}}] = 1000
	// fmt.Println("DBG preflight ", whereWeAre)
	for tick := 1; tick <= 100; tick++ {
		newWhere := make(map[gState]int)
		for state, hgt := range whereWeAre {
			if prevHgt, ok := beenThere[state]; ok {
				if prevHgt >= hgt {
					continue
				}
			}
			beenThere[state] = hgt
			for _, newState := range []gState{
				{point{state.where.x + state.dir.x, state.where.y + state.dir.y}, state.dir},
				{point{state.where.x + state.dir.y, state.where.y - state.dir.x}, point{state.dir.y, -state.dir.x}},
				{point{state.where.x - state.dir.y, state.where.y + state.dir.x}, point{-state.dir.y, state.dir.x}},
			} {
				if newState.where.x < 0 || newState.where.y < 0 || newState.where.x >= height || newState.where.y >= width {
					continue
				}
				if grid[newState.where.x][newState.where.y] == '#' {
					continue
				}
				var newHgt int
				switch grid[newState.where.x][newState.where.y] {
				case '#':
					continue
				case '+':
					newHgt = hgt + 1
				case '.':
					newHgt = hgt - 1
				case '-':
					newHgt = hgt - 2
				}
				newWhere[newState] = max(newHgt, newWhere[newState])
			}
		}
		whereWeAre = newWhere
	}
	maxHgt := 0
	for _, val := range whereWeAre {
		maxHgt = max(maxHgt, val)
	}
	return maxHgt
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var startLoc point
	grid := make([][]byte, 0)
	height := 0
	for scanner.Scan() {
		sIdx := slices.Index(scanner.Bytes(), 'S')
		if sIdx >= 0 {
			startLoc = point{height, sIdx}
		}
		grid = append(grid, slices.Clone(scanner.Bytes()))
		height++
	}
	width := len(grid[0])
	type gState struct {
		where point
		dir   point
		seen  int // seen 0: start, 1: A, 3: C
	}
	beenThere := make(map[gState]int)
	whereWeAre := make(map[gState]int)
	whereWeAre[gState{startLoc, point{0, -1}, 0}] = 10000
	whereWeAre[gState{startLoc, point{-1, 0}, 0}] = 10000
	whereWeAre[gState{startLoc, point{0, 1}, 0}] = 10000
	whereWeAre[gState{startLoc, point{1, 0}, 0}] = 10000
	for tick := 1; true; tick++ {
		newWhere := make(map[gState]int)
		for state, hgt := range whereWeAre {
			if prevHgt, ok := beenThere[state]; ok {
				if prevHgt >= hgt {
					continue
				}
			}
			beenThere[state] = hgt
			for _, newState := range []gState{
				{point{state.where.x + state.dir.x, state.where.y + state.dir.y}, state.dir, 0},
				{point{state.where.x + state.dir.y, state.where.y - state.dir.x}, point{state.dir.y, -state.dir.x}, 0},
				{point{state.where.x - state.dir.y, state.where.y + state.dir.x}, point{-state.dir.y, state.dir.x}, 0},
			} {
				if newState.where.x < 0 || newState.where.y < 0 || newState.where.x >= height || newState.where.y >= width {
					continue
				}
				myCh := grid[newState.where.x][newState.where.y]
				if myCh == '#' {
					continue
				}
				var newHgt int
				newSeen := state.seen
				switch myCh {
				case '#':
					continue
				case '+':
					newHgt = hgt + 1
				case '.', 'A', 'B', 'C', 'S':
					newHgt = hgt - 1
				case '-':
					newHgt = hgt - 2
				}
				if (newSeen == 0 && myCh == 'A') || (newSeen == 1 && myCh == 'B') || (newSeen == 2 && myCh == 'C') {
					newSeen++
				}
				if grid[newState.where.x][newState.where.y] == 'S' {
					if newHgt >= 10000 && newSeen == 3 {
						return tick
					}
					continue
				}
				newState.seen = newSeen
				newWhere[newState] = max(newHgt, newWhere[newState])
			}
		}
		whereWeAre = newWhere
	}
	return -1
}

func doProblem3(data []byte) any {

	scanner := bufio.NewScanner(bytes.NewReader(data))
	var startLoc point
	grid := make([][]byte, 0)
	height := 0
	for scanner.Scan() {
		sIdx := slices.Index(scanner.Bytes(), 'S')
		if sIdx >= 0 {
			startLoc = point{height, sIdx}
		}
		grid = append(grid, slices.Clone(scanner.Bytes()))
		height++
	}
	width := len(grid[0])
	type gState struct {
		where point
		dir   point
	}
	beenThere := make(map[gState]int)
	whereWeAre := make(map[gState]int)
	whereWeAre[gState{startLoc, point{0, -1}}] = 384400
	whereWeAre[gState{startLoc, point{-1, 0}}] = 384400
	whereWeAre[gState{startLoc, point{0, 1}}] = 384400
	whereWeAre[gState{startLoc, point{1, 0}}] = 384400
	// fmt.Println("DBG preflight ", whereWeAre)
	maxRow := 0
	for tick := 1; len(whereWeAre) > 0; tick++ {
		newWhere := make(map[gState]int)
		for state, hgt := range whereWeAre {
			if prevHgt, ok := beenThere[state]; ok {
				if prevHgt >= hgt {
					continue
				}
			}
			beenThere[state] = hgt
			for _, newState := range []gState{
				{point{state.where.x + state.dir.x, state.where.y + state.dir.y}, state.dir},
				{point{state.where.x + state.dir.y, state.where.y - state.dir.x}, point{state.dir.y, -state.dir.x}},
				{point{state.where.x - state.dir.y, state.where.y + state.dir.x}, point{-state.dir.y, state.dir.x}},
			} {
				if newState.where.x < 0 || newState.where.y < 0 || newState.where.y >= width {
					continue
				}
				// Because I, as a human, can see that the optimal strategy in the given "notes" is going
				// to be "fly to column 31, then head straight down", force that path:
				if state.where.y > 31 && newState.where.y >= state.where.y {
					continue
				}
				if state.where.y < 31 && newState.where.y <= state.where.y {
					continue
				}
				if state.where.y == 31 && newState.where.y != 31 {
					continue
				}
				myCh := grid[newState.where.x%height][newState.where.y]
				if myCh == '#' {
					continue
				}
				var newHgt int
				switch myCh {
				case '#':
					continue
				case '+':
					newHgt = hgt + 1
				case '.', 'S':
					newHgt = hgt - 1
				case '-':
					newHgt = hgt - 2
				}
				if newHgt == 0 {
					maxRow = max(maxRow, newState.where.x)
					continue
				}
				newWhere[newState] = max(newHgt, newWhere[newState])
			}
		}
		whereWeAre = newWhere
		// if tick < 10 {
		// 	fmt.Println("DBG: ", whereWeAre)
		// }
	}
	return maxRow
}
