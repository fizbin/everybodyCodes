package main

import (
	"bufio"
	"bytes"
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
			infile = "../input/everybody_codes_e2024_q02_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		// data = []byte("WORDS:THE,OWE,MES,ROD,HER\n\nAWAKEN THE POWER ADORNED WITH THE FLAMES BRIGHT IRE\n")
		scanner := bufio.NewScanner(bytes.NewReader(data))
		scanner.Scan()
		wordLine := scanner.Text()
		words := strings.Split(strings.Split(wordLine, ":")[1], ",")
		wordCheck := make(map[string]bool)
		maxWordLen := 0
		for _, word := range words {
			wordCheck[word] = true
			maxWordLen = max(maxWordLen, len(word))
		}

		scanner.Scan() // blank line
		scanner.Scan()
		haystack := scanner.Text()
		total := 0
		for start := 0; start < len(haystack); start++ {
			found := 0
			for l := 1; l <= maxWordLen; l++ {
				if start+l > len(haystack) {
					break
				}
				if wordCheck[haystack[start:start+l]] {
					found = 1
					break
				}
			}
			total += found
		}

		fmt.Println("Part 1:", total)
	}

	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q02_p2.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		// data = []byte("WORDS:THE,OWE,MES,ROD,HER\n\nAWAKEN THE POWER ADORNED WITH THE FLAMES BRIGHT IRE\n")
		scanner := bufio.NewScanner(bytes.NewReader(data))
		scanner.Scan()
		wordLine := scanner.Text()
		words := strings.Split(strings.Split(wordLine, ":")[1], ",")
		wordCheck := make(map[string]bool)
		maxWordLen := 0
		for _, word := range words {
			wordCheck[word] = true
			maxWordLen = max(maxWordLen, len(word))
			secondStrBytes := []byte(word)
			for j := 0; j < len(secondStrBytes)-j-1; j++ {
				secondStrBytes[j], secondStrBytes[len(secondStrBytes)-j-1] = secondStrBytes[len(secondStrBytes)-j-1], secondStrBytes[j]
			}
			wordCheck[string(secondStrBytes)] = true
		}

		scanner.Scan() // blank line
		total := 0

		for scanner.Scan() {
			haystack := scanner.Text()
			foundchars := make(map[int]bool)
			for start := 0; start < len(haystack); start++ {
				for l := 1; l <= maxWordLen; l++ {
					if start+l > len(haystack) {
						break
					}
					if wordCheck[haystack[start:start+l]] {
						for i := range l {
							foundchars[start+i] = true
						}
					}
				}
			}
			total += len(foundchars)
		}

		fmt.Println("Part 2:", total)
	}

	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q02_p3.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		// data = []byte("WORDS:THE,OWE,MES,ROD,HER\n\nAWAKEN THE POWER ADORNED WITH THE FLAMES BRIGHT IRE\n")
		scanner := bufio.NewScanner(bytes.NewReader(data))
		scanner.Scan()
		wordLine := scanner.Text()
		words := strings.Split(strings.Split(wordLine, ":")[1], ",")
		wordCheck := make(map[string]bool)
		maxWordLen := 0
		for _, word := range words {
			wordCheck[word] = true
			maxWordLen = max(maxWordLen, len(word))
			secondStrBytes := []byte(word)
			for j := 0; j < len(secondStrBytes)-j-1; j++ {
				secondStrBytes[j], secondStrBytes[len(secondStrBytes)-j-1] = secondStrBytes[len(secondStrBytes)-j-1], secondStrBytes[j]
			}
			wordCheck[string(secondStrBytes)] = true
		}

		scanner.Scan() // blank line
		type coord struct {
			x int
			y int
		}
		allfound := make(map[coord]bool)
		row := 0
		var transposed [][]byte
		for scanner.Scan() {
			haystack := scanner.Text()
			if transposed == nil {
				for range []byte(haystack) {
					transposed = append(transposed, make([]byte, 0, len(haystack)))
				}
			}
			for idx, v := range []byte(haystack) {
				transposed[idx] = append(transposed[idx], v)
			}
			foundchars := make(map[int]bool)
			haystackp := haystack + haystack[0:maxWordLen]
			for start := 0; start < len(haystack); start++ {
				for l := 1; l <= maxWordLen; l++ {
					if start+l > len(haystackp) {
						break
					}
					if wordCheck[haystackp[start:start+l]] {
						for i := range l {
							foundchars[(start+i)%len(haystack)] = true
						}
					}
				}
			}
			for yval := range foundchars {
				allfound[coord{row, yval}] = true
			}
			row++
		}

		for yval, column := range transposed {
			foundchars := make(map[int]bool)
			haystack := string(column)
			for start := 0; start < len(haystack); start++ {
				for l := 1; l <= maxWordLen; l++ {
					if start+l > len(haystack) {
						break
					}
					if wordCheck[haystack[start:start+l]] {
						for i := range l {
							foundchars[start+i] = true
						}
					}
				}
			}
			for xval := range foundchars {
				allfound[coord{xval, yval}] = true
			}

		}

		fmt.Println("Part 3:", len(allfound))
	}
}
