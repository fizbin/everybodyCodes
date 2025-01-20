package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Pos struct {
	x int
	y int
}

func (here Pos) Add(dir Pos) Pos {
	return Pos{here.x + dir.x, here.y + dir.y}
}

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q08_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		n, _ := strconv.Atoi(string(data))
		layer := 0
		for layer*layer < n {
			layer++
		}

		fmt.Println("Part 1:", (2*layer-1)*(layer*layer-n))
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q08_p2.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		blocks := 20240000
		priests, _ := strconv.Atoi(string(data))
		acolytes := 1111
		thickness := 1
		totalblocks := 1
		layer := 1
		for totalblocks < blocks {
			layer++
			thickness = (thickness * priests) % acolytes
			totalblocks += thickness * (2*layer - 1)
		}

		fmt.Println("Part 2:", (2*layer-1)*(totalblocks-blocks))
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q08_p3.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		data, _ := os.ReadFile(infile)
		blocks := 202400000
		hpriests, _ := strconv.Atoi(string(data))
		hacolytes := 10
		thickness := 1
		totalblocks := 1
		emptyblocks := 0
		columnHeights := []int{1}
		layer := 1
		for totalblocks-emptyblocks < blocks {
			layer++
			thickness = (thickness*hpriests)%hacolytes + hacolytes
			nColumnHeights := []int{thickness}
			for _, col := range columnHeights {
				nColumnHeights = append(nColumnHeights, col+thickness)
			}
			nColumnHeights = append(nColumnHeights, thickness)
			columnHeights = nColumnHeights
			totalblocks = 0
			emptyblocks = 0
			for idx, ncol := range columnHeights {
				totalblocks += ncol
				if idx == 0 || idx+1 == len(columnHeights) {
					continue
				}
				emptyblocks += (hpriests * len(columnHeights) * ncol) % hacolytes
			}
		}

		fmt.Println("Part 3:", (totalblocks - emptyblocks - blocks))
	}
}
