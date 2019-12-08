package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func addToMap(orbitMap map[string]celestialBody, line string) map[string]celestialBody {
	bodies := strings.Split(strings.TrimSpace(line), ")")
	orbitedName := bodies[0]
	orbiterName := bodies[1]

	orbitedBody := orbitMap[orbitedName]
	orbiterBody := orbitMap[orbiterName]
	orbiterBody.orbited = orbitedName

	orbitMap[orbitedName] = orbitedBody
	orbitMap[orbiterName] = orbiterBody

	return orbitMap
}

func orbitInit(path string) map[string]celestialBody {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	orbitMap := make(map[string]celestialBody)
	for scanner.Scan() {
		orbitMap = addToMap(orbitMap, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return orbitMap
}

type celestialBody struct {
	orbited string
	orbits  int
}

//func newCelestialBody() celestialBody {
//body := celestialBody{}
//body.orbited = ""
//body.orbits = 0

//return body
//}

func (body celestialBody) calculateOrbits(orbitMap map[string]celestialBody) int {
	if body.orbits > 0 {
		return body.orbits
	}
	totalOrbits := 0
	if body.orbited != "" {
		totalOrbits = 1 + orbitMap[body.orbited].calculateOrbits(orbitMap)

	}
	body.orbits = totalOrbits
	return totalOrbits
}

func totalOrbits(orbitMap map[string]celestialBody) int {
	total := 0
	for _, body := range orbitMap {
		total += body.calculateOrbits(orbitMap)
	}
	return total
}

func main() {
	orbitMap := orbitInit("day6/input")
	fmt.Println("Part 1:")
	fmt.Println(totalOrbits(orbitMap))
}
