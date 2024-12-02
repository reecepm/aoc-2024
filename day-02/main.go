package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	arr := make([][]int, 0)

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
		if parts := strings.Split(line, " "); len(parts) > 1 {
			parsed := make([]int, 0)
			for _, part := range parts {
				if num, err := strconv.Atoi(part); err == nil {
					parsed = append(parsed, num)
				}
			}
			arr = append(arr, parsed)
		}
	}

	p1 := partOne(arr)
	p2 := partTwo(arr)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(grid [][]int) int {
	return countSafeSequences(grid, false)
}

func partTwo(grid [][]int) int {
	return countSafeSequences(grid, true)
}

func countSafeSequences(grid [][]int, allowOneRemoval bool) int {
	totalSafe := 0

	for _, row := range grid {
		if isValidSequence(row) {
			totalSafe++
			continue
		}

		if allowOneRemoval {
			for skipIdx := 0; skipIdx < len(row); skipIdx++ {
				testSeq := make([]int, 0, len(row)-1)
				testSeq = append(testSeq, row[:skipIdx]...)
				testSeq = append(testSeq, row[skipIdx+1:]...)

				if isValidSequence(testSeq) {
					totalSafe++
					break
				}
			}
		}
	}

	return totalSafe
}

func isValidSequence(nums []int) bool {
	if len(nums) <= 1 {
		return true
	}

	increasing := nums[1] > nums[0]

	for i := 1; i < len(nums); i++ {
		diff := nums[i] - nums[i-1]

		if increasing && diff <= 0 || !increasing && diff >= 0 {
			return false
		}

		absDiff := abs(diff)
		if absDiff < 1 || absDiff > 3 {
			return false
		}
	}

	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
