package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
)

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q10_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		grid := make([][]byte, 0)
		for scanner.Scan() {
			grid = append(grid, []byte(scanner.Text()))
		}
		word := computeWord(grid, 0, 0)

		fmt.Println("Part 1:", string(word))
	}
	{
		var infile string
		if len(argsWithoutProg) < 2 || argsWithoutProg[1] == "" {
			infile = "../input/everybody_codes_e2024_q10_p2.txt"
		} else {
			infile = argsWithoutProg[1]
		}
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		grid := make([][]byte, 0)
		for scanner.Scan() {
			grid = append(grid, []byte(scanner.Text()))
		}
		total := 0
		sawblank := true
		asterFinder := regexp.MustCompile(`(?:^|[ ])[*]`)
		for idx, row := range grid {
			if sawblank {
				sawblank = false
				for _, match := range asterFinder.FindAllIndex(row, -1) {
					word := computeWord(grid, idx, match[1]-1)
					total += wordPower(word)
				}
			} else if len(row) == 0 {
				sawblank = true
			}
		}

		fmt.Println("Part 2:", total)
	}
	{
		var infile string
		if len(argsWithoutProg) < 3 || argsWithoutProg[2] == "" {
			infile = "../input/everybody_codes_e2024_q10_p3.txt"
		} else {
			infile = argsWithoutProg[2]
		}
		data, _ := os.ReadFile(infile)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		grid := make([][]byte, 0)
		for scanner.Scan() {
			grid = append(grid, []byte(scanner.Text()))
		}
		total := 0
		asterFinder := regexp.MustCompile(`[*][*][A-Z?]`)
		sawNoAster := true
		for idx, row := range grid {
			if idx+3 > len(grid) {
				break
			}
			if sawNoAster {
				for _, match := range asterFinder.FindAllIndex(row, -1) {
					sawNoAster = false
					fillGridPass1(grid, idx, match[0])
				}
			} else if !slices.Contains(row, '*') {
				sawNoAster = true
			}
		}
		for _, row := range grid {
			fmt.Println("1->:", string(row))
		}
		fmt.Println()

		fillGridPass2(grid, 0, 48)

		madeChange := true
		for madeChange {
			madeChange = false
			sawNoAster = true
			for idx, row := range grid {
				if idx+3 > len(grid) {
					break
				}
				if sawNoAster {
					for _, match := range asterFinder.FindAllIndex(row, -1) {
						sawNoAster = false
						if collectWord(grid, idx, match[0]) != nil {
							madeChange = fillGridPass2(grid, idx, match[0]) || madeChange
						}
					}
				} else if !slices.Contains(row, '*') {
					sawNoAster = true
				}
			}
			if madeChange {
				for _, row := range grid {
					fmt.Println("2->:", string(row))
				}
				fmt.Println()
			}
		}

		sawNoAster = true
		for idx, row := range grid {
			if idx+3 > len(grid) {
				break
			}
			if sawNoAster {
				for _, match := range asterFinder.FindAllIndex(row, -1) {
					sawNoAster = false
					word := collectWord(grid, idx, match[0])
					if word != nil {
						total += wordPower(word)
					} else {
						for xpos := 2; xpos < 6; xpos++ {
							for ypos := 2; ypos < 6; ypos++ {
								grid[xpos+idx][ypos+match[0]] = '.'
							}
						}
						fillGridPass1(grid, idx, match[0])
						if collectWord(grid, idx, match[0]) != nil {
							fmt.Println("DBG New word ", string(collectWord(grid, idx, match[0])), " at ", idx, match[0])
						}
						for xpos := 2; xpos < 6; xpos++ {
							for ypos := 2; ypos < 6; ypos++ {
								grid[xpos+idx][ypos+match[0]] = '.'
							}
						}
					}
				}
			} else if !slices.Contains(row, '*') {
				sawNoAster = true
			}
		}

		for _, row := range grid {
			fmt.Println("3->:", string(row))
		}

		fmt.Println("Part 3:", total)
	}
}

func computeWord(grid [][]byte, offsetRow, offsetCol int) []byte {
	word := make([]byte, 0)
	for xpos := 2; xpos < 6; xpos++ {
		charsRow := []byte{grid[xpos+offsetRow][0+offsetCol], grid[xpos+offsetRow][1+offsetCol], grid[xpos+offsetRow][6+offsetCol], grid[xpos+offsetRow][7+offsetCol]}
		for ypos := 2; ypos < 6; ypos++ {
			charsCol := []byte{grid[0+offsetRow][ypos+offsetCol], grid[1+offsetRow][ypos+offsetCol], grid[6+offsetRow][ypos+offsetCol], grid[7+offsetRow][ypos+offsetCol]}
			for _, ch := range charsCol {
				if slices.Contains(charsRow, ch) {
					word = append(word, ch)
					break
				}
			}
		}
	}
	return word
}

