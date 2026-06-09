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

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q09_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q09_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q09_p3.txt`, "the input for part 3")

func main() {
	flag.Parse()
	{
		infile := *input1
		data, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		fmt.Println("p1:", doProblem1(data))
	}
	{
		infile := *input2
		data, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		fmt.Println("p2:", doProblem2(data))
	}
	{
		infile := *input3
		data, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		fmt.Println("p3:", doProblem3(data))
	}
}

var /* const */ lineParser = regexp.MustCompile(`(\d+):([ACTG]+)`)

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scales := make([]string, 0)
	for scanner.Scan() {
		parts := lineParser.FindStringSubmatch(scanner.Text())
		scales = append(scales, parts[2])
	}
	if len(scales) != 3 {
		panic("Didn't get exactly 3 scales!")
	}
	for childIdx := range 3 {
		var match1, match2 int
		for idx, base := range []byte(scales[childIdx]) {
			if (scales[(childIdx+1)%3][idx] != base) && (scales[(childIdx+2)%3][idx] != base) {
				match1 = 0
				match2 = 0
				break
			}
			if scales[(childIdx+1)%3][idx] == base {
				match1++
			}
			if scales[(childIdx+2)%3][idx] == base {
				match2++
			}
		}
		if (match1 > 0) && (match2 > 0) {
			return match1 * match2
		}
	}
	return 0
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scales := make([]string, 0)
	for scanner.Scan() {
		parts := lineParser.FindStringSubmatch(scanner.Text())
		scales = append(scales, parts[2])
	}
	simSum := 0
	for childIdx, child := range scales {
		var match1, match2 int
		for idx1 := 0; idx1 < len(scales)-1; idx1++ {
			if childIdx == idx1 {
				continue
			}
			for idx2 := idx1 + 1; idx2 < len(scales); idx2++ {
				if childIdx == idx2 {
					continue
				}
				match1 = 0
				match2 = 0
				for idx, base := range []byte(child) {
					if (scales[idx1][idx] != base) && (scales[idx2][idx] != base) {
						match1 = 0
						match2 = 0
						break
					}
					if scales[idx1][idx] == base {
						match1++
					}
					if scales[idx2][idx] == base {
						match2++
					}
				}
				if (match1 > 0) && (match2 > 0) {
					break
				}
			}
			if (match1 > 0) && (match2 > 0) {
				break
			}
		}
		if false {
			fmt.Println("p2 debug:", childIdx+1, match1, match2, match1*match2)
		}
		simSum += match1 * match2
	}
	return simSum
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scales := make(map[int]string, 0)
	for scanner.Scan() {
		parts := lineParser.FindStringSubmatch(scanner.Text())
		scaleNum, _ := strconv.Atoi(parts[1])
		scales[scaleNum] = parts[2]
	}
	connections := make([][3]int, 0)
	for childIdx, child := range scales {
		var match1, match2 int
		var idx1, idx2 int
		for idx1 = 1; idx1 <= len(scales)-1; idx1++ {
			if childIdx == idx1 {
				continue
			}
			for idx2 = idx1 + 1; idx2 <= len(scales); idx2++ {
				if childIdx == idx2 {
					continue
				}
				match1 = 0
				match2 = 0
				for idx, base := range []byte(child) {
					if (scales[idx1][idx] != base) && (scales[idx2][idx] != base) {
						match1 = 0
						match2 = 0
						break
					}
					if scales[idx1][idx] == base {
						match1++
					}
					if scales[idx2][idx] == base {
						match2++
					}
				}
				if (match1 > 0) && (match2 > 0) {
					break
				}
			}
			if (match1 > 0) && (match2 > 0) {
				break
			}
		}
		if (match1 > 0) && (match2 > 0) {
			connections = append(connections, [3]int{idx1, idx2, childIdx})
		}
	}
	familyNum := make(map[int]int)
	for key := range scales {
		familyNum[key] = key
	}
	done := false
	for !done {
		done = true
		for _, connection := range connections {
			familyMin := min(familyNum[connection[0]], familyNum[connection[1]], familyNum[connection[2]])
			if familyNum[connection[0]]+familyNum[connection[1]]+familyNum[connection[2]] != 3*familyMin {
				done = false
				familyNum[connection[0]] = familyMin
				familyNum[connection[1]] = familyMin
				familyNum[connection[2]] = familyMin
			}
		}
	}
	famSize := make(map[int]int)
	famScaleTot := make(map[int]int)
	maxFam := 0
	maxFamSize := 0
	for sNum, famNum := range familyNum {
		famSize[famNum]++
		famScaleTot[famNum] += sNum

		if famSize[famNum] > maxFamSize {
			maxFam = famNum
			maxFamSize = famSize[famNum]
		}
	}

	return famScaleTot[maxFam]
}
