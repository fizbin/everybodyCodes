package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q09_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		denominations := []int{1, 3, 5, 10}
		total := 0
		for scanner.Scan() {
			n, _ := strconv.Atoi(scanner.Text())
			for denomIdx := len(denominations); denomIdx > 0; denomIdx-- {
				denom := denominations[denomIdx-1]
				total += n / denom
				n = n % denom
			}
		}
		fmt.Println("Part 1:", total)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q09_p2.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		denominations := []int{1, 3, 5, 10, 15, 16, 20, 24, 25, 30}
		minBeetles := make([]int, 0, 2000)
		for idx := range 2000 {
			if idx == 0 {
				minBeetles = append(minBeetles, 0)
			} else {
				minval := math.MaxInt
				for _, denom := range denominations {
					if denom > idx {
						continue
					}
					minval = min(minval, 1+minBeetles[idx-denom])
				}
				minBeetles = append(minBeetles, minval)
			}
		}
		total := 0
		for scanner.Scan() {
			n, _ := strconv.Atoi(scanner.Text())
			total += minBeetles[n]
		}
		fmt.Println("Part 2:", total)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q09_p3.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		denominations := []int{1, 3, 5, 10, 15, 16, 20, 24, 25, 30, 37, 38, 49, 50, 74, 75, 100, 101}
		minBeetles := make([]int, 0, 200_000)
		for idx := range 200_000 {
			if idx == 0 {
				minBeetles = append(minBeetles, 0)
			} else {
				minval := math.MaxInt
				for _, denom := range denominations {
					if denom > idx {
						continue
					}
					minval = min(minval, 1+minBeetles[idx-denom])
				}
				minBeetles = append(minBeetles, minval)
			}
		}
		total := 0
		for scanner.Scan() {
			n, _ := strconv.Atoi(scanner.Text())
			minBeetlesForN := math.MaxInt
			for splitA := n/2 - 55; splitA <= n/2; splitA++ {
				splitB := n - splitA
				if splitB-splitA > 100 {
					continue
				}
				minBeetlesForN = min(minBeetlesForN, minBeetles[splitA]+minBeetles[splitB])
			}
			total += minBeetlesForN
		}
		fmt.Println("Part 3:", total)
	}
}
