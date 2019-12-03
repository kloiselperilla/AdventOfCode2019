package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

/////////////////////////////
// Structs and Types
/////////////////////////////

type coordinate struct {
	pos     complex64
	elapsed int
}

type coordSet map[complex64]int

func newCoordSet() coordSet {
	return make(coordSet)
}

func (set coordSet) intersect(otherSet coordSet) coordSet {
	intersection := newCoordSet()
	for pos := range set {
		if _, ok := otherSet[pos]; ok {
			intersection.add(pos, set[pos]+otherSet[pos])
		}
	}
	return intersection
}

func (set coordSet) toSlice() []coordinate {
	retval := []coordinate{}
	for pos := range set {
		retval = append(retval, coordinate{pos: pos, elapsed: set[pos]})
	}
	return retval
}

func (set coordSet) add(coord complex64, elapsed int) {
	if _, ok := set[coord]; !ok {
		set[coord] = elapsed
	}
}

/////////////////////////////
// Helper Functions
/////////////////////////////

func manDist(coord coordinate) int {
	dist := math.Abs(float64(real(coord.pos))) + math.Abs(float64(imag(coord.pos)))
	return int(dist)
}

/////////////////////////////
// Comparators
/////////////////////////////

type byManDist []coordinate

func (coord byManDist) Len() int           { return len(coord) }
func (coord byManDist) Less(i, j int) bool { return manDist(coord[i]) < manDist(coord[j]) }
func (coord byManDist) Swap(i, j int)      { coord[i], coord[j] = coord[j], coord[i] }

type byElapsedDist []coordinate

func (coord byElapsedDist) Len() int           { return len(coord) }
func (coord byElapsedDist) Less(i, j int) bool { return coord[i].elapsed < coord[j].elapsed }
func (coord byElapsedDist) Swap(i, j int)      { coord[i], coord[j] = coord[j], coord[i] }

/////////////////////////////
// Main Functions
/////////////////////////////

func directionLists(path string) ([]string, []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	dirListA := strings.Split(strings.TrimSpace(scanner.Text()), ",")
	scanner.Scan()
	dirListB := strings.Split(strings.TrimSpace(scanner.Text()), ",")

	return dirListA, dirListB
}

func populatePath(dirList []string) coordSet {
	path := newCoordSet()
	var pos complex64 = 0
	elapsed := 0

	for _, move := range dirList {
		var dir complex64
		switch dirLetter := string(move[0]); {
		case dirLetter == "R":
			dir = 1
		case dirLetter == "D":
			dir = -1i
		case dirLetter == "L":
			dir = -1
		case dirLetter == "U":
			dir = 1i
		default:
			fmt.Println("Invalid direction")
			fmt.Println(move)
		}
		dist, err := strconv.Atoi(move[1:])
		if err != nil {
			panic(err)
		}

		for i := 0; i < dist; i++ {
			pos += dir
			elapsed++
			path.add(pos, elapsed)
		}

	}
	return path
}

func findClosestIntersect(pathA coordSet, pathB coordSet) coordinate {
	intersections := pathA.intersect(pathB).toSlice()

	sort.Sort(byManDist(intersections))
	return intersections[0]
}

func findShortestIntersect(pathA coordSet, pathB coordSet) coordinate {
	intersections := pathA.intersect(pathB).toSlice()

	sort.Sort(byElapsedDist(intersections))
	return intersections[0]
}

func intersectFindClosestFromLists(listA []string, listB []string) coordinate {
	pathA := populatePath(listA)
	pathB := populatePath(listB)

	return findClosestIntersect(pathA, pathB)
}

func intersectFindShortestFromLists(listA []string, listB []string) coordinate {
	pathA := populatePath(listA)
	pathB := populatePath(listB)

	return findShortestIntersect(pathA, pathB)
}

func main() {
	listA, listB := directionLists("input")
	fmt.Println("Part 1:")
	fmt.Println(manDist(intersectFindClosestFromLists(listA, listB)))

	fmt.Println("Part 2:")
	fmt.Println(intersectFindShortestFromLists(listA, listB).elapsed)
}
