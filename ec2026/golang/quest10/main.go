package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"slices"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q10_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q10_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q10_p3.txt`, "the input for part 3")

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
	r, c int
}

type ghost struct {
	where point
	dir   point
	life  int
}

var ghostDirs map[byte]point

func init() {
	ghostDirs = map[byte]point{
		'<': {0, -1},
		'^': {-1, 0},
		'>': {0, 1},
		'v': {1, 0},
	}
}

func doProblem1(data []byte) any {
	// 	data = []byte(`###########
	// #####.#####
	// ###.....###
	// ##.......##
	// ##.>>.<<.##
	// #v...>...v#
	// ##...v...##
	// ##.......##
	// ###.....###
	// #####.#####
	// ###########`)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	for scanner.Scan() {
		grid = append(grid, bytes.Clone(scanner.Bytes()))
	}
	ghosts := make([]ghost, 0)
	gridMap := make(map[point]byte)
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[r]); c++ {
			myCh := grid[r][c]
			here := point{r, c}
			gridMap[here] = myCh
			gDir, ok := ghostDirs[myCh]
			if ok {
				ghosts = append(ghosts, ghost{where: here, dir: gDir})
				gridMap[here] = '.'
			}
		}
	}
	collisions := 0
	for roundNum := 1; roundNum <= 10; roundNum++ {
		// ocollisions := collisions
		newGhosts := make([]ghost, 0)
		ghostPosMap := make(map[point]int)
		for _, ghost_ := range ghosts {
			nPos := point{ghost_.where.r + ghost_.dir.r, ghost_.where.c + ghost_.dir.c}
			var nGhost ghost
			switch gridMap[nPos] {
			case 0:
				continue
			case '#':
				nGhost = ghost{where: ghost_.where, dir: point{ghost_.dir.c, -ghost_.dir.r}}
				collisions++
			case '.':
				nGhost = ghost{where: nPos, dir: ghost_.dir}
			}
			ghostPosMap[nGhost.where]++
			newGhosts = append(newGhosts, nGhost)
		}
		for idx, nGhost := range newGhosts {
			if ghostPosMap[nGhost.where] > 1 {
				collisions++
				newGhosts[idx].dir.c *= -1
				newGhosts[idx].dir.r *= -1
			}
		}
		// fmt.Println("round ", roundNum, "collisions", collisions-ocollisions)
		ghosts = newGhosts
	}
	return collisions
}

func doProblem2(data []byte) any {
	// 	data = []byte(`###########
	// #####.#####
	// ###.....###
	// ##.......##
	// ##.>>.<<.##
	// #v...>...v#
	// ##...v...##
	// ##.......##
	// ###.....###
	// #####.#####
	// ###########`)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	for scanner.Scan() {
		grid = append(grid, bytes.Clone(scanner.Bytes()))
	}
	ghosts := make([]ghost, 0)
	gridMap := make(map[point]byte)
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[r]); c++ {
			myCh := grid[r][c]
			here := point{r, c}
			gridMap[here] = myCh
			gDir, ok := ghostDirs[myCh]
			if ok {
				ghosts = append(ghosts, ghost{where: here, dir: gDir, life: 10000})
				gridMap[here] = '.'
			}
		}
	}
	return runGhosts(ghosts, gridMap)
}

