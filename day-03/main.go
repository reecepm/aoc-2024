package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	input := strings.ReplaceAll(string(data), "\n", "")
	input = strings.ReplaceAll(input, "\r", "")
	input = strings.ReplaceAll(input, " ", "")

	log.Printf("part1: %d", partOne(input))
	log.Printf("part2: %d", partTwo(input))
}

func extractNumsFromMul(mulStr string) (int, int) {
	values := strings.Split(mulStr[4:len(mulStr)-1], ",")
	num1, _ := strconv.Atoi(values[0])
	num2, _ := strconv.Atoi(values[1])
	return num1, num2
}

func partOne(input string) int {
	re := regexp.MustCompile(`mul\(\d+,\d+\)`)
	sum := 0
	matches := re.FindAllString(input, -1)
	for _, match := range matches {
		num1, num2 := extractNumsFromMul(match)
		sum += num1 * num2
	}
	return sum
}

func partTwo(input string) int {
	re := regexp.MustCompile(`mul\(\d+,\d+\)|do\(\)|don't\(\)`)
	sum := 0
	enabled := true

	matches := re.FindAllString(input, -1)
	for _, match := range matches {
		switch match {
		case "do()":
			enabled = true
		case "don't()":
			enabled = false
		default:
			if enabled {
				num1, num2 := extractNumsFromMul(match)
				sum += num1 * num2
			}
		}
	}
	return sum
}
