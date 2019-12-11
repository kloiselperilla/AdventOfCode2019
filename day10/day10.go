package main

import (
	"AdventOfCode/lib/files"
	"fmt"
	"math"
)

type astMap struct {
	Map    map[point]byte
	Width  int
	Height int
}

type point struct {
	X int
	Y int
}

func buildMap(lines []string) astMap {
	charMap := map[point]byte{}
	for y, line := range lines {
		for x, char := range []byte(line) {
			if char == '#' {
				charMap[point{x, y}] = char
			}
		}
	}
	return astMap{Map: charMap, Width: len(lines[0]), Height: len(lines)}
}

func closeEnough(pt point, x float64, y float64) bool {
	return math.Abs(x-float64(pt.X)) < 1e-9 && math.Abs(y-float64(pt.Y)) < 1e-9
}

func isInt(val float64) bool {
	return val == float64(int64(val))
}

func hasLineOfSight(astA point, astB point, grid astMap) bool {
	xDiff := astB.X - astA.X
	yDiff := astB.Y - astA.Y
	if xDiff == 0 {
		sign := int(math.Abs(float64(yDiff)) / float64(yDiff))
		for yInc := 1; yInc < int(math.Abs(float64(yDiff))); yInc++ {
			dirYInc := sign * yInc
			if grid.Map[point{astA.X, astA.Y + dirYInc}] == '#' {
				return false
			}
		}
	} else {
		yInc := float64(yDiff) / float64(xDiff)
		for xInc := 1; xInc < int(math.Abs(float64(xDiff))); xInc++ {
			sign := int(math.Abs(float64(xDiff)) / float64(xDiff))
			dirXInc := sign * xInc
			y := float64(astA.Y) + float64(dirXInc)*yInc
			if isInt(y) {
				if grid.Map[point{astA.X + dirXInc, int(y)}] == '#' {
					return false
				}
			}
		}
	}
	return true
}

func numVisibleAsteroids(station point, grid astMap) int {
	total := 0
	for coord, char := range grid.Map {
		if char == '#' && coord != station && hasLineOfSight(station, coord, grid) {
			total++
		}
	}
	return total
}

func bestStation(grid astMap) (point, int) {
	max := -1
	var maxPt point
	for coord, char := range grid.Map {
		if char == '#' {
			numVisible := numVisibleAsteroids(coord, grid)
			if numVisible > max {
				max = numVisible
				maxPt = coord
			}
		}
	}
	return maxPt, max
}

func main() {
	lines := files.ReadLines("day10/input")
	grid := buildMap(lines)
	fmt.Println("Part 1:")
	fmt.Println(bestStation(grid))
	fmt.Println()
}
