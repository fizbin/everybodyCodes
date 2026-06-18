package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var input1 = flag.String("p1", `..\input\everybody_codes_e2025_q18_p1.txt`, "the input for part 1")
var input2 = flag.String("p2", `..\input\everybody_codes_e2025_q18_p2.txt`, "the input for part 2")
var input3 = flag.String("p3", `..\input\everybody_codes_e2025_q18_p3.txt`, "the input for part 3")

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
	scanner := bufio.NewScanner(bytes.NewReader(data))
	plantEnergies := make(map[string]int)
	plantLineParser := regexp.MustCompile(`(Plant \d+) with thickness (\d+):`)
	branchParser := regexp.MustCompile(`- (?:free branch with thickness (\d+)|branch to (Plant \d+) with thickness (\d+))`)
	plantName := ""
	for scanner.Scan() {
		plantLineParts := plantLineParser.FindStringSubmatch(scanner.Text())
		if len(plantLineParts) == 0 {
			panic("Bad Plant line: " + scanner.Text())
		}
		plantName = plantLineParts[1]
		plantThickness, _ := strconv.Atoi(plantLineParts[2])
		incomingEnergy := 0
		for scanner.Scan() {
			if scanner.Text() == "" {
				break
			}
			parsed := branchParser.FindStringSubmatch(scanner.Text())
			if parsed[1] != "" {
				thick, _ := strconv.Atoi(parsed[1])
				incomingEnergy += thick
			} else {
				incomingSrc := parsed[2]
				thick, _ := strconv.Atoi(parsed[3])
				incomingEnergy += plantEnergies[incomingSrc] * thick
			}
		}
		if plantThickness <= incomingEnergy {
			plantEnergies[plantName] = incomingEnergy
		} else {
			plantEnergies[plantName] = 0
		}
	}
	return plantEnergies[plantName]
}

type plantAction func(map[string]int, map[int]bool)

func doProblem2(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	actions, lastPlantName, _ := parseActions(scanner)
	sum := 0
	for scanner.Scan() {
		vals := strings.Split(scanner.Text(), " ")
		testCases := make(map[int]bool)
		for idx, v := range vals {
			testCases[idx] = (v == "1")
		}

		plantEnergies := make(map[string]int)
		for _, act := range actions {
			act(plantEnergies, testCases)
		}
		sum += plantEnergies[lastPlantName]
	}
	return sum
}

func parseActions(scanner *bufio.Scanner) ([]plantAction, string, map[int]int) {
	plantLineParser := regexp.MustCompile(`(Plant \d+) with thickness (\d+):`)
	branchParser := regexp.MustCompile(`- (?:free branch with thickness (-?\d+)|branch to (Plant \d+) with thickness (-?\d+))`)
	actions := make([]plantAction, 0)
	plantName := ""
	inputNumber := 0
	plantNames := make(map[string]int)
	inputCharacters := make(map[int]int) // 0 - nothing; 1 - always positive; 2 - always negative; 3 - both
	for scanner.Scan() {
		plantLineParts := plantLineParser.FindStringSubmatch(scanner.Text())
		if len(plantLineParts) == 0 {
			if scanner.Text() == "" {
				break // finished plants, now read testcases
			}
			panic("Bad Plant line: " + scanner.Text())
		}
		plantName = plantLineParts[1]
		plantThickness, _ := strconv.Atoi(plantLineParts[2])
		for scanner.Scan() {
			if scanner.Text() == "" {
				break
			}
			parsed := branchParser.FindStringSubmatch(scanner.Text())
			// fmt.Print("scanner.Text: '", scanner.Text(), "'")
			// fmt.Println()
			var action plantAction
			if parsed[1] != "" {
				myInputNumber := inputNumber
				myDstName := plantName
				thick, _ := strconv.Atoi(parsed[1])
				action = func(plantEnergies map[string]int, testCases map[int]bool) {
					if testCases[myInputNumber] {
						plantEnergies[myDstName] += thick
					}
				}
				plantNames[myDstName] = myInputNumber
				inputNumber++
			} else {
				myDstName := plantName
				incomingSrc := parsed[2]
				thick, _ := strconv.Atoi(parsed[3])
				action = func(plantEnergies map[string]int, testCases map[int]bool) {
					plantEnergies[myDstName] += plantEnergies[incomingSrc] * thick
				}
				if thick > 0 {
					inputCharacters[plantNames[incomingSrc]] |= 1
				}
				if thick < 0 {
					inputCharacters[plantNames[incomingSrc]] |= 2
				}
			}
			actions = append(actions, action)
		}
		myDstName := plantName
		actions = append(actions, func(plantEnergies map[string]int, testCases map[int]bool) {
			if plantEnergies[myDstName] < plantThickness {
				plantEnergies[myDstName] = 0
			}
		})
	}
	return actions, plantName, inputCharacters
}

func doProblem3(data []byte) any {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	actions, lastPlantName, inputCharacters := parseActions(scanner)
	foundValues := make([]int, 0)
	nInputs := 0
	maxLastPlantEnergy := 0
	for scanner.Scan() {
		vals := strings.Split(scanner.Text(), " ")
		testCases := make(map[int]bool)
		nInputs = max(nInputs, len(vals))
		for idx, v := range vals {
			if v == "1" {
				testCases[idx] = true
			}
		}

		plantEnergies := make(map[string]int)
		for _, act := range actions {
			act(plantEnergies, testCases)
		}
		if plantEnergies[lastPlantName] != 0 {
			foundValues = append(foundValues, plantEnergies[lastPlantName])
			maxLastPlantEnergy = max(maxLastPlantEnergy, plantEnergies[lastPlantName])
			// fmt.Println("Found", plantEnergies[lastPlantName])
		}
	}

	{
		testCases := make(map[int]bool)
		for idx, val := range inputCharacters {
			testCases[idx] = ((val & 2) == 0)
		}
		keepGoing := true
		for keepGoing {
			keepGoing = false
			plantEnergies := make(map[string]int)
			for _, act := range actions {
				act(plantEnergies, testCases)
			}
			maxLastPlantEnergy = max(maxLastPlantEnergy, plantEnergies[lastPlantName])
			// fmt.Println("Found2:", plantEnergies[lastPlantName])
			for upIdx := range nInputs {
				if inputCharacters[upIdx] == 3 && !testCases[upIdx] {
					keepGoing = true
					testCases[upIdx] = true
					for downIdx := upIdx - 1; downIdx >= 0; downIdx-- {
						if inputCharacters[downIdx] == 3 {
							testCases[downIdx] = false
						}
					}
					break
				}
			}
		}
	}

	// fmt.Println("max is", maxLastPlantEnergy)
	sum := 0
	for _, val := range foundValues {
		sum += maxLastPlantEnergy - val
	}
	return sum
}
