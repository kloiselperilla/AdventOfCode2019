package intcode

import (
	"strconv"
	"strings"
)

const (
	positionModeCode  = 0
	immediateModeCode = 1
	relativeModeCode  = 2
)

func opcodeParse(opVal int) (int, []int) {
	opcode := opVal % 100
	opVal /= 100

	modes := make([]int, 5)
	for i := 0; opVal > 0; i++ {
		modes[i] = opVal % 10
		opVal /= 10
	}

	return opcode, modes[:]
}

func parameterIndex(eng *Engine, paramNum int) *int {
	var index *int
	value := eng.Code[eng.Ip+paramNum+1]
	switch eng.Modes[paramNum] {
	case positionModeCode:
		index = &eng.Code[value]
	case immediateModeCode:
		index = &value
	case relativeModeCode:
		index = &eng.Code[eng.RelBase+value]
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
		if s == "-" {
			break
		}
		val, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		intArr = append(intArr, val)
	}
	intArrCap := make([]int, len(intArr)*20)
	copy(intArrCap, intArr)
	return intArrCap
}
