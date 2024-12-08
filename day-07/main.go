package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	numMap, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(numMap)
	p2 := partTwo(numMap)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(numMap map[int][][]int) int {
	return sumOfTargets(numMap, canFormWithPlusMultiply)
}

func partTwo(numMap map[int][][]int) int {
	return sumOfTargets(numMap, canFormWithPlusMultiplyConcat)
}

func sumOfTargets(numMap map[int][][]int, checkFn func(int, [][]int) bool) int {
	total := 0
	for target, equations := range numMap {
		if checkFn(target, equations) {
			total += target
		}
	}
	return total
}

func canFormWithPlusMultiply(target int, equations [][]int) bool {
	for _, nums := range equations {
		if canFormTargetOps(nums, target, false) {
			return true
		}
	}
	return false
}

func canFormWithPlusMultiplyConcat(target int, equations [][]int) bool {
	for _, nums := range equations {
		if canFormTargetOps(nums, target, true) {
			return true
		}
	}
	return false
}

func canFormTargetOps(nums []int, target int, allowConcat bool) bool {
	if len(nums) == 1 {
		return nums[0] == target
	}

	a, b := nums[0], nums[1]
	rest := nums[2:]

	if canFormTargetOps(append([]int{a + b}, rest...), target, allowConcat) {
		return true
	}

	if canFormTargetOps(append([]int{a * b}, rest...), target, allowConcat) {
		return true
	}

	if allowConcat {
		concatVal := concatInts(a, b)
		if canFormTargetOps(append([]int{concatVal}, rest...), target, allowConcat) {
			return true
		}
	}

	return false
}

func concatInts(a, b int) int {
	s := strconv.Itoa(a) + strconv.Itoa(b)
	val, _ := strconv.Atoi(s)
	return val
}

func parseInput(inputPath string) (map[int][][]int, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	numMap := make(map[int][][]int)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid input line: %s", line)
		}

		targetStr := strings.TrimSpace(parts[0])
		numStrs := strings.Fields(parts[1])

		target, err := strconv.Atoi(targetStr)
		if err != nil {
			return nil, fmt.Errorf("invalid target %q: %w", targetStr, err)
		}

		ints := make([]int, 0, len(numStrs))
		for _, ns := range numStrs {
			n, err := strconv.Atoi(strings.TrimSpace(ns))
			if err != nil {
				return nil, fmt.Errorf("invalid number %q: %w", ns, err)
			}
			ints = append(ints, n)
		}

		numMap[target] = append(numMap[target], ints)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return numMap, nil
}
