package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type bodyMap map[string]*celestialBody

func addToMap(orbitMap bodyMap, line string) bodyMap {
	bodies := strings.Split(strings.TrimSpace(line), ")")
	orbitedName := bodies[0]
	orbiterName := bodies[1]

	//orbitedBody := *orbitMap[orbitedName]
	if orbitedPtr := orbitMap[orbitedName]; orbitedPtr == nil {
		orbitMap[orbitedName] = &celestialBody{bodyName: orbitedName}
	}
	if orbiterPtr := orbitMap[orbiterName]; orbiterPtr == nil {
		orbitMap[orbiterName] = &celestialBody{bodyName: orbiterName}
	}
	orbitMap[orbiterName].orbited = orbitedName

	//orbitMap[orbitedName] = orbitedBody
	//orbitMap[orbiterName] = orbiterBody

	return orbitMap
}

func orbitInit(path string) bodyMap {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	orbitMap := make(bodyMap)
	for scanner.Scan() {
		orbitMap = addToMap(orbitMap, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return orbitMap
}

type celestialBody struct {
	bodyName string
	orbited  string
	orbits   int
}

// JESUS CHRIST! I was having trouble saving the 'orbits' field in the body
// structs. BUT it turns out that you have to make the map point to pointers of
// celestialBodies AND you ALSO have to make the receiver type a pointer
// If you want to not use pointers you have to set the key's value again:
// a = map[thing]; a.blah = blah; map[thing] = a
func (body *celestialBody) calculateOrbits(orbitMap bodyMap) int {
	if body.orbits > 0 {
		return body.orbits
	}
	totalOrbits := 0
	if body.orbited != "" {
		var orbitedBody celestialBody = *orbitMap[body.orbited]
		totalOrbits = 1 + orbitedBody.calculateOrbits(orbitMap)
	}
	body.orbits = totalOrbits
	return totalOrbits
}

func (body celestialBody) findCommonOrbited(otherBody celestialBody, orbitMap bodyMap) celestialBody {
	for desc1 := body; desc1.orbited != ""; desc1 = *orbitMap[desc1.orbited] {
		for desc2 := otherBody; desc2.orbited != ""; desc2 = *orbitMap[desc2.orbited] {
			if desc1 == desc2 {
				return desc1
			}
		}
	}
	fmt.Println("ERRROR: they should all orbit COM")
	return celestialBody{}
}

func findDistance(bodyA celestialBody, bodyB celestialBody, orbitMap bodyMap) int {
	orbitedA := *orbitMap[bodyA.orbited]
	orbitedB := *orbitMap[bodyB.orbited]
	common := orbitedA.findCommonOrbited(orbitedB, orbitMap)
	return (orbitedA.orbits - common.orbits) + (orbitedB.orbits - common.orbits)
}

func totalOrbits(orbitMap bodyMap) int {
	total := 0
	for bodyName := range orbitMap {
		body := orbitMap[bodyName]
		total += body.calculateOrbits(orbitMap)
	}
	return total
}

func main() {
	orbitMap := orbitInit("day6/input")
	fmt.Println("Part 1:")
	fmt.Println(totalOrbits(orbitMap))

	fmt.Println()

	fmt.Println("Part 2:")
	fmt.Println(findDistance(*orbitMap["YOU"], *orbitMap["SAN"], orbitMap))
}
