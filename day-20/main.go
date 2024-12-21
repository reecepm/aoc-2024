package main

import (
	"bufio"
	"log"
	"os"
)

type Coordinate struct {
	x, y int
}

var directions = []Coordinate{
	{0, -1}, // up
	{0, 1},  // down
	{-1, 0}, // left
	{1, 0},  // right
}

type Maze struct {
	grid  [][]rune
	start Coordinate
	end   Coordinate
}

func NewMaze() *Maze {
	return &Maze{}
}

func (m *Maze) FindCheats(maxCheatDist, minSaving int) int {
	path := m.findPath()
	count := 0

	for i := 0; i < len(path)-2; i++ {
		for j := i + 2; j < len(path); j++ {
			cheatDist := manhattan(path[i], path[j])
			if cheatDist <= maxCheatDist {
				saving := (j - i) - cheatDist
				if saving >= minSaving {
					count++
				}
			}
		}
	}
	return count
}

func (m *Maze) findPath() []Coordinate {
	queue := []pathState{{
		pos:   m.start,
		steps: 0,
		path:  []Coordinate{m.start},
	}}
	visited := make(map[Coordinate]bool)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.pos == m.end {
			return current.path
		}

		if visited[current.pos] {
			continue
		}
		visited[current.pos] = true

		for _, d := range directions {
			next := Coordinate{current.pos.x + d.x, current.pos.y + d.y}
			if m.isValid(next) {
				newPath := make([]Coordinate, len(current.path))
				copy(newPath, current.path)
				queue = append(queue, pathState{
					pos:   next,
					steps: current.steps + 1,
					path:  append(newPath, next),
				})
			}
		}
	}
	return nil
}

func (m *Maze) isValid(pos Coordinate) bool {
	return pos.x >= 0 && pos.x < len(m.grid[0]) &&
		pos.y >= 0 && pos.y < len(m.grid) &&
		m.grid[pos.y][pos.x] != '#'
}

type pathState struct {
	pos   Coordinate
	steps int
	path  []Coordinate
}

func manhattan(a, b Coordinate) int {
	return abs(a.y-b.y) + abs(a.x-b.x)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	maze, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(maze)
	p2 := partTwo(maze)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(m *Maze) int {
	return m.FindCheats(2, 100)
}

func partTwo(m *Maze) int {
	return m.FindCheats(20, 100)
}

func parseInput(inputPath string) (*Maze, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	maze := NewMaze()
	scanner := bufio.NewScanner(file)
	var grid [][]rune

	for scanner.Scan() {
		line := scanner.Text()
		row := []rune(line)
		for col, ch := range row {
			if ch == 'S' {
				maze.start = Coordinate{col, len(grid)}
			} else if ch == 'E' {
				maze.end = Coordinate{col, len(grid)}
			}
		}
		grid = append(grid, row)
	}

	maze.grid = grid
	return maze, scanner.Err()
}