// fills innards we know right off
func fillGridPass1(grid [][]byte, offsetRow, offsetCol int) {
	for xpos := 2; xpos < 6; xpos++ {
		charsRow := []byte{grid[xpos+offsetRow][0+offsetCol], grid[xpos+offsetRow][1+offsetCol], grid[xpos+offsetRow][6+offsetCol], grid[xpos+offsetRow][7+offsetCol]}
		for ypos := 2; ypos < 6; ypos++ {
			charsCol := []byte{grid[0+offsetRow][ypos+offsetCol], grid[1+offsetRow][ypos+offsetCol], grid[6+offsetRow][ypos+offsetCol], grid[7+offsetRow][ypos+offsetCol]}
			foundSpot := false
			for _, ch := range charsCol {
				if ch == '?' {
					continue
				}
				if slices.Contains(charsRow, ch) {
					grid[xpos+offsetRow][ypos+offsetCol] = ch
					foundSpot = true
					break
				}
			}
			if !foundSpot && !slices.Contains(charsCol, '?') && !slices.Contains(charsRow, '?') {
				grid[xpos+offsetRow][ypos+offsetCol] = '%'
			}
		}
	}
}

func fillGridPass2(grid [][]byte, offsetRow, offsetCol int) bool {
	madeChange := false
	for xpos := 2; xpos < 6; xpos++ {
		charsRow := []byte{grid[xpos+offsetRow][0+offsetCol], grid[xpos+offsetRow][1+offsetCol], grid[xpos+offsetRow][6+offsetCol], grid[xpos+offsetRow][7+offsetCol]}
		for ypos := 2; ypos < 6; ypos++ {
			charsCol := []byte{grid[0+offsetRow][ypos+offsetCol], grid[1+offsetRow][ypos+offsetCol], grid[6+offsetRow][ypos+offsetCol], grid[7+offsetRow][ypos+offsetCol]}
			me := grid[xpos+offsetRow][ypos+offsetCol]
			if me == '.' {
				nQuestionMarks := 0
				for _, ch := range slices.Concat(charsRow, charsCol) {
					if ch == '?' {
						nQuestionMarks++
					}
				}
				if nQuestionMarks > 1 {
					continue
				}
				if nQuestionMarks == 0 {
					foundSpot := false
					for _, ch := range charsCol {
						if slices.Contains(charsRow, ch) {
							grid[xpos+offsetRow][ypos+offsetCol] = ch
							foundSpot = true
							break
						}
					}
					if !foundSpot {
						grid[xpos+offsetRow][ypos+offsetCol] = '%'
					}
				} else {
					if slices.Contains(charsRow, '?') {
						colcontents := make([]byte, 0)
						for xpos2 := 2; xpos2 < 6; xpos2++ {
							colcontents = append(colcontents, grid[xpos2+offsetRow][ypos+offsetCol])
						}
						var replacementChar byte
						for _, ch := range charsCol {
							if !slices.Contains(colcontents, ch) {
								if replacementChar == 0 {
									replacementChar = ch
								} else {
									replacementChar = 0
									break
								}
							}
						}
						if replacementChar == 0 {
							continue
						}
						grid[xpos+offsetRow][ypos+offsetCol] = replacementChar
						for _, ypos2 := range []int{0, 1, 6, 7} {
							if grid[xpos+offsetRow][ypos2+offsetCol] == '?' {
								grid[xpos+offsetRow][ypos2+offsetCol] = replacementChar
							}
						}
					} else {
						// slices.Contains(charsCol, '?')
						rowcontents := grid[xpos+offsetRow][2+offsetCol : 6+offsetCol]
						var replacementChar byte
						for _, ch := range charsRow {
							if !slices.Contains(rowcontents, ch) {
								if replacementChar == 0 {
									replacementChar = ch
								} else {
									replacementChar = 0
									break
								}
							}
						}
						if replacementChar == 0 {
							continue
						}
						grid[xpos+offsetRow][ypos+offsetCol] = replacementChar
						for _, xpos2 := range []int{0, 1, 6, 7} {
							if grid[xpos2+offsetRow][ypos+offsetCol] == '?' {
								grid[xpos2+offsetRow][ypos+offsetCol] = replacementChar
							}
						}
					}
				}
				madeChange = true
			}
		}
	}
	return madeChange
}

func collectWord(grid [][]byte, offsetRow, offsetCol int) []byte {
	retval := make([]byte, 0)
	for xpos := 2; xpos < 6; xpos++ {
		for ypos := 2; ypos < 6; ypos++ {
			me := grid[xpos+offsetRow][ypos+offsetCol]
			if me == '%' {
				return nil
			}
			retval = append(retval, me)
		}
	}
	return retval
}

