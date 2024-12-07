package main

import (
	"log"
	"os"
	"strings"
)

const (
	Wall     = '#'
	Empty    = '.'
	StartPos = '^'
)

type Coordinate struct {
	x, y int
}

func (c Coordinate) move(d Direction) Coordinate {
	return Coordinate{c.x + d.dx, c.y + d.dy}
}

type Direction struct {
	dx, dy int
}

var directions = []Direction{
	{0, -1}, // up
	{1, 0},  // right
	{0, 1},  // down
	{-1, 0}, // left
}

type GuardMap [][]bool

func (m GuardMap) isWall(pos Coordinate) bool {
	return m[pos.y][pos.x]
}

func (m GuardMap) withinBounds(pos Coordinate) bool {
	return pos.x >= 0 && pos.x < len(m[0]) &&
		pos.y >= 0 && pos.y < len(m)
}

type Guard struct {
	pos      Coordinate
	dirIndex int
	visited  map[Coordinate]bool
	path     map[Coordinate]bool
}

func NewGuard(startPos Coordinate) *Guard {
	return &Guard{
		pos:      startPos,
		dirIndex: 0,
		visited:  make(map[Coordinate]bool),
		path:     make(map[Coordinate]bool),
	}
}

func (g *Guard) currentDirection() Direction {
	return directions[g.dirIndex]
}

func (g *Guard) turnRight() {
	g.dirIndex = (g.dirIndex + 1) % len(directions)
}

func (g *Guard) move(guardMap GuardMap, trackVisited bool) bool {
	if !guardMap.withinBounds(g.pos) {
		return false
	}

	if trackVisited {
		g.visited[g.pos] = true
	} else {
		g.path[g.pos] = true
	}

	next := g.pos.move(g.currentDirection())
	if !guardMap.withinBounds(next) {
		return false
	}

	if guardMap.isWall(next) {
		g.turnRight()
	} else {
		g.pos = next
	}
	return true
}

func main() {
	guardMap, initPos := parseInput("input.txt")

	p1 := partOne(guardMap, initPos)
	p2 := partTwo(guardMap, initPos)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(guardMap GuardMap, initPos Coordinate) int {
	guard := NewGuard(initPos)
	for guard.move(guardMap, true) {
	}
	return len(guard.visited)
}

func partTwo(guardMap GuardMap, initPos Coordinate) int {
	guard := NewGuard(initPos)
	for guard.move(guardMap, false) {
	}

	loopCount := 0
	for pos := range guard.path {
		if pos == initPos {
			continue
		}

		guardMap[pos.y][pos.x] = true
		if checkLoop(guardMap, initPos) {
			loopCount++
		}
		guardMap[pos.y][pos.x] = false
	}
	return loopCount
}

func checkLoop(guardMap GuardMap, initPos Coordinate) bool {
	guard := NewGuard(initPos)
	visited := make(map[Coordinate]struct{})

	for guard.move(guardMap, true) {
		key := Coordinate{guard.pos.x*4 + guard.dirIndex, guard.pos.y}
		if _, exists := visited[key]; exists {
			return true
		}
		visited[key] = struct{}{}
	}
	return false
}

func parseInput(input string) (GuardMap, Coordinate) {
	content, err := os.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	grid := make(GuardMap, len(lines))
	var init Coordinate

	for y, line := range lines {
		grid[y] = make([]bool, len(line))
		for x, char := range line {
			switch char {
			case Wall:
				grid[y][x] = true
			case StartPos:
				init = Coordinate{x, y}
				grid[y][x] = false
			case Empty:
				grid[y][x] = false
			}
		}
	}

	return grid, init
}
