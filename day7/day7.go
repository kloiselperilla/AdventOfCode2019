package main

import (
	"AdventOfCode/lib/intcode"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

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

func calculateSignal(code []int, phase []int) int {
	var wg sync.WaitGroup
	sig := 0
	var inChan [5]chan int
	var ready [5]chan bool
	for i := 0; i < 5; i++ {
		inChan[i] = make(chan int, 2)
		ready[i] = make(chan bool)
	}

	for i := 0; i < 5; i++ {
		codeCopy := make([]int, len(code))
		copy(codeCopy, code)

		wg.Add(1)
		go intcode.EvaluateIntcode(codeCopy, inChan[i], inChan[(i+1)%5], ready[i], &wg)
		inChan[i] <- phase[i]

	}
	inChan[0] <- 0
	for i := 0; i < 5; i++ {
		ready[i] <- true
	}
	wg.Wait()
	sig = <-inChan[0]

	return sig
}

func findMaxSignal(code []int, possiblePhases []int) int {
	max := -1

	phaseSequence := make([]int, 5)
	for i := 0; i < 5; i++ {
		possibleCopy0 := make([]int, len(possiblePhases))
		copy(possibleCopy0, possiblePhases)
		phaseSequence[0] = possibleCopy0[i]
		possibleCopy0 = append(possibleCopy0[:i], possibleCopy0[i+1:]...)
		for j := 0; j < 4; j++ {
			possibleCopy1 := make([]int, len(possibleCopy0))
			copy(possibleCopy1, possibleCopy0)
			phaseSequence[1] = possibleCopy1[j]
			possibleCopy1 = append(possibleCopy1[:j], possibleCopy1[j+1:]...)
			for k := 0; k < 3; k++ {
				possibleCopy2 := make([]int, len(possibleCopy1))
				copy(possibleCopy2, possibleCopy1)
				phaseSequence[2] = possibleCopy2[k]
				possibleCopy2 = append(possibleCopy2[:k], possibleCopy2[k+1:]...)
				for l := 0; l < 2; l++ {
					phaseSequence[3] = possibleCopy2[l]

					phaseSequence[4] = possibleCopy2[1-l]
					sig := calculateSignal(code, phaseSequence)
					if sig > max {
						max = sig
					}
				}
			}
		}
	}
	return max
}

func main() {
	originalCode := intcodeArray(readIntcode("day7/input"))
	possiblePhasesPt1 := []int{0, 1, 2, 3, 4}
	fmt.Println("Part 1:")
	fmt.Println(findMaxSignal(originalCode, possiblePhasesPt1))
	possiblePhasesPt2 := []int{5, 6, 7, 8, 9}
	fmt.Println("Part 2:")
	fmt.Println(findMaxSignal(originalCode, possiblePhasesPt2))

}
