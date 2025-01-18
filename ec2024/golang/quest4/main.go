package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func convolve[E any, F any](in [][]E, def E, comb func(nb [4]E, mid E) F) [][]F {
	retval := make([][]F, 0, len(in))
	for ridx, row := range in {
		newrow := make([]F, 0, len(row))
		for cidx, val := range row {
			up, dn, lt, rt := def, def, def, def
			if ridx > 0 {
				up = in[ridx-1][cidx]
			}
			if cidx > 0 {
				lt = in[ridx][cidx-1]
			}
			if ridx < len(in)-1 {
				dn = in[ridx+1][cidx]
			}
			if cidx < len(row)-1 {
				rt = in[ridx][cidx+1]
			}
			newrow = append(newrow, comb([4]E{up, lt, rt, dn}, val))
		}
		retval = append(retval, newrow)
	}
	return retval
}

func convolveMore[E any, F any](in [][]E, def E, comb func(nb [8]E, mid E) F) [][]F {
	retval := make([][]F, 0, len(in))
	for ridx, row := range in {
		newrow := make([]F, 0, len(row))
		for cidx, val := range row {
			var nbs [8]E
			nidx := 0
			for rowoff := -1; rowoff <= 1; rowoff++ {
				for coloff := -1; coloff <= 1; coloff++ {
					if rowoff == 0 && coloff == 0 {
						continue
					}
					if ridx+rowoff < 0 || ridx+rowoff >= len(in) {
						nbs[nidx] = def
					} else if cidx+coloff < 0 || cidx+coloff >= len(row) {
						nbs[nidx] = def
					} else {
						nbs[nidx] = in[ridx+rowoff][cidx+coloff]
					}
					nidx++
				}
			}
			newrow = append(newrow, comb(nbs, val))
		}
		retval = append(retval, newrow)
	}
	return retval
}

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q04_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		numStrs := strings.Fields(string(data))
		numSum := 0
		minNum := math.MaxUint32
		for _, numStr := range numStrs {
			n, _ := strconv.Atoi(numStr)
			numSum += n
			minNum = min(minNum, n)
		}
		total := numSum - len(numStrs)*minNum

		fmt.Println("Part 1:", total)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q04_p2.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		numStrs := strings.Fields(string(data))
		numSum := 0
		minNum := math.MaxUint32
		for _, numStr := range numStrs {
			n, _ := strconv.Atoi(numStr)
			numSum += n
			minNum = min(minNum, n)
		}
		total := numSum - len(numStrs)*minNum

		fmt.Println("Part 2:", total)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q04_p3.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		numStrs := strings.Fields(string(data))
		nums := make([]int, 0, 1000)
		for _, numStr := range numStrs {
			n, _ := strconv.Atoi(numStr)
			nums = append(nums, n)
		}
		slices.Sort(nums)
		median := nums[len(nums)/2]
		total := 0
		for _, n := range nums {
			if n < median {
				total += median - n
			} else {
				total += n - median
			}
		}

		fmt.Println("Part 3:", total)
	}
}
