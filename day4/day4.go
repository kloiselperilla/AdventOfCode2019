package main

import (
	"fmt"
	"runtime"
	"strconv"
)

const threadRange = 10000

func checkRepeating(number int) bool {
	prev := -1
	for number > 0 {
		curr := number % 10
		if curr == prev {
			return true
		}
		prev = curr
		number /= 10
	}
	return false
}

func checkMonotonicIncr(number int) bool {
	prev := 10
	for number > 0 {
		curr := number % 10
		if curr > prev {
			return false
		}
		prev = curr
		number /= 10
	}
	return true
}

func validPass(candidate int, digLen int, min int, max int) bool {
	valid := true
	// In range
	if candidate < min || candidate > max {
		valid = false
	}
	// Correct length
	if len(strconv.Itoa(candidate)) != digLen {
		valid = false
	}
	// Does not decrease
	if !checkMonotonicIncr(candidate) {
		valid = false
	}
	// Has two repeating
	if !checkRepeating(candidate) {
		valid = false
	}
	return valid
}

func checkRangeForValid(min int, max int, digLen int, passMin int, passMax int, validChan chan []int) {
	validNums := []int{}
	for i := min; i <= max; i++ {
		if validPass(i, digLen, passMin, passMax) {
			validNums = append(validNums, i)
		}
	}
	validChan <- validNums
}

func findValidPasswords(digLen int, passMin int, passMax int) []int {
	runtime.GOMAXPROCS(12)
	validPasswords := []int{}

	validChan := make(chan []int)
	for i := passMin; i <= passMax; i += threadRange {
		rangeMax := i + threadRange - 1
		if rangeMax > passMax {
			rangeMax = passMax
		}
		go checkRangeForValid(i, rangeMax, digLen, passMin, passMax, validChan)
	}

	for i := passMin; i <= passMax; i += threadRange {
		validPasswords = append(validPasswords, <-validChan...)
	}

	return validPasswords
}

///////////////////
// Part 2 Functions
///////////////////

func revisedCheckRepeating(number int) bool {
	repeats := make(map[int]int)

	prev := -1
	for number > 0 {
		curr := number % 10
		if curr == prev {
			if _, ok := repeats[curr]; ok {
				repeats[curr]++
			} else {
				repeats[curr] = 1
			}
		}
		prev = curr
		number /= 10
	}
	for k := range repeats {
		if repeats[k] == 1 {
			return true
		}
	}
	return false
}

func revisedValidPass(candidate int, digLen int, min int, max int) bool {
	valid := true
	// In range
	if candidate < min || candidate > max {
		valid = false
	}
	// Correct length
	if len(strconv.Itoa(candidate)) != digLen {
		valid = false
	}
	// Does not decrease
	if !checkMonotonicIncr(candidate) {
		valid = false
	}
	// Has two repeating
	if !revisedCheckRepeating(candidate) {
		valid = false
	}
	return valid
}

func revisedCheckRangeForValid(min int, max int, digLen int, passMin int, passMax int, validChan chan []int) {
	validNums := []int{}
	for i := min; i <= max; i++ {
		if revisedValidPass(i, digLen, passMin, passMax) {
			validNums = append(validNums, i)
		}
	}
	validChan <- validNums
}

func revisedFindValidPasswords(digLen int, passMin int, passMax int) []int {
	runtime.GOMAXPROCS(12)
	validPasswords := []int{}

	validChan := make(chan []int)
	for i := passMin; i <= passMax; i += threadRange {
		rangeMax := i + threadRange - 1
		if rangeMax > passMax {
			rangeMax = passMax
		}
		go revisedCheckRangeForValid(i, rangeMax, digLen, passMin, passMax, validChan)
	}

	for i := passMin; i <= passMax; i += threadRange {
		validPasswords = append(validPasswords, <-validChan...)
	}

	return validPasswords
}

////////////////////

func main() {
	passMin := 372304
	passMax := 847060
	digLen := 6
	fmt.Println("Part 1:")
	fmt.Println(len(findValidPasswords(digLen, passMin, passMax)))

	fmt.Println("Part 2:")
	fmt.Println(len(revisedFindValidPasswords(digLen, passMin, passMax)))
	fmt.Println(revisedCheckRepeating(111123))
}
