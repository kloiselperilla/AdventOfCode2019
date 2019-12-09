package intcode

import (
	"fmt"
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

func parameterIndex(mode int, value int, intcode []int) int {
	var index int
	switch mode {
	case positionModeCode:
		index = intcode[value]
	case immediateModeCode:
		index = value
	}
	return index
}

func add(pos int, intcode []int, modes []int) ([]int, int) {
	operand1 := parameterIndex(modes[0], intcode[pos+1], intcode)
	operand2 := parameterIndex(modes[1], intcode[pos+2], intcode)
	sum := operand1 + operand2

	if modes[2] != 0 {
		fmt.Println("No write param should be immediate mode")
	}
	intcode[intcode[pos+3]] = sum
	return intcode, 4
}

func mult(pos int, intcode []int, modes []int) ([]int, int) {
	operand1 := parameterIndex(modes[0], intcode[pos+1], intcode)
	operand2 := parameterIndex(modes[1], intcode[pos+2], intcode)
	prod := operand1 * operand2

	if modes[2] != 0 {
		fmt.Println("No write param should be immediate mode")
	}
	intcode[intcode[pos+3]] = prod
	return intcode, 4
}

func lessThan(pos int, intcode []int, modes []int) ([]int, int) {
	operand1 := parameterIndex(modes[0], intcode[pos+1], intcode)
	operand2 := parameterIndex(modes[1], intcode[pos+2], intcode)
	ans := operand1 < operand2

	if modes[2] != 0 {
		fmt.Println("No write param should be immediate mode")
	}
	var ansInt int
	if ans {
		ansInt = 1
	}
	intcode[intcode[pos+3]] = ansInt
	return intcode, 4
}

func equals(pos int, intcode []int, modes []int) ([]int, int) {
	operand1 := parameterIndex(modes[0], intcode[pos+1], intcode)
	operand2 := parameterIndex(modes[1], intcode[pos+2], intcode)
	ans := operand1 == operand2

	if modes[2] != 0 {
		fmt.Println("No write param should be immediate mode")
	}
	var ansInt int
	if ans {
		ansInt = 1
	}
	intcode[intcode[pos+3]] = ansInt
	return intcode, 4
}

func input(input int, pos int, intcode []int) int {
	intcode[intcode[pos+1]] = input

	return 2
}

func output(pos int, intcode []int, modes []int) (int, int) {
	outVal := parameterIndex(modes[0], intcode[pos+1], intcode)
	//fmt.Println("Output: ", outVal)
	return outVal, 2
}

func jumpIfTrue(pos int, intcode []int, modes []int, ip *int) int {
	operand := parameterIndex(modes[0], intcode[pos+1], intcode)
	ipIncr := 3
	if operand != 0 {
		jumpTo := parameterIndex(modes[1], intcode[pos+2], intcode)
		*ip = jumpTo
		ipIncr = 0
	}
	return ipIncr
}

func jumpIfFalse(pos int, intcode []int, modes []int, ip *int) int {
	operand := parameterIndex(modes[0], intcode[pos+1], intcode)
	ipIncr := 3
	if operand == 0 {
		jumpTo := parameterIndex(modes[1], intcode[pos+2], intcode)
		*ip = jumpTo
		ipIncr = 0
	}
	return ipIncr
}

// EvaluateIntcode evaluates and runs a given intcode
func EvaluateIntcode(intcode []int, inChan chan int, outChan chan int, thrusterOut chan int) {
	outVal := -1
	ip := 0
loop:
	for ip < len(intcode) {
		var ipIncr int
		switch op, modes := opcodeParse(intcode[ip]); op {
		case addCode:
			intcode, ipIncr = add(ip, intcode, modes)
		case multCode:
			intcode, ipIncr = mult(ip, intcode, modes)
		case inputCode:
			var inputVal int
			inputVal = <-inChan
			ipIncr = input(inputVal, ip, intcode)
		case outputCode:
			outVal, ipIncr = output(ip, intcode, modes)
			select {
			case outChan <- outVal:
				//fmt.Println("sending")
			default:
				//fmt.Println("Does it get here?")
				thrusterOut <- outVal
			}
		case jumpIfTrueCode:
			ipIncr = jumpIfTrue(ip, intcode, modes, &ip)
		case jumpIfFalseCode:
			ipIncr = jumpIfFalse(ip, intcode, modes, &ip)
		case lessThanCode:
			intcode, ipIncr = lessThan(ip, intcode, modes)
		case equalsCode:
			intcode, ipIncr = equals(ip, intcode, modes)
		case stopCode:
			break loop
		default:
			fmt.Println("Something went wrong, invalid opcode")
			fmt.Println("ip", ip)
			fmt.Println("op", op)
		}
		//if num == 4 {
		//fmt.Println("did a command: ", ip, inChan, outChan)
		//}
		ip += ipIncr
	}
	//fmt.Println("One done: ", inChan, outChan)
	//return intcode, outVal
}

// ResetMemory changes the intcode at indexes 1 and 2 with the given noun and
// verb
func ResetMemory(intcode []int, noun int, verb int) []int {
	intcode[1] = noun
	intcode[2] = verb
	return intcode
}

func main() {
	fmt.Println(opcodeParse(1101))
}
