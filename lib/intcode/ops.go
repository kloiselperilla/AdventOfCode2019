package intcode

import (
	"fmt"
)

func add(eng *Engine) int {
	operand1 := parameterIndex(eng, 0)
	operand2 := parameterIndex(eng, 1)
	sum := operand1 + operand2

	if eng.Modes[2] != 0 {
		fmt.Println("No write param should be immediate mode")
	}
	eng.Code[eng.Code[eng.Ip+3]] = sum
	return 4
}

func mult(eng *Engine) int {
	operand1 := parameterIndex(eng, 0)
	operand2 := parameterIndex(eng, 1)
	prod := operand1 * operand2

	if eng.Modes[2] != 0 {
		fmt.Println("No write param should be immediate mode")
	}
	eng.Code[eng.Code[eng.Ip+3]] = prod
	return 4
}

func lessThan(eng *Engine) int {
	operand1 := parameterIndex(eng, 0)
	operand2 := parameterIndex(eng, 1)
	ans := operand1 < operand2

	if eng.Modes[2] != 0 {
		fmt.Println("No write param should be immediate mode")
	}
	var ansInt int
	if ans {
		ansInt = 1
	}
	eng.Code[eng.Code[eng.Ip+3]] = ansInt
	return 4
}

func equals(eng *Engine) int {
	operand1 := parameterIndex(eng, 0)
	operand2 := parameterIndex(eng, 1)
	ans := operand1 == operand2

	if eng.Modes[2] != 0 {
		fmt.Println("No write param should be immediate mode")
	}
	var ansInt int
	if ans {
		ansInt = 1
	}
	eng.Code[eng.Code[eng.Ip+3]] = ansInt
	return 4
}

func input(eng *Engine) int {
	eng.InputVal = eng.Inputs.Dequeue()
	eng.Code[eng.Code[eng.Ip+1]] = eng.InputVal

	return 2
}

func output(eng *Engine) int {
	outVal := parameterIndex(eng, 0)
	eng.Outputs.Enqueue(outVal)
	return 2
}

func jumpIfTrue(eng *Engine) int {
	operand := parameterIndex(eng, 0)
	ipIncr := 3
	if operand != 0 {
		jumpTo := parameterIndex(eng, 1)
		eng.Ip = jumpTo
		ipIncr = 0
	}
	return ipIncr
}

func jumpIfFalse(eng *Engine) int {
	operand := parameterIndex(eng, 0)
	ipIncr := 3
	if operand == 0 {
		jumpTo := parameterIndex(eng, 1)
		eng.Ip = jumpTo
		ipIncr = 0
	}
	return ipIncr
}
