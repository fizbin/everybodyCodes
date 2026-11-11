package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q12_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q12_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q12_p3.txt`, "the input for part 3")

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

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make([][]byte, 0)
	for scanner.Scan() {
		grid = append(grid, bytes.Clone(scanner.Bytes()))
	}
	height := len(grid)
	width := len(grid[0])
	ans := 0
	for top := range height {
		for left := range width {
			if grid[top][left] == '#' {
				// now find the right
				var right, bottom int
				for right = left; (right < width) && grid[top][right] == '#'; right++ {
				}
				for bottom = top; (bottom < height) && grid[bottom][left] == '#'; bottom++ {
				}
				ans += min(bottom-top, right-left)
				for x := left; x < right; x++ {
					for y := top; y < bottom; y++ {
						grid[y][x] = '*'
					}
				}
			}
		}
	}
	return ans
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make(map[point]byte)
	width := 0
	height := 0
	for scanner.Scan() {
		width = max(width, len(scanner.Bytes()))
		for col, bval := range scanner.Bytes() {
			grid[point{height, col}] = bval
		}
		height++
	}
	ans := 0
	for top := range height {
		for left := range width {
			if grid[point{top, left}] == '#' {
				// now find all connected
				myPoints := make(map[point]bool)
				myPoints[point{top, left}] = true
				border := []point{{top, left}}
				for len(border) > 0 {
					oldBorder := border
					border = nil
					for _, pt := range oldBorder {
						nbrs := neighbors(pt, height, width)
						for _, nbr := range nbrs {
							if grid[nbr] == '#' && !myPoints[nbr] {
								myPoints[nbr] = true
								border = append(border, nbr)
							}
						}
					}
				}
				ans += spellsNeededP2(myPoints)
				for pt := range myPoints {
					grid[pt] = '*'
				}
			}
		}
	}
	return ans
}

func spellsNeededP2(targets map[point]bool) int {
	maxX, maxY, minX, minY := -1, -1, 9999999, 9999999
	for pt := range targets {
		maxX = max(maxX, pt.x)
		maxY = max(maxY, pt.y)
		minX = min(minX, pt.x)
		minY = min(minY, pt.y)
	}
	var workingTargets map[point]bool
	if !targets[point{maxX, maxY}] {
		// fmt.Print("case 1 ")
		workingTargets = targets
	} else if !targets[point{maxX, minY}] {
		// fmt.Print("case 2 ")
		workingTargets = make(map[point]bool)
		for pt := range targets {
			workingTargets[point{pt.x, -pt.y}] = true
		}
		maxY, minY = -minY, -maxY
	} else if !targets[point{minX, maxY}] {
		// fmt.Print("case 3 ")
		workingTargets = make(map[point]bool)
		for pt := range targets {
			workingTargets[point{-pt.x, pt.y}] = true
		}
		maxX, minX = -minX, -maxX
	} else if !targets[point{minX, minY}] {
		// fmt.Print("case 4 ")
		workingTargets = make(map[point]bool)
		for pt := range targets {
			workingTargets[point{-pt.x, -pt.y}] = true
		}
		maxY, minY = -minY, -maxY
		maxX, minX = -minX, -maxX
	} else {
		// rectangle!
		// fmt.Println("Rectangle @", minX, minY, ":", maxX, maxY, "ans", min(maxX-minX+1, maxY-minY+1))
		return min(maxX-minX+1, maxY-minY+1)
	}
	holeLeft, holeTop := 0, 0
	for holeLeft = minY; (holeLeft <= maxY) && workingTargets[point{maxX, holeLeft}]; holeLeft++ {
	}
	for holeTop = minX; (holeTop <= maxX) && workingTargets[point{holeTop, maxY}]; holeTop++ {
	}
	ans := min(
		// shoot top-left block horizontally
		(holeTop-minX)+min(maxX-holeTop+1, holeLeft-minY),
		// shoot top-left block vertically
		(holeLeft-minY)+min(maxY-holeLeft+1, holeTop-minX),
	)
	// fmt.Println("L shape @", minX, minY, ":", maxX, maxY, "ans", ans)
	return ans
}

func neighbors(pt point, height, width int) []point {
	ans := make([]point, 0)
	if pt.x > 0 {
		ans = append(ans, point{pt.x - 1, pt.y})
	}
	if pt.y > 0 {
		ans = append(ans, point{pt.x, pt.y - 1})
	}
	if pt.x+1 < height {
		ans = append(ans, point{pt.x + 1, pt.y})
	}
	if pt.y+1 < width {
		ans = append(ans, point{pt.x, pt.y + 1})
	}
	return ans
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid := make(map[point]byte)
	width := 0
	height := 0
	for scanner.Scan() {
		width = max(width, len(scanner.Bytes()))
		for col, bval := range scanner.Bytes() {
			grid[point{height, col}] = bval
		}
		height++
	}
	ans := 0
	for top := range height {
		for left := range width {
			if grid[point{top, left}] == '#' {
				// now find all connected
				myPoints := make(map[point]bool)
				myPoints[point{top, left}] = true
				border := []point{{top, left}}
				for len(border) > 0 {
					oldBorder := border
					border = nil
					for _, pt := range oldBorder {
						nbrs := neighbors(pt, height, width)
						for _, nbr := range nbrs {
							if grid[nbr] == '#' && !myPoints[nbr] {
								myPoints[nbr] = true
								border = append(border, nbr)
							}
						}
					}
				}
				ans += spellsNeededP3(myPoints)
				for pt := range myPoints {
					grid[pt] = '*'
				}
			}
		}
	}
	return ans
}

func spellsNeededP3(targets map[point]bool) int {
	maxX, maxY, minX, minY := -1, -1, 9999999, 9999999
	for pt := range targets {
		maxX = max(maxX, pt.x)
		maxY = max(maxY, pt.y)
		minX = min(minX, pt.x)
		minY = min(minY, pt.y)
	}
	holeTop, holeLeft, holeBottom, holeRight := maxX+1, maxY+1, minX-1, minY-1
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			if !targets[point{x, y}] {
				holeTop = min(holeTop, x)
				holeLeft = min(holeLeft, y)
				holeBottom = max(holeBottom, x)
				holeRight = max(holeRight, y)
			}
		}
	}

	// double-check we only have "rectangle + rectangular hole"
	isNormal := true
	for x := minX; isNormal && x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			inHole := (holeLeft <= y) && (y <= holeRight) && (holeTop <= x) && (x <= holeBottom)
			if inHole == targets[point{x, y}] {
				isNormal = false
				break
			}
		}
	}
	if !isNormal {
		// fmt.Println("Not a normal thing @", minX, minY, maxX, maxY, holeTop, holeLeft, holeBottom, holeRight)
		// panic("Not normal")

		// Can we fix it?
		for x := minX; x <= maxX; x++ {
			for y := minY; y <= maxY; y++ {
				inHole := (holeLeft <= y) && (y <= holeRight) && (holeTop <= x) && (x <= holeBottom)
				if _, ok := targets[point{x, y}]; inHole == ok {
					// we got a problem. Rip it out
					delete(targets, point{x, y})
				}
			}
		}

	}
	// fmt.Println("Rect @ ", minX, minY, ":", maxX, maxY, "with hole", holeTop, holeLeft, ":", holeBottom, holeRight)
	return spellsKnowingHole(minX, minY, maxX, maxY, holeTop, holeLeft, holeBottom, holeRight)
}

func spellsKnowingHole(minX, minY, maxX, maxY, holeTop, holeLeft, holeBottom, holeRight int) int {
	// fmt.Println("Calling sKH @ ", minX, minY, maxX, maxY, holeTop, holeLeft, holeBottom, holeRight)

	if minX > maxX || minY > maxY {
		return 0
	}
	if holeTop > maxX {
		// no hole; rectangle!
		return min(maxX-minX+1, maxY-minY+1)
	}
	// so we know the box, and we know the hole
	// now we can think about eight regions:
	//  1 2 3
	//  4 H 6
	//  7 8 9
	// each region must be eliminated

	// if region 4 or region 2 doesn't exist, flip coordinates so that it does
	if holeLeft == minY {
		maxY, minY = -minY, -maxY
		holeLeft, holeRight = -holeRight, -holeLeft
	}
	if holeTop == minX {
		maxX, minX = -minX, -maxX
		holeTop, holeBottom = -holeBottom, -holeTop
	}
	// so now, regions 1, 2, and 4 always exist

	// Is the hole aligned at bottom-right? If so, L shape
	if holeBottom == maxX && holeRight == maxY {
		return min(
			// shoot block 1 horiz
			(holeTop-minX)+min(maxX-holeTop+1, holeLeft-minY),
			// shoot block 1 vertically
			(holeLeft-minY)+min(maxY-holeLeft+1, holeTop-minX),
		)
	}

	// so possibly region 6 or 8 doesn't exist, but all others do
	ans := 999999999

	// possible strategies:
	// 1 horiz.
	// with 1 wiped out horizontally, we have:
	// - - -
	// 4 H 6
	// 7 8 9
	// We can then treat the rest in one of four ways:
	// L shape of 4 7 8, rectangle of 6 9
	// L shape of 4 7 8 9, rectangle of 6
	// L shape of 7 8 9 6, rectangle of 4
	// L shape of 8 9 6, rectangle of 4 7
	{
		oneHoriz := holeTop - minX
		strategy1 := (spellsKnowingHole(holeTop, minY, maxX, holeRight, holeTop, holeLeft, holeBottom, holeRight) +
			spellsKnowingHole(holeTop, holeRight+1, maxX, maxY, maxX+1, maxY+1, maxX+1, maxY+1))
		strategy2 := (spellsKnowingHole(holeTop, minY, maxX, maxY, holeTop, holeLeft, holeBottom, maxY) +
			spellsKnowingHole(holeTop, holeRight+1, holeBottom, maxY, holeBottom+1, maxY+1, holeBottom+1, maxY+1))
		strategy3 := (spellsKnowingHole(holeTop, minY, maxX, maxY, holeTop, minY, holeBottom, holeRight) +
			spellsKnowingHole(holeTop, minY, holeBottom, holeLeft-1, holeBottom+1, holeLeft, holeBottom+1, holeLeft))
		strategy4 := (spellsKnowingHole(holeTop, holeLeft, maxX, maxY, holeTop, holeLeft, holeBottom, holeRight) +
			spellsKnowingHole(holeTop, minY, maxX, holeLeft-1, maxX+1, holeLeft, maxX+1, holeLeft))
		ans = min(ans, oneHoriz+strategy1, oneHoriz+strategy2, oneHoriz+strategy3, oneHoriz+strategy4)
	}

	// Be lazy and swap variables, then do the same thing again:
	{
		minY, minX, maxY, maxX, holeLeft, holeTop, holeRight, holeBottom := minX, minY, maxX, maxY, holeTop, holeLeft, holeBottom, holeRight
		oneHoriz := holeTop - minX
		strategy1 := (spellsKnowingHole(holeTop, minY, maxX, holeRight, holeTop, holeLeft, holeBottom, holeRight) +
			spellsKnowingHole(holeTop, holeRight+1, maxX, maxY, maxX+1, maxY+1, maxX+1, maxY+1))
		strategy2 := (spellsKnowingHole(holeTop, minY, maxX, maxY, holeTop, holeLeft, holeBottom, maxY) +
			spellsKnowingHole(holeTop, holeRight+1, holeBottom, maxY, holeBottom+1, maxY+1, holeBottom+1, maxY+1))
		strategy3 := (spellsKnowingHole(holeTop, minY, maxX, maxY, holeTop, minY, holeBottom, holeRight) +
			spellsKnowingHole(holeTop, minY, holeBottom, holeLeft-1, holeBottom+1, holeLeft, holeBottom+1, holeLeft))
		strategy4 := (spellsKnowingHole(holeTop, holeLeft, maxX, maxY, holeTop, holeLeft, holeBottom, holeRight) +
			spellsKnowingHole(holeTop, minY, maxX, holeLeft-1, maxX+1, holeLeft, maxX+1, holeLeft))
		ans = min(ans, oneHoriz+strategy1, oneHoriz+strategy2, oneHoriz+strategy3, oneHoriz+strategy4)
	}

	return ans
}
