package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q11_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		rules := make(map[string]map[string]int)
		for scanner.Scan() {
			splitLine := strings.Split(scanner.Text(), ":")
			res := strings.Split(splitLine[1], ",")
			rules[splitLine[0]] = make(map[string]int)
			for _, child := range res {
				rules[splitLine[0]][child] += 1
			}
		}
		termites := map[string]int{"A": 1}
		for range 4 {
			ntermites := make(map[string]int)
			for ttyp, num := range termites {
				for otyp, onum := range rules[ttyp] {
					ntermites[otyp] += num * onum
				}
			}
			termites = ntermites
		}
		total := 0
		for _, v := range termites {
			total += v
		}
		fmt.Println("Part 1:", total)
	}
	{
		var infile string
		if len(argsWithoutProg) < 2 || argsWithoutProg[1] == "" {
			infile = "../input/everybody_codes_e2024_q11_p2.txt"
		} else {
			infile = argsWithoutProg[1]
		}
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		rules := make(map[string]map[string]int)
		for scanner.Scan() {
			splitLine := strings.Split(scanner.Text(), ":")
			res := strings.Split(splitLine[1], ",")
			rules[splitLine[0]] = make(map[string]int)
			for _, child := range res {
				rules[splitLine[0]][child] += 1
			}
		}
		termites := map[string]int{"Z": 1}
		for range 10 {
			ntermites := make(map[string]int)
			for ttyp, num := range termites {
				for otyp, onum := range rules[ttyp] {
					ntermites[otyp] += num * onum
				}
			}
			termites = ntermites
		}
		total := 0
		for _, v := range termites {
			total += v
		}
		fmt.Println("Part 2:", total)
	}
	{
		var infile string
		if len(argsWithoutProg) < 3 || argsWithoutProg[2] == "" {
			infile = "../input/everybody_codes_e2024_q11_p3.txt"
		} else {
			infile = argsWithoutProg[2]
		}
		data, _ := os.ReadFile(infile)
		// data = []byte("A:B,C\nB:C,A,A\nC:A")
		scanner := bufio.NewScanner(bytes.NewReader(data))
		rules := make(map[string]map[string]int)
		for scanner.Scan() {
			splitLine := strings.Split(scanner.Text(), ":")
			res := strings.Split(splitLine[1], ",")
			rules[splitLine[0]] = make(map[string]int)
			for _, child := range res {
				rules[splitLine[0]][child] += 1
			}
		}
		termiteTotal := func(intype string) int {
			termites := map[string]int{intype: 1}
			for range 20 {
				ntermites := make(map[string]int)
				for ttyp, num := range termites {
					for otyp, onum := range rules[ttyp] {
						ntermites[otyp] += num * onum
					}
				}
				termites = ntermites
			}
			total := 0
			for _, v := range termites {
				total += v
			}
			return total
		}
		minTot := math.MaxInt
		maxTot := 0
		for ttyp := range rules {
			v := termiteTotal(ttyp)
			// fmt.Println("DBG", ttyp, v)
			minTot = min(minTot, v)
			maxTot = max(maxTot, v)
		}
		// fmt.Println(minTot, maxTot)
		fmt.Println("Part 3:", maxTot-minTot)
	}
}
