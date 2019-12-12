package main

import (
	"AdventOfCode/lib/files"
	"fmt"
	"math"
	"runtime"
	"sort"
	"sync"
)

type astMap struct {
	Map              map[point]mapPosition
	Width            int
	Height           int
	VisibleAsteroids map[point][]point
}

type point struct {
	X int
	Y int
}

type mapPosition struct {
	Pt    point
	Type  byte
	Angle float64
}

type byAngle []mapPosition

func (a byAngle) Len() int           { return len(a) }
func (a byAngle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byAngle) Less(i, j int) bool { return a[i].Angle < a[j].Angle }

func buildMap(lines []string) astMap {
	charMap := map[point]mapPosition{}
	for y, line := range lines {
		for x, char := range []byte(line) {
			if char == '#' {
				pt := point{X: x, Y: y}
				charMap[pt] = mapPosition{Pt: pt, Type: char}
			}
		}
	}
	return astMap{Map: charMap, Width: len(lines[0]), Height: len(lines)}
}

func isInt(val float64) bool {
	return val == float64(int64(val))
}

func hasLineOfSight(astA point, astB point, grid *astMap, out *bool, wg *sync.WaitGroup) {
	defer wg.Done()

	hasLine := true
	xDiff := astB.X - astA.X
	yDiff := astB.Y - astA.Y
	if xDiff == 0 {
		sign := int(math.Abs(float64(yDiff)) / float64(yDiff))
		for yInc := 1; yInc < int(math.Abs(float64(yDiff))); yInc++ {
			dirYInc := sign * yInc
			if grid.Map[point{X: astA.X, Y: astA.Y + dirYInc}].Type == '#' {
				hasLine = false
				break
			}
		}
	} else {
		yInc := float64(yDiff) / float64(xDiff)
		for xInc := 1; xInc < int(math.Abs(float64(xDiff))); xInc++ {
			sign := int(math.Abs(float64(xDiff)) / float64(xDiff))
			dirXInc := sign * xInc
			y := float64(astA.Y) + float64(dirXInc)*yInc
			if isInt(y) {
				if grid.Map[point{X: astA.X + dirXInc, Y: int(y)}].Type == '#' {
					hasLine = false
					break
				}
			}
		}
	}
	*out = hasLine
}

// A concurrent paradigm: For embarrasingly parallel tasks, give each task a
// pointer to an output. That pointer belongs to an array of outputs.
// Then if you need to, you can have a helper map to map relevant info to an
// index of the output array
func numVisibleAsteroids(station point, grid *astMap, out *int, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	total := 0
	var myWg sync.WaitGroup
	visibleArr := make([]bool, len(grid.Map))
	helperMap := map[point]int{}
	i := 0
	for coord, pos := range grid.Map {
		if pos.Type == '#' && coord != station {
			myWg.Add(1)
			hasLineOfSight(station, coord, grid, &visibleArr[i], &myWg)
			helperMap[coord] = i
			i++
		}
	}
	myWg.Wait()

	for coord, ix := range helperMap {
		isVisible := visibleArr[ix]
		if isVisible {
			total++
			visAst := grid.VisibleAsteroids[station]
			if visAst == nil {
				grid.VisibleAsteroids = map[point][]point{}
				grid.VisibleAsteroids[station] = []point{}
				visAst = grid.VisibleAsteroids[station]
			}
			grid.VisibleAsteroids[station] = append(visAst, coord)
			//fmt.Println(grid.VisibleAsteroids[point{26, 28}])
			//if wg == nil {
			//fmt.Println(len(grid.VisibleAsteroids[station]))
			//}
		}

	}
	*out = total
}

func bestStation(grid *astMap) (point, int) {
	max := -1
	var maxPt point
	var myWg sync.WaitGroup
	stationArr := make([]int, len(grid.Map))
	helperMap := map[point]int{}
	i := 0
	for coord, pos := range grid.Map {
		if pos.Type == '#' {
			myWg.Add(1)
			numVisibleAsteroids(coord, grid, &stationArr[i], &myWg)
			helperMap[coord] = i
			i++
		}
	}
	myWg.Wait()

	for coord, ix := range helperMap {
		numVisible := stationArr[ix]
		if numVisible > max {
			max = numVisible
			maxPt = coord
		}
	}
	return maxPt, max
}

func populateAngles(station point, grid *astMap) {
	for coord, pos := range grid.Map {
		if pos.Type == '#' {
			xDiff := float64(coord.X - station.X)
			yDiff := float64(coord.Y - station.Y)
			angRad := math.Atan2(xDiff, yDiff)
			if angRad < 0 {
				angRad += 1 * math.Pi
			}
			pos := grid.Map[coord]
			pos.Angle = math.Atan2(-1*xDiff, yDiff)
			grid.Map[coord] = pos
			fmt.Println(coord, angRad)
		}
	}
}

func rotateLaser(station point, grid *astMap) []point {
	vaporized := []point{}

	visibleAsteroidsPt := grid.VisibleAsteroids[station]
	visibleAsteroids := []mapPosition{}
	for _, coord := range visibleAsteroidsPt {
		visibleAsteroids = append(visibleAsteroids, grid.Map[coord])
	}

	sort.Sort(byAngle(visibleAsteroids))

	for _, pos := range visibleAsteroids {
		pos.Type = '.'
		grid.Map[pos.Pt] = pos
		vaporized = append(vaporized, pos.Pt)
	}

	//fmt.Println(vaporized)
	return vaporized
}

func destroyAsteroids(station point, grid *astMap) []point {
	vaporized := []point{}
	newVisNum := 0
	numVisibleAsteroids(station, grid, &newVisNum, nil)
	//fmt.Println(grid.VisibleAsteroids[station])
	for len(grid.VisibleAsteroids[station]) > 0 {
		vapOneRound := rotateLaser(station, grid)
		vaporized = append(vaporized, vapOneRound...)
		grid.VisibleAsteroids[station] = []point{}
		newVisNum = 0
		numVisibleAsteroids(station, grid, &newVisNum, nil)
	}

	return vaporized
}

func main() {
	runtime.GOMAXPROCS(12)
	lines := files.ReadLines("day10/input")
	grid := buildMap(lines)
	fmt.Println("Part 1:")
	station, visible := bestStation(&grid)
	fmt.Println(station, visible)
	fmt.Println()
	populateAngles(station, &grid)
	fmt.Println("Part 2:")
	fmt.Println(destroyAsteroids(station, &grid)[199])
}