// computeWordP3 can return nil; also, it will adjust ? values in grid as needed
func computeWordP3(grid [][]byte, offsetRow, offsetCol int) []byte {
	dotsx := make([]int, 0)
	dotsy := make([]int, 0)
	for xpos := 2; xpos < 6; xpos++ {
		charsRow := []byte{grid[xpos+offsetRow][0+offsetCol], grid[xpos+offsetRow][1+offsetCol], grid[xpos+offsetRow][6+offsetCol], grid[xpos+offsetRow][7+offsetCol]}
		for ypos := 2; ypos < 6; ypos++ {
			charsCol := []byte{grid[0+offsetRow][ypos+offsetCol], grid[1+offsetRow][ypos+offsetCol], grid[6+offsetRow][ypos+offsetCol], grid[7+offsetRow][ypos+offsetCol]}
			foundSpot := false
			for _, ch := range charsCol {
				if ch == '?' {
					continue
				}
				if slices.Contains(charsRow, ch) {
					grid[xpos+offsetRow][ypos+offsetCol] = ch
					foundSpot = true
					break
				}
			}
			if !foundSpot && !slices.Contains(charsCol, '?') && !slices.Contains(charsRow, '?') {
				grid[xpos+offsetRow][ypos+offsetCol] = '%'
				return nil
			} else if !foundSpot {
				dotsx = append(dotsx, xpos)
				dotsy = append(dotsy, ypos)
			}
		}
	}
	if len(dotsx) > 0 {

		// fmt.Println("DBG-cwp3")
		// for xpos := 0; xpos < 8; xpos++ {
		// 	for ypos := 0; ypos < 8; ypos++ {
		// 		fmt.Print(string(grid[xpos+offsetRow][ypos+offsetCol : ypos+offsetCol+1]))
		// 	}
		// 	fmt.Println()
		// }
		// fmt.Println()
		for idx, xpos := range dotsx {
			ypos := dotsy[idx]
			charsRow := []byte{grid[xpos+offsetRow][0+offsetCol], grid[xpos+offsetRow][1+offsetCol], grid[xpos+offsetRow][6+offsetCol], grid[xpos+offsetRow][7+offsetCol]}
			charsCol := []byte{grid[0+offsetRow][ypos+offsetCol], grid[1+offsetRow][ypos+offsetCol], grid[6+offsetRow][ypos+offsetCol], grid[7+offsetRow][ypos+offsetCol]}
			// if slices.Contains(charsRow, '?') == slices.Contains(charsCol, '?') {
			// 	log.Fatal("Expected exactly one question mark", string(charsRow), string(charsCol), offsetRow, offsetCol, "@", xpos, ypos)
			// }
			// fmt.Println("DBG-cwp3.charsRow: ", string(charsRow))
			// fmt.Println("DBG-cwp3.charsCol: ", string(charsCol))
			if slices.Contains(charsRow, '?') {
				colcontents := make([]byte, 0)
				for xpos2 := 2; xpos2 < 6; xpos2++ {
					colcontents = append(colcontents, grid[xpos2+offsetRow][ypos+offsetCol])
				}
				// fmt.Println("DBG-cwp3.colcontents: ", string(colcontents))
				for _, ch := range charsCol {
					if !slices.Contains(colcontents, ch) {
						grid[xpos+offsetRow][ypos+offsetCol] = ch
						for _, ypos2 := range []int{0, 1, 6, 7} {
							if grid[xpos+offsetRow][ypos2+offsetCol] == '?' {
								grid[xpos+offsetRow][ypos2+offsetCol] = ch
							}
						}
						break
					}
				}
			} else {
				// slices.Contains(charsCol, '?')
				rowcontents := grid[xpos+offsetRow][2+offsetCol : 6+offsetCol]
				for _, ch := range charsRow {
					if !slices.Contains(rowcontents, ch) {
						grid[xpos+offsetRow][ypos+offsetCol] = ch
						for _, xpos2 := range []int{0, 1, 6, 7} {
							if grid[xpos2+offsetRow][ypos+offsetCol] == '?' {
								grid[xpos2+offsetRow][ypos+offsetCol] = ch
							}
						}
						break
					}
				}
			}
		}
	}
	word := make([]byte, 0)
	for xpos := 2; xpos < 6; xpos++ {
		word = slices.Concat(word, grid[xpos+offsetRow][2+offsetCol:6+offsetCol])
	}
	return word
}

func wordPower(word []byte) int {
	retval := 0
	for idx, v := range word {
		if v == '.' {
			log.Fatal("Bad word: ", string(word))
		}
		retval += (idx + 1) * int(v-'A'+1)
	}
	return retval
}
