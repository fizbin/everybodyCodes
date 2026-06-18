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

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q16_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q16_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q16_p3.txt`, "the input for part 3")

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

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	dataStr := scanner.Text()
	datas := strings.Split(dataStr, ",")
	// colHeights := make([]int, 91)
	total := 0
	for _, data := range datas {
		val, _ := strconv.Atoi(data)
		for idx := val; idx <= 90; idx += val {
			total += 1
			// colHeights[idx]++
		}
	}
	return total
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	dataStr := scanner.Text()
	datas := strings.Split(dataStr, ",")
	colHeights := make([]int, len(datas)+1)
	for idx, data := range datas {
		val, _ := strconv.Atoi(data)
		colHeights[idx+1] = val
	}
	totProduct := uint64(1)
	for {
		spellFactor := 0
		for ; spellFactor < len(colHeights); spellFactor++ {
			if colHeights[spellFactor] != 0 {
				break
			}
		}
		if spellFactor >= len(colHeights) {
			return totProduct
		}
		mult := colHeights[spellFactor]
		for range mult {
			totProduct *= uint64(spellFactor)
		}
		for x := spellFactor; x < len(colHeights); x += spellFactor {
			colHeights[x] -= mult
		}
	}
}

func blocksNeeded(cols int, spell []int) int {
	retval := 0
	for _, x := range spell {
		retval += cols / x
	}
	return retval
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	dataStr := scanner.Text()
	datas := strings.Split(dataStr, ",")
	colHeights := make([]int, len(datas)+1)
	for idx, data := range datas {
		val, _ := strconv.Atoi(data)
		colHeights[idx+1] = val
	}
	spell := make([]int, 0)
	for {
		spellFactor := 0
		for ; spellFactor < len(colHeights); spellFactor++ {
			if colHeights[spellFactor] != 0 {
				break
			}
		}
		if spellFactor >= len(colHeights) {
			break
		}
		mult := colHeights[spellFactor]
		for range mult {
			spell = append(spell, spellFactor)
		}
		for x := spellFactor; x < len(colHeights); x += spellFactor {
			colHeights[x] -= mult
		}
	}
	blocksTarget := 202520252025000
	low := 1
	high := 100
	for ; blocksNeeded(high, spell) < blocksTarget; high *= 2 {
	}
	for high-low > 1 {
		mid := high/2 + low/2
		if (high%2 == 1) && (low%2 == 1) {
			mid++
		}
		val := blocksNeeded(mid, spell)
		if val < blocksTarget {
			low = mid
		} else if val > blocksTarget {
			high = mid
		} else {
			return mid
		}
	}
	return low
}
