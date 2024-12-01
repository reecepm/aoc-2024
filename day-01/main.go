package main

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	arr1, arr2 := make([]int, 0), make([]int, 0)

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if parts := strings.Split(line, "   "); len(parts) == 2 {
			num1, err1 := strconv.Atoi(parts[0])
			num2, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil {
				arr1 = append(arr1, num1)
				arr2 = append(arr2, num2)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	p1 := part1(arr1, arr2)
	p2 := part2(arr1, arr2)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func part1(arr1 []int, arr2 []int) int {
	sorted1 := make([]int, len(arr1))
	sorted2 := make([]int, len(arr2))
	copy(sorted1, arr1)
	copy(sorted2, arr2)

	slices.Sort(sorted1)
	slices.Sort(sorted2)

	total := 0
	for i := range sorted1 {
		total += abs(sorted1[i] - sorted2[i])
	}
	return total
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func part2(arr1 []int, arr2 []int) int {
	secondItems := make(map[int]int, len(arr2))

	for _, v := range arr2 {
		secondItems[v]++
	}

	total := 0

	for _, v := range arr1 {
		if count, exists := secondItems[v]; exists {
			total += count * v
		}
	}

	return total
}
