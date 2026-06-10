package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q10_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q10_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q10_p3.txt`, "the input for part 3")

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

type coord struct {
	x, y int
}

var jumps = []coord{coord{-1, -2}, coord{-2, -1}, coord{2, -1}, coord{-1, 2}, coord{-2, 1}, coord{1, -2}, coord{1, 2}, coord{2, 1}}

func doProblem1(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	board := make(map[coord]byte)
	var dragon coord
	row := 0
	for scanner.Scan() {
		for col, val := range scanner.Bytes() {
			board[coord{row, col}] = val
			if val == 'D' {
				dragon = coord{row, col}
			}
		}
		row++
	}
	possibleMoves := make(map[coord]bool)
	possibleMoves[dragon] = true
	sheepEaten := make(map[coord]bool)
	for range 4 {
		newMoves := make(map[coord]bool)
		for _, jump := range jumps {
			for start := range possibleMoves {
				dst := coord{jump.x + start.x, jump.y + start.y}
				if board[dst] == 'S' {
					sheepEaten[dst] = true
				}
				newMoves[dst] = true
			}
		}
		possibleMoves = newMoves
	}
	// fmt.Println(sheepEaten)
	return len(sheepEaten)
}

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	board := make(map[coord]byte)
	var dragon coord
	row := 0
	for scanner.Scan() {
		for col, val := range scanner.Bytes() {
			board[coord{row, col}] = val
			if val == 'D' {
				dragon = coord{row, col}
			}
		}
		row++
	}
	pDragons := make(map[coord]bool)
	pDragons[dragon] = true
	sheepEaten := make(map[coord]bool)
	for pSheepMoves := range 20 {
		newDragons := make(map[coord]bool)
		for _, jump := range jumps {
			for start := range pDragons {
				dst := coord{jump.x + start.x, jump.y + start.y}
				if dst.x >= row {
					continue
				}
				sheepSpot0 := coord{dst.x - pSheepMoves, dst.y}
				sheepSpot1 := coord{dst.x - pSheepMoves - 1, dst.y}
				if board[dst] != '#' && board[sheepSpot0] == 'S' {
					sheepEaten[sheepSpot0] = true
				}
				if board[dst] != '#' && board[sheepSpot1] == 'S' {
					sheepEaten[sheepSpot1] = true
				}
				newDragons[dst] = true
			}
		}
		pDragons = newDragons
	}
	// fmt.Println(sheepEaten)
	return len(sheepEaten)
}

const MaxCols = 10

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	board := make(map[coord]byte)
	var dragon coord
	var sheep [MaxCols]int
	for i := range MaxCols {
		sheep[i] = -1
	}
	row := 0
	col := 0
	for scanner.Scan() {
		for c, val := range scanner.Bytes() {
			board[coord{row, c}] = val
			if val == 'D' {
				dragon = coord{row, c}
			} else if val == 'S' {
				sheep[c] = row
			}
			col = max(col, c)
		}
		row++
	}
	col++

	memo := make(map[memoKey]uint64)
	// so now row and col give the limits of the board
	return nMovesSheep(memoKey{dragon, sheep}, memo, board, row, col)
}

type memoKey struct {
	dragon coord
	sheep  [MaxCols]int
}

func nMovesSheep(state memoKey, memo map[memoKey]uint64, board map[coord]byte, row, col int) uint64 {
	if ans, ok := memo[state]; ok {
		return ans
	}
	ans := uint64(0)
	anySheep := false
	allBlocked := true
	for c := range col {
		if state.sheep[c] != -1 {
			anySheep = true
			if state.sheep[c]+1 == row {
				// moving this sheep means sheep escape
				allBlocked = false
				continue
			}
			if (state.dragon.y == c) && (state.sheep[c]+1 == state.dragon.x) && (board[state.dragon] != '#') {
				// moving this sheep means dragon would eat it, so don't move it
			} else {
				allBlocked = false
				nState := memoKey{state.dragon, state.sheep}
				nState.sheep[c] += 1
				subAns := nMovesDragon(nState, memo, board, row, col)
				ans += subAns
			}
		}
	}
	if !anySheep {
		// dragon must have eaten the last sheep prior to this move
		ans = 1
	}
	if ans == 0 && allBlocked {
		ans = nMovesDragon(state, memo, board, row, col)
	}
	if ans > 0 {
		// fmt.Println("p3 debug (S)", state.dragon, state.sheep[0:col], "->", ans)
	}
	memo[state] = ans
	return ans
}

func nMovesDragon(state memoKey, memo map[memoKey]uint64, board map[coord]byte, row, col int) uint64 {
	ans := uint64(0)
	for _, jump := range jumps {
		nDragon := coord{state.dragon.x + jump.x, state.dragon.y + jump.y}
		nSheep := state.sheep
		if nDragon.x < 0 {
			continue
		}
		if nDragon.y < 0 {
			continue
		}
		if nDragon.x >= row {
			continue
		}
		if nDragon.y >= col {
			continue
		}
		if board[nDragon] != '#' {
			if nSheep[nDragon.y] == nDragon.x {
				nSheep[nDragon.y] = -1
			}
		}
		ans += nMovesSheep(memoKey{nDragon, nSheep}, memo, board, row, col)
	}
	return ans
}
