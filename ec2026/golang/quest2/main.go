package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q02_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q02_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q02_p3.txt`, "the input for part 3")

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
	ans := 0
	prsr := regexp.MustCompile(`(\d+)([<>])(\d+)`)
	// data = []byte("4<3 2<5 3>1 1>2 1<5 4<1 3>2 4<2 3>5 4<5\n4>1 2>1 3>2 1<5 4>5 2<5 3>1 4>3 2<4 5<3\n5>1 5>2 4>1 2<4 3>2 3<4 1>2 3>5 1<3 4>5\n")
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		notes := prsr.FindAllStringSubmatch(scanner.Text(), -1)
		alloys := make(map[string]bool)
		for _, note := range notes {
			alloys[note[1]] = true
			alloys[note[3]] = true
		}
		alloysList := make([]string, len(alloys))
		for alloy := range alloys {
			alloysList = append(alloysList, alloy)
		}
		slices.SortFunc(alloysList, func(a, b string) int {
			for _, note := range notes {
				if note[1] == a && note[3] == b {
					if note[2] == ">" {
						return -1
					} else {
						return 1
					}
				} else if note[3] == a && note[1] == b {
					if note[2] == ">" {
						return 1
					} else {
						return -1
					}
				}
			}
			return 0
		})
		thisval, _ := strconv.Atoi(strings.Join(alloysList, ""))
		// fmt.Println(thisval)
		ans += thisval
	}

	return ans
}
func doProblem2(data []byte) any {
	ans := 0
	// scanner.Scan()
	// specStr := scanner.Text()
	// if scanner.Err() != nil {
	// 	panic("Scanner error")
	// }
	// specBits := strings.Split(specStr, ",")
	prsr := regexp.MustCompile(`(\d+)([<>])(\d+)`)
	// 	data = []byte(`2<5 3>1 1>2 1<5 4<1 3>2 4<2 3>5 4<5
	// 2>1 3>2 1<5 4>5 2<5 3>1 4>3 2<4 5<3
	// 5>1 5>2 4>1 3>2 3<4 1>2 3>5 1<3 4>5`)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		notes := prsr.FindAllStringSubmatch(scanner.Text(), -1)
		alloys := make(map[string]bool)
		for _, note := range notes {
			alloys[note[1]] = true
			alloys[note[3]] = true
		}
		alloysList := make([]string, 0, len(alloys))
		for alloy := range alloys {
			alloysList = append(alloysList, alloy)
		}
		// fmt.Println(alloysList)
		slices.SortFunc(alloysList, func(a, b string) int {
			arels := make(map[string]int)
			brels := make(map[string]int)
			for _, note := range notes {
				cmpOrder := make([]string, 0, 2)
				dir := -1
				if note[2] == "<" {
					dir = 1
				}
				if note[1] == b || note[3] == a {
					cmpOrder = append(cmpOrder, note[3])
					cmpOrder = append(cmpOrder, note[1])
					dir = -dir
				} else if note[1] == a || note[3] == b {
					cmpOrder = append(cmpOrder, note[1])
					cmpOrder = append(cmpOrder, note[3])
				} else {
					continue
				}
				if cmpOrder[0] == a && cmpOrder[1] == b {
					return dir
				} else if cmpOrder[0] == a {
					arels[cmpOrder[1]] = dir
					if arels[cmpOrder[1]] == brels[cmpOrder[1]] {
						return dir
					}
				} else if cmpOrder[1] == b {
					brels[cmpOrder[0]] = dir
					if arels[cmpOrder[0]] == brels[cmpOrder[0]] {
						return dir
					}
				}
			}
			fmt.Println("Couldn't order ", a, "and", b, "given", scanner.Text())
			return 0
		})
		thisval, _ := strconv.Atoi(strings.Join(alloysList, ""))
		// fmt.Println(thisval)
		ans += thisval
	}
	return ans
}

func doProblem3(data []byte) any {
	ans := 0
	prsr := regexp.MustCompile(`(\d+)([<>])(\d+)`)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		notes := prsr.FindAllStringSubmatch(scanner.Text(), -1)
		alloyGraph := make(map[string]map[string]int)
		alloys := make(map[string]bool)
		for _, note := range notes {
			alloys[note[1]] = true
			alloys[note[3]] = true
			d := 1
			if note[2] == ">" {
				d = -1
			}
			if alloyGraph[note[1]] == nil {
				alloyGraph[note[1]] = make(map[string]int)
			}
			if alloyGraph[note[3]] == nil {
				alloyGraph[note[3]] = make(map[string]int)
			}
			alloyGraph[note[1]][note[3]] = d
			alloyGraph[note[3]][note[1]] = -d
		}
		alloysList := make([]string, 0, len(alloys))
		for alloy := range alloys {
			alloysList = append(alloysList, alloy)
		}
		// fmt.Println(alloysList)
		slices.SortFunc(alloysList, func(a, b string) int {
			if alloyGraph[a][b] != 0 {
				return alloyGraph[a][b]
			}
			reached := make(map[string]int)
			reached[a] = 2
			edge := make(map[string]bool)
			for p, d := range alloyGraph[a] {
				reached[p] = d
				edge[p] = true
			}
			for len(edge) > 0 {
				oldEdge := edge
				edge = make(map[string]bool)
				for alloy := range oldEdge {
					for p, d := range alloyGraph[alloy] {
						if d == reached[alloy] && reached[p] == 0 {
							reached[p] = d
							edge[p] = true
							if p == b {
								return d
							}
						}
					}
				}
			}
			fmt.Println("Couldn't order ", a, "and", b, "given", scanner.Text())
			return 0
		})
		thisval, _ := strconv.Atoi(strings.Join(alloysList, ""))
		ans += thisval
	}
	return ans
}
