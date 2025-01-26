package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
)

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q17_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 1:", doProblem1(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 2 || argsWithoutProg[1] == "" {
			infile = "../input/everybody_codes_e2024_q17_p2.txt"
		} else {
			infile = argsWithoutProg[1]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 2:", doProblem1(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 3 || argsWithoutProg[2] == "" {
			infile = "../input/everybody_codes_e2024_q17_p3.txt"
		} else {
			infile = argsWithoutProg[2]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 3:", doProblem3(data))
	}
}

type point struct{ x, y int }

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func doProblem1(data []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	stars := make([]point, 0)
	height := 0
	for scanner.Scan() {
		for colIdx, ch := range scanner.Text() {
			if ch == '*' {
				stars = append(stars, point{height, colIdx})
			}
		}
		height++
	}

	return strconv.Itoa(ConstellationSize(stars))
}

func ConstellationSize(stars []point) int {
	distances := make([][]int, 0, len(stars))
	for _, stari := range stars {
		row := make([]int, 0, len(stars))
		for _, starj := range stars {
			row = append(row, abs(stari.x-starj.x)+abs(stari.y-starj.y))
		}
		distances = append(distances, row)
	}
	pointsInTree := make(map[int]bool)
	pointsNotInTree := make(map[int]bool)
	for i := range stars {
		pointsNotInTree[i] = true
	}
	pointsInTree[0] = true
	spanningTreeEdgeLen := 0
	delete(pointsNotInTree, 0)
	for len(pointsNotInTree) > 0 {
		// minInIdx := -1
		minOutIdx := -1
		minDist := math.MaxInt
		for innie := range pointsInTree {
			for outie := range pointsNotInTree {
				if distances[innie][outie] < minDist {
					// minInIdx = innie
					minOutIdx = outie
					minDist = distances[innie][outie]
				}
			}
		}
		delete(pointsNotInTree, minOutIdx)
		pointsInTree[minOutIdx] = true
		spanningTreeEdgeLen += minDist
		// fmt.Println("Edge from", stars[minInIdx], "to", stars[minOutIdx], "length", minDist)
	}
	ans := spanningTreeEdgeLen + len(pointsInTree)
	return ans
}

func doProblem3(data []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	stars := make([]point, 0)
	height := 0
	for scanner.Scan() {
		for colIdx, ch := range scanner.Text() {
			if ch == '*' {
				stars = append(stars, point{height, colIdx})
			}
		}
		height++
	}
	starCategory := make(map[point]int)
	for idx, star := range stars {
		starCategory[star] = idx
	}
	done := false
	for !done {
		done = true
		for idx, starA := range stars {
			for _, starB := range stars[idx+1:] {
				starDist := abs(starA.x-starB.x) + abs(starA.y-starB.y)
				if starDist < 6 && starCategory[starA] != starCategory[starB] {
					starCategory[starA] = min(starCategory[starA], starCategory[starB])
					starCategory[starB] = starCategory[starA]
					done = false
				}
			}
		}
	}
	constellationHolder := make(map[int][]point)
	for pt, cat := range starCategory {
		constellationHolder[cat] = append(constellationHolder[cat], pt)
	}
	sizes := make([]int, 0, len(starCategory))
	for _, starArray := range constellationHolder {
		sizes = append(sizes, ConstellationSize(starArray))
	}
	slices.Sort(sizes)
	ans := sizes[len(sizes)-1] * sizes[len(sizes)-2] * sizes[len(sizes)-3]
	return strconv.Itoa(ans)
}
