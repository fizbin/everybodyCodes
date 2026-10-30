package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2026_q04_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2026_q04_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2026_q04_p3.txt`, "the input for part 3")

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

func doProblem1(data []byte) any {
	ans := 0
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`id=(\d+) quest=(\d+) part=(\d+) time=(\d+)`)
	times := make(map[string]int)
	for scanner.Scan() {
		matched := prsr.FindStringSubmatch(scanner.Text())
		timeNum, _ := strconv.Atoi(matched[4])
		times[matched[1]] += timeNum
	}
	idvals := make([]string, 0, len(times))
	for idval := range times {
		idvals = append(idvals, idval)
	}
	slices.SortFunc(idvals, func(a, b string) int {
		return times[a] - times[b]
	})
	for idx, idval := range idvals {
		ans += (idx + 1) * times[idval]
	}
	return ans
}
func doProblem2(data []byte) any {
	ans := 0
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`id=(\d+) quest=(\d+) part=(\d+) time=(\d+)`)
	times := make(map[string]int)
	parts := make(map[string]int)
	for scanner.Scan() {
		matched := prsr.FindStringSubmatch(scanner.Text())
		timeNum, _ := strconv.Atoi(matched[4])
		times[matched[1]] += timeNum
		parts[matched[1]] += 1
	}
	idvals := make([]string, 0, len(times))
	for idval := range times {
		idvals = append(idvals, idval)
	}
	slices.SortFunc(idvals, func(a, b string) int {
		if parts[a] != parts[b] {
			return parts[b] - parts[a]
		}
		return times[a] - times[b]
	})
	for idx, idval := range idvals {
		ans += (idx + 1) * times[idval]
	}
	return ans
}
func doProblem3(data []byte) any {
	ans := 0
	// data = []byte(`id=2 quest=1 part=1 time=90
	// id=4 quest=2 part=1 time=150
	// id=1 quest=2 part=1 time=130
	// id=1 quest=1 part=2 time=180
	// id=2 quest=2 part=1 time=130
	// id=4 quest=1 part=1 time=110
	// id=1 quest=2 part=2 time=210
	// id=2 quest=2 part=2 time=220
	// id=1 quest=2 part=3 time=330
	// id=2 quest=1 part=2 time=170
	// id=3 quest=1 part=1 time=90
	// id=4 quest=1 part=2 time=170
	// id=3 quest=2 part=1 time=120
	// id=1 quest=1 part=1 time=120
	// id=3 quest=2 part=2 time=220
	// id=2 quest=2 part=3 time=350
	// id=3 quest=1 part=2 time=170`)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	prsr := regexp.MustCompile(`id=(\d+) quest=(\d+) part=(\d+) time=(\d+)`)
	times := make(map[int]map[int]map[string]int) // quest -> part -> id -> time
	idvalMap := make(map[string]bool)
	for scanner.Scan() {
		matched := prsr.FindStringSubmatch(scanner.Text())
		questNum, _ := strconv.Atoi(matched[2])
		partNum, _ := strconv.Atoi(matched[3])
		timeNum, _ := strconv.Atoi(matched[4])
		if times[questNum] == nil {
			times[questNum] = make(map[int]map[string]int)
		}
		if times[questNum][partNum] == nil {
			times[questNum][partNum] = make(map[string]int)
		}
		times[questNum][partNum][matched[1]] = timeNum
		idvalMap[matched[1]] = true
	}
	scores := scoreQuest(times)
	idvals := make([]string, 0, len(times))
	for idval := range idvalMap {
		idvals = append(idvals, idval)
	}
	slices.SortFunc(idvals, func(a, b string) int {
		return scores[b] - scores[a]
	})
	pointRank := make(map[int]int)
	for idx, idval := range idvals {
		if pointRank[scores[idval]] == 0 {
			pointRank[scores[idval]] = idx + 1
		}
		ans += pointRank[scores[idval]] * scores[idval]
	}
	return ans
}

func scoreQuest(times map[int]map[int]map[string]int) map[string]int {
	ans := make(map[string]int)
	for quest, qMap := range times {
		miniAns := make(map[string]int)
		p1TimeToScore := make(map[int]int)
		p1Times := make([]int, 0)
		p1Map := qMap[1]
		for _, time := range p1Map {
			p1Times = append(p1Times, time)
		}
		slices.Sort(p1Times)
		for idx, time := range p1Times {
			p1TimeToScore[time] = max(p1TimeToScore[time], 50-idx)
		}
		for idval, time := range p1Map {
			miniAns[idval] += p1TimeToScore[time]
		}
		if quest < 0 {
			stred := ""
			for idval, score := range miniAns {
				if score > 0 {
					stred += fmt.Sprintf(" %s:%d", idval, score)
				}
			}
			fmt.Println("For quest", quest, "part 1 looks like [", stred, "]")
		}

		p2TimeToScore := make(map[int]int)
		p2Times := make([]int, 0)
		p2Map := qMap[2]
		for idval, time2 := range p2Map {
			time1 := p1Map[idval]
			if time1 == 0 {
				panic("wait, part2")
			}
			p2Times = append(p2Times, time1+time2)
		}
		slices.Sort(p2Times)
		for idx, time := range p2Times {
			p2TimeToScore[time] = max(p2TimeToScore[time], 100-idx)
		}
		for idval, time2 := range p2Map {
			time1 := p1Map[idval]
			miniAns[idval] += p2TimeToScore[time1+time2]
		}

		p3TimeToScore := make(map[int]int)
		p3Times := make([]int, 0)
		p3Map := qMap[3]
		for idval, time3 := range p3Map {
			time1 := p1Map[idval]
			if time1 == 0 {
				panic("wait, part3")
			}
			time2 := p2Map[idval]
			if time2 == 0 {
				panic("wait, part3, 2")
			}
			p3Times = append(p3Times, time1+time2+time3)
		}
		slices.Sort(p3Times)
		for idx, time := range p3Times {
			p3TimeToScore[time] = max(p3TimeToScore[time], 150-idx)
		}
		for idval, time3 := range p3Map {
			time1 := p1Map[idval]
			time2 := p2Map[idval]
			miniAns[idval] += p3TimeToScore[time1+time2+time3]
		}
		if quest < 0 {
			fmt.Println("For quest ", quest, " answers are ", miniAns)
		}
		for idval, score := range miniAns {
			ans[idval] += score
		}
	}
	return ans
}
