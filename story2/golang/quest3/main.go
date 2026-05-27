package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Die struct {
	pulse      int
	seed       int
	faces      []int
	lastIdx    int
	rollNumber int
}
type coord struct {
	x int
	y int
}

func Roll(die *Die) int {
	pulse := die.pulse
	rollNumber := die.rollNumber
	spin := pulse * rollNumber
	resultIdx := (die.lastIdx + spin) % len(die.faces)
	result := die.faces[resultIdx]
	pulse = pulse + spin
	pulse = pulse % die.seed
	pulse = pulse + 1 + rollNumber + die.seed
	// now store everything back
	die.pulse = pulse
	die.lastIdx = resultIdx
	die.rollNumber += 1
	return result
}

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e2_q03_p1.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()

		// testDie := Die{pulse: 3, seed: 3, faces: []int{1, 2, 4, -1, 5, 7, 9}, lastIdx: 0, rollNumber: 1}
		// for range 100 {
		// 	rollNum := testDie.rollNumber
		// 	result := Roll(&testDie)
		// 	pulseAfter := testDie.pulse
		// 	fmt.Println("Roll", rollNum, "result", result, "pulseAfter", pulseAfter)
		// }

		// return

		dice := make(map[int]*Die)
		scanner := bufio.NewScanner(file)
		parser := regexp.MustCompile(`(\d+): *faces=\[([-\d,]+)\] *seed=(\d+)`)
		for scanner.Scan() {
			parseResult := parser.FindStringSubmatch(scanner.Text())
			if len(parseResult) == 0 {
				panic(scanner.Text())
			}
			diceId, err := strconv.ParseInt(parseResult[1], 10, 32)
			if err != nil {
				panic("Bad Input " + parseResult[1] + " in " + parseResult[0])
			}
			seed, err := strconv.ParseInt(parseResult[3], 10, 32)
			if err != nil {
				panic("Bad Input " + parseResult[3] + " in " + parseResult[0])
			}
			facesStr := strings.Split(parseResult[2], ",")
			faces := make([]int, 0, len(facesStr))
			for _, faceStr := range facesStr {
				face, err := strconv.ParseInt(faceStr, 10, 32)
				if err != nil {
					panic("Bad Input " + faceStr + " in " + parseResult[0])
				}
				faces = append(faces, (int)(face))
			}
			dice[(int)(diceId)] = &Die{pulse: (int)(seed), seed: (int)(seed), faces: faces, lastIdx: 0, rollNumber: 1}
		}
		if scanner.Err() != nil {
			panic(scanner.Err().Error())
		}
		rollsTaken := 0
		pointsTotal := 0
		for pointsTotal < 10000 {
			for dId, pdie := range dice {
				roll := Roll(pdie)
				if false {
					fmt.Println("roll", dId, "->", roll)
				}
				pointsTotal += roll
			}
			rollsTaken++
			// fmt.Println("p1 dbg after", rollsTaken, "have", pointsTotal)
		}
		// fmt.Println("p1 dbg pointstotal", pointsTotal)
		fmt.Println("p1", rollsTaken)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e2_q03_p2.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()

		dice := make(map[int]*Die)
		scanner := bufio.NewScanner(file)
		parser := regexp.MustCompile(`(\d+): *faces=\[([-\d,]+)\] *seed=(\d+)`)
		for scanner.Scan() {
			parseResult := parser.FindStringSubmatch(scanner.Text())
			if len(parseResult) == 0 {
				if len(scanner.Text()) > 0 {
					panic(scanner.Text())
				}
				break
			}
			diceId, err := strconv.ParseInt(parseResult[1], 10, 32)
			if err != nil {
				panic("Bad Input " + parseResult[1] + " in " + parseResult[0])
			}
			seed, err := strconv.ParseInt(parseResult[3], 10, 32)
			if err != nil {
				panic("Bad Input " + parseResult[3] + " in " + parseResult[0])
			}
			facesStr := strings.Split(parseResult[2], ",")
			faces := make([]int, 0, len(facesStr))
			for _, faceStr := range facesStr {
				face, err := strconv.ParseInt(faceStr, 10, 32)
				if err != nil {
					panic("Bad Input " + faceStr + " in " + parseResult[0])
				}
				faces = append(faces, (int)(face))
			}
			dice[(int)(diceId)] = &Die{pulse: (int)(seed), seed: (int)(seed), faces: faces, lastIdx: 0, rollNumber: 1}
		}
		trackStr := "1"
		if scanner.Scan() {
			trackStr = scanner.Text()
		}
		if scanner.Err() != nil {
			panic(scanner.Err().Error())
		}
		finishers := make([]int, 0, len(dice))
		players := make(map[int]int)
		for dId, _ := range dice {
			players[dId] = 0
		}
		for len(players) > 0 {
			for dId, current := range players {
				diep := dice[dId]
				result := Roll(diep)
				if trackStr[current] == '0'+(byte)(result) {
					current += 1
					if current >= len(trackStr) {
						delete(players, dId)
						finishers = append(finishers, dId)
					} else {
						players[dId] += 1
					}
				}
			}
		}
		fmt.Print("p2 ", finishers[0])
		for _, player := range finishers[1:] {
			fmt.Printf(",%d", player)
		}
		fmt.Println()
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e2_q03_p3.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()

		dice := make(map[int]*Die)
		scanner := bufio.NewScanner(file)
		parser := regexp.MustCompile(`(\d+): *faces=\[([-\d,]+)\] *seed=(\d+)`)
		for scanner.Scan() {
			parseResult := parser.FindStringSubmatch(scanner.Text())
			if len(parseResult) == 0 {
				if len(scanner.Text()) > 0 {
					panic(scanner.Text())
				}
				break
			}
			diceId, err := strconv.ParseInt(parseResult[1], 10, 32)
			if err != nil {
				panic("Bad Input " + parseResult[1] + " in " + parseResult[0])
			}
			seed, err := strconv.ParseInt(parseResult[3], 10, 32)
			if err != nil {
				panic("Bad Input " + parseResult[3] + " in " + parseResult[0])
			}
			facesStr := strings.Split(parseResult[2], ",")
			faces := make([]int, 0, len(facesStr))
			for _, faceStr := range facesStr {
				face, err := strconv.ParseInt(faceStr, 10, 32)
				if err != nil {
					panic("Bad Input " + faceStr + " in " + parseResult[0])
				}
				faces = append(faces, (int)(face))
			}
			dice[(int)(diceId)] = &Die{pulse: (int)(seed), seed: (int)(seed), faces: faces, lastIdx: 0, rollNumber: 1}
		}
		grid := make([]string, 0)
		for scanner.Scan() {
			grid = append(grid, scanner.Text())
		}
		covered := make(map[coord]bool)
		allmap := make(map[coord]bool)
		for x := range len(grid) {
			for y := range len(grid[x]) {
				allmap[coord{x, y}] = true
			}
		}
		for _, diep := range dice {
			candidates := allmap
			for len(candidates) > 0 {
				newCandidates := make(map[coord]bool)
				roll := Roll(diep)
				for c := range candidates {
					if grid[c.x][c.y] == '0'+(byte)(roll) {
						covered[c] = true
						newCandidates[c] = true
						c1 := coord{c.x - 1, c.y}
						c2 := coord{c.x + 1, c.y}
						c3 := coord{c.x, c.y - 1}
						c4 := coord{c.x, c.y + 1}
						if allmap[c1] {
							newCandidates[c1] = true
						}
						if allmap[c2] {
							newCandidates[c2] = true
						}
						if allmap[c3] {
							newCandidates[c3] = true
						}
						if allmap[c4] {
							newCandidates[c4] = true
						}
					}
				}
				candidates = newCandidates
			}
		}
		fmt.Println("p3", len(covered))
		// Uncomment for image
		// for x := range len(grid) {
		// 	for y := range len(grid[x]) {
		// 		if covered[coord{x, y}] {
		// 			fmt.Print("*")
		// 		} else {
		// 			fmt.Print(" ")
		// 		}
		// 	}
		// 	fmt.Println()
		// }
	}
}
