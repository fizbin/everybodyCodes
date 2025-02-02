package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
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

	// Something special about my input is that the grid is divided into three
	// vertical parts, that are connected only in one spot at the bottom that's
	// occupied by a "K" plant. Furthermore, plants that appear in one of the
	// three regions only appear in that region. Therefore, we can divide the
	// task of traveling salesman into thirds solving it for the left piece and
	// right piece as one chunk. The key will be that later we know to
	// immediately bail on any potential solution that gets some of the "left"
	// (bzw. "right") third but not all of it, and then gets something not
	// from the "left" ("right") third.

	// Also, the "K" plants have no other plants in the same vertical column
	// of the map.

	kPower := plantSpec['K']
	kPlantCols := []int{}
	kPlantIdx := []int{}
	plantPowerLeft := uint64(0)
	plantPowerRight := uint64(0)
	for idx, plantLoc := range plantLocs {
		if plantLocPower[idx] == kPower {
			kPlantCols = append(kPlantCols, plantLoc.y)
			kPlantIdx = append(kPlantIdx, idx)
		}
	}
	if len(kPlantCols) != 2 {
		log.Fatal("Expected exactly two K plants")
	}

	// Now we classify plants as "left", "right" or neither based on their
	// column relative to the two "K" plant columns.
	for idx, plantLoc := range plantLocs {
		if plantLoc.y < kPlantCols[0] {
			// this plant is a type you see in the left
			plantPowerLeft |= plantLocPower[idx]
		}
		if plantLoc.y > kPlantCols[1] {
			plantPowerRight |= plantLocPower[idx]
		}
	}

	plMap := make(map[point]int)
	for idx, plant := range plantLocs {
		plMap[plant] = idx
	}

	// map (point, point) -> dist; this gives the distance between two plants
	plantGraph := make([][]int, len(plantLocs))
	for idx := range plantGraph {
		plantGraph[idx] = make([]int, len(plantLocs))
	}

	// Now I compute the distance between each pair of plants.
	// But here I also take advantage of the "divide in thirds" thing
	// and don't bother to map out distances through "K" plants, since
	// I know that if there are two plants and the only path is through
	// a "K" then I can just get the distance from "A" to "B" by adding
	// "A to K" and "K to B"
	for startPIdx, startPlant := range plantLocs {
		plantName := gridMap[startPlant.x][startPlant.y]
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
					if plantName == 'K' || gridMap[nbr.x][nbr.y] != 'K' {
						q = append(q, nbr)
					}
				}
			}
		}
		for idx, plant := range plantLocs {
			plantGraph[startPIdx][idx] = beenThere[plant]
		}
	}
	// Now add in the distances around the "K" spots
	for rowIdx, row := range plantGraph {
		for colIdx, dist := range row {
			if rowIdx == colIdx {
				continue
			} else if dist == 0 {
				row[colIdx] = min(
					plantGraph[kPlantIdx[0]][rowIdx]+plantGraph[kPlantIdx[0]][colIdx],
					plantGraph[kPlantIdx[1]][rowIdx]+plantGraph[kPlantIdx[1]][colIdx])
			}
		}
	}
	// fmt.Println("DBG: Got all dists")

	plantGoal := (uint64(1) << len(plantSpec)) - 1

	// type bestDistCacheType struct {
	// 	goal   uint64
	// 	ending int
	// }
	bestDistCache := make([][]int, 0, len(plantGraph))
	for range plantGraph {
		bestDistCache = append(bestDistCache, make([]int, plantGoal+1))
	}

	const disconnected = 999999

	// bestDist is a function that determines the smallest possible total distance if you
	// start at the start spot, collect together the plant power bitset in "goal" and
	// then end at the spot "ending"
	var bestDist func(uint64, int) int
	bestDist = func(goal uint64, ending int) int {
		if plantLocPower[ending]&goal == 0 {
			return disconnected
		}
		if ((goal & plantPowerLeft) != 0) && ((goal & plantPowerLeft) != plantPowerLeft) {
			// so some, but not all, of the left third.
			if plantLocPower[ending]&plantPowerLeft == 0 {
				// Collecting some, but not all, of the left third and then ending
				// *outside* the left third? Not allowed.
				return disconnected
			}
		}
		if ((goal & plantPowerRight) != 0) && ((goal & plantPowerRight) != plantPowerRight) {
			// same but for right
			if plantLocPower[ending]&plantPowerRight == 0 {
				return disconnected
			}
		}
		val := bestDistCache[ending][goal]
		if val != 0 {
			return val
		}
		preGoal := goal ^ plantLocPower[ending]
		best := disconnected
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
		bestDistCache[ending][goal] = best
		return best
	}

	best := disconnected
	for plantIdx := 1; plantIdx < len(plantGraph); plantIdx++ {
		best = min(best, bestDist(plantGoal, plantIdx)+plantGraph[0][plantIdx])
	}
	return best
}
