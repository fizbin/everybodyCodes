package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x, y int
}

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q12_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		// 		data = []byte(`.............
		// .C...........
		// .B......T....
		// .A......T.T..
		// =============`)
		towerSpots := make(map[rune]point)
		targetSpots := make([]point, 0)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		rowIdx := 0
		for scanner.Scan() {
			for colIdx, ch := range scanner.Text() {
				if ch == 'A' || ch == 'B' || ch == 'C' {
					towerSpots[ch] = point{rowIdx, colIdx}
				}
				if ch == 'T' {
					targetSpots = append(targetSpots, point{rowIdx, colIdx})
				}
			}
			rowIdx++
		}
		total := 0
		for _, target := range targetSpots {
			for idx, tower := range " ABC" {
				if idx == 0 {
					continue
				}
				q := (target.y - towerSpots[tower].y) - (target.x - towerSpots[tower].x)
				if q%3 == 0 {
					// fmt.Println("DBG", "towerloc", towerSpots[tower], "toweri ", ti, " tgti ", idx, " power ", q/3)
					total += (q / 3) * idx
				}
			}
		}
		fmt.Println("Part 1:", total)
	}
	{
		var infile string
		if len(argsWithoutProg) < 2 || argsWithoutProg[1] == "" {
			infile = "../input/everybody_codes_e2024_q12_p2.txt"
		} else {
			infile = argsWithoutProg[1]
		}
		data, _ := os.ReadFile(infile)
		towerSpots := make(map[rune]point)
		targetSpots := make([]point, 0)
		hardSpots := make([]point, 0)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		rowIdx := 0
		for scanner.Scan() {
			for colIdx, ch := range scanner.Text() {
				if ch == 'A' || ch == 'B' || ch == 'C' {
					towerSpots[ch] = point{rowIdx, colIdx}
				}
				if ch == 'T' {
					targetSpots = append(targetSpots, point{rowIdx, colIdx})
				}
				if ch == 'H' {
					hardSpots = append(hardSpots, point{rowIdx, colIdx})
				}
			}
			rowIdx++
		}
		total := 0
		for _, target := range targetSpots {
			for idx, tower := range " ABC" {
				if idx == 0 {
					continue
				}
				q := (target.y - towerSpots[tower].y) - (target.x - towerSpots[tower].x)
				if q%3 == 0 {
					// fmt.Println("DBG", "towerloc", towerSpots[tower], "toweri ", ti, " tgti ", idx, " power ", q/3)
					total += (q / 3) * idx
				}
			}
		}
		for _, target := range hardSpots {
			for idx, tower := range " ABC" {
				if idx == 0 {
					continue
				}
				q := (target.y - towerSpots[tower].y) - (target.x - towerSpots[tower].x)
				if q%3 == 0 {
					// fmt.Println("DBG", "towerloc", towerSpots[tower], "toweri ", ti, " tgti ", idx, " power ", q/3)
					total += (q / 3) * idx * 2
				}
			}
		}
		fmt.Println("Part 2:", total)
	}
	{
		var infile string
		if len(argsWithoutProg) < 3 || argsWithoutProg[2] == "" {
			infile = "../input/everybody_codes_e2024_q12_p3.txt"
		} else {
			infile = argsWithoutProg[2]
		}
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		towers := []point{{0, 0}, {0, 1}, {0, 2}}
		trajArray := make([][][]point, 0)
		for ti, tower := range towers {
			trajArray = append(trajArray, make([][]point, 0, 2000))
			for power := range 2000 {
				trajArray[ti] = append(trajArray[ti], getTrajectory(tower, power))
			}
		}
		total := 0
		for scanner.Scan() {
			numStrs := strings.Fields(scanner.Text())
			col, _ := strconv.Atoi(numStrs[0])
			hgt, _ := strconv.Atoi(numStrs[1])
			solHeight := -1
			solScore := math.MaxInt
			for tidx := range towers {
				for power := range hgt/2 + 5 {
					trajectory := trajArray[tidx][power]
					for time, sp := range trajectory {
						if sp.x-sp.y == col-hgt && sp.x <= col-time {
							if sp.y > solHeight {
								solHeight = sp.y
								solScore = (tidx + 1) * power
							} else if sp.y == solHeight {
								solScore = min(solScore, (tidx+1)*power)
							}
						}
					}
				}
			}
			if solHeight < 0 {
				log.Fatal("Couldn't solve for ", col, hgt)
			}
			// fmt.Println("DBG", col, hgt, solScore)
			total += solScore
		}
		fmt.Println("Part 3:", total)
	}
}

func getTrajectory(start point, power int) []point {
	retval := make([]point, 0, 3*(power+2))
	for range power {
		retval = append(retval, start)
		start = point{start.x + 1, start.y + 1}
	}
	for range power {
		retval = append(retval, start)
		start = point{start.x + 1, start.y}
	}
	for start.y >= 0 {
		retval = append(retval, start)
		start = point{start.x + 1, start.y - 1}
	}
	return retval
}
