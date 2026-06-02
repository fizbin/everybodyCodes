package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q05_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q05_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q05_p3.txt`, "the input for part 3")

type fishRow struct {
	left, spine, right *int
}

type sword struct {
	id               uint64
	fish             []fishRow
	spineVal         uint64
	comparisonValues []uint64
}

func main() {
	flag.Parse()
	{
		infile := *input1
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		specStr := scanner.Text()
		if scanner.Err() != nil {
			panic("Scanner error")
		}
		specBits := strings.Split(specStr, ":")
		numsStr := specBits[1]
		fish := parseFish(numsStr)
		fmt.Print("p1: ")
		for _, frow := range fish {
			fmt.Print(*frow.spine)
		}
		fmt.Println()
	}
	{
		infile := *input2
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		maxFishVal := uint64(0)
		minFishVal := uint64(math.MaxUint64)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			specStr := scanner.Text()
			specBits := strings.Split(specStr, ":")
			numsStr := specBits[1]
			fish := parseFish(numsStr)
			fishVal := uint64(0)
			for _, fishRow := range fish {
				fishVal *= 10
				fishVal += uint64(*fishRow.spine)
			}
			maxFishVal = max(fishVal, maxFishVal)
			minFishVal = min(fishVal, minFishVal)
		}
		if scanner.Err() != nil {
			panic("Scanner error")
		}
		fmt.Println("p2:", maxFishVal-minFishVal)
	}
	{
		infile := *input3
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		swords := make([]sword, 0)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			specStr := scanner.Text()
			specBits := strings.Split(specStr, ":")
			numsStr := specBits[1]
			fish := parseFish(numsStr)
			fishVal := uint64(0)
			cmpvals := make([]uint64, 0)
			for _, fishRow := range fish {
				fishVal *= 10
				fishVal += uint64(*fishRow.spine)

				rval := uint64(*fishRow.spine)
				if fishRow.left != nil {
					rval += uint64(10 * (*fishRow.left))
				}
				if fishRow.right != nil {
					rval *= 10
					rval += uint64(*fishRow.right)
				}
				cmpvals = append(cmpvals, rval)
			}
			idVal, _ := strconv.ParseUint(specBits[0], 10, 64)
			swords = append(swords, sword{id: idVal, fish: fish, spineVal: fishVal, comparisonValues: cmpvals})
		}
		if scanner.Err() != nil {
			panic("Scanner error")
		}
		slices.SortFunc(swords, func(a, b sword) int {
			if a.spineVal > b.spineVal {
				return -1
			}
			if b.spineVal > a.spineVal {
				return 1
			}
			for idx := range a.comparisonValues {
				if a.comparisonValues[idx] > b.comparisonValues[idx] {
					return -1
				}
				if b.comparisonValues[idx] > a.comparisonValues[idx] {
					return 1
				}
			}
			// last resort, by id
			return int(b.id) - int(a.id)
		})
		checksum := uint64(0)
		for idx, sword := range swords {
			checksum += uint64(idx+1) * sword.id
		}
		fmt.Println("p3:", checksum)
	}
}

func parseFish(numsStr string) []fishRow {
	nums := make([]int, 0)
	for _, numStr := range strings.Split(numsStr, ",") {
		num, _ := strconv.Atoi(numStr)
		nums = append(nums, num)
	}
	fish := make([]fishRow, 0)
	for idx := range nums {
		val := nums[idx]
		rowIdx := 0
		for ; rowIdx < len(fish); rowIdx++ {
			if (val < *fish[rowIdx].spine) && fish[rowIdx].left == nil {
				fish[rowIdx].left = &val
				break
			}
			if (val > *fish[rowIdx].spine) && fish[rowIdx].right == nil {
				fish[rowIdx].right = &val
				break
			}
		}
		if rowIdx == len(fish) {
			fish = append(fish, fishRow{spine: &val})
		}
	}
	return fish
}
