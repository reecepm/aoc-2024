package main

import (
	"bufio"
	"log"
	"os"
)

const (
	MinHeight = 0
	MaxHeight = 9
)

type Coordinate struct {
	x, y int
}

func (c Coordinate) add(other Coordinate) Coordinate {
	return Coordinate{c.x + other.x, c.y + other.y}
}

var directions = []Coordinate{
	{0, -1}, // up
	{1, 0},  // right
	{0, 1},  // down
	{-1, 0}, // left
}

type Grid map[Coordinate]int

type HikingTrails struct {
	grid Grid
}

func (h *HikingTrails) findTrailheadScores(countPaths bool) int {
	totalScore := 0
	for pos, height := range h.grid {
		if height == MinHeight {
			if countPaths {
				totalScore += h.countTrailPaths(pos)
			} else {
				totalScore += len(h.findReachableNines(pos))
			}
		}
	}
	return totalScore
}

func (h *HikingTrails) countTrailPaths(start Coordinate) int {
	visited := make(map[Coordinate]bool)
	return h.exploreTrailPaths(start, visited)
}

func (h *HikingTrails) findReachableNines(start Coordinate) map[Coordinate]bool {
	reachableNines := make(map[Coordinate]bool)
	visited := make(map[Coordinate]bool)
	h.exploreTrails(start, visited, reachableNines)
	return reachableNines
}

func (h *HikingTrails) exploreTrails(current Coordinate, visited map[Coordinate]bool, nines map[Coordinate]bool) {
	currentHeight := h.grid[current]
	visited[current] = true
	defer func() { visited[current] = false }()

	for _, dir := range directions {
		next := current.add(dir)
		nextHeight, exists := h.grid[next]
		if !exists || visited[next] {
			continue
		}

		if nextHeight == MaxHeight && currentHeight == MaxHeight-1 {
			nines[next] = true
			continue
		}

		if nextHeight == currentHeight+1 {
			h.exploreTrails(next, visited, nines)
		}
	}
}

func (h *HikingTrails) exploreTrailPaths(current Coordinate, visited map[Coordinate]bool) int {
	currentHeight := h.grid[current]
	if currentHeight == MaxHeight {
		return 1
	}

	visited[current] = true
	defer func() { visited[current] = false }()

	pathCount := 0
	for _, dir := range directions {
		next := current.add(dir)
		nextHeight, exists := h.grid[next]
		if !exists || visited[next] {
			continue
		}

		if nextHeight == currentHeight+1 {
			pathCount += h.exploreTrailPaths(next, visited)
		}
	}
	return pathCount
}

func main() {
	trails, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(trails)
	p2 := partTwo(trails)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(h *HikingTrails) int {
	return h.findTrailheadScores(false)
}

func partTwo(h *HikingTrails) int {
	return h.findTrailheadScores(true)
}

func parseInput(inputPath string) (*HikingTrails, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	grid := make(Grid)
	scanner := bufio.NewScanner(file)
	y := 0

	for scanner.Scan() {
		line := scanner.Text()
		for x, char := range line {
			height := int(char - '0')
			grid[Coordinate{x, y}] = height
		}
		y++
	}

	return &HikingTrails{grid: grid}, scanner.Err()
}
