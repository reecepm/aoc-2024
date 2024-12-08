package main

import (
	"bufio"
	"log"
	"os"
)

type Coordinate struct {
	x, y int
}

func (c Coordinate) add(other Coordinate) Coordinate {
	return Coordinate{x: c.x + other.x, y: c.y + other.y}
}

func (c Coordinate) subtract(other Coordinate) Coordinate {
	return Coordinate{x: c.x - other.x, y: c.y - other.y}
}

func (c Coordinate) isWithinBounds(bounds Coordinate) bool {
	return c.x >= 0 && c.x < bounds.x && c.y >= 0 && c.y < bounds.y
}

func (c Coordinate) multiply(factor int) Coordinate {
	return Coordinate{x: c.x * factor, y: c.y * factor}
}

type LocationMap map[rune][]Coordinate

func main() {
	locations, bounds, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(locations, *bounds)
	p2 := partTwo(locations, *bounds)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func (lm LocationMap) findAntinodes(bounds Coordinate, extrapolate bool) map[Coordinate]bool {
	antinodes := make(map[Coordinate]bool)

	for _, coords := range lm {
		for i := 0; i < len(coords); i++ {
			for j := i + 1; j < len(coords); j++ {
				if extrapolate {
					antinodes[coords[i]] = true
					antinodes[coords[j]] = true
				}

				diff := coords[i].subtract(coords[j])

				processDirection(coords[i], diff, bounds, antinodes, extrapolate)
				processDirection(coords[j], diff.multiply(-1), bounds, antinodes, extrapolate)
			}
		}
	}

	return antinodes
}

func processDirection(start Coordinate, diff Coordinate, bounds Coordinate, antinodes map[Coordinate]bool, extrapolate bool) {
	current := start.add(diff)
	if !extrapolate {
		if current.isWithinBounds(bounds) {
			antinodes[current] = true
		}
		return
	}

	for current.isWithinBounds(bounds) {
		antinodes[current] = true
		current = current.add(diff)
	}
}

func partOne(lmap LocationMap, bounds Coordinate) int {
	return len(lmap.findAntinodes(bounds, false))
}

func partTwo(lmap LocationMap, bounds Coordinate) int {
	return len(lmap.findAntinodes(bounds, true))
}

func parseInput(inputPath string) (LocationMap, *Coordinate, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	locations := make(LocationMap)
	scanner := bufio.NewScanner(f)

	y := 0
	x := 0

	for scanner.Scan() {
		line := scanner.Text()
		x = len(line)

		for x, char := range line {
			if char == '.' {
				continue
			}

			coord := Coordinate{x: x, y: y}
			locations[char] = append(locations[char], coord)
		}

		y++
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return locations, &Coordinate{x: x, y: y}, nil
}
