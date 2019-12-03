package main

import (
	"bufio"
	"fmt"
	"github.com/deckarep/golang-set"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type coord struct {
	x int
	y int
}

func manDist(coord complex64) int {
	dist := math.Abs(float64(real(coord))) + math.Abs(float64(imag(coord)))
	return int(dist)
}

type byManDist []complex64

func (coord byManDist) Len() int           { return len(coord) }
func (coord byManDist) Less(i, j int) bool { return manDist(coord[i]) < manDist(coord[j]) }
func (coord byManDist) Swap(i, j int)      { coord[i], coord[j] = coord[j], coord[i] }

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

func populatePath(dirList []string) mapset.Set {
	path := mapset.NewSet()
	var pos complex64 = 0

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
			path.Add(pos)
		}

	}
	return path
}

func findClosestIntersect(pathA mapset.Set, pathB mapset.Set) complex64 {
	intersections := pathA.Intersect(pathB).ToSlice()
	interList := make([]complex64, len(intersections))
	for i := range intersections {
		interList[i] = intersections[i].(complex64)
	}

	sort.Sort(byManDist(interList))
	return interList[0]
}

func intersectFindFromLists(listA []string, listB []string) complex64 {
	pathA := populatePath(listA)
	pathB := populatePath(listB)

	return findClosestIntersect(pathA, pathB)
}

func main() {
	listA, listB := directionLists("input")
	fmt.Println("Part 1:")
	fmt.Println(manDist(intersectFindFromLists(listA, listB)))
}
