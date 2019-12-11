package main

import (
	"AdventOfCode/lib/files"
	"AdventOfCode/lib/intcode"
	"fmt"
	"strconv"
	"strings"
)

func intcodeArray(intcode string) []int {
	var intArr = []int{}
	for _, s := range strings.Split(strings.TrimSpace(intcode), ",") {
		if s == "-" {
			break
		}
		val, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		intArr = append(intArr, val)
	}
	return intArr
}

func readIntcodes(path string) [][]int {
	lines := files.ReadLines(path)
	icodes := [][]int{}
	for _, line := range lines {
		// string to int
		icode := intcodeArray(line)
		icodes = append(icodes, icode)
	}

	return icodes
}

func main() {
	//icode := []int{3, 0, 4, 0, 99}
	//ans := intcode.EvaluateIntcode(icode, 69)
	//fmt.Println(ans)
	fmt.Println("Part 1:")

	icodes := readIntcodes("day5/input")
	for _, icode := range icodes {
		intcode.EvaluateIntcode(icode, 1)
	}

	fmt.Println()
	fmt.Println("Part 2:")
	icodes2 := readIntcodes("day5/input")
	for _, icode := range icodes2 {
		intcode.EvaluateIntcode(icode, 5)
	}

}
