package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"runtime/pprof"
	"slices"
	"strconv"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q18_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q18_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q18_p3.txt`, "the input for part 3")

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

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	stacks := make(map[int]int)
	prsr := regexp.MustCompile(`x=(\d+) +y=(\d+)`)
	minX := math.MaxInt
	maxX := math.MinInt
	for scanner.Scan() {
		matched := prsr.FindStringSubmatch(scanner.Text())
		xval, _ := strconv.Atoi(matched[1])
		yval, _ := strconv.Atoi(matched[2])
		stacks[xval] = yval
		minX = min(xval, minX)
		maxX = max(xval, maxX)
	}
	heightArray := make([]int, maxX-minX+1)
	heightArrayL := make([]int, maxX-minX+1)
	heightArrayR := make([]int, maxX-minX+1)
	maxL := 0
	for i := 0; i < maxX-minX+1; i++ {
		heightArray[i] = stacks[i-minX]
		heightArrayL[i] = max(maxL, stacks[i-minX])
		maxL = heightArrayL[i]
	}
	maxR := 0
	for i := maxX - minX; i >= 0; i-- {
		heightArrayR[i] = max(maxR, stacks[i-minX])
		maxR = heightArrayR[i]
	}
	ans := 0
	for i := 0; i < maxX-minX+1; i++ {
		ans += min(heightArrayL[i], heightArrayR[i])
	}
	return ans
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	stacks := make(map[int]int)
	prsr := regexp.MustCompile(`x=(\d+) +y=(\d+) +width=(\d+)`)
	minX := math.MaxInt
	maxX := math.MinInt
	for scanner.Scan() {
		matched := prsr.FindStringSubmatch(scanner.Text())
		xval, _ := strconv.Atoi(matched[1])
		yval, _ := strconv.Atoi(matched[2])
		wval, _ := strconv.Atoi(matched[3])
		for w := range wval {
			stacks[xval+w] = yval
			minX = min(xval+w, minX)
			maxX = max(xval+w, maxX)
		}
	}
	heightArray := make([]int, maxX-minX+1)
	heightArrayL := make([]int, maxX-minX+1)
	heightArrayR := make([]int, maxX-minX+1)
	maxL := 0
	for i := 0; i < maxX-minX+1; i++ {
		heightArray[i] = stacks[i-minX]
		heightArrayL[i] = max(maxL, stacks[i-minX])
		maxL = heightArrayL[i]
	}
	maxR := 0
	for i := maxX - minX; i >= 0; i-- {
		heightArrayR[i] = max(maxR, stacks[i-minX])
		maxR = heightArrayR[i]
	}
	ans := 0
	for i := 0; i < maxX-minX+1; i++ {
		ans += min(heightArrayL[i], heightArrayR[i])
		//	fmt.Println("for p2:", i-minX, ans)
	}
	return ans
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	stacks := make(map[int]int)
	widths := make(map[int]int)
	prsr := regexp.MustCompile(`x=(\d+) +y=(\d+) +width=(\d+)`)
	xvals := make([]int, 0)
	for scanner.Scan() {
		matched := prsr.FindStringSubmatch(scanner.Text())
		xval, _ := strconv.Atoi(matched[1])
		yval, _ := strconv.Atoi(matched[2])
		wval, _ := strconv.Atoi(matched[3])
		stacks[xval] = yval
		widths[xval] = wval
		xvals = append(xvals, xval)
	}
	slices.Sort(xvals)
	heightArrayL := make([]int, len(xvals))
	heightArrayR := make([]int, len(xvals))
	toTheLeft := 0
	for i := range len(xvals) {
		heightArrayL[i] = toTheLeft
		toTheLeft = max(toTheLeft, stacks[xvals[i]])
	}
	toTheRight := 0
	for i := len(xvals) - 1; i >= 0; i-- {
		heightArrayR[i] = toTheRight
		toTheRight = max(toTheRight, stacks[xvals[i]])
	}

	ans := 0
	prevXend := xvals[0] - 1
	for i := range len(xvals) {
		// height of region in-between
		ibHeight := min(heightArrayL[i], max(heightArrayR[i], stacks[xvals[i]]))
		ibWidth := xvals[i] - prevXend - 1
		myHeight := max(stacks[xvals[i]], ibHeight)
		ans += ibHeight*ibWidth + myHeight*widths[xvals[i]]
		prevXend = xvals[i] + widths[xvals[i]] - 1
	}
	return ans
}
