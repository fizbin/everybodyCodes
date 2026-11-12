package main

import (
	"bufio"
	"bytes"
	"cmp"
	"flag"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q13_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q13_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q13_p3.txt`, "the input for part 3")

const doingPart = 3
const ClockDebug = false

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
	hands, handDir, handPos, handMove := parsePortal(scanner)
	return workPortal(hands, handDir, handPos, handMove)
}

func parsePortal(scanner *bufio.Scanner) ([]int, []int, []int, []int) {
	var hands []int
	var handDir []int // 1 or -1 for CW or CCW
	var handPos []int
	var handMove []int
	scanner.Scan()
	{
		pieces := strings.Split(scanner.Text(), " ")
		if pieces[0] != "hands" {
			panic("Bad first string " + scanner.Text())
		}
		for _, valStr := range strings.Split(pieces[1], ":") {
			handN, _ := strconv.Atoi(valStr)
			hands = append(hands, handN)
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
			handMove = append(handMove, handN)
		}
	}
	return hands, handDir, handPos, handMove
}

func workPortal(hands []int, handDir []int, handPos []int, handMove []int) int {
	handLCM := lcm(lcm(hands[0], hands[1]), hands[2])
	handSpeed := make([]int, len(hands))
	closeVal := 0
	for idx, hval := range hands {
		baseMul := handLCM / hval
		handSpeed[idx] = baseMul * handDir[idx]
		handPos[idx] *= baseMul

		closeVal = max(closeVal, 2*baseMul)
	}

	if ClockDebug {
		fmt.Println("Using", handLCM, "uniform positions; hands at ", handPos, "with speed", handSpeed)
	}
	ans := 0
	prevRound := 0
	prevAns := 0
	var prevPos []int
	var prevSpeed []int
	for roundNum := 0; roundNum < max(handMove[0], handMove[1], handMove[2]); roundNum++ {
		for idx, moveVal := range handMove {
			if roundNum >= moveVal && handSpeed[idx] != 0 {
				// fmt.Println("Invalidating cache at", roundNum)
				handSpeed[idx] = 0
				prevPos = nil
			}
		}
		if prevPos != nil && slices.Equal(prevPos, handPos) && slices.Equal(prevSpeed, handSpeed) {
			// fast forward!
			stepSize := roundNum - prevRound
			roomToMove := math.MaxInt
			for idx := range 3 {
				if handSpeed[idx] != 0 {
					roomToMove = min(roomToMove, handMove[idx]-roundNum-2)
				}
			}
			if roomToMove > stepSize {
				nsteps := roomToMove / stepSize
				// fmt.Println("Stepping from ", roundNum, "to", roundNum+nsteps*stepSize, "von", roomToMove, stepSize, handMove, roundNum)
				ansJump := ans - prevAns
				roundNum += nsteps * stepSize
				ans += nsteps * ansJump
			}
		}
		if prevPos == nil {
			// fmt.Println("Recording at", roundNum)
			prevPos = slices.Clone(handPos)
			prevSpeed = slices.Clone(handSpeed)
			prevAns = ans
			prevRound = roundNum
		}
		handPosNew := []int{
			(handPos[0] + handSpeed[0] + handLCM) % handLCM,
			(handPos[1] + handSpeed[1] + handLCM) % handLCM,
			(handPos[2] + handSpeed[2] + handLCM) % handLCM,
		}
		for _, handNums := range [][]int{{0, 1}, {1, 2}, {0, 2}} {
			handA := handNums[0]
			handB := handNums[1]
			posDiff0 := handPos[handA] - handPos[handB]
			if posDiff0 < 0 {
				posDiff0 = -posDiff0
			}
			posDiff0 %= handLCM
			if posDiff0 != 0 && (posDiff0 <= closeVal || posDiff0+closeVal >= handLCM) {
				// close enough that we need to check
				posA0 := handPos[handA]
				posB0 := handPos[handB]
				posA1 := handPosNew[handA]
				posB1 := handPosNew[handB]

				for (handSpeed[handA] != 0 && (posA1-posA0)/handSpeed[handA] != 1) ||
					(handSpeed[handB] != 0 && (posB1-posB0)/handSpeed[handB] != 1) {
					adjustment := handLCM - max(posA0, posB0, posA1, posB1)
					posA0 += adjustment
					posB0 += adjustment
					posA1 += adjustment
					posB1 += adjustment
					posA0 %= handLCM
					posA1 %= handLCM
					posB0 %= handLCM
					posB1 %= handLCM
				}

				if cmp.Compare(posA0, posB0) != cmp.Compare(posA1, posB1) {
					if ClockDebug {
						fmt.Println("Meeting of hands", handA, handB, "Using normed positions", posA0, "->", posA1, "/", posB0, "->", posB1)
					}
					ans++
				}
			}
		}
		if ClockDebug {
			fmt.Println("After second", roundNum+1, "have met", ans, "times; hands at ", handPosNew, "record of", prevPos)
		}
		handPos = handPosNew
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
		hands, handDir, handPos, handMove := parsePortal(scanner)
		clockAns := workPortal(hands, handDir, handPos, handMove)
		// fmt.Println("clock", clockNum, "meetings", clockAns)
		ans += clockAns
	}
	return ans
}
func doProblem3(data []byte) any {
	return doProblem2(data)
}
