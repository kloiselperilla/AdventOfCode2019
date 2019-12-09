package main

import (
	"AdventOfCode/lib/intcode"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
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
	//fmt.Println(phase)
	//totalSig := 0
	//prevSig := 0
	sig := 0
	var inChan [5]chan int
	for i := 0; i < 5; i++ {
		inChan[i] = make(chan int)
	}
	thrusterChan := make(chan int)

	for i := 0; i < 5; i++ {
		codeCopy := make([]int, len(code))
		copy(codeCopy, code)

		//wg.Add(1)
		go intcode.EvaluateIntcode(codeCopy, inChan[i], inChan[(i+1)%5], thrusterChan)
		inChan[i] <- phase[i]
		if i == 0 {
			inChan[0] <- 0
		}

	}
	sig = <-thrusterChan

	return sig
}

func findMaxSignal(code []int, possiblePhases []int) int {
	max := -1

	phaseSequence := make([]int, 5)
	for i := 0; i < 5; i++ {
		//fmt.Println("LOOP ONE: ", i)
		//fmt.Println(possiblePhases)
		possibleCopy0 := make([]int, len(possiblePhases))
		copy(possibleCopy0, possiblePhases)
		phaseSequence[0] = possibleCopy0[i]
		possibleCopy0 = append(possibleCopy0[:i], possibleCopy0[i+1:]...)
		//fmt.Println(possibleCopy0)
		for j := 0; j < 4; j++ {
			//fmt.Println("LOOP TWO: ", j)
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
						//fmt.Println("max: ", max, phaseSequence)
					}
				}
			}
		}
	}
	fmt.Println(max, phaseSequence)
	return max
}

func main() {
	originalCode := intcodeArray(readIntcode("day7/input"))
	possiblePhasesPt1 := []int{0, 1, 2, 3, 4}
	fmt.Println(findMaxSignal(originalCode, possiblePhasesPt1))
	possiblePhasesPt2 := []int{5, 6, 7, 8, 9}
	fmt.Println(findMaxSignal(originalCode, possiblePhasesPt2))

}
