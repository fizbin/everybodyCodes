package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q06_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q06_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q06_p3.txt`, "the input for part 3")

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
	myMap := make(map[string][]string)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`(\S+):(\S+)`)
	for scanner.Scan() {
		match := prsr.FindStringSubmatch(scanner.Text())
		target := match[1]
		sources := strings.Split(match[2], ",")
		myMap[target] = sources
	}
	nullSrc := make(map[string]bool)
	boundary := []string{"NULL"}
	for len(boundary) > 0 {
		oldBoundary := boundary
		boundary = make([]string, 0)
		for _, start := range oldBoundary {
			startSrcs := myMap[start]
			for _, s := range startSrcs {
				if !nullSrc[s] {
					nullSrc[s] = true
					boundary = append(boundary, s)
				}
			}
		}
	}
	return len(nullSrc)
}

func doProblem2(data []byte) any {
	revMap := make(map[string][]string)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`(\S+):(\S+)`)
	for scanner.Scan() {
		match := prsr.FindStringSubmatch(scanner.Text())
		target := match[1]
		sources := strings.Split(match[2], ",")
		for _, src := range sources {
			revMap[src] = append(revMap[src], target)
		}
	}
	nullSrc := make(map[string]bool)
	boundary := []string{"NULL"}
	for len(boundary) > 0 {
		oldBoundary := boundary
		boundary = make([]string, 0)
		for _, start := range oldBoundary {
			startSrcs := revMap[start]
			for _, s := range startSrcs {
				if !nullSrc[s] {
					nullSrc[s] = true
					boundary = append(boundary, s)
				}
			}
		}
	}
	return len(nullSrc)
}

func doProblem3(data []byte) any {
	revMap := make(map[string][]string)
	allTargets := make(map[string]bool)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`(\S+):(\S+)`)
	for scanner.Scan() {
		match := prsr.FindStringSubmatch(scanner.Text())
		target := match[1]
		allTargets[target] = true
		sources := strings.Split(match[2], ",")
		for _, src := range sources {
			revMap[src] = append(revMap[src], target)
		}
	}
	findContaminated := func(contaminant string) map[string]bool {
		nullSrc := make(map[string]bool)
		boundary := []string{contaminant}
		for len(boundary) > 0 {
			oldBoundary := boundary
			boundary = make([]string, 0)
			for _, start := range oldBoundary {
				startSrcs := revMap[start]
				for _, s := range startSrcs {
					if !nullSrc[s] {
						nullSrc[s] = true
						boundary = append(boundary, s)
					}
				}
			}
		}
		return nullSrc
	}
	contaminated1 := findContaminated("NULLv1")
	contaminated2 := findContaminated("NULLv2")
	contaminated3 := findContaminated("NULLv3")

	atLeastTwo := make(map[string]bool)
	for tgt := range allTargets {
		v := 0
		if contaminated1[tgt] {
			v++
		}
		if contaminated2[tgt] {
			v++
		}
		if contaminated3[tgt] {
			v++
		}
		if v >= 2 {
			atLeastTwo[tgt] = true
		}
	}
	return len(atLeastTwo)
}
