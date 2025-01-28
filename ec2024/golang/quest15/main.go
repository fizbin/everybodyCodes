package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"unicode"
)

type point struct{ x, y int }

func neighbors(current point, gridMap [][]byte, height, width int) []point {
	retval := make([]point, 0, 4)
	sx := current.x
	sy := current.y
	for _, loc := range []point{{sx + 1, sy}, {sx - 1, sy}, {sx, sy + 1}, {sx, sy - 1}} {
		if loc.x >= 0 && loc.x < height && loc.y >= 0 && loc.y < width {
			contents := gridMap[loc.x][loc.y]
			if contents != '#' && contents != '~' {
				retval = append(retval, loc)
			}
		}
	}
	return retval
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

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
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q15_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 1:", doProblem(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 2 || argsWithoutProg[1] == "" {
			infile = "../input/everybody_codes_e2024_q15_p2.txt"
		} else {
			infile = argsWithoutProg[1]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 2:", doProblem(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 3 || argsWithoutProg[2] == "" {
			infile = "../input/everybody_codes_e2024_q15_p3.txt"
		} else {
			infile = argsWithoutProg[2]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 3:", doProblem3(data))
	}
}

func doProblem(data []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var startLoc point
	plantSpec := make(map[byte]uint64)
	gridMap := make([][]byte, 0)
	var width, height int
	for rowIdx := 0; scanner.Scan(); rowIdx++ {
		row := make([]byte, 0, width)
		for colIdx, ch := range scanner.Bytes() {
			here := point{rowIdx, colIdx}
			if rowIdx == 0 && ch == '.' {
				startLoc = here
			}
			row = append(row, ch)
			if unicode.IsLetter(rune(ch)) && plantSpec[ch] == 0 {
				plantSpec[ch] = (1 << len(plantSpec))
			}
		}
		gridMap = append(gridMap, row)
		width = len(row)
	}
	height = len(gridMap)
	plantGoal := (uint64(1) << len(plantSpec)) - 1
	type beenThereType struct {
		where    point
		carrying uint64
	}
	beenThere := make(map[beenThereType]int)
	type queueType struct {
		dist     int
		where    point
		carrying uint64
	}
	q := []queueType{{dist: 0, where: startLoc, carrying: 0}}
	for len(q) > 0 {
		current := q[0]
		q = q[1:]
		btKey := beenThereType{current.where, current.carrying}
		if val, ok := beenThere[btKey]; ok && val <= current.dist {
			continue
		}
		beenThere[btKey] = current.dist
		if current.where == startLoc && current.carrying > 0 {
			continue
		}
		newCarrying := current.carrying | plantSpec[gridMap[current.where.x][current.where.y]]
		if val, ok := beenThere[beenThereType{current.where, plantGoal ^ newCarrying}]; ok {
			q = append(q, queueType{val + current.dist, startLoc, plantGoal})
			continue
		}
		for _, nbr := range neighbors(current.where, gridMap, height, width) {
			q = append(q, queueType{current.dist + 1, nbr, newCarrying})
		}
	}
	return beenThere[beenThereType{startLoc, plantGoal}]
}