func runGhosts(ghosts []ghost, gridMap map[point]byte) any {
	// Basic idea: run the ghosts, tracking collisions and removing ghosts
	// as they fade to 0 life. When we see a situation with all the ghost
	// position and direction exactly as it was before, "fast forward" by
	// taking big steps, as many steps as we can without actually having
	// a ghost fade to 0. (Save all the "fading to 0" stuff for slow,
	// careful steps)
	roundNum := 0
	lastGhostLostRound := 0
	recordedGW := make([]point, 0, len(ghosts))
	recordedGL := make([]int, 0)
	recordedRoundNum := 0
	for ; len(ghosts) > 0; roundNum++ {
		// oGhostCount := len(ghosts)

		// Start checking for whether we've seen this ghost configuration before
		myGW := make([]point, 0, len(ghosts))
		for _, ghost_ := range ghosts {
			myGW = append(myGW, ghost_.where, ghost_.dir)
		}
		if slices.Equal(myGW, recordedGW) {
			// fmt.Println("Got a hit at ", roundNum)

			// We have seen it before! Determine how long it's been and
			// use that as a step size
			stepSize := roundNum - recordedRoundNum
			stepsToTake := roundNum
			for idx, ghost_ := range ghosts {
				lifeLost := recordedGL[idx] - ghost_.life
				stepsToTake = min(stepsToTake, (ghost_.life-1)/lifeLost)
			}
			// fmt.Println("before:", ghosts)
			// minGhostLife := ghosts[0].life
			for idx, ghost_ := range ghosts {
				lifeLost := recordedGL[idx] - ghost_.life
				ghost_.life -= lifeLost * stepsToTake
				// minGhostLife = min(minGhostLife, ghost_.life)
				ghosts[idx] = ghost_
			}
			// fmt.Print("Taking ", maxSteps, " steps of size ", stepSize)
			// fmt.Println(" min ghost life", minGhostLife, "to land at", roundNum+stepSize*maxSteps)
			// fmt.Println("after:", ghosts)
			roundNum += stepSize * stepsToTake
			lastGhostLostRound = roundNum
			recordedGW = make([]point, 0)
		}
		if (len(recordedGW) != len(myGW)) && roundNum > 5000+lastGhostLostRound {
			// fmt.Println("Recording at", roundNum, "with", len(ghosts), "ghosts")

			// The "5000" is an arbitrary warm-up constant because it might take a bit after a
			// ghost disappears from the scene for the ghost interaction to settle into a
			// pattern.
			recordedGW = myGW
			recordedGL = make([]int, 0)
			for _, ghost_ := range ghosts {
				recordedGL = append(recordedGL, ghost_.life)
			}
			recordedRoundNum = roundNum
		}

		newGhosts := make([]ghost, 0)
		ghostPosMap := make(map[point]int)
		for _, ghost_ := range ghosts {
			nPos := point{ghost_.where.r + ghost_.dir.r, ghost_.where.c + ghost_.dir.c}
			var nGhost ghost
			switch gridMap[nPos] {
			case 0:
				continue
			case '#':
				nGhost = ghost{where: ghost_.where, dir: point{ghost_.dir.c, -ghost_.dir.r}, life: ghost_.life - 1}
				if nGhost.life <= 0 {
					lastGhostLostRound = roundNum
					continue
				}
			case '.':
				nGhost = ghost{where: nPos, dir: ghost_.dir, life: ghost_.life}
			}
			ghostPosMap[nGhost.where]++
			newGhosts = append(newGhosts, nGhost)
		}
		newGhosts2 := make([]ghost, 0, len(newGhosts))
		for _, nGhost := range newGhosts {
			if ghostPosMap[nGhost.where] > 1 {
				nGhost.dir = point{r: -nGhost.dir.r, c: -nGhost.dir.c}
				nGhost.life--
			}
			if nGhost.life > 0 {
				newGhosts2 = append(newGhosts2, nGhost)
			} else {
				lastGhostLostRound = roundNum
			}
		}

		ghosts = newGhosts2
		// if len(ghosts) != oGhostCount {
		// 	fmt.Println("Ghost death at round", roundNum, "went from", oGhostCount, "to", len(ghosts))
		// }
	}
	return roundNum
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	for scanner.Scan() {
		grid = append(grid, bytes.Clone(scanner.Bytes()))
	}
	ghosts := make([]ghost, 0)
	gridMap := make(map[point]byte)
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[r]); c++ {
			myCh := grid[r][c]
			here := point{r, c}
			gridMap[here] = myCh
			gDir, ok := ghostDirs[myCh]
			if ok {
				ghosts = append(ghosts, ghost{where: here, dir: gDir, life: 1000000000000})
				gridMap[here] = '.'
			}
		}
	}
	return runGhosts(ghosts, gridMap)
}
