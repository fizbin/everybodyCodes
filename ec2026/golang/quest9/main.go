package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q09_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q09_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q09_p3.txt`, "the input for part 3")

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

type treeSpec struct {
	id, growth, limit int
}

func countTreeSecs(spec treeSpec) int {
	componentCount := 0
	branches := []int{0}
	secs := 0
	for ; componentCount < spec.limit; secs++ {
		newBranches := make([]int, 0, len(branches))
		for _, branch := range branches {
			newBranch := branch + 1
			componentCount++
			if newBranch < spec.growth {
				newBranches = append(newBranches, newBranch)
			} else {
				newBranches = append(newBranches, newBranch/2, newBranch/2, 1)
				componentCount++
			}
		}
		branches = newBranches
	}
	return secs
}

func doProblem1(data []byte) any {
	// 	data = []byte(`id=1, growth=3, limit=10
	// id=2, growth=5, limit=20
	// id=3, growth=7, limit=30`)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`id=(\d+), growth=(\d+), limit=(\d+)`)
	var specs []treeSpec
	for scanner.Scan() {
		matched := prsr.FindStringSubmatch(scanner.Text())
		idVal, _ := strconv.Atoi(matched[1])
		growthVal, _ := strconv.Atoi(matched[2])
		limitVal, _ := strconv.Atoi(matched[3])
		specs = append(specs, treeSpec{id: idVal, growth: growthVal, limit: limitVal})
	}
	ans := 0
	for _, spec := range specs {
		ans += countTreeSecs(spec)
	}

	return ans
}

func countTreeSecs2(spec treeSpec) int {
	record := make([]int, 0, 1000)
	record = append(record, 1)

	leafCount := 1
	branches := make(map[int]int)
	branches[0] = 1
	secs := 0
	for leafCount < spec.limit {
		newBranches := make(map[int]int)
		topLen := 0
		for branch, mult := range branches {
			newBranch := branch + 1
			if newBranch < spec.growth {
				newBranches[newBranch] += mult
				topLen = max(topLen, newBranch)
			} else {
				newBranches[newBranch/2] += mult * 2
				topLen = max(topLen, newBranch/2)
				newBranches[1] += mult
				leafCount += mult
			}
		}
		jumpAmount := max(0, spec.growth-1-topLen-5)
		if leafCount >= spec.limit {
			jumpAmount = 0
		}
		branches = make(map[int]int)
		maxnew := 0
		for branch, mult := range newBranches {
			branches[branch+jumpAmount] = mult
			maxnew = max(maxnew, branch+jumpAmount)
		}
		secs += jumpAmount + 1
		// fmt.Println("jumped", jumpAmount, "maxnew", maxnew, "growth", spec.growth)
		// fmt.Println("After ", secs, "seconds, have", leafCount, " and went from ", topLen, "to", maxnew)
	}
	// fmt.Println("tree id", spec.id, "secs", secs)
	return max(1, secs)
}

func doProblem2(data []byte) any {
	// data = []byte(`id=1, growth=3, limit=10
	// id=2, growth=5, limit=20
	// id=3, growth=7, limit=30`)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`id=(\d+), growth=(\d+), limit=(\d+)`)
	var specs []treeSpec
	for scanner.Scan() {
		matched := prsr.FindStringSubmatch(scanner.Text())
		idVal, _ := strconv.Atoi(matched[1])
		growthVal, _ := strconv.Atoi(matched[2])
		limitVal, _ := strconv.Atoi(matched[3])
		specs = append(specs, treeSpec{id: idVal, growth: growthVal, limit: limitVal})
	}
	ans := 0
	for _, spec := range specs {
		ans += countTreeSecs2(spec)
	}

	return ans
}

func doProblem3(data []byte) any {
	return doProblem2(data)
}
