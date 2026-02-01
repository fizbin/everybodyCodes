package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type node struct {
	left    *node
	right   *node
	label   rune
	nodeVal int64
	parent  **node
}

func addNode(tree, newNode *node) (*node, int) {
	if tree == nil {
		return newNode, 0
	}
	var lvl int
	if tree.nodeVal <= newNode.nodeVal {
		tree.right, lvl = addNode(tree.right, newNode)
	} else {
		tree.left, lvl = addNode(tree.left, newNode)
	}
	return tree, (lvl + 1)
}

func addNode2(tree, newNode *node, parentSlot **node) {
	if tree == nil {
		newNode.parent = parentSlot
		*parentSlot = newNode
		return
	}
	if tree.nodeVal <= newNode.nodeVal {
		addNode2(tree.right, newNode, &(tree.right))
	} else {
		addNode2(tree.left, newNode, &(tree.left))
	}
}

func collect(tree *node, lvl int, res map[int]string) {
	if tree != nil {
		res[lvl] += string([]rune{tree.label})
		collect(tree.left, lvl+1, res)
		collect(tree.right, lvl+1, res)
	}
}

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e1_q02_p1.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		parser := regexp.MustCompile(`ADD id=(\d+) left=\[(\d+),(.)\] right=\[(\d+),(.)\]`)
		scanner := bufio.NewScanner(file)
		var leftTree *node
		var rightTree *node
		lvlcountLeft := make(map[int]int)
		lvlcountRight := make(map[int]int)
		for scanner.Scan() {
			match := parser.FindStringSubmatch(scanner.Text())
			// addId, _ := strconv.ParseInt(match[1], 10, 64)
			leftId, _ := strconv.ParseInt(match[2], 10, 64)
			var leftSym rune
			for _, c := range match[3] {
				leftSym = c
				break
			}
			rightId, _ := strconv.ParseInt(match[4], 10, 64)
			var rightSym rune
			for _, c := range match[5] {
				rightSym = c
				break
			}
			var lvl int
			leftTree, lvl = addNode(leftTree, &node{nil, nil, leftSym, leftId, nil})
			lvlcountLeft[lvl]++
			rightTree, lvl = addNode(rightTree, &node{nil, nil, rightSym, rightId, nil})
			lvlcountRight[lvl]++
		}
		maxlevel := 0
		maxlevelVal := 0
		for lvl, lvlVal := range lvlcountLeft {
			if lvlVal > maxlevelVal {
				maxlevel, maxlevelVal = lvl, lvlVal
			}
		}
		fmt.Print(traverseToLevel(leftTree, maxlevel))
		maxlevel = 0
		maxlevelVal = 0
		for lvl, lvlVal := range lvlcountRight {
			if lvlVal > maxlevelVal {
				maxlevel, maxlevelVal = lvl, lvlVal
			}
		}
		fmt.Print(traverseToLevel(rightTree, maxlevel))
		fmt.Println()
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e1_q02_p2.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		parser := regexp.MustCompile(`ADD id=(\d+) left=\[(\d+),(.)\] right=\[(\d+),(.)\]|SWAP (\d+)`)
		scanner := bufio.NewScanner(file)
		var leftTree *node
		var rightTree *node
		lvlCountLeft := make(map[int]int)
		lvlCountRight := make(map[int]int)
		nodeIdLeft := make(map[int](*node))
		nodeIdRight := make(map[int](*node))
		for scanner.Scan() {
			match := parser.FindStringSubmatch(scanner.Text())
			if match[1] != "" {
				addId, _ := strconv.ParseInt(match[1], 10, 64)
				leftId, _ := strconv.ParseInt(match[2], 10, 64)
				var leftSym rune
				for _, c := range match[3] {
					leftSym = c
					break
				}
				rightId, _ := strconv.ParseInt(match[4], 10, 64)
				var rightSym rune
				for _, c := range match[5] {
					rightSym = c
					break
				}
				var lvl int
				newNodeL := node{nil, nil, leftSym, leftId, nil}
				leftTree, lvl = addNode(leftTree, &newNodeL)
				lvlCountLeft[lvl]++
				newNodeR := node{nil, nil, rightSym, rightId, nil}
				rightTree, lvl = addNode(rightTree, &newNodeR)
				lvlCountRight[lvl]++
				nodeIdLeft[int(addId)] = &newNodeL
				nodeIdRight[int(addId)] = &newNodeR
			} else {
				swapIdx, _ := strconv.ParseInt(match[6], 10, 64)
				nodeL := nodeIdLeft[int(swapIdx)]
				nodeR := nodeIdRight[int(swapIdx)]
				nodeL.label, nodeR.label = nodeR.label, nodeL.label
				nodeL.nodeVal, nodeR.nodeVal = nodeR.nodeVal, nodeL.nodeVal
			}
		}
		maxlevel := 0
		maxlevelVal := 0
		for lvl, lvlVal := range lvlCountLeft {
			if lvlVal > maxlevelVal {
				maxlevel, maxlevelVal = lvl, lvlVal
			}
		}
		fmt.Print(traverseToLevel(leftTree, maxlevel))
		maxlevel = 0
		maxlevelVal = 0
		for lvl, lvlVal := range lvlCountRight {
			if lvlVal > maxlevelVal {
				maxlevel, maxlevelVal = lvl, lvlVal
			}
		}
		fmt.Print(traverseToLevel(rightTree, maxlevel))
		fmt.Println()
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e1_q02_p3.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		parser := regexp.MustCompile(`ADD id=(\d+) left=\[(\d+),(.)\] right=\[(\d+),(.)\]|SWAP (\d+)`)
		scanner := bufio.NewScanner(file)
		var leftTree *node
		var rightTree *node
		nodeIdLeft := make(map[int](*node))
		nodeIdRight := make(map[int](*node))
		for scanner.Scan() {
			match := parser.FindStringSubmatch(scanner.Text())
			if match[1] != "" {
				addId, _ := strconv.ParseInt(match[1], 10, 64)
				leftId, _ := strconv.ParseInt(match[2], 10, 64)
				var leftSym rune
				for _, c := range match[3] {
					leftSym = c
					break
				}
				rightId, _ := strconv.ParseInt(match[4], 10, 64)
				var rightSym rune
				for _, c := range match[5] {
					rightSym = c
					break
				}
				newNodeL := node{nil, nil, leftSym, leftId, nil}
				addNode2(leftTree, &newNodeL, &leftTree)
				newNodeR := node{nil, nil, rightSym, rightId, nil}
				addNode2(rightTree, &newNodeR, &rightTree)
				nodeIdLeft[int(addId)] = &newNodeL
				nodeIdRight[int(addId)] = &newNodeR
			} else {
				swapIdx, _ := strconv.ParseInt(match[6], 10, 64)
				nodeL := nodeIdLeft[int(swapIdx)]
				nodeR := nodeIdRight[int(swapIdx)]
				*(nodeL.parent), *(nodeR.parent) = nodeR, nodeL
				nodeL.parent, nodeR.parent = nodeR.parent, nodeL.parent
			}
		}
		messages := make(map[int]string)
		collect(leftTree, 0, messages)
		maxlevel := 0
		maxlevelVal := ""
		for lvl, lvlVal := range messages {
			if len(lvlVal) > len(maxlevelVal) || (len(lvlVal) == len(maxlevelVal) && lvl < maxlevel) {
				maxlevel, maxlevelVal = lvl, lvlVal
			}
		}
		fmt.Print(maxlevelVal)
		messages = make(map[int]string)
		collect(rightTree, 0, messages)
		maxlevel = 0
		maxlevelVal = ""
		for lvl, lvlVal := range messages {
			if len(lvlVal) > len(maxlevelVal) || (len(lvlVal) == len(maxlevelVal) && lvl < maxlevel) {
				maxlevel, maxlevelVal = lvl, lvlVal
			}
		}
		fmt.Print(maxlevelVal)
		fmt.Println()
	}
}

func traverseToLevel(tree *node, level int) string {
	if level == 0 {
		return string([]rune{tree.label})
	}
	var retval string
	if tree.left != nil {
		retval += traverseToLevel(tree.left, level-1)
	}
	if tree.right != nil {
		retval += traverseToLevel(tree.right, level-1)
	}
	return retval
}
