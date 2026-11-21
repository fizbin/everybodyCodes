package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"strconv"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q20_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q20_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q20_p3.txt`, "the input for part 3")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
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

var figSegs map[byte]int

var prsr *regexp.Regexp

func init() {
	figSegs = make(map[byte]int)
	for idx, val := range []int{6, 2, 5, 5, 4, 5, 6, 3, 7, 6} {
		figSegs['0'+byte(idx)] = val
	}
	figSegs['+'] = 2
	figSegs['-'] = 1
	figSegs['='] = 2
	figSegs['A'] = 6
	prsr = regexp.MustCompile(`([\dA]+) +([-+]) +([\dA]+) *= *([\dA]+) +([-+]) +([\dA]+)`)
}

func evalLine(line []byte, base int) bool {
	matched := prsr.FindSubmatch(line)
	val1, _ := strconv.ParseInt(string(matched[1]), base, 32)
	val2, _ := strconv.ParseInt(string(matched[3]), base, 32)
	val3, _ := strconv.ParseInt(string(matched[4]), base, 32)
	val4, _ := strconv.ParseInt(string(matched[6]), base, 32)
	var lft, rgt int64
	switch matched[2][0] {
	case '+':
		lft = val1 + val2
	case '-':
		lft = val1 - val2
	}
	switch matched[5][0] {
	case '+':
		rgt = val3 + val4
	case '-':
		rgt = val3 - val4
	}
	// if lft == rgt {
	// 	fmt.Println("Evaluating", string(line), "yields", lft, "=", rgt)
	// }
	return lft == rgt
}

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	ans := 0
	for scanner.Scan() {
		if evalLine(scanner.Bytes(), 10) {
			segs := 0
			for _, ch := range scanner.Bytes() {
				segs += figSegs[ch]
			}
			ans += segs
		}
	}
	return ans
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	ans := 0
	for scanner.Scan() {
		if evalLine(scanner.Bytes(), 11) {
			segs := 0
			for _, ch := range scanner.Bytes() {
				segs += figSegs[ch]
			}
			ans += segs
		}
	}

	return strconv.FormatInt(int64(ans), 11)
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	movesBoth, movesUp, movesDown := make(map[byte][]byte), make(map[byte][]byte), make(map[byte][]byte)
	movesBoth['0'] = []byte{'6', '9', 'A'}
	movesBoth['6'] = []byte{'0', '9', 'A'}
	movesBoth['9'] = []byte{'6', '0', 'A'}
	movesBoth['A'] = []byte{'6', '9', '0'}
	movesBoth['3'] = []byte{'2', '5'}
	movesBoth['2'] = []byte{'3'}
	movesBoth['5'] = []byte{'3'}
	movesUp['1'] = []byte{'7'}
	movesUp['3'] = []byte{'9'}
	movesUp['5'] = []byte{'9', '6'}
	movesUp['0'] = []byte{'8'}
	movesUp['9'] = []byte{'8'}
	movesUp['6'] = []byte{'8'}
	movesUp['A'] = []byte{'8'}
	movesUp['-'] = []byte{'+'}
	movesDown['8'] = []byte{'0', '9', '6', 'A'}
	movesDown['+'] = []byte{'-'}
	movesDown['6'] = []byte{'5'}
	movesDown['9'] = []byte{'5', '3'}
	movesDown['7'] = []byte{'1'}
	ans := 0
	for scanner.Scan() {
		line := bytes.Clone(scanner.Bytes())
		accept := false
		if evalLine(line, 11) {
			accept = true
		}
		for i := 0; i < len(line) && !accept; i++ {
			oldVal := line[i]
			for _, repl := range movesBoth[oldVal] {
				if figSegs[oldVal] != figSegs[repl] {
					panic("Bad1")
				}
				line[i] = repl
				if evalLine(line, 11) {
					accept = true
				}
			}
			line[i] = oldVal
		}
		for i := 0; i < len(line) && !accept; i++ {
			oldVal1 := line[i]
			for _, repl1 := range movesDown[oldVal1] {
				if figSegs[oldVal1]-1 != figSegs[repl1] {
					panic("Bad2")
				}
				line[i] = repl1
				for j := 0; j < len(line) && !accept; j++ {
					if i == j {
						continue
					}
					oldVal2 := line[j]
					for _, repl2 := range movesUp[oldVal2] {
						if figSegs[oldVal2]+1 != figSegs[repl2] {
							panic("Bad3")
						}
						line[j] = repl2
						if evalLine(line, 11) {
							accept = true
						}
					}
					line[j] = oldVal2
				}
			}
			line[i] = oldVal1
		}
		if accept {
			segs := 0
			for _, ch := range line {
				segs += figSegs[ch]
			}
			ans += segs
		}
	}
	return strconv.FormatInt(int64(ans), 11)
}

const doingPart = 3
