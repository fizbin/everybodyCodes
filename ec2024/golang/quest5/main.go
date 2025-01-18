package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func grabInts(line string) []int {
	fields := strings.Fields(line)
	retval := make([]int, 0, len(line))
	for _, field := range fields {
		n, _ := strconv.Atoi(field)
		retval = append(retval, n)
	}
	return retval
}

func danceround(dancecols [][]int, roundNum int, digmul int) int {
	// do one round, tell me the number shouted at the end of the round
	clapper := dancecols[roundNum%len(dancecols)][0]
	dancecols[roundNum%len(dancecols)] = dancecols[roundNum%len(dancecols)][1:]
	targetCol := dancecols[(roundNum+1)%len(dancecols)]
	insertSpot := -1
	realClap := clapper % (2 * len(targetCol))
	if realClap == 0 {
		realClap = 2 * len(targetCol)
	}
	if realClap <= len(targetCol) {
		insertSpot = realClap - 1
	} else {
		insertSpot = 1 + (2*len(targetCol) - realClap)
	}
	targetCol = slices.Concat(targetCol[0:insertSpot], []int{clapper}, targetCol[insertSpot:])
	dancecols[(roundNum+1)%len(dancecols)] = targetCol
	finalnum := 0
	for _, col := range dancecols {
		finalnum = finalnum*digmul + col[0]
	}
	return finalnum
}

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q05_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		var dancecols [][]int
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			linenums := grabInts(scanner.Text())
			if dancecols == nil {
				dancecols = make([][]int, len(linenums))
			}
			for idx, n := range linenums {
				dancecols[idx] = append(dancecols[idx], n)
			}
		}

		var shouted int
		for round := range 10 {
			shouted = danceround(dancecols, round, 10)
		}

		fmt.Println("Part 1:", shouted)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q05_p2.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		var dancecols [][]int
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			linenums := grabInts(scanner.Text())
			if dancecols == nil {
				dancecols = make([][]int, len(linenums))
			}
			for idx, n := range linenums {
				dancecols[idx] = append(dancecols[idx], n)
			}
		}

		shoutrecord := make(map[int]int)
		for round := 0; true; round++ {
			shouted := danceround(dancecols, round, 100)
			if shoutrecord[shouted] == 2023 {
				fmt.Println("Part 2:", shouted*(round+1))
				break
			}
			shoutrecord[shouted] += 1
		}
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q05_p3.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		var dancecols [][]int
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			linenums := grabInts(scanner.Text())
			if dancecols == nil {
				dancecols = make([][]int, len(linenums))
			}
			for idx, n := range linenums {
				dancecols[idx] = append(dancecols[idx], n)
			}
		}

		shoutmax := 0
		dancecolsSlow := make([][]int, len(dancecols))
		for idx, col := range dancecols {
			dancecolsSlow[idx] = slices.Clone(col)
		}
		for round := 0; true; round += 2 {
			shouted := danceround(dancecols, round, 10000)
			if shouted > shoutmax {
				shoutmax = shouted
			}
			shouted = danceround(dancecols, round+1, 10000)
			if shouted > shoutmax {
				shoutmax = shouted
			}
			danceround(dancecolsSlow, round/2, 10000)
			// now break if dancecols and dancecolsSlow are equal
			if slices.EqualFunc(dancecols, dancecolsSlow, slices.Equal) {
				break
			}
		}
		fmt.Println("Part 3:", shoutmax)
	}
}
