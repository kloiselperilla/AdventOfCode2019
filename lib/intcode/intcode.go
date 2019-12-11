package intcode

import (
	"fmt"
	"sync"
)

// SignalQueue gives an implementation of a connection between engines
type SignalQueue struct {
	Queue []int
	ready chan bool
	mux   *sync.Mutex
	cond  *sync.Cond
}

func newSignalQueue() SignalQueue {
	mux := sync.Mutex{}
	q := SignalQueue{ready: make(chan bool), mux: &mux, cond: sync.NewCond(&mux)}
	return q
}

// Enqueue adds to end of queue
func (q *SignalQueue) Enqueue(val int) {
	q.cond.L.Lock()
	q.Queue = append(q.Queue, val)
	q.cond.Broadcast()
	q.cond.L.Unlock()
}

// Dequeue removes from beginning of queue
func (q *SignalQueue) Dequeue() int {
	q.cond.L.Lock()
	// Wait for not empty
	for len(q.Queue) == 0 {
		q.cond.Wait()
	}
	retval := q.Queue[0]
	q.Queue[0] = 0
	q.Queue = q.Queue[1:]
	q.cond.L.Unlock()

	return retval
}

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

// Engine encapsulates an intcode consumer
type Engine struct {
	code    []int
	Inputs  *SignalQueue
	Outputs *SignalQueue
}

// NewEngine makes a new Engine with a copy of an intcode
func NewEngine(code []int) Engine {
	codeCopy := make([]int, len(code))
	copy(codeCopy, code)

	q := newSignalQueue()
	e := Engine{code: codeCopy, Inputs: &q}
	return e
}

// ConnectOutput sets an output queue
func (e *Engine) ConnectOutput(outputs *SignalQueue) {
	e.Outputs = outputs
}

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
func (e *Engine) EvaluateIntcode(wg *sync.WaitGroup) {
	defer wg.Done()
	outVal := -1
	ip := 0
loop:
	for ip < len(e.code) {
		var ipIncr int
		switch op, modes := opcodeParse(e.code[ip]); op {
		case addCode:
			e.code, ipIncr = add(ip, e.code, modes)
		case multCode:
			e.code, ipIncr = mult(ip, e.code, modes)
		case inputCode:
			inputVal := e.Inputs.Dequeue()
			ipIncr = input(inputVal, ip, e.code)
		case outputCode:
			outVal, ipIncr = output(ip, e.code, modes)
			e.Outputs.Enqueue(outVal)
		case jumpIfTrueCode:
			ipIncr = jumpIfTrue(ip, e.code, modes, &ip)
		case jumpIfFalseCode:
			ipIncr = jumpIfFalse(ip, e.code, modes, &ip)
		case lessThanCode:
			e.code, ipIncr = lessThan(ip, e.code, modes)
		case equalsCode:
			e.code, ipIncr = equals(ip, e.code, modes)
		case stopCode:
			break loop
		default:
			fmt.Println("Something went wrong, invalid opcode")
			fmt.Println("ip", ip)
			fmt.Println("op", op)
		}
		ip += ipIncr
	}
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
