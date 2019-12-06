package main

import (
	"AdventOfCode/lib/intcode"
	"bufio"
	"fmt"
	"log"
	"os"
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
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	icodes := [][]int{}
	for scanner.Scan() {
		// string to int
		icode := intcodeArray(scanner.Text())
		icodes = append(icodes, icode)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
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
