package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
)

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e2_q02_p1.txt`
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
		if !scanner.Scan() {
			panic("Empty input file")
		}
		balloons := scanner.Text()
		// balloons = "GRBGGGBBBRRRRRRRR"
		balloonStart := 0
		shot := 0
		for ; balloonStart < len(balloons); shot++ {
			var shotColor byte
			switch shot % 3 {
			case 0:
				shotColor = 'R'
			case 1:
				shotColor = 'G'
			case 2:
				shotColor = 'B'
			}
			for (balloonStart < len(balloons)) && (balloons[balloonStart] == shotColor) {
				balloonStart++
			}
			balloonStart++
		}
		fmt.Println("p1", shot)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e2_q02_p2.txt`
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
		if !scanner.Scan() {
			panic("Empty input file")
		}
		balloons := scanner.Text()
		// balloons = "BBRGGRRGBBRGGBRGBBRRBRRRBGGRRRBGBGG"
		multiplicity := 100
		circle := make([]byte, 0, multiplicity*len(balloons))
		for range multiplicity {
			circle = append(circle, []byte(balloons)...)
		}
		shot := 0
		for ; len(circle) > 0; shot++ {
			cLen := len(circle)
			var shotColor byte
			switch shot % 3 {
			case 0:
				shotColor = 'R'
			case 1:
				shotColor = 'G'
			case 2:
				shotColor = 'B'
			}
			if circle[0] == shotColor && (cLen%2 == 0) {
				circle = slices.Concat(circle[1:cLen/2], circle[cLen/2+1:])
			} else {
				circle = circle[1:]
			}
			// fmt.Println("After shot", shot+1, "have", len(circle), "balloons")
		}
		fmt.Println("p2", shot)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e2_q02_p3.txt`
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
		if !scanner.Scan() {
			panic("Empty input file")
		}
		balloons := scanner.Text()
		// balloons = "BBRGGRRGBBRGGBRGBBRRBRRRBGGRRRBGBGG"
		multiplicity := 100000
		circle := make([]byte, 0, multiplicity*len(balloons))
		for range multiplicity {
			circle = append(circle, []byte(balloons)...)
		}
		shot := 0
		for len(circle) > 0 {
			cLen := len(circle)
			if cLen > 300 {
				origCLen := cLen
				origShot := shot
				// going to shoot directly from 0 to cLen/2, but if cLen is even, don't shoot cLen/2
				// so shoot balloon 0 and (cLen-1)/2, and everything in between
				// so shoot (cLen+1)/2 times
				onesShotOnOtherSide := 0
				for range (cLen + 1) / 2 {
					var shotColor byte
					switch shot % 3 {
					case 0:
						shotColor = 'R'
					case 1:
						shotColor = 'G'
					case 2:
						shotColor = 'B'
					}
					if circle[shot-origShot] == shotColor && (cLen%2 == 0) {
						// so the thing at circle[origCLen-1] is the end, and we want
						// the thing at circle[origCLen-(cLen/2)] to pop
						circle[origCLen-(cLen/2)] = 'X'
						onesShotOnOtherSide++
						cLen--
					}
					circle[shot-origShot] = 'X'
					shot++
					cLen--
				}
				circleNew := make([]byte, 0, origCLen/2-onesShotOnOtherSide+1)
				for idx := origCLen / 2; idx < origCLen; idx++ {
					b := circle[idx]
					if b != 'X' {
						circleNew = append(circleNew, b)
					}
				}
				circle = circleNew
			} else {
				var shotColor byte
				switch shot % 3 {
				case 0:
					shotColor = 'R'
				case 1:
					shotColor = 'G'
				case 2:
					shotColor = 'B'
				}
				if circle[0] == shotColor && (cLen%2 == 0) {
					circle = slices.Concat(circle[1:cLen/2], circle[cLen/2+1:])
				} else {
					circle = circle[1:]
				}
				shot++
			}
		}
		fmt.Println("p3", shot)
	}

}
