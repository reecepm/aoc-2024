package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Pin struct {
	heights []int
}

type LockAndKey struct {
	locks []Pin
	keys  []Pin
}

func NewLockAndKey() *LockAndKey {
	return &LockAndKey{
		locks: make([]Pin, 0),
		keys:  make([]Pin, 0),
	}
}

func (lk *LockAndKey) addLock(heights []int) {
	lk.locks = append(lk.locks, Pin{heights: heights})
}

func (lk *LockAndKey) addKey(heights []int) {
	lk.keys = append(lk.keys, Pin{heights: heights})
}

func (lk *LockAndKey) countValidPairs() int {
	count := 0
	for _, lock := range lk.locks {
		for _, key := range lk.keys {
			if lk.isValidPair(lock, key) {
				count++
			}
		}
	}
	return count
}

func (lk *LockAndKey) isValidPair(lock, key Pin) bool {
	if len(lock.heights) != len(key.heights) {
		return false
	}

	for i := 0; i < len(lock.heights); i++ {
		if lock.heights[i]+key.heights[i] > 5 {
			return false
		}
	}
	return true
}

func calculateLockHeights(lines []string) []int {
	var heights []int
	width := len(lines[0])

	for col := 0; col < width; col++ {
		height := 0
		for row := 0; row < len(lines); row++ {
			if lines[row][col] == '#' {
				height++
			} else {
				break
			}
		}
		heights = append(heights, height-1)
		if heights[len(heights)-1] < 0 {
			heights[len(heights)-1] = 0
		}
	}
	return heights
}

func calculateKeyHeights(lines []string) []int {
	var heights []int
	width := len(lines[0])

	for col := 0; col < width; col++ {
		height := 0
		for row := len(lines) - 1; row >= 0; row-- {
			if lines[row][col] == '#' {
				height++
			} else {
				break
			}
		}
		heights = append(heights, height-1)
		if heights[len(heights)-1] < 0 {
			heights[len(heights)-1] = 0
		}
	}
	return heights
}

func parseInput(path string) (*LockAndKey, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lk := NewLockAndKey()

	var currentSchematic []string

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if len(currentSchematic) > 0 {
				if strings.Count(currentSchematic[0], "#") == 5 {
					heights := calculateLockHeights(currentSchematic)
					lk.addLock(heights)
				} else {
					heights := calculateKeyHeights(currentSchematic)
					lk.addKey(heights)
				}
				currentSchematic = nil
			}
			continue
		}

		currentSchematic = append(currentSchematic, line)
	}

	if len(currentSchematic) > 0 {
		if strings.Count(currentSchematic[0], "#") == 5 {
			heights := calculateLockHeights(currentSchematic)
			lk.addLock(heights)
		} else {
			heights := calculateKeyHeights(currentSchematic)
			lk.addKey(heights)
		}
	}

	return lk, nil
}

func main() {
	lk, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println(lk.locks)
	fmt.Println(lk.keys)

	p1 := partOne(lk)
	log.Printf("part1: %d", p1)
}

func partOne(lk *LockAndKey) int {
	return lk.countValidPairs()
}
