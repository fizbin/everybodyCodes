package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e2_q01_p1.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		board := make([][]byte, 0)
		for scanner.Scan() {
			if scanner.Text() == "" {
				break
			}
			board = append(board, bytes.Clone(scanner.Bytes()))
		}
		tokspecs := make([][]byte, 0)
		for scanner.Scan() {
			tokspecs = append(tokspecs, bytes.Clone(scanner.Bytes()))
		}
		if err := scanner.Err(); err != nil {
			panic("Oops")
		}
		tot_score := 0
		for idx, tokspec := range tokspecs {
			out_col := RunToken(board, 2*idx, tokspec, false)
			in_slot := 1 + idx
			out_slot := out_col/2 + 1
			tot_score += max(out_slot*2-in_slot, 0)
		}
		fmt.Println("p1 Total score: ", tot_score)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			// infile = `..\input\everybody_codes_e2_q01_p2_samp.txt`
			infile = `..\input\everybody_codes_e2_q01_p2.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		board := make([][]byte, 0)
		for scanner.Scan() {
			if scanner.Text() == "" {
				break
			}
			board = append(board, bytes.Clone(scanner.Bytes()))
		}
		tokspecs := make([][]byte, 0)
		for scanner.Scan() {
			tokspecs = append(tokspecs, bytes.Clone(scanner.Bytes()))
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			panic("Oops")
		}
		tot_score := 0
		for idx, tokspec := range tokspecs {
			max_win := 0
			chosen := 0
			chosen_out := 0
			for col_half := range (len(board[0]) + 1) / 2 {
				out_col := RunToken(board, 2*col_half, tokspec, idx < 0)
				in_slot := 1 + col_half
				out_slot := out_col/2 + 1
				spot_score := max(out_slot*2-in_slot, 0)
				if spot_score > max_win {
					chosen = in_slot
					chosen_out = out_slot
					max_win = max(spot_score, max_win)
				}
			}
			if max_win < 0 {
				fmt.Println(string(tokspec), chosen, "->", chosen_out, "=", max_win)
			}
			tot_score += max_win
		}
		fmt.Println("p2 Total score: ", tot_score)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			// infile = `..\input\everybody_codes_e2_q01_p3_samp.txt`
			infile = `..\input\everybody_codes_e2_q01_p3.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		board := make([][]byte, 0)
		for scanner.Scan() {
			if scanner.Text() == "" {
				break
			}
			board = append(board, bytes.Clone(scanner.Bytes()))
		}
		tokspecs := make([][]byte, 0)
		for scanner.Scan() {
			tokspecs = append(tokspecs, bytes.Clone(scanner.Bytes()))
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			panic("Oops")
		}
		slot_max := (len(board[0]) + 1) / 2
		score_array := make([][]int, 0, len(tokspecs))
		for range len(tokspecs) {
			score_array = append(score_array, make([]int, slot_max))
		}
		for idx, tokspec := range tokspecs {
			for col_half := range slot_max {
				out_col := RunToken(board, 2*col_half, tokspec, idx < 0)
				in_slot := 1 + col_half
				out_slot := out_col/2 + 1
				spot_score := max(out_slot*2-in_slot, 0)
				score_array[idx][col_half] = spot_score
			}
		}
		max_score := 0
		min_score := slot_max * len(tokspecs)
		// min_slots := make([]int, 0)
		// max_slots := make([]int, 0)
		initial_slots := make([]int, len(tokspecs))
		for idx := range initial_slots {
			initial_slots[idx] = idx
		}
		// fmt.Println("...reticulating splines...")
		for {
			score := 0
			for tokidx, slot := range initial_slots {
				score += score_array[tokidx][slot]
			}
			// if score > max_score {
			// 	max_slots = slices.Clone(initial_slots)
			// }
			// if score < min_score {
			// 	min_slots = slices.Clone(initial_slots)
			// }
			max_score = max(max_score, score)
			min_score = min(min_score, score)

			// now inc
			inc_col := len(initial_slots) - 1
			for {
				initial_slots[inc_col] += 1
				if initial_slots[inc_col] < slot_max {
					good_val := true
					for pcol := 0; pcol < inc_col; pcol++ {
						if initial_slots[pcol] == initial_slots[inc_col] {
							good_val = false
						}
					}
					if good_val {
						break // We were able to increment initial_slots[inc_col] to a good value
					} else {
						continue // Try to increment initial_slots[inc_col] again
					}
				} else {
					// We ran out of room to increment initial_slots[inc_col]
					inc_col -= 1
					if inc_col < 0 {
						break
					}
				}
			}
			if inc_col < 0 {
				// We're done
				break
			}
			// Okay, we incremented initial_slots[inc_col]. Now we set all the
			// initial_slots[ocol] for ocol > inc_col to the minimum value that's legal
			for ocol := inc_col + 1; ocol < len(initial_slots); ocol++ {
				initial_slots[ocol] = 0
				for {
					good_val := true
					for pcol := 0; pcol < ocol; pcol++ {
						if initial_slots[pcol] == initial_slots[ocol] {
							good_val = false
						}
					}
					if good_val {
						break
					}
					initial_slots[ocol] += 1
				}
			}
			// if inc_col == 0 {
			// 	fmt.Println("incremented to", initial_slots)
			// }
		}
		fmt.Println("p3", min_score, max_score)
		// fmt.Println("p3 at", min_slots, max_slots)
	}
}

func RunToken(board [][]byte, i int, tokspec []byte, debug bool) int {
	ins_idx := 0
	cur_col := i
	for rowidx, row := range board {
		if debug {
			fmt.Println("At row", rowidx, string(row), "have", cur_col)
		}
		if row[cur_col] == '*' {
			instruction := tokspec[ins_idx]
			ins_idx += 1
			new_col := cur_col
			if instruction == 'L' {
				new_col = cur_col - 1
			} else {
				new_col = cur_col + 1
			}
			if new_col < 0 {
				new_col = cur_col + 1
			}
			if new_col >= len(row) {
				new_col = cur_col - 1
			}
			cur_col = new_col
		}
	}
	if debug {
		fmt.Println("Finally hav column", cur_col)
	}
	return cur_col
}
