package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Coordinate struct {
	x, y int64
}

func (c Coordinate) move(dir Direction) Coordinate {
	switch dir {
	case Up:
		return Coordinate{c.x, c.y - 1}
	case Down:
		return Coordinate{c.x, c.y + 1}
	case Left:
		return Coordinate{c.x - 1, c.y}
	case Right:
		return Coordinate{c.x + 1, c.y}
	}
	return c
}

type Cell rune

const (
	Empty Cell = '.'
	Wall  Cell = '#'
	Box   Cell = 'O'
	Robot Cell = '@'
	LBox  Cell = '['
	RBox  Cell = ']'
)

type Direction rune

const (
	Up    Direction = '^'
	Down  Direction = 'v'
	Left  Direction = '<'
	Right Direction = '>'
)

type Warehouse struct {
	Grid     [][]Cell
	Moves    []Direction
	RobotPos Coordinate
	Width    int
	Height   int
}

func NewWarehouse() *Warehouse {
	return &Warehouse{}
}

func (w *Warehouse) Copy() *Warehouse {
	copy := *w
	copy.Grid = make([][]Cell, len(w.Grid))
	for i, row := range w.Grid {
		copy.Grid[i] = make([]Cell, len(row))
		copy.Grid[i] = append([]Cell(nil), row...)
	}
	return &copy
}

func (w *Warehouse) Display() {
	for _, rows := range w.Grid {
		for _, col := range rows {
			fmt.Print(string(col))
		}
		fmt.Println()
	}
	fmt.Println()
}

func (w *Warehouse) CalculateScore(boxType Cell) int {
	total := 0
	for r, rows := range w.Grid {
		for c, col := range rows {
			if col == boxType {
				total += (100 * r) + c
			}
		}
	}
	return total
}

func main() {
	wh, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(wh.Copy())
	p2 := partTwo(wh.Copy())

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func (w *Warehouse) canMove(pos Coordinate, dir Direction) bool {
	next := pos.move(dir)
	switch w.Grid[next.y][next.x] {
	case Wall:
		return false
	case LBox:
		if dir == Up || dir == Down {
			nextRight := Coordinate{next.x + 1, next.y}
			return w.canMove(next, dir) && w.canMove(nextRight, dir)
		}
		return w.canMove(next, dir)
	case RBox:
		if dir == Up || dir == Down {
			nextLeft := Coordinate{next.x - 1, next.y}
			return w.canMove(nextLeft, dir) && w.canMove(next, dir)
		}
		return w.canMove(next, dir)
	default:
		return true
	}
}

func (w *Warehouse) moveSingleBox(pos Coordinate, dir Direction) bool {
	var boxPositions []Coordinate
	curr := pos

	for w.Grid[curr.y][curr.x] == Box {
		boxPositions = append(boxPositions, curr)
		next := curr.move(dir)
		if w.Grid[next.y][next.x] == Wall {
			return false
		}
		curr = next
	}

	for i := len(boxPositions) - 1; i >= 0; i-- {
		pos := boxPositions[i]
		newBoxPos := pos.move(dir)
		w.Grid[newBoxPos.y][newBoxPos.x] = Box
		w.Grid[pos.y][pos.x] = Empty
	}

	return true
}

func (w *Warehouse) moveDoubleBox(pos Coordinate, dir Direction) {
	next := pos.move(dir)
	switch w.Grid[next.y][next.x] {
	case LBox:
		if dir == Up || dir == Down {
			w.moveDoubleBox(Coordinate{next.x + 1, next.y}, dir)
		}
		w.moveDoubleBox(next, dir)
	case RBox:
		if dir == Up || dir == Down {
			w.moveDoubleBox(Coordinate{next.x - 1, next.y}, dir)
		}
		w.moveDoubleBox(next, dir)
	}
	w.Grid[next.y][next.x], w.Grid[pos.y][pos.x] = w.Grid[pos.y][pos.x], w.Grid[next.y][next.x]
}

func (w *Warehouse) processMove(move Direction, isDoubleWidth bool) bool {
	newPos := w.RobotPos.move(move)

	if isDoubleWidth {
		if !w.canMove(w.RobotPos, move) {
			return false
		}
		w.moveDoubleBox(w.RobotPos, move)
		w.RobotPos = newPos
		return true
	}

	if w.Grid[newPos.y][newPos.x] == Wall {
		return false
	}

	if w.Grid[newPos.y][newPos.x] == Empty ||
		w.moveSingleBox(newPos, move) {
		w.Grid[w.RobotPos.y][w.RobotPos.x] = Empty
		w.Grid[newPos.y][newPos.x] = Robot
		w.RobotPos = newPos
		return true
	}

	return false
}

func partOne(wh *Warehouse) int {
	for _, move := range wh.Moves {
		wh.processMove(move, false)
	}
	return wh.CalculateScore(Box)
}

func partTwo(warehouse *Warehouse) int {
	wh := warehouse.Double()
	for _, move := range wh.Moves {
		wh.processMove(move, true)
	}
	return wh.CalculateScore(LBox)
}

func (w *Warehouse) Double() *Warehouse {
	newWidth := w.Width * 2
	newGrid := make([][]Cell, w.Height)

	robotPos := w.RobotPos

	for y := 0; y < w.Height; y++ {
		newGrid[y] = make([]Cell, newWidth)

		for x := 0; x < w.Width; x++ {
			newX := x * 2

			switch w.Grid[y][x] {
			case Wall:
				newGrid[y][newX] = Wall
				newGrid[y][newX+1] = Wall
			case Box:
				newGrid[y][newX] = LBox
				newGrid[y][newX+1] = RBox
			case Empty:
				newGrid[y][newX] = Empty
				newGrid[y][newX+1] = Empty
			case Robot:
				newGrid[y][newX] = Robot
				newGrid[y][newX+1] = Empty
				robotPos = Coordinate{int64(newX), w.RobotPos.y}
			}
		}
	}

	return &Warehouse{
		Grid:     newGrid,
		Moves:    w.Moves,
		RobotPos: robotPos,
		Width:    newWidth,
		Height:   w.Height,
	}
}

func parseInput(inputPath string) (*Warehouse, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}
	defer file.Close()

	warehouse := NewWarehouse()
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
			if cell == Robot {
				warehouse.RobotPos = Coordinate{int64(i), int64(len(grid))}
			}
		}
		grid = append(grid, row)
	}

	warehouse.Grid = grid
	warehouse.Height = len(grid)
	if len(grid) > 0 {
		warehouse.Width = len(grid[0])
	}

	var moves []Direction
	for scanner.Scan() {
		line := scanner.Text()
		for _, ch := range line {
			if ch == '^' || ch == 'v' || ch == '<' || ch == '>' {
				moves = append(moves, Direction(ch))
			}
		}
	}
	warehouse.Moves = moves

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return warehouse, nil
}
