package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q06_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q06_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q06_p3.txt`, "the input for part 3")

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
		knightList := scanner.Text()
		if scanner.Err() != nil {
			panic("scanner error")
		}
		seenCaps := 0
		pairingTotal := 0
		for _, ch := range []byte(knightList) {
			switch ch {
			case 'A':
				seenCaps++
			case 'a':
				pairingTotal += seenCaps
			}
		}
		fmt.Println("p1:", pairingTotal)
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
		knightList := scanner.Text()
		if scanner.Err() != nil {
			panic("scanner error")
		}
		var seenA, seenB, seenC int
		pairingTotal := 0
		for _, ch := range []byte(knightList) {
			switch ch {
			case 'A':
				seenA++
			case 'B':
				seenB++
			case 'C':
				seenC++
			case 'a':
				pairingTotal += seenA
			case 'b':
				pairingTotal += seenB
			case 'c':
				pairingTotal += seenC
			}
		}
		fmt.Println("p2:", pairingTotal)
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
		knightList := scanner.Bytes()
		if scanner.Err() != nil {
			panic("scanner error")
		}
		var pairingTotal0, pairingTotal1, pairingTotal999 int
		// first repetition
		for idx, ch := range []byte(knightList) {
			switch ch {
			case 'a':
				pairingTotal0 += countCh('A', knightList[max(0, idx-1000):min(len(knightList), idx+1001)])
				pairingTotal0 += countCh('A', knightList[0:max(0, idx+1001-len(knightList))])
			case 'b':
				pairingTotal0 += countCh('B', knightList[max(0, idx-1000):min(len(knightList), idx+1001)])
				pairingTotal0 += countCh('B', knightList[0:max(0, idx+1001-len(knightList))])
			case 'c':
				pairingTotal0 += countCh('C', knightList[max(0, idx-1000):min(len(knightList), idx+1001)])
				pairingTotal0 += countCh('C', knightList[0:max(0, idx+1001-len(knightList))])
			}
		}
		// second repetition
		for idx, ch := range []byte(knightList) {
			switch ch {
			case 'a':
				pairingTotal1 += countCh('A', knightList[min(idx-1000+len(knightList), len(knightList)):])
				pairingTotal1 += countCh('A', knightList[max(0, idx-1000):min(len(knightList), idx+1001)])
				pairingTotal1 += countCh('A', knightList[0:max(0, idx+1001-len(knightList))])
			case 'b':
				pairingTotal1 += countCh('B', knightList[min(idx-1000+len(knightList), len(knightList)):])
				pairingTotal1 += countCh('B', knightList[max(0, idx-1000):min(len(knightList), idx+1001)])
				pairingTotal1 += countCh('B', knightList[0:max(0, idx+1001-len(knightList))])
			case 'c':
				pairingTotal1 += countCh('C', knightList[min(idx-1000+len(knightList), len(knightList)):])
				pairingTotal1 += countCh('C', knightList[max(0, idx-1000):min(len(knightList), idx+1001)])
				pairingTotal1 += countCh('C', knightList[0:max(0, idx+1001-len(knightList))])
			}
		}
		// last repetition
		for idx, ch := range []byte(knightList) {
			switch ch {
			case 'a':
				pairingTotal999 += countCh('A', knightList[min(idx-1000+len(knightList), len(knightList)):])
				pairingTotal999 += countCh('A', knightList[max(0, idx-1000):min(len(knightList), idx+1001)])
				// pairingTotal1 += countCh('A', knightList[0:max(0, idx+1001-len(knightList))])
			case 'b':
				pairingTotal999 += countCh('B', knightList[min(idx-1000+len(knightList), len(knightList)):])
				pairingTotal999 += countCh('B', knightList[max(0, idx-1000):min(len(knightList), idx+1001)])
				// pairingTotal1 += countCh('B', knightList[0:max(0, idx+1001-len(knightList))])
			case 'c':
				pairingTotal999 += countCh('C', knightList[min(idx-1000+len(knightList), len(knightList)):])
				pairingTotal999 += countCh('C', knightList[max(0, idx-1000):min(len(knightList), idx+1001)])
				// pairingTotal1 += countCh('C', knightList[0:max(0, idx+1001-len(knightList))])
			}
		}
		fmt.Println("p3:", pairingTotal0+998*pairingTotal1+pairingTotal999)
	}
}

func countCh(ch byte, haystack []byte) int {
	retval := 0
	for _, x := range haystack {
		if x == ch {
			retval++
		}
	}
	return retval
}
