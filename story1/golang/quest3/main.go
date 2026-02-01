package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
)

func extendedGCD(a, b int64) (g, m, n int64) {
	if a < 0 {
		g, m, n = extendedGCD(-a, b)
		m = -m
		return
	}
	if b < 0 {
		g, m, n = extendedGCD(a, -b)
		n = -n
		return
	}
	if a > b {
		g, n, m = extendedGCD(b, a)
		return
	}
	if a == 0 {
		return b, 0, 1
	}
	c := b % a
	g, m1, n1 := extendedGCD(c, a)
	// m1 * c + n1 * a == g
	// m1 * (b - a*(b/a)) + n1*a == g
	// m1 * b + (n1 - m1*(b/a))*a == g
	return g, (n1 - m1*(b/a)), m1
}

func doCRTCalc(off1, mod1, off2, mod2 int64) (off3, mod3 int64) {
	g, m, n := extendedGCD(mod1, mod2)
	// m*mod1 + n*mod2 == g
	if (off1-off2)%g != 0 {
		panic("Can't align")
	}
	if g != 1 {
		fmt.Println("Interesting! ", mod1, mod2)
	}
	// Need to do intermediate calculations with math/big
	var mod3w, off3w, scratch1, scratch2 big.Int
	//mod3 = (mod1 / g) * mod2
	mod3w.SetInt64(mod1 / g)
	mod3w.Mul(&mod3w, scratch1.SetInt64(mod2))
	//off3 = (off1/g)*n*mod2 + (off2/g)*m*mod1
	off3w.SetInt64(off1 / g)
	off3w.Mul(&off3w, scratch1.SetInt64(n))
	off3w.Mul(&off3w, scratch1.SetInt64(mod2))
	scratch1.SetInt64(off2 / g)
	scratch1.Mul(&scratch1, scratch2.SetInt64(m))
	scratch1.Mul(&scratch1, scratch2.SetInt64(mod1))
	off3w.Add(&off3w, &scratch1)
	// off3 += off1 % g
	// for off3 < 0 {
	// 	off3 += mod3
	// }
	// off3 %= mod3
	off3w.Add(&off3w, scratch1.SetInt64(off1%g))
	off3w.Mod(&off3w, &mod3w)
	if !mod3w.IsInt64() {
		panic("Non i64 mod")
	}
	if !off3w.IsInt64() {
		panic("Non i64 offset")
	}
	mod3 = mod3w.Int64()
	off3 = off3w.Int64()
	return
}

func main() {
	flag.Parse()
	argsWithoutProg := flag.Args()
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e1_q03_p1.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		parser := regexp.MustCompile(`x=(\d+) y=(\d+)`)
		scanner := bufio.NewScanner(file)
		retval := int64(0)
		for scanner.Scan() {
			match := parser.FindStringSubmatch(scanner.Text())
			x, _ := strconv.ParseInt(match[1], 10, 64)
			y, _ := strconv.ParseInt(match[2], 10, 64)
			mod := x + y - 1
			finalx := (x+99)%mod + 1
			finaly := mod + 1 - finalx
			retval += finalx + 100*finaly
		}
		fmt.Println(retval)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e1_q03_p2.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		parser := regexp.MustCompile(`x=(\d+) y=(\d+)`)
		scanner := bufio.NewScanner(file)
		totaloff := int64(0)
		totalmod := int64(1)
		for scanner.Scan() {
			match := parser.FindStringSubmatch(scanner.Text())
			x, _ := strconv.ParseInt(match[1], 10, 64)
			y, _ := strconv.ParseInt(match[2], 10, 64)
			mod := x + y - 1
			off := y - 1
			totaloff, totalmod = doCRTCalc(off, mod, totaloff, totalmod)
		}
		fmt.Println(totaloff, "out of", totalmod)
	}
	{
		var infile string
		if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "" {
			infile = `..\input\everybody_codes_e1_q03_p3.txt`
		} else {
			infile = argsWithoutProg[0]
		}
		file, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			panic("Couldn't open input file")
		}
		defer file.Close()
		parser := regexp.MustCompile(`x=(\d+) y=(\d+)`)
		scanner := bufio.NewScanner(file)
		totaloff := int64(0)
		totalmod := int64(1)
		for scanner.Scan() {
			match := parser.FindStringSubmatch(scanner.Text())
			x, _ := strconv.ParseInt(match[1], 10, 64)
			y, _ := strconv.ParseInt(match[2], 10, 64)
			mod := x + y - 1
			off := y - 1
			totaloff, totalmod = doCRTCalc(off, mod, totaloff, totalmod)
		}
		fmt.Println(totaloff, "out of", totalmod)
	}
}
