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

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q16_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 1:", doProblem1(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 2 || argsWithoutProg[1] == "" {
			infile = "../input/everybody_codes_e2024_q16_p2.txt"
		} else {
			infile = argsWithoutProg[1]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 2:", doProblem2(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 3 || argsWithoutProg[2] == "" {
			infile = "../input/everybody_codes_e2024_q16_p3.txt"
		} else {
			infile = argsWithoutProg[2]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 3:", doProblem3(data))
	}
}

func doProblem1(data []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	wheelMovesStr := strings.Split(scanner.Text(), ",")
	wheelMoves := make([]int, 0)
	for _, wheelStr := range wheelMovesStr {
		n, _ := strconv.Atoi(wheelStr)
		wheelMoves = append(wheelMoves, n)
	}
	scanner.Scan() // blank line
	wheels := make([][]string, len(wheelMoves))
	for scanner.Scan() {
		t := scanner.Text()
		for wheelIdx := range wheels {
			if 4*wheelIdx+3 > len(t) {
				break
			}
			wheelStr := t[4*wheelIdx : 3+4*wheelIdx]
			if strings.Trim(wheelStr, " ") != "" {
				wheels[wheelIdx] = append(wheels[wheelIdx], t[4*wheelIdx:3+4*wheelIdx])
			}
		}
	}
	res := make([]string, 0)
	for idx, catFaces := range wheels {
		res = append(res, catFaces[(wheelMoves[idx]*100)%len(catFaces)])
	}
	return strings.Join(res, " ")
}

func gcd(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	if a > b {
		return gcd(b, a)
	}
	return gcd(b%a, a)
}

func lcm(a, b int) int {
	return a * (b / gcd(a, b))
}

// Like a % b but with python/perl semantics
func mod(a, b int) int {
	res := a % b
	if (res == 0) || ((res < 0) == (b < 0)) {
		return res
	}
	return res + b
}

func doProblem2(data []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	wheelMovesStr := strings.Split(scanner.Text(), ",")
	wheelMoves := make([]int, 0)
	for _, wheelStr := range wheelMovesStr {
		n, _ := strconv.Atoi(wheelStr)
		wheelMoves = append(wheelMoves, n)
	}
	scanner.Scan() // blank line
	wheels := make([][]string, len(wheelMoves))
	for scanner.Scan() {
		t := scanner.Text()
		for wheelIdx := range wheels {
			if 4*wheelIdx+3 > len(t) {
				break
			}
			wheelStr := t[4*wheelIdx : 3+4*wheelIdx]
			if strings.Trim(wheelStr, " ") != "" {
				wheels[wheelIdx] = append(wheels[wheelIdx], t[4*wheelIdx:3+4*wheelIdx])
			}
		}
	}
	wheelLCM := 1
	for _, wheel := range wheels {
		wheelLCM = lcm(wheelLCM, len(wheel))
	}
	scorePull := func(pullN int) int {
		has := make(map[byte]int)
		for idx, catFaces := range wheels {
			face := catFaces[(wheelMoves[idx]*pullN)%len(catFaces)]
			has[face[0]] += 1
			has[face[2]] += 1
		}
		tot := 0
		for _, v := range has {
			if v >= 3 {
				tot += v - 2
			}
		}
		return tot
	}
	totCoins := 0
	for pull := 1; pull <= wheelLCM; pull++ {
		totCoins += scorePull(pull)
	}
	lcmCycles := 202420242024 / wheelLCM
	lcmCycleEnds := lcmCycles * wheelLCM
	totCoins *= lcmCycles
	for pull := lcmCycleEnds + 1; pull <= 202420242024; pull++ {
		totCoins += scorePull(pull)
	}
	return totCoins
}

func doProblem3(data []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	wheelMovesStr := strings.Split(scanner.Text(), ",")
	wheelMoves := make([]int, 0)
	for _, wheelStr := range wheelMovesStr {
		n, _ := strconv.Atoi(wheelStr)
		wheelMoves = append(wheelMoves, n)
	}
	scanner.Scan() // blank line
	wheels := make([][]string, len(wheelMoves))
	for scanner.Scan() {
		t := scanner.Text()
		for wheelIdx := range wheels {
			if 4*wheelIdx+3 > len(t) {
				break
			}
			wheelStr := t[4*wheelIdx : 3+4*wheelIdx]
			if strings.Trim(wheelStr, " ") != "" {
				wheels[wheelIdx] = append(wheels[wheelIdx], t[4*wheelIdx:3+4*wheelIdx])
			}
		}
	}
	scorePull := func(pullN int, offset int) int {
		has := make(map[byte]int)
		for idx, catFaces := range wheels {
			face := catFaces[mod(wheelMoves[idx]*pullN+offset, len(catFaces))]
			has[face[0]] += 1
			has[face[2]] += 1
		}
		tot := 0
		for _, v := range has {
			if v >= 3 {
				tot += v - 2
			}
		}
		return tot
	}
	computeAnswer := func(combine func(int, int) int) int {
		// we'll have current pull go from 256 to 0

		// answers will be the number of coins that can be gotten *after*
		// "current pull" if current pull ended with the total offset being
		// [idx] at the very end, we'll use answers[0]
		answers := make(map[int]int)
		// after pull 256, you can have offset from -256 to 256, but can't get any more coins
		for offset := -256; offset <= 256; offset++ {
			answers[offset] = 0
		}
		for pull := 256; pull > 0; pull-- {
			thisPullScore := make(map[int]int)
			for offset := -pull; offset <= pull; offset++ {
				thisPullScore[offset] = scorePull(pull, offset)
			}
			// newAnswers will be the number of coins that can be gotten *after* round "pull-1"
			newAnswers := make(map[int]int)
			for offset := -pull + 1; offset <= pull-1; offset++ {
				candidate1 := thisPullScore[offset-1] + answers[offset-1]
				candidate2 := thisPullScore[offset] + answers[offset]
				candidate3 := thisPullScore[offset+1] + answers[offset+1]
				newAnswers[offset] = combine(candidate1, combine(candidate2, candidate3))
			}
			answers = newAnswers // answers is now correct for *after* round "pull-1"
		}
		return answers[0]
	}

	return strconv.Itoa(computeAnswer(func(a, b int) int { return max(a, b) })) + " " +
		strconv.Itoa(computeAnswer(func(a, b int) int { return min(a, b) }))
}
