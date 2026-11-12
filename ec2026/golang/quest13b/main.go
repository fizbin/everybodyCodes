package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q13_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q13_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q13_p3.txt`, "the input for part 3")

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

// type point struct {
// 	x, y int
// }

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	handCount, handDir, handPos, handMoves := parsePortal(scanner)
	return workPortal(handCount, handDir, handPos, handMoves)
}

func parsePortal(scanner *bufio.Scanner) ([]int, []int, []int, []int) {
	var handCount []int
	var handDir []int // 1 or -1 for CW or CCW
	var handPos []int
	var handMoves []int
	scanner.Scan()
	{
		pieces := strings.Split(scanner.Text(), " ")
		if pieces[0] != "hands" {
			panic("Bad first string " + scanner.Text())
		}
		for _, valStr := range strings.Split(pieces[1], ":") {
			handN, _ := strconv.Atoi(valStr)
			handCount = append(handCount, handN)
		}
	}
	scanner.Scan()
	{
		pieces := strings.Split(scanner.Text(), " ")
		if pieces[0] != "direction" {
			panic("Bad 2nd string " + scanner.Text())
		}
		for _, valStr := range strings.Split(pieces[1], ":") {
			switch valStr {
			case "CW":
				handDir = append(handDir, 1)
			case "CCW":
				handDir = append(handDir, -1)
			default:
				panic("Bad direction " + valStr + " in " + scanner.Text())
			}
		}
	}
	scanner.Scan()
	{
		pieces := strings.Split(scanner.Text(), " ")
		if pieces[0] != "current" {
			panic("Bad 3rd string " + scanner.Text())
		}
		for _, valStr := range strings.Split(pieces[1], ":") {
			handN, _ := strconv.Atoi(valStr)
			handPos = append(handPos, handN)
		}
	}
	scanner.Scan()
	{
		pieces := strings.Split(scanner.Text(), " ")
		if pieces[0] != "move" {
			panic("Bad 4th string " + scanner.Text())
		}
		for _, valStr := range strings.Split(pieces[1], ":") {
			handN, _ := strconv.Atoi(valStr)
			handMoves = append(handMoves, handN)
		}
	}
	return handCount, handDir, handPos, handMoves
}

func workPortal(handCount []int, handDir []int, handPos []int, handMoves []int) int {
	ans := 0
	// Take each pair of hands and find the intersections for that pair
	for idx1 := range 3 {
		for idx2 := idx1 + 1; idx2 < 3; idx2++ {
			var handA, handB int
			// choose handA so that it stops first
			if handMoves[idx2] < handMoves[idx1] {
				handA = idx2
				handB = idx1
			} else {
				handB = idx2
				handA = idx1
			}

			// Convert to integer coordinates; "raw" is for using a
			// clock-relative reference frame.
			handLCM := lcm(handCount[handA], handCount[handB])
			speedAraw := handLCM / handCount[handA] * handDir[handA]
			speedBraw := handLCM / handCount[handB] * handDir[handB]
			posAraw := handLCM / handCount[handA] * handPos[handA]
			posBraw := handLCM / handCount[handB] * handPos[handB]

			// Now work two parts: when both hands are moving, and after
			// handA stops.
			// Change coordinates to use a reference frame trick: hold
			// handA still at 0 and normalize movement direction to always be
			// clockwise/positive
			{
				posB := posBraw - posAraw
				speedB := speedBraw - speedAraw
				if speedB < 0 {
					speedB = -speedB
					posB = -posB
				}
				posB = (posB%handLCM + handLCM) % handLCM
				finalPos := posB + speedB*handMoves[handA]
				ans += finalPos / handLCM
			}
			{
				posB := posBraw - posAraw + (speedBraw-speedAraw)*handMoves[handA]
				speedB := speedBraw
				if speedB < 0 {
					speedB = -speedB
					posB = -posB
				}
				posB = (posB%handLCM + handLCM) % handLCM
				finalPos := posB + speedB*(handMoves[handB]-handMoves[handA])
				ans += finalPos / handLCM
			}
		}
	}

	return ans
}

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	if a > b {
		a, b = b, a
	}
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

func lcm(a, b int) int {
	return b * (a / gcd(a, b))
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	isFirst := true
	ans := 0
	clockNum := 0
	for isFirst || scanner.Scan() {
		isFirst = false
		clockNum++
		handCount, handDir, handPos, handMoves := parsePortal(scanner)
		clockAns := workPortal(handCount, handDir, handPos, handMoves)
		// fmt.Println("clock", clockNum, "meetings", clockAns)
		ans += clockAns
	}
	return ans
}
func doProblem3(data []byte) any {
	return doProblem2(data)
}
