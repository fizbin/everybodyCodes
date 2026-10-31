package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q05_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q05_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q05_p3.txt`, "the input for part 3")

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

type coord struct {
	x, y int
}

func getNeighbors(where coord, height int, width int) []coord {
	ans := make([]coord, 0, 8)
	for xdiff := -1; xdiff < 2; xdiff++ {
		for ydiff := -1; ydiff < 2; ydiff++ {
			if xdiff == 0 && ydiff == 0 {
				continue
			}
			newWhere := coord{where.x + xdiff, where.y + ydiff}
			if 0 <= newWhere.x && newWhere.x < height && 0 <= newWhere.y && newWhere.y < width {
				ans = append(ans, newWhere)
			}
		}
	}
	return ans
}

func doProblem1core(data []byte) (int, [][]byte) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	myMap := make([][]byte, 0)
	width := 0
	for scanner.Scan() {
		currBytes := scanner.Bytes()
		end := len(currBytes)
		for currBytes[end-1] != '0' && currBytes[end-1] != '1' {
			end--
		}
		myMap = append(myMap, bytes.Clone(scanner.Bytes()[:end]))
		if width == 0 {
			width = end
		}
		if width != end {
			panic(fmt.Sprintf("Bad end at %d vs %d", width, end))
		}
	}
	height := len(myMap)
	midX := height / 2
	midY := width / 2
	myMap[midX][midY] = '#'
	landCount := 1
	border := getNeighbors(coord{midX, midY}, height, width)
	for len(border) > 0 {
		oldBorder := border
		border = make([]coord, 0, len(oldBorder))
		for _, spot := range oldBorder {
			if myMap[spot.x][spot.y] == '#' {
				continue
			}
			spotBorder := getNeighbors(spot, height, width)
			count := 0
			for _, nbr := range spotBorder {
				if myMap[nbr.x][nbr.y] == '1' || myMap[nbr.x][nbr.y] == '#' {
					count++
				}
			}
			if count >= 5 {
				myMap[spot.x][spot.y] = '#'
				border = append(border, spotBorder...)
				landCount++
			}
		}
	}

	return landCount, myMap
}

func doProblem1(data []byte) any {
	landCount, _ := doProblem1core(data)
	return landCount
}

func doProblem2core(data []byte) (int, [][]byte) {
	_, myMap := doProblem1core(data)
	width := len(myMap[0])
	height := len(myMap)
	sandCount := 0
	sandBorder := make([]coord, 0)
	for x, row := range myMap {
		for y, spotCh := range row {
			if spotCh == '#' {
				spotNbrs := getNeighbors(coord{x, y}, height, width)
				if len(spotNbrs) < 8 {
					spotCh = '='
				} else {
					for _, nbr := range spotNbrs {
						if myMap[nbr.x][nbr.y] == '0' || myMap[nbr.x][nbr.y] == '1' {
							spotCh = '='
						}
					}
				}
				if spotCh == '=' {
					sandCount++
					myMap[x][y] = '='
					sandBorder = append(sandBorder, spotNbrs...)
				}
			}
		}
	}
	for len(sandBorder) > 0 {
		oldBorder := sandBorder
		sandBorder = make([]coord, 0)
		for _, spot := range oldBorder {
			if myMap[spot.x][spot.y] != '#' {
				continue
			}
			spotNbrs := getNeighbors(spot, height, width)
			count := 0
			for _, nbr := range spotNbrs {
				if myMap[nbr.x][nbr.y] == '=' {
					count++
				}
			}
			if count >= 5 {
				myMap[spot.x][spot.y] = '='
				sandCount++
				sandBorder = append(sandBorder, spotNbrs...)
			}
		}
	}
	return sandCount, myMap
}

func doProblem2(data []byte) any {
	sandCount, _ := doProblem2core(data)
	return sandCount
}

func doProblem3(data []byte) any {
	_, myMap := doProblem2core(data)
	width := len(myMap[0])
	height := len(myMap)
	treeCount := 0
	treeBorder := make([]coord, 0)
	for centerX := 2; centerX+2 < height; centerX++ {
		for centerY := 2; centerY+2 < width; centerY++ {
			if myMap[centerX][centerY] != '#' {
				continue
			}
			foundSpots := 0
			for xdiff := -2; xdiff <= 2; xdiff++ {
				for ydiff := -2; ydiff <= 2; ydiff++ {
					if (xdiff == ydiff || xdiff == -ydiff) && (xdiff == 2 || xdiff == -2) {
						continue
					}
					found := myMap[centerX+xdiff][centerY+ydiff]
					if found == '#' || found == '@' {
						foundSpots++
					}
				}
			}
			if foundSpots == 21 {
				myMap[centerX][centerY] = '@'
				treeCount++
				treeBorder = append(treeBorder, getNeighbors(coord{centerX, centerY}, height, width)...)
			}
		}
	}
	for growthStep := byte('A'); growthStep < 'C'; growthStep++ {
		oldBorder := treeBorder
		treeBorder = make([]coord, 0)
		for _, spot := range oldBorder {
			if myMap[spot.x][spot.y] != '#' {
				continue
			}
			localTrees := 0
			nbrs := getNeighbors(spot, height, width)
			for _, nbr := range nbrs {
				nVal := myMap[nbr.x][nbr.y]
				if nVal == '=' {
					localTrees = -20
				} else if nVal >= '@' && nVal < growthStep {
					localTrees++
				}
			}
			if localTrees >= 4 {
				myMap[spot.x][spot.y] = growthStep
				treeCount++
				treeBorder = append(treeBorder, nbrs...)
			}
		}
	}
	// for _, row := range myMap {
	// 	fmt.Println(string(row))
	// }
	return treeCount
}
