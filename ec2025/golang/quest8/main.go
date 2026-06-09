package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q08_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q08_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q08_p3.txt`, "the input for part 3")

func main() {
	flag.Parse()
	{
		infile := *input1
		data, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		fmt.Println("p1:", doProblem1(data))
	}
	{
		infile := *input2
		data, err := os.ReadFile(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		fmt.Println("p2:", doProblem2(data))
	}
	{
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
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	spotString := scanner.Text()
	spots := strings.Split(spotString, ",")
	crossings := 0
	pspot := 999
	for _, spot := range spots {
		spotNum, _ := strconv.Atoi(spot)
		if (spotNum+16 == pspot) || (spotNum-16 == pspot) {
			crossings++
		}
		pspot = spotNum
	}
	return crossings
}

type LineSeg struct {
	a, b int
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	spotString := scanner.Text()
	spots := strings.Split(spotString, ",")
	crossings := 0
	prevSegs := make([]LineSeg, 0, len(spots))
	pspot := 0
	for _, spotStr := range spots {
		spot, _ := strconv.Atoi(spotStr)
		// count crossings here - TODO
		for _, seg := range prevSegs {
			easyReject := (spot == seg.a) || (spot == seg.b) || (pspot == seg.a) || (pspot == seg.b)
			if easyReject {
				continue
			}
			spotInside := (spot > seg.a) != (spot > seg.b)
			pspotInside := (pspot > seg.a) != (pspot > seg.b)
			if spotInside != pspotInside {
				crossings++
			}
		}
		if pspot != 0 {
			prevSegs = append(prevSegs, LineSeg{pspot, spot})
		}
		pspot = spot
	}
	return crossings
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Scan()
	spotString := scanner.Text()
	spots := strings.Split(spotString, ",")
	prevSegs := make([]LineSeg, 0, len(spots))
	{
		pspot := 0
		for _, spotStr := range spots {
			spot, _ := strconv.Atoi(spotStr)
			if pspot != 0 {
				prevSegs = append(prevSegs, LineSeg{min(spot, pspot), max(spot, pspot)})
			}
			pspot = spot
		}
	}

	maxCrossings := 0
	for spotA := 1; spotA <= 255; spotA++ {
		for spotB := spotA + 1; spotB <= 256; spotB++ {
			crossings := 0
			for _, seg := range prevSegs {
				if (spotA == seg.a) && (spotB == seg.b) {
					crossings++
					continue
				}
				easyReject := (spotA == seg.a) || (spotA == seg.b) || (spotB == seg.a) || (spotB == seg.b)
				if easyReject {
					continue
				}
				spotInside := (spotA > seg.a) != (spotA > seg.b)
				pspotInside := (spotB > seg.a) != (spotB > seg.b)
				if spotInside != pspotInside {
					crossings++
				}
			}
			maxCrossings = max(maxCrossings, crossings)
		}
	}
	return maxCrossings
}
