package intcode

import (
	"strconv"
	"strings"
)

const (
	addCode         = 1
	multCode        = 2
	inputCode       = 3
	outputCode      = 4
	jumpIfTrueCode  = 5
	jumpIfFalseCode = 6
	lessThanCode    = 7
	equalsCode      = 8
	stopCode        = 99

	positionModeCode  = 0
	immediateModeCode = 1
)

func opcodeParse(opVal int) (int, []int) {
	opcode := opVal % 100
	opVal /= 100

	modes := [5]int{}
	for i := 0; opVal > 0; i++ {
		modes[i] = opVal % 10
		opVal /= 10
	}

	return opcode, modes[:]
}

func parameterIndex(eng *Engine, paramNum int) int {
	var index int
	value := eng.Code[eng.Ip+paramNum+1]
	switch eng.Modes[paramNum] {
	case positionModeCode:
		index = eng.Code[value]
	case immediateModeCode:
		index = value
	}
	return index
}

// ResetMemory changes the intcode at indexes 1 and 2 with the given noun and
// verb
func ResetMemory(intcode []int, noun int, verb int) []int {
	intcode[1] = noun
	intcode[2] = verb
	return intcode
}

// StringToCode converts a string of intcode to slice
func StringToCode(intcode string) []int {
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
