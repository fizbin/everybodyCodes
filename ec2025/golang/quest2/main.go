package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q02_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q02_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q02_p3.txt`, "the input for part 3")

func CxAdd(a, b CxInt) CxInt {
	return CxInt{a.x + b.x, a.y + b.y}
}
func CxMul(a, b CxInt) CxInt {
	return CxInt{a.x*b.x - a.y*b.y, a.x*b.y + a.y*b.x}
}
func CxDiv(a, b CxInt) CxInt {
	return CxInt{a.x / b.x, a.y / b.y}
}

func shouldEngrave(test CxInt) bool {
	check := CxInt{0, 0}
	for range 100 {
		check = CxMul(check, check)
		check = CxDiv(check, CxInt{100000, 100000})
		check = CxAdd(check, test)

		if check.x > 1000000 || check.x < -1000000 {
			return false
		}
		if check.y > 1000000 || check.y < -1000000 {
			return false
		}
	}
	return true
}

type CxInt struct {
	x, y int
}

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
		noteStr := scanner.Text()
		parser := regexp.MustCompile(`A=\[ *(-?\d+), *(-?\d+)\]`)
		parsed := parser.FindStringSubmatch(noteStr)
		valX, errX := strconv.ParseInt(parsed[1], 10, 32)
		valY, errY := strconv.ParseInt(parsed[2], 10, 32)
		if (errX != nil) || (errY != nil) {
			panic("Bad input: " + noteStr)
		}
		valA := CxInt{(int)(valX), (int)(valY)}
		valR := CxInt{0, 0}

		valR = CxMul(valR, valR)
		valR = CxDiv(valR, CxInt{10, 10})
		valR = CxAdd(valR, valA)

		valR = CxMul(valR, valR)
		valR = CxDiv(valR, CxInt{10, 10})
		valR = CxAdd(valR, valA)

		valR = CxMul(valR, valR)
		valR = CxDiv(valR, CxInt{10, 10})
		valR = CxAdd(valR, valA)

		fmt.Printf("p1: [%d,%d]\n", valR.x, valR.y)
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
		noteStr := scanner.Text()
		parser := regexp.MustCompile(`A=\[ *(-?\d+), *(-?\d+)\]`)
		parsed := parser.FindStringSubmatch(noteStr)
		valX, errX := strconv.ParseInt(parsed[1], 10, 32)
		valY, errY := strconv.ParseInt(parsed[2], 10, 32)
		if (errX != nil) || (errY != nil) {
			panic("Bad input: " + noteStr)
		}
		valA := CxInt{(int)(valX), (int)(valY)}
		// valA = CxInt{35300, -64910}
		endCorner := CxAdd(valA, CxInt{1000, 1000})
		engraved := 0
		for testX := valA.x; testX <= endCorner.x; testX += 10 {
			for testY := valA.y; testY <= endCorner.y; testY += 10 {
				if shouldEngrave(CxInt{testX, testY}) {
					engraved++
				}
			}
		}
		fmt.Println("p2:", engraved)
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
		noteStr := scanner.Text()
		parser := regexp.MustCompile(`A=\[ *(-?\d+), *(-?\d+)\]`)
		parsed := parser.FindStringSubmatch(noteStr)
		valX, errX := strconv.ParseInt(parsed[1], 10, 32)
		valY, errY := strconv.ParseInt(parsed[2], 10, 32)
		if (errX != nil) || (errY != nil) {
			panic("Bad input: " + noteStr)
		}
		valA := CxInt{(int)(valX), (int)(valY)}
		// valA = CxInt{35300, -64910}
		endCorner := CxAdd(valA, CxInt{1000, 1000})
		engraved := 0
		for testX := valA.x; testX <= endCorner.x; testX += 1 {
			for testY := valA.y; testY <= endCorner.y; testY += 1 {
				if shouldEngrave(CxInt{testX, testY}) {
					engraved++
				}
			}
		}
		fmt.Println("p2:", engraved)
	}
}
