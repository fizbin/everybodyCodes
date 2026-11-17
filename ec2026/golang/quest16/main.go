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
	"slices"
	"unicode"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q16_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q16_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q16_p3.txt`, "the input for part 3")

const doingPart = 3

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

// I call the coords x and y, but they're really row and column - "x" is
// vertical and increases as you go down, "y" is horizontal
// type point struct {
// 	x, y int
// }

// func addpt(a, b point) point {
// 	return point{a.x + b.x, a.y + b.y}
// }

// func leftpt(a point) point {
// 	return point{-a.y, a.x}
// }

// func rightpt(a point) point {
// 	return point{a.y, -a.x}
// }

func doProblem1(data []byte) any {
	// data = []byte(`RRGGBBGRRGGBBG`)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	workingShelf := bytes.Clone(scanner.Bytes())
	flips := 0
	for true {
		recentLocations := make(map[byte]int)
		startFlip := math.MaxInt
		var stopFlip int
		for i := 0; i < len(workingShelf); i++ {
			ch := workingShelf[i]
			prevLoc, ok := recentLocations[ch]
			if i == prevLoc+1 || !ok {
				recentLocations[ch] = i
			} else {
				if prevLoc+1 < startFlip {
					startFlip = prevLoc + 1
					for stopFlip = len(workingShelf) - 1; workingShelf[stopFlip] != ch; stopFlip-- {
					}
				}
			}
		}
		if startFlip > len(workingShelf) {
			break
		}
		slices.Reverse(workingShelf[startFlip : stopFlip+1])
		flips++
	}
	return flips
}

func doProblem2(data []byte) any {
	return doProblem3guts(data, 10)
}

func eliminateOneColor(data []byte) ([]byte, int) {
	if len(data) == 0 {
		return nil, 0
	}
	targetColor := data[0]
	// so we have something like (where T is target color)
	// TTT other1 TT other2 TT other3 TT .... TT otherN
	// we're gonna end up with
	// .... 3rehto otherN-2 2rehto otherN-1 1rehto otherN
	// I think that the easiest way to build that up is in reverse,
	// then reverse in the end before returning
	retN := 0
	retData := make([]byte, 0, len(data))
	idxFront := 0
	idxBack := len(data) - 1
	for ; idxFront <= idxBack; idxFront++ {
		if data[idxFront] != targetColor {
			break
		}
	}
	for idxFront <= idxBack {
		// idxFront does not point at targetColor
		// idxBack might point at targetColor, but likely not
		for idxFront <= idxBack && data[idxBack] != targetColor {
			retData = append(retData, data[idxBack])
			idxBack--
		}
		// idxBack now points at targetColor
		retN += idxBack - idxFront + 1
		for idxFront <= idxBack && data[idxBack] == targetColor {
			idxBack--
		}
		// idxBack now does not point at targetColor
		for idxFront <= idxBack && data[idxFront] != targetColor {
			retData = append(retData, data[idxFront])
			idxFront++
		}
		// idxFront now points at targetColor
		retN += idxBack - idxFront + 1
		for idxFront <= idxBack && data[idxFront] == targetColor {
			idxFront++
		}
	}
	slices.Reverse(retData)
	return retData, retN
}

func doProblem3guts(data []byte, nCopies int) int {
	// data = []byte(`RRYGGYBBYG`)
	dL := len(data)
	for unicode.IsSpace(rune(data[dL-1])) {
		dL--
	}
	myData := slices.Repeat(data, nCopies)
	retval := 0
	for myData != nil {
		newData, newFlips := eliminateOneColor(myData)
		myData = newData
		retval += newFlips
	}
	return retval
}

func doProblem3(data []byte) any {
	return doProblem3guts(data, 100)
}
