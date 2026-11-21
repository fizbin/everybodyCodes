package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"slices"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q19_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q19_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q19_p3.txt`, "the input for part 3")

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

func combineColors(color1, color2 byte) byte {
	if color1 == color2 {
		return color1
	}
	if color1 == 'K' || color2 == 'K' {
		return 'K'
	}
	if color1 > color2 {
		color1, color2 = color2, color1
	}
	if color1 == 'C' && color2 == 'M' {
		return 'B'
	}
	if color1 == 'C' && color2 == 'Y' {
		return 'G'
	}
	return 'R'
}

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid1 := make([][]byte, 0)
	grid2 := make([][]byte, 0)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			break
		}
		grid1 = append(grid1, bytes.Clone(scanner.Bytes()))
	}
	for scanner.Scan() {
		grid2 = append(grid2, bytes.Clone(scanner.Bytes()))
	}
	colorCounts := make(map[byte]int)
	for rowidx, row1 := range grid1 {
		row2 := grid2[rowidx]
		for cidx, ch1 := range row1 {
			ch2 := row2[cidx]
			colorCounts[combineColors(ch1, ch2)]++
		}
	}
	allCounts := make([]int, 0, len(colorCounts))
	for _, x := range colorCounts {
		allCounts = append(allCounts, x)
	}
	slices.Sort(allCounts)
	lac := len(allCounts)
	return allCounts[lac-1] * allCounts[lac-2] * allCounts[lac-3]
}

// call this the first time with all colors in order
func advanceToNextPermutation(colors []byte) bool {
	var idx int
	for idx = len(colors) - 2; idx >= 0; idx-- {
		if colors[idx] < colors[idx+1] {
			break
		}
	}
	if idx < 0 {
		return false
	}
	idx2 := 0
	for idx2 = len(colors) - 1; idx2 > idx; idx2-- {
		if colors[idx2] > colors[idx] {
			break
		}
	}
	colors[idx], colors[idx2] = colors[idx2], colors[idx]
	slices.Reverse(colors[idx+1:])
	return true
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	grid1 := make([][]byte, 0)
	grid2 := make([][]byte, 0)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			break
		}
		grid1 = append(grid1, bytes.Clone(scanner.Bytes()))
	}
	for scanner.Scan() {
		grid2 = append(grid2, bytes.Clone(scanner.Bytes()))
	}

	permutation := []byte("CKMY")
	isFirst := true
	ans := 0
	for isFirst || advanceToNextPermutation(permutation) {
		isFirst = false
		// fmt.Println(string(permutation))
		colorCounts := make(map[byte]int)
		for rowidx, row1 := range grid1 {
			row2 := grid2[rowidx]
			for cidx, ch1 := range row1 {
				ch2 := permutation[row2[cidx]-'1']
				colorCounts[combineColors(ch1, ch2)]++
			}
		}
		allCounts := make([]int, 0, len(colorCounts))
		for _, x := range colorCounts {
			allCounts = append(allCounts, x)
		}
		slices.Sort(allCounts)
		lac := len(allCounts)
		ans = max(ans, allCounts[lac-1]*allCounts[lac-2]*allCounts[lac-3])
	}
	return ans
}

func doProblem3(data []byte) any {
	countStuff := make(map[byte]int)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			break
		}
		for _, ch := range scanner.Bytes() {
			countStuff[ch]++
		}
	}
	totalFreeChoices := countStuff['C'] + countStuff['M'] + countStuff['Y']
	count1 := totalFreeChoices / 3
	count2 := (totalFreeChoices - count1) / 2
	count3 := totalFreeChoices - count2 - count1
	return count1 * count2 * count3
}

const doingPart = 3
