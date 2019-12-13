package main

import (
	"AdventOfCode/lib/files"
	"AdventOfCode/lib/intcode"
	"fmt"
	"sync"
)

// RobMap is the grid
type RobMap map[complex64]MapPos

// MapPos is a spot on the map
type MapPos struct {
	Color   int
	Painted bool
}

// PaintRobot controls the robot functionality
type PaintRobot struct {
	Engine intcode.Engine
	Dir    complex64
	Pos    complex64
}

func newPaintRobot(eng *intcode.Engine) PaintRobot {
	rob := PaintRobot{Engine: *eng, Dir: 1i, Pos: 0 + 0i}
	return rob
}

func (rob *PaintRobot) runRobot(robMap RobMap) {
	q := intcode.NewSignalQueue()
	rob.Engine.ConnectOutput(&q)

	var wg sync.WaitGroup
	wg.Add(1)
	rob.Engine.Inputs.Enqueue(robMap[rob.Pos].Color)
	go rob.Engine.EvaluateIntcode(&wg)
	isDone := false

	for !isDone {
		color := rob.Engine.Outputs.Dequeue()
		dir := rob.Engine.Outputs.Dequeue()
		pos := robMap[rob.Pos]
		pos.Painted = true
		pos.Color = color
		robMap[rob.Pos] = pos
		if dir == 0 {
			rob.Dir *= 1i
		} else if dir == 1 {
			rob.Dir *= -1i
		} else {
			fmt.Println("Something went wrong: ", dir)
		}
		rob.Pos += rob.Dir
		rob.Engine.Inputs.Enqueue(robMap[rob.Pos].Color)
		rob.Engine.WaitForReady()
		isDone = rob.Engine.IsDone
	}
}

func numPainted(rm RobMap) int {
	total := 0
	for _, pos := range rm {
		if pos.Painted {
			total++
		}
	}
	return total
}

func initializeMap(rm RobMap) {
	origin := rm[0+0i]
	origin.Color = 1
	rm[0+0i] = origin
}

// minX, maxX, minY, maxY
func findMapRange(rm RobMap) (int, int, int, int) {
	var minX, maxX, minY, maxY int
	for coord := range rm {
		if int(real(coord)) < minX {
			minX = int(real(coord))
		}
		if int(real(coord)) > maxX {
			maxX = int(real(coord))
		}
		if int(imag(coord)) < minY {
			minY = int(imag(coord))
		}
		if int(imag(coord)) > maxY {
			maxY = int(imag(coord))
		}
	}
	return minX, maxX, minY, maxY
}

func colorPosition(pos MapPos) string {
	if pos.Color == 1 {
		return "#"
	}
	return "."
}

func printMap(rm RobMap) {
	minX, maxX, minY, maxY := findMapRange(rm)
	//yRange := maxY - minY
	//xRange := maxX - minX

	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			coord := complex64(complex(float64(x), float64(y)))
			fmt.Print(colorPosition(rm[coord]))
		}
		fmt.Print("\n")
	}
}

func main() {
	code := intcode.StringToCode(files.Read("day11/input"))
	eng := intcode.NewEngine(code)
	rob := newPaintRobot(&eng)
	robMap := RobMap{}
	fmt.Println("Part 1:")
	rob.runRobot(robMap)
	fmt.Println(numPainted(robMap))

	fmt.Println()
	fmt.Println("Part 2:")

	eng2 := intcode.NewEngine(code)
	rob2 := newPaintRobot(&eng2)
	robMap2 := RobMap{}
	initializeMap(robMap2)

	rob2.runRobot(robMap2)
	printMap(robMap2)
}
