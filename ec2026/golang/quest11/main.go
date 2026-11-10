package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"maps"
	"os"
	"regexp"
	"strconv"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q11_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q11_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q11_p3.txt`, "the input for part 3")

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

func minpt(pt point) point {
	if pt.y < pt.x {
		return point{pt.y, pt.x}
	}
	return pt
}

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr1 := regexp.MustCompile(`portals: (\d+)`)
	prsr2 := regexp.MustCompile(`(\d+) *(\d+)`)
	scanner.Scan()
	matched1 := prsr1.FindStringSubmatch(scanner.Text())
	portalCount, _ := strconv.Atoi(matched1[1])
	scanner.Scan()
	side1 := make([]int, 0)
	side2 := make([]int, 0)
	for scanner.Scan() {
		matched2 := prsr2.FindStringSubmatch(scanner.Text())
		v1, _ := strconv.Atoi(matched2[1])
		v2, _ := strconv.Atoi(matched2[2])
		side1 = append(side1, v1)
		side2 = append(side2, v2)
	}

	connections := make(map[point]bool)
	for start := range portalCount {
		connections[minpt(point{start, (start + 1) % portalCount})] = true
	}
	// oPD := computePortalDist(portalCount, connections)
	// originalTotal := 0
	// for portal := range portalCount {
	// 	for end := range portalCount {
	// 		originalTotal += oPD[portal][end]
	// 	}
	// }
	originalTotal := (portalCount * portalCount * portalCount) / 4
	for idx := range side1 {
		connections[minpt(point{side1[idx], side2[idx]})] = true
	}

	portalDist := computePortalDist(portalCount, connections, make(map[point]int))

	sum := sumPortalDist(portalDist)
	return originalTotal - sum
}

func sumPortalDist(portalDist map[point]int) int {
	sum := 0
	for _, dist := range portalDist {
		sum += dist
	}
	return sum * 2
}

func computePortalDist(portalCount int, connections map[point]bool, portalDist map[point]int) map[point]int {
	portalsToWatch := make(map[int]bool)
	for start := range portalCount {
		for end := start + 1; end < portalCount; end++ {
			if _, ok := portalDist[point{start, end}]; !ok {
				portalDist[point{start, end}] = portalCount
			}
		}
		portalDist[point{start, start}] = 0
		portalsToWatch[start] = true
	}
	for len(portalsToWatch) > 0 {
		oPortalsToWatch := portalsToWatch
		portalsToWatch = make(map[int]bool)
		for portal := range portalCount {
			for other := range oPortalsToWatch {
				if connections[minpt(point{portal, other})] {
					for end := range portalCount {
						newval := 1 + portalDist[minpt(point{other, end})]
						if newval < portalDist[minpt(point{portal, end})] {
							portalsToWatch[portal] = true
							portalDist[minpt(point{portal, end})] = newval
						}
					}
				}
			}
		}
	}
	return portalDist
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr1 := regexp.MustCompile(`portals: (\d+)`)
	prsr2 := regexp.MustCompile(`(\d+) *(\d+)`)
	scanner.Scan()
	matched1 := prsr1.FindStringSubmatch(scanner.Text())
	portalCount, _ := strconv.Atoi(matched1[1])
	scanner.Scan()
	side1 := make([]int, 0)
	side2 := make([]int, 0)
	for scanner.Scan() {
		matched2 := prsr2.FindStringSubmatch(scanner.Text())
		v1, _ := strconv.Atoi(matched2[1])
		v2, _ := strconv.Atoi(matched2[2])
		side1 = append(side1, v1, v2)
		side2 = append(side2, v2, v1)
	}

	connections := make(map[point]bool)
	for start := range portalCount {
		connections[point{start, (start + 1) % portalCount}] = true
		connections[point{start, (start + portalCount - 1) % portalCount}] = true
		connections[point{start, (start + 2) % portalCount}] = true
		connections[point{start, (start + portalCount - 2) % portalCount}] = true
		connections[point{start, (start + 3) % portalCount}] = true
		connections[point{start, (start + portalCount - 3) % portalCount}] = true
	}
	partialDist := make(map[point]int)
	originalTotal := sumPortalDist(computePortalDist(portalCount, connections, partialDist))
	// fmt.Println("originalTotal", originalTotal)
	for idx := range side1 {
		connections[point{side1[idx] - 1, side2[idx] - 1}] = true
	}

	resonance := 0
	for start := range portalCount {
		for end := start + 1; end < portalCount; end++ {
			if !connections[point{start, end}] {
				connections[point{start, end}] = true
				connections[point{end, start}] = true

				cloned := maps.Clone(partialDist)

				newTotal := sumPortalDist(computePortalDist(portalCount, connections, cloned))
				// fmt.Println("newTotal", newTotal)
				resonance = max(resonance, originalTotal-newTotal)
				delete(connections, point{start, end})
				delete(connections, point{end, start})
			}
		}
	}

	return resonance
}

func doProblem3(data []byte) any {
	// 	data = []byte(`portals: 32

	// 1 10
	// 20 30
	// 15 29
	// 13 10`)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr1 := regexp.MustCompile(`portals: (\d+)`)
	prsr2 := regexp.MustCompile(`(\d+) *(\d+)`)
	scanner.Scan()
	matched1 := prsr1.FindStringSubmatch(scanner.Text())
	portalCount, _ := strconv.Atoi(matched1[1])
	scanner.Scan()
	side1 := make([]int, 0)
	side2 := make([]int, 0)
	for scanner.Scan() {
		matched2 := prsr2.FindStringSubmatch(scanner.Text())
		v1, _ := strconv.Atoi(matched2[1])
		v2, _ := strconv.Atoi(matched2[2])
		side1 = append(side1, v1, v2)
		side2 = append(side2, v2, v1)
	}

	connections := make(map[point]bool)
	for start := range portalCount {
		connections[point{start, (start + 1) % portalCount}] = true
		connections[point{start, (start + portalCount - 1) % portalCount}] = true
		connections[point{start, (start + 2) % portalCount}] = true
		connections[point{start, (start + portalCount - 2) % portalCount}] = true
		connections[point{start, (start + 3) % portalCount}] = true
		connections[point{start, (start + portalCount - 3) % portalCount}] = true
	}
	partialDist := make(map[point]int)
	originalTotal := sumPortalDist(computePortalDist(portalCount, connections, partialDist))
	// fmt.Println("originalTotal", originalTotal)
	for idx := range side1 {
		connections[point{side1[idx] - 1, side2[idx] - 1}] = true
	}
	resonance := 0
	for topGem := range portalCount / 8 {
		gems := make([]int, 0, 8)
		for gem := topGem; gem < portalCount; gem += portalCount / 8 {
			gems = append(gems, gem)
		}
		for _, arrangement := range getGemPairs(gems) {
			added := make([]bool, len(arrangement))
			for idx, pt := range arrangement {
				if !connections[pt] {
					connections[pt] = true
					connections[point{pt.y, pt.x}] = true
					added[idx] = true
				}
			}
			cloned := maps.Clone(partialDist)
			newTotal := sumPortalDist(computePortalDist(portalCount, connections, cloned))
			// fmt.Println("newTotal", newTotal, "for", arrangement)
			resonance = max(resonance, originalTotal-newTotal)
			for idx, pt := range arrangement {
				if added[idx] {
					delete(connections, pt)
					delete(connections, point{pt.y, pt.x})
				}
			}
		}
	}
	return resonance
}

func getGemPairs(gems []int) [][]point {
	return getGemPairs1(gems, make([]bool, len(gems)))
}

func getGemPairs1(gems []int, used []bool) [][]point {
	retval := make([][]point, 0)
	foundSomething := false
	for lowidx := 0; lowidx < len(gems); lowidx++ {
		if used[lowidx] {
			continue
		}
		for hiidx := lowidx + 1; hiidx < len(gems); hiidx++ {
			if used[hiidx] {
				continue
			}
			foundSomething = true
			mypt := point{gems[lowidx], gems[hiidx]}
			used[lowidx] = true
			used[hiidx] = true
			downStream := getGemPairs1(gems, used)
			for _, val := range downStream {
				val = append(val, mypt)
				retval = append(retval, val)
			}
			used[lowidx] = false
			used[hiidx] = false
		}
		break
	}
	if !foundSomething {
		return [][]point{make([]point, 0)}
	}
	return retval
}
