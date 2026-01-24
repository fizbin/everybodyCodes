package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func powmod(a, b, m uint64) uint64 {
	a %= m
	if b == 0 {
		return 1
	}
	working := powmod(a, b/2, m)
	if b%2 == 1 {
		return (working * working * a) % m
	}
	return (working * working) % m
}

func eni(a, b, m, firstPow uint64) uint64 {
	remainders := make([]uint64, 0, b-firstPow)
	working := powmod(a, firstPow, m)
	for idx := firstPow + 1; idx <= b; idx++ {
		working *= a
		working %= m
		remainders = append(remainders, working)
	}
	var fullval uint64
	powmul := uint64(1)
	for _, val := range remainders {
		for powmul <= fullval {
			powmul *= 10
		}
		fullval += powmul * val
	}
	// fmt.Printf("Dbg: %d, %d, %d, (%d) -> %d [%v]\n", a, b, m, firstPow, fullval, remainders)
	return fullval
}

func eni3(a, b, m uint64) uint64 {
	var retval uint64
	remainders := make([]uint64, 0, m)
	working := uint64(1)
	for idx := uint64(1); idx <= b; idx++ {
		working *= a
		working %= m
		prevIdx := slices.Index(remainders, working)
		if idx < b && prevIdx > -1 {
			// so (a**idx) % m == (a**(prevIdx+1)) % m
			bigleapSize := uint64(len(remainders) - prevIdx)
			bigleapVal := uint64(0)
			for _, val := range remainders[prevIdx:] {
				bigleapVal += val
			}
			if (b-idx)%bigleapSize == 0 {
				return retval + bigleapVal*((b-idx)/bigleapSize)
			}
			nBigLeaps := (b - idx) / bigleapSize
			retval += bigleapVal * nBigLeaps
			idx += bigleapSize * nBigLeaps
			remainders = remainders[0:prevIdx]
		}
		remainders = append(remainders, working)
		retval += working
	}
	return retval
}

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e1_q01_p1.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		var maxval uint64
		parser := regexp.MustCompile(`A=(\d+) B=(\d+) C=(\d+) X=(\d+) Y=(\d+) Z=(\d+) M=(\d+)`)
		for scanner.Scan() {
			mtch := parser.FindStringSubmatch(scanner.Text())
			if mtch == nil {
				panic("Bad line! " + scanner.Text())
			}
			a, _ := strconv.ParseUint(mtch[1], 10, 64)
			b, _ := strconv.ParseUint(mtch[2], 10, 64)
			c, _ := strconv.ParseUint(mtch[3], 10, 64)
			x, _ := strconv.ParseUint(mtch[4], 10, 64)
			y, _ := strconv.ParseUint(mtch[5], 10, 64)
			z, _ := strconv.ParseUint(mtch[6], 10, 64)
			m, _ := strconv.ParseUint(mtch[7], 10, 64)
			val := eni(a, x, m, 0) + eni(b, y, m, 0) + eni(c, z, m, 0)
			// fmt.Printf("DBG: %d\n", val)
			maxval = max(val, maxval)
		}
		fmt.Println("Part 1:", maxval)
	}

	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e1_q01_p2.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		var maxval uint64
		parser := regexp.MustCompile(`A=(\d+) B=(\d+) C=(\d+) X=(\d+) Y=(\d+) Z=(\d+) M=(\d+)`)
		for scanner.Scan() {
			mtch := parser.FindStringSubmatch(scanner.Text())
			if mtch == nil {
				panic("Bad line! " + scanner.Text())
			}
			a, _ := strconv.ParseUint(mtch[1], 10, 64)
			b, _ := strconv.ParseUint(mtch[2], 10, 64)
			c, _ := strconv.ParseUint(mtch[3], 10, 64)
			x, _ := strconv.ParseUint(mtch[4], 10, 64)
			y, _ := strconv.ParseUint(mtch[5], 10, 64)
			z, _ := strconv.ParseUint(mtch[6], 10, 64)
			m, _ := strconv.ParseUint(mtch[7], 10, 64)
			val := eni(a, x, m, max(x, 5)-5) + eni(b, y, m, max(y, 5)-5) + eni(c, z, m, max(z, 5)-5)
			// fmt.Printf("DBG: %d\n", val)
			maxval = max(val, maxval)
		}
		fmt.Println("Part 2:", maxval)
	}

	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e1_q01_p3.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		var maxval uint64
		parser := regexp.MustCompile(`A=(\d+) B=(\d+) C=(\d+) X=(\d+) Y=(\d+) Z=(\d+) M=(\d+)`)
		for scanner.Scan() {
			mtch := parser.FindStringSubmatch(scanner.Text())
			if mtch == nil {
				panic("Bad line! " + scanner.Text())
			}
			a, _ := strconv.ParseUint(mtch[1], 10, 64)
			b, _ := strconv.ParseUint(mtch[2], 10, 64)
			c, _ := strconv.ParseUint(mtch[3], 10, 64)
			x, _ := strconv.ParseUint(mtch[4], 10, 64)
			y, _ := strconv.ParseUint(mtch[5], 10, 64)
			z, _ := strconv.ParseUint(mtch[6], 10, 64)
			m, _ := strconv.ParseUint(mtch[7], 10, 64)
			val := eni3(a, x, m) + eni3(b, y, m) + eni3(c, z, m)
			// fmt.Printf("DBG: %d\n", val)
			maxval = max(val, maxval)
		}
		fmt.Println("Part 3:", maxval)
	}
}
