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

type point struct {
	row int
	col int
}

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
		points := make([]point, 0)
		asterFinder := regexp.MustCompile(`[*][*][A-Z?]`)
		sawNoAster := true
		for idx, row := range grid {
			if idx+3 > len(grid) {
				break
			}
			if sawNoAster {
				for _, match := range asterFinder.FindAllIndex(row, -1) {
					points = append(points, point{idx, match[0]})
					sawNoAster = false
				}
			} else if !slices.Contains(row, '*') {
				sawNoAster = true
			}
		}
		// for _, row := range grid {
		// 	fmt.Println("1->:", string(row))
		// }
		// fmt.Println()

		// rand.Shuffle(len(points), func(i, j int) { points[i], points[j] = points[j], points[i] })

		madeChange := true
		for madeChange {
			madeChange = false
			for _, spot := range points {
				if collectWord(grid, spot.row, spot.col) != nil {
					before := ""
					for idx := range 8 {
						before += string(grid[spot.row+idx][spot.col:spot.col+8]) + "\n"
					}
					fillGridPass1(grid, spot.row, spot.col)
					madeChange = fillGridPass2(grid, spot.row, spot.col) || madeChange
					if !crossCheckBlock(grid, spot.row, spot.col) {
						after := ""
						for idx := range 8 {
							after += string(grid[spot.row+idx][spot.col:spot.col+8]) + "\n"
						}
						fmt.Println("ERR!", "\n"+before, "\n"+after)
						panic("Inconsistent!")
					}
				}
			}
			// Debugging
			// if madeChange {
			// 	for _, row := range grid {
			// 		fmt.Println("2->:", string(row))
			// 	}
			// 	fmt.Println()
			// }
		}

		total := 0
		for _, spot := range points {
			word := collectWord(grid, spot.row, spot.col)
			wordPowr := 0
			if word != nil {
				wordPowr = wordPower(word)
			}
			if wordPowr != 0 {
				total += wordPowr
			} else {
				for xpos := 2; xpos < 6; xpos++ {
					for ypos := 2; ypos < 6; ypos++ {
						grid[xpos+spot.row][ypos+spot.col] = '.'
					}
				}
				fillGridPass1(grid, spot.row, spot.col)
				if collectWord(grid, spot.row, spot.col) != nil {
					fmt.Println("DBG New word ", string(collectWord(grid, spot.row, spot.col)), " at ", spot)
				}
				for xpos := 2; xpos < 6; xpos++ {
					for ypos := 2; ypos < 6; ypos++ {
						grid[xpos+spot.row][ypos+spot.col] = '.'
					}
				}
			}
		}

		// Debugging
		// for _, row := range grid {
		// 	fmt.Println("3->:", string(row))
		// }

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

func crossCheckBlock(grid [][]byte, offsetRow, offsetCol int) bool {
	seen := make(map[byte]bool)
	for xpos := 2; xpos < 6; xpos++ {
		charsRow := []byte{grid[xpos+offsetRow][0+offsetCol], grid[xpos+offsetRow][1+offsetCol], grid[xpos+offsetRow][6+offsetCol], grid[xpos+offsetRow][7+offsetCol]}
		for ypos := 2; ypos < 6; ypos++ {
			charsCol := []byte{grid[0+offsetRow][ypos+offsetCol], grid[1+offsetRow][ypos+offsetCol], grid[6+offsetRow][ypos+offsetCol], grid[7+offsetRow][ypos+offsetCol]}
			me := grid[xpos+offsetRow][ypos+offsetCol]
			if me >= 'A' && me <= 'Z' {
				if seen[me] {
					return false // log.Panic("Duplicate letter at ", xpos, ypos, " within block at ", offsetRow, offsetCol)
				}
				if !slices.Contains(charsRow, me) {
					return false // log.Panic("Letter not in charsRow at ", xpos, ypos, " within block at ", offsetRow, offsetCol)
				}
				if !slices.Contains(charsCol, me) {
					return false // log.Panic("Letter not in charsCol at ", xpos, ypos, " within block at ", offsetRow, offsetCol)
				}
			}
		}
	}
	return true
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
	changes := make(map[point]byte)
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
					return false
				}
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
						return false
					}
					changes[point{xpos + offsetRow, ypos + offsetCol}] = replacementChar
					if !slices.Contains(charsRow, replacementChar) {
						for _, ypos2 := range []int{0, 1, 6, 7} {
							if grid[xpos+offsetRow][ypos2+offsetCol] == '?' {
								changes[point{xpos + offsetRow, ypos2 + offsetCol}] = replacementChar
							}
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
						return false
					}
					changes[point{xpos + offsetRow, ypos + offsetCol}] = replacementChar
					if !slices.Contains(charsCol, replacementChar) {
						for _, xpos2 := range []int{0, 1, 6, 7} {
							if grid[xpos2+offsetRow][ypos+offsetCol] == '?' {
								changes[point{xpos2 + offsetRow, ypos + offsetCol}] = replacementChar
							}
						}
					}
				}
			}
		}
	}
	for spot, val := range changes {
		grid[spot.row][spot.col] = val
	}
	return len(changes) > 0
}

func collectWord(grid [][]byte, offsetRow, offsetCol int) []byte {
	retval := make([]byte, 0)
	seenLetters := make(map[byte]bool)
	for xpos := 2; xpos < 6; xpos++ {
		for ypos := 2; ypos < 6; ypos++ {
			me := grid[xpos+offsetRow][ypos+offsetCol]
			if me == '%' {
				return nil
			}
			if me != '.' {
				if seenLetters[me] {
					return nil
				}
				seenLetters[me] = true
			}
			retval = append(retval, me)
		}
	}
	return retval
}

func wordPower(word []byte) int {
	retval := 0
	for idx, v := range word {
		if v == '.' {
			log.Fatal("Bad word ", string(word))
		}
		retval += (idx + 1) * int(v-'A'+1)
	}
	return retval
}
