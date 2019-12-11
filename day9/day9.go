package main

import (
	"AdventOfCode/lib/files"
	"AdventOfCode/lib/intcode"
	"fmt"
)

func runBoost(code []int, inputVal int) int {
	var out int
	eng := intcode.NewEngine(code)
	eng.Inputs.Enqueue(inputVal)
	q := intcode.NewSignalQueue()
	eng.ConnectOutput(&q)
	eng.EvaluateIntcode()
	if len(eng.Outputs.Queue) > 0 {
		out = eng.Outputs.Dequeue()
	}
	return out
}

func main() {
	code := intcode.StringToCode(files.ReadFile("day9/input"))
	fmt.Println("Part 1:")
	fmt.Println(runBoost(code, 1))
	fmt.Println()

	fmt.Println("Part 2:")
	fmt.Println(runBoost(code, 2))
}
