package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"maps"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q07_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q07_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q07_p3.txt`, "the input for part 3")

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

type bolt struct {
	dest  int
	power int
}

func doProblem1Guts(data []byte, initialPulse int) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`id=(\d+), energy=(\d+), connections=([\d,]+)`)
	mostConnected := 0
	mostConnections := 0
	energies := make(map[int]int)
	network := make(map[int][]int)
	for scanner.Scan() {
		matched := prsr.FindStringSubmatch(scanner.Text())
		id, _ := strconv.Atoi(matched[1])
		energy, _ := strconv.Atoi(matched[2])
		connectionStrs := strings.SplitSeq(matched[3], ",")
		var outgoing []int
		for cStr := range connectionStrs {
			cVal, _ := strconv.Atoi(cStr)
			outgoing = append(outgoing, cVal)
		}
		network[id] = outgoing
		energies[id] = energy
		if len(outgoing) > mostConnections {
			mostConnections = len(outgoing)
			mostConnected = id
		}
	}
	boltMap := make(map[bolt]int)
	boltMap[bolt{dest: mostConnected, power: initialPulse}] = 1
	for len(boltMap) > 0 {
		oldBolts := boltMap
		boltMap = make(map[bolt]int)
		for b, mult := range oldBolts {
			energies[b.dest] += mult * b.power
			if b.power > 1 {
				for _, d := range network[b.dest] {
					boltMap[bolt{dest: d, power: b.power - 1}] += mult
				}
			}
		}
	}
	total := 0
	for _, val := range energies {
		total += val
	}
	return total
}
func doProblem1(data []byte) any {
	return doProblem1Guts(data, 5)
}
func doProblem2(data []byte) any {
	return doProblem1Guts(data, 10)
}
func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`id=(\d+), energy=(\d+), connections=([\d,]+)`)
	energies := make(map[int]int)
	network := make(map[int][]int)
	for scanner.Scan() {
		matched := prsr.FindStringSubmatch(scanner.Text())
		id, _ := strconv.Atoi(matched[1])
		energy, _ := strconv.Atoi(matched[2])
		connectionStrs := strings.SplitSeq(matched[3], ",")
		var outgoing []int
		for cStr := range connectionStrs {
			cVal, _ := strconv.Atoi(cStr)
			outgoing = append(outgoing, cVal)
		}
		network[id] = outgoing
		energies[id] = energy
	}
	ans := 0
	for cloudId := range network {
		energiesCopy := maps.Clone(energies)
		ans = max(ans, handleTopShot(network, energiesCopy, cloudId, 10))
	}
	return ans
}

func handleTopShot(network map[int][]int, energies map[int]int, dest, power int) int {
	bolts := []bolt{{dest: dest, power: power}}
	handleRound(network, energies, bolts)
	total := 0
	for _, val := range energies {
		total += val
	}
	return total
}

func handleRound(network map[int][]int, energies map[int]int, bolts []bolt) {
	boltMap := make(map[bolt]int)
	for _, b := range bolts {
		boltMap[b] += 1
	}
	for len(boltMap) > 0 {
		oldBolts := boltMap
		boltMap = make(map[bolt]int)
		for b, mult := range oldBolts {
			if energies[b.dest] == 0 {
				continue
			}
			energies[b.dest] += b.power * mult
			if b.power > 1 {
				for _, d := range network[b.dest] {
					boltMap[bolt{dest: d, power: b.power - 1}] += mult
				}
			}
		}

		for tgt, val := range energies {
			if val >= 1000000 {
				energies[tgt] = 0
				for _, d := range network[tgt] {
					boltMap[bolt{dest: d, power: 10}] += 1
				}
			}
		}
	}
}
