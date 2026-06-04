package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q07_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q07_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q07_p3.txt`, "the input for part 3")

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
		nameListStr := scanner.Text()
		names := strings.Split(nameListStr, ",")
		scanner.Scan() // blank line
		parser := regexp.MustCompile(`(\S) *> ([a-zA-Z,]+)`)
		ruleMap := make(map[byte][]byte)
		for scanner.Scan() {
			ruleBits := parser.FindStringSubmatch(scanner.Text())
			res := make([]byte, 0, len(ruleBits[2]))
			for _, following := range strings.Split(ruleBits[2], ",") {
				res = append(res, following[0])
			}
			ruleMap[ruleBits[1][0]] = res
		}
		if scanner.Err() != nil {
			panic("scanner error")
		}
		for _, name := range names {
			prevChar := name[0]
			valid := true
			for _, ch := range ([]byte(name))[1:] {
				if slices.Contains(ruleMap[prevChar], ch) {
					prevChar = ch
					continue
				} else {
					// fmt.Println(name, "was rejected at character", ch, " with rule ", (string)([]byte{prevChar}), ">", (string)(ruleMap[prevChar]))
					valid = false
					break
				}
			}
			if valid {
				fmt.Println("p1:", name)
			}
		}
		// knightList := scanner.Text()
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
		scanner.Scan()
		nameListStr := scanner.Text()
		names := strings.Split(nameListStr, ",")
		scanner.Scan() // blank line
		parser := regexp.MustCompile(`(\S) *> ([a-zA-Z,]+)`)
		ruleMap := make(map[byte][]byte)
		for scanner.Scan() {
			ruleBits := parser.FindStringSubmatch(scanner.Text())
			res := make([]byte, 0, len(ruleBits[2]))
			for _, following := range strings.Split(ruleBits[2], ",") {
				res = append(res, following[0])
			}
			ruleMap[ruleBits[1][0]] = res
		}
		if scanner.Err() != nil {
			panic("scanner error")
		}
		acceptSum := 0
		for nameIdx, name := range names {
			prevChar := name[0]
			valid := true
			for _, ch := range ([]byte(name))[1:] {
				if slices.Contains(ruleMap[prevChar], ch) {
					prevChar = ch
					continue
				} else {
					valid = false
					break
				}
			}
			if valid {
				acceptSum += nameIdx + 1
			}
		}
		fmt.Println("p2:", acceptSum)
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
		scanner.Scan()
		nameListStr := scanner.Text()
		names := strings.Split(nameListStr, ",")
		scanner.Scan() // blank line
		parser := regexp.MustCompile(`(\S) *> ([a-zA-Z,]+)`)
		ruleMap := make(map[byte][]byte)
		for scanner.Scan() {
			ruleBits := parser.FindStringSubmatch(scanner.Text())
			res := make([]byte, 0, len(ruleBits[2]))
			for _, following := range strings.Split(ruleBits[2], ",") {
				res = append(res, following[0])
			}
			ruleMap[ruleBits[1][0]] = res
		}
		if scanner.Err() != nil {
			panic("scanner error")
		}
		memo := make(map[memoParam]uint64)
		var acceptSum uint64
		slices.SortFunc(names, func(a, b string) int {
			return len(a) - len(b)
		})
		seenNames := make(map[string]bool)
		for _, name := range names {
			prevChar := name[0]
			valid := true
			for idx, ch := range ([]byte(name))[1:] {
				if seenNames[name[0:idx+1]] {
					valid = false
					break
				}
				if slices.Contains(ruleMap[prevChar], ch) {
					prevChar = ch
					continue
				} else {
					valid = false
					break
				}
			}
			if valid {
				seenNames[name] = true
				acceptSum += namesPossible(len(name), prevChar, ruleMap, memo)
			}
		}
		fmt.Println("p3:", acceptSum)
	}
}

type memoParam struct {
	lengthSoFar int
	prevChar    byte
}

func namesPossible(lengthSoFar int, prevChar byte, ruleMap map[byte][]byte, memo map[memoParam]uint64) uint64 {
	if computed, ok := memo[memoParam{lengthSoFar, prevChar}]; ok {
		return computed
	}
	if lengthSoFar == 11 {
		return 1
	}
	if lengthSoFar > 11 {
		return 0
	}
	var retval uint64
	if lengthSoFar >= 7 {
		retval = 1
	}
	for _, nextch := range ruleMap[prevChar] {
		retval += namesPossible(lengthSoFar+1, nextch, ruleMap, memo)
	}
	memo[memoParam{lengthSoFar, prevChar}] = retval
	return retval
}
