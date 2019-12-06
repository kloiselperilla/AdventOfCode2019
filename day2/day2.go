package day2

import (
	//"bufio"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// FEEDBACK: Define the opCodes as constants to make them more readable.

func intcodeArray(intcode string) []int {
	var intArr = []int{}
	for _, s := range strings.Split(strings.TrimSpace(intcode), ",") {
		val, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		intArr = append(intArr, val)
	}
	return intArr
}

func readIntcode(path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(dat)
}

func add(pos int, intcode []int) []int {
	sum := intcode[intcode[pos+1]] + intcode[intcode[pos+2]]
	intcode[intcode[pos+3]] = sum
	return intcode
}

func mult(pos int, intcode []int) []int {
	prod := intcode[intcode[pos+1]] * intcode[intcode[pos+2]]
	intcode[intcode[pos+3]] = prod
	return intcode
}

func EvaluateIntcode(intcode []int) []int {
	ix := 0
loop:
	for ix < len(intcode) {
		switch op := intcode[ix]; {
		case op == 1:
			intcode = add(ix, intcode)
		case op == 2:
			intcode = mult(ix, intcode)
		case op == 99:
			break loop
		default:
			fmt.Println("Something went wrong, invalid opcode")
			fmt.Println("ix", ix)
			fmt.Println("op", op)
		}
		ix += 4
	}
	return intcode
}

func ResetMemory(intcode []int, noun int, verb int) []int {
	intcode[1] = noun
	intcode[2] = verb
	return intcode
}
func findNounVerb(intcode []int, goal int) int {
	candidate := make([]int, len(intcode))
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			copy(candidate, intcode)
			candidate = ResetMemory(candidate, noun, verb)

			result := EvaluateIntcode(candidate)
			if result[0] == goal {
				return 100*noun + verb
			}
		}
	}
	fmt.Println("Reached the end")
	return -1
}

func main() {
	ic := readIntcode("input")
	arr := intcodeArray(ic)
	pt1 := make([]int, len(arr))
	copy(pt1, arr)
	pt1 = ResetMemory(pt1, 12, 2)
	pt1 = EvaluateIntcode(pt1)

	fmt.Println("Part 1:")
	fmt.Println(pt1[0])

	fmt.Println("Part 2:")
	pt2 := make([]int, len(arr))
	copy(pt2, arr)
	fmt.Println(findNounVerb(pt2, 19690720))

}
