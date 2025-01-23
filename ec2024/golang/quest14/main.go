package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q14_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 1:", doProblem1(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 2 || argsWithoutProg[1] == "" {
			infile = "../input/everybody_codes_e2024_q14_p2.txt"
		} else {
			infile = argsWithoutProg[1]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 2:", doProblem2(data))
	}
	{
		var infile string
		if len(argsWithoutProg) < 3 || argsWithoutProg[2] == "" {
			infile = "../input/everybody_codes_e2024_q14_p3.txt"
		} else {
			infile = argsWithoutProg[2]
		}
		data, _ := os.ReadFile(infile)
		fmt.Println("Part 3:", doProblem3(data))
	}
}

func doProblem1(data []byte) int {
	strdata := string(data)
	maxHeight := 0
	height := 0
	for _, fld := range strings.Split(strdata, ",") {
		dir := fld[0]
		n, _ := strconv.Atoi(fld[1:])
		switch dir {
		case 'U':
			height += n
		case 'D':
			height -= n
		}
		maxHeight = max(maxHeight, height)
	}
	return maxHeight
}

type point struct{ x, y, z int }

func (p point) add(q point) point {
	return point{p.x + q.x, p.y + q.y, p.z + q.z}
}

func doProblem2(data []byte) int {
	dirMap := map[byte]point{
		'U': {0, 0, 1}, 'D': {0, 0, -1},
		'R': {0, 1, 0}, 'L': {0, -1, 0},
		'F': {1, 0, 0}, 'B': {-1, 0, 0}}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	tree := make(map[point]bool)
	for scanner.Scan() {
		where := point{0, 0, 0}
		for _, fld := range strings.Split(scanner.Text(), ",") {
			dirch := fld[0]
			n, _ := strconv.Atoi(fld[1:])
			dir := dirMap[dirch]
			for range n {
				where = where.add(dir)
				tree[where] = true
			}
		}
	}
	return len(tree)
}

func doProblem3(data []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	tree := make(map[point]bool)
	leaves := make([]point, 0)
	maxHeight := 0
	dirMap := map[byte]point{
		'U': {0, 0, 1}, 'D': {0, 0, -1},
		'R': {0, 1, 0}, 'L': {0, -1, 0},
		'F': {1, 0, 0}, 'B': {-1, 0, 0}}
	for scanner.Scan() {
		where := point{0, 0, 0}
		for _, fld := range strings.Split(scanner.Text(), ",") {
			dirch := fld[0]
			n, _ := strconv.Atoi(fld[1:])
			dir := dirMap[dirch]
			for range n {
				where = where.add(dir)
				tree[where] = true
			}
			maxHeight = max(maxHeight, where.z+1)
		}
		leaves = append(leaves, where)
	}
	murkiness := make(map[point]int)
	type queueType struct {
		where point
		dist  int
	}
	for _, leaf := range leaves {
		beenThere := make(map[point]bool)
		q := []queueType{{leaf, 0}}
		for len(q) > 0 {
			current := q[0]
			q = q[1:]
			if beenThere[current.where] {
				continue
			}
			beenThere[current.where] = true
			murkiness[current.where] += current.dist
			for _, dir := range dirMap {
				npos := current.where.add(dir)
				if tree[npos] && !beenThere[npos] {
					q = append(q, queueType{npos, current.dist + 1})
				}
			}
		}
	}
	minMurky := math.MaxInt
	for idx := range maxHeight {
		if murk, ok := murkiness[point{0, 0, idx}]; ok {
			minMurky = min(minMurky, murk)
		}
	}
	return minMurky
}