func doProblem3(data []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	plantSpec := make(map[byte]uint64)
	gridMap := make([][]byte, 0)
	plantLocs := make([]point, 0)
	plantLocPower := make([]uint64, 0)
	var width, height int
	for rowIdx := 0; scanner.Scan(); rowIdx++ {
		row := make([]byte, 0, width)
		for colIdx, ch := range scanner.Bytes() {
			here := point{rowIdx, colIdx}
			if rowIdx == 0 && ch == '.' {
				plantLocs = append(plantLocs, here)
				plantLocPower = append(plantLocPower, 0)
			}
			row = append(row, ch)
			if unicode.IsLetter(rune(ch)) {
				if plantSpec[ch] == 0 {
					plantSpec[ch] = (1 << len(plantSpec))
				}
				plantLocs = append(plantLocs, here)
				plantLocPower = append(plantLocPower, plantSpec[ch])
			}
		}
		gridMap = append(gridMap, row)
		width = len(row)
	}
	height = len(gridMap)

	kPower := plantSpec['K']
	kPowerCols := []int{}
	plantPowerLeft := uint64(0)
	plantPowerRight := uint64(0)
	for idx, plantLoc := range plantLocs {
		if plantLocPower[idx] == kPower {
			kPowerCols = append(kPowerCols, plantLoc.y)
		}
	}
	if len(kPowerCols) != 2 {
		log.Fatal("Expected exactly two K plants")
	}
	for idx, plantLoc := range plantLocs {
		if plantLoc.y < kPowerCols[0] {
			plantPowerLeft |= plantLocPower[idx]
		}
		if plantLoc.y > kPowerCols[1] {
			plantPowerRight |= plantLocPower[idx]
		}
	}

	plMap := make(map[point]int)
	for idx, plant := range plantLocs {
		plMap[plant] = idx
	}

	// map (point, point) -> dist
	plantGraph := make([][]int, 0, len(plantLocs))

	for _, startPlant := range plantLocs {
		beenThere := make(map[point]int)
		beenThere[startPlant] = 0
		q := []point{startPlant}
		for len(q) > 0 {
			current := q[0]
			q = q[1:]
			mydist := beenThere[current]
			for _, nbr := range neighbors(current, gridMap, height, width) {
				if _, ok := beenThere[nbr]; !ok {
					beenThere[nbr] = 1 + mydist
					q = append(q, nbr)
				}
			}
		}
		outMap := make([]int, len(plantLocs))
		for idx, plant := range plantLocs {
			outMap[idx] = beenThere[plant]
		}
		plantGraph = append(plantGraph, outMap)
	}
	// fmt.Println("DBG: Got all dists")

	plantGoal := (uint64(1) << len(plantSpec)) - 1

	type bestDistCacheType struct {
		goal   uint64
		ending int
	}
	bestDistCache := make(map[bestDistCacheType]int)
	var bestDist func(uint64, int) int
	bestDist = func(goal uint64, ending int) int {
		if plantLocPower[ending]&goal == 0 {
			// would return just MaxInt, but I don't want overflow on addition later
			return math.MaxInt - 3*(width+height)
		}
		if ((goal & plantPowerLeft) != 0) && ((goal & plantPowerLeft) != plantPowerLeft) {
			if plantLocPower[ending]&plantPowerLeft == 0 {
				return math.MaxInt - 3*(width+height)
			}
		}
		if ((goal & plantPowerRight) != 0) && ((goal & plantPowerRight) != plantPowerRight) {
			if plantLocPower[ending]&plantPowerRight == 0 {
				return math.MaxInt - 3*(width+height)
			}
		}
		if val, ok := bestDistCache[bestDistCacheType{goal, ending}]; ok {
			return val
		}
		preGoal := goal ^ plantLocPower[ending]
		best := math.MaxInt
		if preGoal == 0 {
			best = plantGraph[ending][0]
		} else {
			for butEnd := 1; butEnd < len(plantGraph); butEnd++ {
				if butEnd == ending || plantLocPower[butEnd]&preGoal == 0 {
					continue
				}
				newDist := plantGraph[ending][butEnd] + bestDist(preGoal, butEnd)
				best = min(newDist, best)
			}
		}
		bestDistCache[bestDistCacheType{goal, ending}] = best
		return best
	}
	for plantIdx := 1; plantIdx < len(plantGraph); plantIdx++ {
		for goal := uint64(1); goal <= plantGoal; goal = 2*goal + 1 {
			bestDist(goal, plantIdx)
		}
	}
	best := math.MaxInt
	for plantIdx := 1; plantIdx < len(plantGraph); plantIdx++ {
		best = min(best, bestDist(plantGoal, plantIdx)+plantGraph[0][plantIdx])
	}
	return best
}
