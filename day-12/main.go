package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Coordinate struct {
	x, y int
}

func (c Coordinate) add(dir Coordinate) Coordinate {
	return Coordinate{c.x + dir.x, c.y + dir.y}
}

func (c Coordinate) rotateLeft() Coordinate {
	return Coordinate{c.y, -c.x}
}

func (c Coordinate) rotateRight() Coordinate {
	return Coordinate{-c.y, c.x}
}

var directions = []Coordinate{
	{0, -1}, // up
	{1, 0},  // right
	{0, 1},  // down
	{-1, 0}, // left
}

type Garden struct {
	grid   [][]rune
	width  int
	height int
}

func (g *Garden) isInBounds(c Coordinate) bool {
	return c.x >= 0 && c.x < g.width && c.y >= 0 && c.y < g.height
}

func (g *Garden) at(c Coordinate) (rune, bool) {
	if !g.isInBounds(c) {
		return 0, false
	}
	return g.grid[c.y][c.x], true
}

type Region struct {
	coords map[Coordinate]bool
	area   int
}

func (r *Region) contains(c Coordinate) bool {
	return r.coords[c]
}

func (r *Region) add(c Coordinate) {
	r.coords[c] = true
	r.area++
}

type OrientedEdge struct {
	pos Coordinate
	dir Coordinate
}

func (g *Garden) findAllRegions() []Region {
	visited := make([][]bool, g.height)
	for i := range visited {
		visited[i] = make([]bool, g.width)
	}

	var regions []Region
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if !visited[y][x] {
				coord := Coordinate{x, y}
				plantType := g.grid[y][x]
				region := g.floodFill(coord, plantType, visited)
				regions = append(regions, region)
			}
		}
	}
	return regions
}

func (g *Garden) floodFill(start Coordinate, plantType rune, visited [][]bool) Region {
	region := Region{coords: make(map[Coordinate]bool)}
	queue := []Coordinate{start}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if visited[curr.y][curr.x] {
			continue
		}

		visited[curr.y][curr.x] = true
		region.add(curr)

		for _, dir := range directions {
			next := curr.add(dir)
			if plant, ok := g.at(next); ok && plant == plantType && !visited[next.y][next.x] {
				queue = append(queue, next)
			}
		}
	}
	return region
}

func (g *Garden) calcEdges(r Region) int {
	edges := 0
	for coord := range r.coords {
		exposed := 4
		for _, dir := range directions {
			next := coord.add(dir)
			if r.contains(next) {
				exposed--
			}
		}
		edges += exposed
	}
	return edges
}

func (g *Garden) calcSides(r Region) int {
	processed := make(map[OrientedEdge]bool)
	sides := 0

	for coord := range r.coords {
		for _, dir := range directions {
			next := coord.add(dir)
			if !r.contains(next) {
				edge := OrientedEdge{coord, dir}
				if !processed[edge] {
					sides++

					for p1, p2 := coord, next; r.contains(p1) && !r.contains(p2); {
						processed[OrientedEdge{p1, dir}] = true
						p1 = p1.add(dir.rotateLeft())
						p2 = p2.add(dir.rotateLeft())
					}

					for p1, p2 := coord, next; r.contains(p1) && !r.contains(p2); {
						processed[OrientedEdge{p1, dir}] = true
						p1 = p1.add(dir.rotateRight())
						p2 = p2.add(dir.rotateRight())
					}
				}
			}
		}
	}
	return sides
}

func main() {
	f, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(f)
	p2 := partTwo(f)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(g *Garden) int {
	total := 0
	for _, region := range g.findAllRegions() {
		total += region.area * g.calcEdges(region)
	}
	return total
}

func partTwo(g *Garden) int {
	total := 0
	for _, region := range g.findAllRegions() {
		total += region.area * g.calcSides(region)
	}
	return total
}

func parseInput(inputPath string) (*Garden, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}
	defer file.Close()

	var grid [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, []rune(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return &Garden{
		grid:   grid,
		height: len(grid),
		width:  len(grid[0]),
	}, nil
}
