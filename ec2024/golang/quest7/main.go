package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
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
			infile = "../input/everybody_codes_e2024_q07_p1.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		file, _ := os.Open(infile)
		scanner := bufio.NewScanner(file)
		planScores := make(map[string]int)
		plans := make([]string, 0)
		for scanner.Scan() {
			splitLine := strings.Split(scanner.Text(), ":")
			splitDests := strings.Split(splitLine[1], ",")
			planName := splitLine[0]
			total := 0
			power := 10
			for idx := range 10 {
				switch splitDests[idx%len(splitDests)] {
				case "=":
					power += 0
				case "+":
					power += 1
				case "-":
					power -= 1
				default:
					log.Fatal("Bad spec", splitDests[idx%len(splitDests)])
				}
				total += power
			}
			planScores[planName] = total
			plans = append(plans, planName)
		}
		slices.SortFunc(plans, func(plan1 string, plan2 string) int {
			return planScores[plan2] - planScores[plan1]
		})
		ans := strings.Join(plans, "")
		fmt.Println("Part 1:", ans)
	}
	{
		raceTrack :=
			`S-=++=-==++=++=-=+=-=+=+=--=-=++=-==++=-+=-=+=-=+=+=++=-+==++=++=-=-=--
-                                                                     -
=                                                                     =
+                                                                     +
=                                                                     +
+                                                                     =
=                                                                     =
-                                                                     -
--==++++==+=+++-=+=-=+=-+-=+-=+-=+=-=+=--=+++=++=+++==++==--=+=++==+++-`
		raceTrackLines := strings.Split(raceTrack, "\n")
		raceTrackFlat := raceTrackLines[0]
		for _, line := range raceTrackLines[1 : len(raceTrackLines)-1] {
			raceTrackFlat += line[len(line)-1:]
		}
		raceTrackB := []byte(raceTrackLines[len(raceTrackLines)-1])
		slices.Reverse(raceTrackB)
		raceTrackFlat += string(raceTrackB)
		for idx := len(raceTrackLines) - 2; idx > 0; idx-- {
			raceTrackFlat += raceTrackLines[idx][0:1]
		}
		raceTrackFlat = raceTrackFlat[1:] + "S"
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q07_p2.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		file, _ := os.Open(infile)
		scanner := bufio.NewScanner(file)
		planScores := make(map[string]int)
		plans := make([]string, 0)
		for scanner.Scan() {
			splitLine := strings.Split(scanner.Text(), ":")
			splitDests := strings.Split(splitLine[1], ",")
			planName := splitLine[0]
			total := 0
			power := 10
			specIdx := 0
			for range 10 {
				for _, trackval := range raceTrackFlat {
					switch splitDests[specIdx] + string(trackval) {
					case "==", "=S":
						power += 0
					case "+=", "+S", "-+", "=+", "++":
						power += 1
					case "-=", "-S", "--", "=-", "+-":
						power -= 1
					default:
						log.Fatal("Bad spec", splitDests[specIdx]+string(trackval))
					}
					total += power
					specIdx += 1
					specIdx %= len(splitDests)
				}
			}
			planScores[planName] = total
			plans = append(plans, planName)
		}
		slices.SortFunc(plans, func(plan1 string, plan2 string) int {
			return planScores[plan2] - planScores[plan1]
		})
		ans := strings.Join(plans, "")
		fmt.Println("Part 2:", ans)
	}
	{
		raceTrack :=
			`S+= +=-== +=++=     =+=+=--=    =-= ++=     +=-  =+=++=-+==+ =++=-=-=--
- + +   + =   =     =      =   == = - -     - =  =         =-=        -
= + + +-- =-= ==-==-= --++ +  == == = +     - =  =    ==++=    =++=-=++
+ + + =     +         =  + + == == ++ =     = =  ==   =   = =++=
= = + + +== +==     =++ == =+=  =  +  +==-=++ =   =++ --= + =
+ ==- = + =   = =+= =   =       ++--          +     =   = = =--= ==++==
=     ==- ==+-- = = = ++= +=--      ==+ ==--= +--+=-= ==- ==   =+=    =
-               = = = =   +  +  ==+ = = +   =        ++    =          -
-               = + + =   +  -  = + = = +   =        +     =          -
--==++++==+=+++-= =-= =-+-=  =+-= =-= =--   +=++=+++==     -=+=++==+++-`
		raceTrackMap := make(map[Pos]rune)
		for rowI, line := range strings.Split(raceTrack, "\n") {
			for colI, bval := range line {
				raceTrackMap[Pos{rowI, colI}] = bval
			}
		}
		raceTrackFlat := ""
		here := Pos{0, 1}
		looking := Pos{0, 1}
		for raceTrackMap[here] != 'S' {
			raceTrackFlat += string(raceTrackMap[here])
			glance0 := raceTrackMap[here.Add(looking)]
			glance1 := raceTrackMap[here.Add(Pos{looking.y, -looking.x})]
			glance2 := raceTrackMap[here.Add(Pos{-looking.y, looking.x})]
			if strings.Contains("+-=S", string(glance0)) {
				here = here.Add(looking)
			} else if strings.Contains("+-=S", string(glance1)) {
				looking = Pos{looking.y, -looking.x}
				here = here.Add(looking)
			} else if strings.Contains("+-=S", string(glance2)) {
				looking = Pos{-looking.y, looking.x}
				here = here.Add(looking)
			}
		}
		raceTrackFlat += "S"

		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = "../input/everybody_codes_e2024_q07_p3.txt"
		} else {
			infile = argsWithoutProg[0]
		}
		file, _ := os.Open(infile)
		scanner := bufio.NewScanner(file)
		planScores := make(map[string]int)
		plans := make([]string, 0)
		for scanner.Scan() {
			splitLine := strings.Split(scanner.Text(), ":")
			splitDests := strings.Split(splitLine[1], ",")
			planName := splitLine[0]
			total := 0
			power := 10
			specIdx := 0
			for range 11 {
				for _, trackval := range raceTrackFlat {
					switch splitDests[specIdx] + string(trackval) {
					case "==", "=S":
						power += 0
					case "+=", "+S", "-+", "=+", "++":
						power += 1
					case "-=", "-S", "--", "=-", "+-":
						power -= 1
					default:
						log.Fatal("Bad spec", splitDests[specIdx]+string(trackval))
					}
					total += power
					specIdx += 1
					specIdx %= len(splitDests)
				}
			}
			planScores[planName] = total
			plans = append(plans, planName)
		}
		otherKnight := planScores[plans[0]]

		actionPlan := "+++++---==="
		winCount := 0
		for actionPlan != "" {
			total := 0
			power := 10
			specIdx := 0
			for range 11 {
				for _, trackval := range raceTrackFlat {
					switch actionPlan[specIdx:specIdx+1] + string(trackval) {
					case "==", "=S":
						power += 0
					case "+=", "+S", "-+", "=+", "++":
						power += 1
					case "-=", "-S", "--", "=-", "+-":
						power -= 1
					default:
						log.Fatal("Bad spec", actionPlan[specIdx:specIdx+1]+string(trackval))
					}
					total += power
					specIdx += 1
					specIdx %= len(actionPlan)
				}
			}
			if total > otherKnight {
				winCount++
			}
			actionPlan = nextActionPlan(actionPlan)
		}
		fmt.Println("Part 3:", winCount)
	}
}

