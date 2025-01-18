package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q01_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		total := 0
		for _, ch := range data {
			switch ch {
			case 'B':
				total += 1
			case 'C':
				total += 3
				// else ignore
			}
		}
		fmt.Println("Part 1:", total)
	}
	{
		var infile string
		if len(argsWithoutProg) < 2 || argsWithoutProg[1] == "" {
			infile = "../input/everybody_codes_e2024_q01_p2.txt"
		} else {
			infile = argsWithoutProg[1]
		}
		data, _ := os.ReadFile(infile)
		datastr := strings.Trim(string(data), " \r\n\t")
		total := 0

		for start := 0; start < len(datastr); start += 2 {
			encounter := datastr[start : start+2]
			if encounter == "xx" {
				continue
			}
			for _, ch := range encounter {
				switch ch {
				case 'A':
					total += 1
				case 'B':
					total += 2
				case 'C':
					total += 4
				case 'D':
					total += 6
				case 'x':
					total -= 1
				}
			}
		}
		fmt.Println("Part 2:", total)
	}
	{
		var infile string
		if len(argsWithoutProg) < 3 || argsWithoutProg[2] == "" {
			infile = "../input/everybody_codes_e2024_q01_p3.txt"
		} else {
			infile = argsWithoutProg[2]
		}
		data, _ := os.ReadFile(infile)
		datastr := strings.Trim(string(data), " \r\n\t")
		total := 0

		for start := 0; start < len(datastr); start += 3 {
			encounter := datastr[start : start+3]
			if encounter == "xxx" {
				continue
			}
			xdrop := 2
			for _, ch := range encounter {
				switch ch {
				case 'A':
					total += 2
				case 'B':
					total += 3
				case 'C':
					total += 5
				case 'D':
					total += 7
				case 'x':
					total -= xdrop
					xdrop = 0
				}
			}
		}
		fmt.Println("Part 3:", total)
	}
}
