package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
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

type Cell rune

const (
	Empty    Cell = '.'
	Wall     Cell = '#'
	Reindeer Cell = 'S'
	End      Cell = 'E'
)

type State struct {
	pos      Coordinate
	dir      Direction
	cost     int
	estimate int
	parent   *State
}

type PriorityQueue []*State

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].estimate < pq[j].estimate }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*State)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func manhattan(a, b Coordinate) int {
	return abs(b.x-a.x) + abs(b.y-a.y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Maze struct {
	Grid     [][]Cell
	StartPos Coordinate
	EndPos   Coordinate
	StartDir Direction
	Width    int
	Height   int
}

func NewMaze() *Maze {
	return &Maze{}
}

func (m *Maze) FindOptimalPath() int {
	return m.astar(false)
}

func (m *Maze) collectOptimalPaths(minCost int) map[Coordinate]bool {
	optimalTiles := make(map[Coordinate]bool)
	m.astar(true, minCost, optimalTiles)
	return optimalTiles
}

func (m *Maze) astar(collectPaths bool, params ...interface{}) int {
	start := &State{
		pos:      m.StartPos,
		dir:      m.StartDir,
		cost:     0,
		estimate: manhattan(m.StartPos, m.EndPos),
	}

	openSet := &PriorityQueue{start}
	heap.Init(openSet)
	visited := make(map[string]int)

	var minCost int
	var optimalTiles map[Coordinate]bool
	if collectPaths {
		minCost = params[0].(int)
		optimalTiles = params[1].(map[Coordinate]bool)
	}

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*State)

		if current.pos == m.EndPos {
			if !collectPaths {
				return current.cost
			}
			if current.cost == minCost {
				for node := current; node != nil; node = node.parent {
					optimalTiles[node.pos] = true
				}
			}
			continue
		}

		if collectPaths && current.cost+manhattan(current.pos, m.EndPos) > minCost {
			continue
		}

		key := fmt.Sprintf("%d,%d,%d,%d",
			current.pos.x, current.pos.y,
			current.dir.dx, current.dir.dy)

		if score, exists := visited[key]; exists {
			if !collectPaths || score < current.cost {
				continue
			}
		}
		visited[key] = current.cost

		if next := current.pos.move(current.dir); m.isValid(next) {
			heap.Push(openSet, &State{
				pos:      next,
				dir:      current.dir,
				cost:     current.cost + 1,
				estimate: current.cost + 1 + manhattan(next, m.EndPos),
				parent:   current,
			})
		}

		dirs := []Direction{
			{-current.dir.dy, current.dir.dx},
			{current.dir.dy, -current.dir.dx},
		}

		for _, newDir := range dirs {
			heap.Push(openSet, &State{
				pos:      current.pos,
				dir:      newDir,
				cost:     current.cost + 1000,
				estimate: current.cost + 1000 + manhattan(current.pos, m.EndPos),
				parent:   current,
			})
		}
	}

	return -1
}

func (m *Maze) isValid(pos Coordinate) bool {
	return pos.x >= 0 && pos.x < m.Width &&
		pos.y >= 0 && pos.y < m.Height &&
		m.Grid[pos.y][pos.x] != Wall
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
	return m.FindOptimalPath()
}

func partTwo(m *Maze) int {
	minCost := m.FindOptimalPath()
	optimalPaths := m.collectOptimalPaths(minCost)
	return len(optimalPaths)
}

func parseInput(inputPath string) (*Maze, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}
	defer file.Close()

	maze := NewMaze()
	scanner := bufio.NewScanner(file)

	var grid [][]Cell
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		row := make([]Cell, len(line))
		for i, ch := range line {
			cell := Cell(ch)
			row[i] = cell
			if cell == Reindeer {
				maze.StartPos = Coordinate{i, len(grid)}
			}
			if cell == End {
				maze.EndPos = Coordinate{i, len(grid)}
			}
		}
		grid = append(grid, row)
	}

	maze.Grid = grid
	maze.Height = len(grid)
	if len(grid) > 0 {
		maze.Width = len(grid[0])
	}
	maze.StartDir = Direction{1, 0}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return maze, nil
}
