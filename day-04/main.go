package main

import (
	"bufio"
	"log"
	"os"
)

type Direction struct {
	dx, dy int
}

type Position struct {
	row, col int
	dir      Direction
}

var (
	directions = []Direction{
		{0, 1},   // right
		{1, 0},   // down
		{1, 1},   // diagonal down-right
		{-1, 1},  // diagonal up-right
		{0, -1},  // left
		{-1, 0},  // up
		{-1, -1}, // diagonal up-left
		{1, -1},  // diagonal down-left
	}

	xCorners = []struct{ row, col int }{
		{-1, -1}, // topLeft
		{-1, 1},  // topRight
		{1, -1},  // bottomLeft
		{1, 1},   // bottomRight
	}
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var grid [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, []rune(scanner.Text()))
	}

	p1 := partOne(grid)
	p2 := partTwo(grid)
	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(grid [][]rune) int {
	return len(findWord(grid, "XMAS"))
}

func partTwo(grid [][]rune) int {
	if len(grid) < 3 { // X pattern needs at least 3x3 space
		return 0
	}

	count := 0
	for row := 1; row < len(grid)-1; row++ {
		for col := 1; col < len(grid[0])-1; col++ {
			if grid[row][col] == 'A' && isValidXMAS(grid, row, col) {
				count++
			}
		}
	}
	return count
}

func findWord(grid [][]rune, word string) []Position {
	if len(grid) == 0 {
		return nil
	}

	searchRunes := []rune(word)
	var results []Position

	for row := range grid {
		for col := range grid[row] {
			for _, dir := range directions {
				if checkWord(grid, searchRunes, row, col, dir) {
					results = append(results, Position{row, col, dir})
				}
			}
		}
	}
	return results
}

func checkWord(grid [][]rune, word []rune, startRow, startCol int, dir Direction) bool {
	endRow, endCol := startRow+dir.dx*(len(word)-1), startCol+dir.dy*(len(word)-1)

	if endRow < 0 || endRow >= len(grid) || endCol < 0 || endCol >= len(grid[0]) {
		return false
	}

	for i := range word {
		if grid[startRow+dir.dx*i][startCol+dir.dy*i] != word[i] {
			return false
		}
	}
	return true
}

func isValidXMAS(grid [][]rune, centerRow, centerCol int) bool {
	var d1, d2 bool

	tl, br := grid[centerRow-1][centerCol-1], grid[centerRow+1][centerCol+1]
	d1 = (tl == 'M' && br == 'S') || (tl == 'S' && br == 'M')

	tr, bl := grid[centerRow-1][centerCol+1], grid[centerRow+1][centerCol-1]
	d2 = (tr == 'M' && bl == 'S') || (tr == 'S' && bl == 'M')

	return d1 && d2
}