func nextActionPlan(actionPlan string) string {
	// viewing actionPlan as a number, sort of, with '+' < '-' < '='

	actionPlanSB := []byte(actionPlan)

	// possible values at "spot" given what's to the left of spot
	possibilities := func(spot int) []byte {
		gots := make(map[byte]int)
		for _, v := range actionPlanSB[:spot] {
			gots[v] += 1
		}
		retval := make([]byte, 0, 3)
		if gots['-'] < 3 {
			retval = append(retval, '-')
		}
		if gots['='] < 3 {
			retval = append(retval, '=')
		}
		if gots['+'] < 5 {
			retval = append(retval, '+')
		}
		return retval
	}

	// main algorithm
	idxUp := len(actionPlanSB) - 1
	for idxUp >= 0 && !slices.ContainsFunc(possibilities(idxUp), func(x byte) bool { return x > actionPlanSB[idxUp] }) {
		idxUp--
	}
	if idxUp < 0 {
		return ""
	}
	p := possibilities(idxUp)
	actionPlanSB[idxUp] += 1
	for !slices.Contains(p, actionPlanSB[idxUp]) {
		actionPlanSB[idxUp] += 1
	}
	for idx := idxUp + 1; idx < len(actionPlanSB); idx++ {
		actionPlanSB[idx] = slices.Min(possibilities(idx))
	}

	return string(actionPlanSB)
}
