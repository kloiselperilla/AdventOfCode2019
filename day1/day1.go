package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func moduleFuelRequired(mass int) int {
	return (mass / 3) - 2
}

func moduleFuelRequiredRec(mass int) int {
	fuel := (mass / 3) - 2
	if fuel <= 0 {
		return 0
	}
	return fuel + moduleFuelRequiredRec(fuel)
}

func moduleMassList(path string) []int {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var masses []int
	for scanner.Scan() {
		// string to int
		mass, err := strconv.Atoi(scanner.Text())
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		masses = append(masses, mass)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return masses
}

func sumOfFuels(masses []int) int {
	total := 0
	for _, mass := range masses {
		total += moduleFuelRequired(mass)
	}
	return total
}

func sumOfFuelsRec(masses []int) int {
	total := 0
	for _, mass := range masses {
		total += moduleFuelRequiredRec(mass)
	}
	return total
}

func main() {
	fmt.Println("Part 1:")
	fmt.Println(sumOfFuels(moduleMassList("input")))
	fmt.Println()

	fmt.Println("Part 2:")
	fmt.Println(sumOfFuelsRec(moduleMassList("input")))
}
