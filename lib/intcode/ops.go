package intcode

import (
	"fmt"
)

func add(eng *Engine) int {
	operand1 := parameterIndex(eng, 0)
	operand2 := parameterIndex(eng, 1)
	sum := *operand1 + *operand2

	if eng.Modes[2] == 1 {
		fmt.Println("No write param should be immediate mode")
	}
	dst := parameterIndex(eng, 2)
	*dst = sum
	return 4
}

func mult(eng *Engine) int {
	operand1 := parameterIndex(eng, 0)
	operand2 := parameterIndex(eng, 1)
	prod := *operand1 * *operand2

	if eng.Modes[2] == 1 {
		fmt.Println("No write param should be immediate mode")
	}
	dst := parameterIndex(eng, 2)
	*dst = prod
	return 4
}

func lessThan(eng *Engine) int {
	operand1 := parameterIndex(eng, 0)
	operand2 := parameterIndex(eng, 1)
	ans := *operand1 < *operand2

	if eng.Modes[2] == 1 {
		fmt.Println("No write param should be immediate mode")
	}
	var ansInt int
	if ans {
		ansInt = 1
	}
	dst := parameterIndex(eng, 2)
	*dst = ansInt
	return 4
}

func equals(eng *Engine) int {
	operand1 := parameterIndex(eng, 0)
	operand2 := parameterIndex(eng, 1)
	ans := *operand1 == *operand2

	if eng.Modes[2] == 1 {
		fmt.Println("No write param should be immediate mode")
	}
	var ansInt int
	if ans {
		ansInt = 1
	}
	dst := parameterIndex(eng, 2)
	*dst = ansInt
	return 4
}

func input(eng *Engine) int {
	eng.InputVal = eng.Inputs.Dequeue()
	dst := parameterIndex(eng, 0)
	*dst = eng.InputVal

	return 2
}

func output(eng *Engine) int {
	outVal := parameterIndex(eng, 0)
	eng.Outputs.Enqueue(*outVal)
	return 2
}

func jumpIfTrue(eng *Engine) int {
	operand := parameterIndex(eng, 0)
	ipIncr := 3
	if *operand != 0 {
		jumpTo := parameterIndex(eng, 1)
		eng.Ip = *jumpTo
		ipIncr = 0
	}
	return ipIncr
}

func jumpIfFalse(eng *Engine) int {
	operand := parameterIndex(eng, 0)
	ipIncr := 3
	if *operand == 0 {
		jumpTo := parameterIndex(eng, 1)
		eng.Ip = *jumpTo
		ipIncr = 0
	}
	return ipIncr
}

func relBaseOffset(eng *Engine) int {
	operand := parameterIndex(eng, 0)
	eng.RelBase += *operand

	return 2
}
