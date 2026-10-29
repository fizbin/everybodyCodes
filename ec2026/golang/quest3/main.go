package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q03_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q03_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q03_p3.txt`, "the input for part 3")

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
	ans := big.NewInt(0)
	// 	data = []byte(`ISO-3210 7654
	// ISO-8742 2468
	// ISO-1111 1111`)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`ISO-(\d+) +(\d+)`)
	for scanner.Scan() {
		match := prsr.FindStringSubmatch(scanner.Text())
		stdNum := match[1]
		voteNum := match[2]
		for idx := range len(stdNum) {
			voteBit := big.NewInt(int64(stdNum[idx] - '0'))
			voteBit.Exp(voteBit, big.NewInt(int64(len(stdNum)-idx-1)), big.NewInt(0))
			voteBit.Mul(voteBit, big.NewInt(int64(voteNum[idx]-'0')))
			ans.Add(ans, voteBit)
		}
	}
	return ans
}

func doProblem2(data []byte) any {
	ans := big.NewInt(0)
	// 	data = []byte(`ISO-3210 7654
	// ISO-8742 2468
	// ISO-1 1111
	// ISO-13 777
	// ISO-34 140856
	// ISO-789 2465112
	// ISO-2341 945`)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`ISO-(\d+) +(\d+)`)
	for scanner.Scan() {
		match := prsr.FindStringSubmatch(scanner.Text())
		stdNum := match[1]
		voteNum := match[2]
		adj := len(stdNum) - (len(voteNum) % len(stdNum))
		for idx := range len(voteNum) {
			voteBit := big.NewInt(int64(stdNum[(idx+adj)%len(stdNum)] - '0'))
			voteBit.Exp(voteBit, big.NewInt(int64(len(voteNum)-idx-1)), big.NewInt(0))
			voteBit.Mul(voteBit, big.NewInt(int64(voteNum[idx]-'0')))
			ans.Add(ans, voteBit)
		}
	}
	return ans
}

func evalVote(stdNum, voteNum string) *big.Int {
	ans := big.NewInt(0)
	adj := len(stdNum) - (len(voteNum) % len(stdNum))
	for idx := range len(voteNum) {
		voteBit := big.NewInt(int64(stdNum[(idx+adj)%len(stdNum)] - '0'))
		voteBit.Exp(voteBit, big.NewInt(int64(len(voteNum)-idx-1)), big.NewInt(0))
		voteBit.Mul(voteBit, big.NewInt(int64(voteNum[idx]-'0')))
		ans.Add(ans, voteBit)
	}
	return ans
}

func evalIso2Vote(stdNum, voteNum string) *big.Int {
	voteI, _ := strconv.Atoi(voteNum)
	baseVal := evalVote(stdNum, voteNum)
	for tst := 0; tst < voteI; tst++ {
		tstStr := strconv.Itoa(tst)
		tstVal := evalVote(stdNum, tstStr)
		cmpVal := tstVal.Cmp(baseVal)
		nTst := tst + 9 - (tst % 10)
		if cmpVal > 0 {
			tst = nTst
		} else if cmpVal < 0 {
			if nTst < voteI {
				nTstVal := evalVote(stdNum, strconv.Itoa(nTst))
				if nTstVal.Cmp(baseVal) >= 0 {
					return big.NewInt(0)
				} else {
					tst = nTst
				}
			}
		} else {
			return big.NewInt(0)
		}
	}
	return baseVal
}

func doProblem3(data []byte) any {
	ans := big.NewInt(0)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`ISO2-(\d+) +(\d+)`)
	for scanner.Scan() {
		match := prsr.FindStringSubmatch(scanner.Text())
		stdNum := match[1]
		voteNum := match[2]
		ans.Add(ans, evalIso2Vote(stdNum, voteNum))
	}
	return ans
}
