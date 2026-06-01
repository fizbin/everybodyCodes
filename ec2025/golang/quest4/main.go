package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q04_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q04_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q04_p3.txt`, "the input for part 3")

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
		// sizeStr := scanner.Text()
		if scanner.Err() != nil {
			panic("Scanner error")
		}
		// did p1 "by hand" basically
	}
	{
		infile := *input2
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		gears := make([]*big.Int, 0)
		for scanner.Scan() {
			asInt64, _ := strconv.ParseInt(scanner.Text(), 10, 64)
			gears = append(gears, big.NewInt(asInt64))
		}
		target := big.NewInt(10000000000000)
		working := big.NewInt(0)
		working.Mul(target, gears[len(gears)-1])
		working.DivMod(working, gears[0], target)
		var ans uint64
		if target.Cmp(&big.Int{}) > 0 {
			ans = working.Uint64() + 1
		} else {
			ans = working.Uint64()
		}
		// sizeStr := scanner.Text()
		if scanner.Err() != nil {
			panic("Scanner error")
		}
		fmt.Println("p2", ans)

	}
	{
		infile := *input3
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		ans := big.NewRat(100, 1)
		scanner.Scan()
		gear1, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		ans.Mul(ans, big.NewRat(gear1, 1))
		parser := regexp.MustCompile(` *(\d+) *\| *(\d+)`)
		for scanner.Scan() {
			lineBits := parser.FindStringSubmatch(scanner.Text())
			if len(lineBits) == 0 {
				final, _ := strconv.ParseInt(scanner.Text(), 10, 64)
				ans.Mul(ans, big.NewRat(1, final))
				break
			}
			denom, _ := strconv.ParseInt(lineBits[1], 10, 64)
			numer, _ := strconv.ParseInt(lineBits[2], 10, 64)
			ans.Mul(ans, big.NewRat(numer, denom))
		}
		if scanner.Err() != nil {
			panic("Scanner error")
		}
		numerI := ans.Num()
		denomI := ans.Denom()
		numerI.Div(numerI, denomI)
		fmt.Println("p3", numerI)
	}
}
