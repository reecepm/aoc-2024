package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Coordinate struct {
	x, y int
}

func (c Coordinate) add(other Coordinate) Coordinate {
	return Coordinate{c.x + other.x, c.y + other.y}
}

var directions = []Coordinate{
	{1, 0},  // right
	{-1, 0}, // left
	{0, 1},  // down
	{0, -1}, // up
}

type MemorySpace struct {
	corrupted map[Coordinate]struct{}
	bounds    Coordinate
}

func NewMemorySpace() *MemorySpace {
	return &MemorySpace{
		corrupted: make(map[Coordinate]struct{}),
	}
}

func (m *MemorySpace) addCorruption(coord Coordinate) {
	m.corrupted[coord] = struct{}{}
	if coord.x > m.bounds.x {
		m.bounds.x = coord.x
	}
	if coord.y > m.bounds.y {
		m.bounds.y = coord.y
	}
}

func (m *MemorySpace) isValid(pos Coordinate) bool {
	return pos.x >= 0 && pos.x <= m.bounds.x &&
		pos.y >= 0 && pos.y <= m.bounds.y &&
		!m.isCorrupted(pos)
}

func (m *MemorySpace) isCorrupted(pos Coordinate) bool {
	_, exists := m.corrupted[pos]
	return exists
}

func (m *MemorySpace) findPath(start, end Coordinate) int {
	visited := make(map[Coordinate]bool)
	visited[start] = true

	current := []Coordinate{start}
	next := []Coordinate{}
	steps := 0

	for len(current) > 0 {
		next = next[:0]

		for _, pos := range current {
			if pos == end {
				return steps
			}

			for _, dir := range directions {
				nextPos := pos.add(dir)
				if !visited[nextPos] && m.isValid(nextPos) {
					visited[nextPos] = true
					next = append(next, nextPos)
				}
			}
		}

		current, next = next, current
		steps++
	}

	return -1
}

func (m *MemorySpace) copyWithCorruption(corruptions []Coordinate, limit int) *MemorySpace {
	copy := NewMemorySpace()
	copy.bounds = m.bounds

	for i := 0; i < limit && i < len(corruptions); i++ {
		copy.addCorruption(corruptions[i])
	}

	return copy
}

func main() {
	memory, corruptions, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(memory, corruptions)
	p2 := partTwo(memory, corruptions)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(memory *MemorySpace, corruptions []Coordinate) int {
	mem := memory.copyWithCorruption(corruptions, 1024)
	return mem.findPath(Coordinate{0, 0}, mem.bounds)
}

func partTwo(memory *MemorySpace, corruptions []Coordinate) Coordinate {
	index := sort.Search(len(corruptions), func(i int) bool {
		mem := memory.copyWithCorruption(corruptions, i+1)
		return mem.findPath(Coordinate{0, 0}, mem.bounds) == -1
	})
	return corruptions[index]
}

func parseInput(inputPath string) (*MemorySpace, []Coordinate, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	memory := NewMemorySpace()
	var corruptions []Coordinate

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var coord Coordinate
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 2 {
			continue
		}
		fmt.Sscanf(parts[0], "%d", &coord.x)
		fmt.Sscanf(parts[1], "%d", &coord.y)

		corruptions = append(corruptions, coord)
		memory.addCorruption(coord)
	}

	return memory, corruptions, scanner.Err()
}
