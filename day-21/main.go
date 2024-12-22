package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Coordinate struct {
	x, y int
}

var directions = map[rune]Coordinate{
	'^': {0, -1},
	'v': {0, 1},
	'<': {-1, 0},
	'>': {1, 0},
}

type Keypad struct {
	buttons map[rune]Coordinate
	start   Coordinate
	width   int
	height  int
}

func NewKeypad(isNumeric bool) *Keypad {
	if isNumeric {
		return &Keypad{
			buttons: map[rune]Coordinate{
				'7': {0, 0}, '8': {1, 0}, '9': {2, 0},
				'4': {0, 1}, '5': {1, 1}, '6': {2, 1},
				'1': {0, 2}, '2': {1, 2}, '3': {2, 2},
				'0': {1, 3}, 'A': {2, 3},
			},
			start:  Coordinate{2, 3},
			width:  3,
			height: 4,
		}
	}
	return &Keypad{
		buttons: map[rune]Coordinate{
			'^': {1, 0}, 'A': {2, 0},
			'<': {0, 1}, 'v': {1, 1}, '>': {2, 1},
		},
		start:  Coordinate{2, 0},
		width:  3,
		height: 2,
	}
}

type DoorCode struct {
	numeric     string
	numericPart int
}

type MemoKey struct {
	input   string
	keypads int
}

type MemoValue struct {
	length int
	ok     bool
}

func (k *Keypad) isValidMove(pos Coordinate) bool {
	if pos.x < 0 || pos.x >= k.width || pos.y < 0 || pos.y >= k.height {
		return false
	}
	_, exists := k.buttons[k.buttonAt(pos)]
	return exists
}

func (k *Keypad) buttonAt(p Coordinate) rune {
	for btn, pos := range k.buttons {
		if pos == p {
			return btn
		}
	}
	return 0
}

func generatePath(from, to Coordinate, moveOrder []rune) string {
	var path string
	dx := to.x - from.x
	dy := to.y - from.y

	for _, dir := range moveOrder {
		switch dir {
		case '^':
			for dy < 0 {
				path += "^"
				dy++
			}
		case 'v':
			for dy > 0 {
				path += "v"
				dy--
			}
		case '<':
			for dx < 0 {
				path += "<"
				dx++
			}
		case '>':
			for dx > 0 {
				path += ">"
				dx--
			}
		}
	}
	return path + "A"
}

func getPossiblePaths(from, to Coordinate) []string {
	moveOrders := [][]rune{
		{'^', 'v', '<', '>'}, // hor
		{'<', '>', '^', 'v'}, //ver
	}

	var paths []string
	for _, order := range moveOrders {
		path := generatePath(from, to, order)
		paths = append(paths, path)
	}

	if len(paths) == 2 && paths[0] == paths[1] {
		return paths[:1]
	}
	return paths
}

func validatePath(start Coordinate, path string, k *Keypad) bool {
	pos := start
	for _, move := range path {
		if move == 'A' {
			continue
		}
		pos.x += directions[rune(move)].x
		pos.y += directions[rune(move)].y
		if !k.isValidMove(pos) {
			return false
		}
	}
	return true
}

func solve(input string, keypads []*Keypad, memo map[MemoKey]MemoValue) (int, bool) {
	key := MemoKey{input: input, keypads: len(keypads)}
	if cached, ok := memo[key]; ok {
		return cached.length, cached.ok
	}

	if len(keypads) == 0 {
		memo[key] = MemoValue{length: len(input), ok: true}
		return len(input), true
	}

	k := keypads[0]
	pos := k.start
	var length int

	for _, c := range input {
		target := k.buttons[c]
		paths := getPossiblePaths(pos, target)
		var shortest int
		var found bool

		for _, path := range paths {
			if !validatePath(pos, path, k) {
				continue
			}

			if len(keypads) == 1 {
				shortest = len(path)
				found = true
				break
			}

			if subLength, ok := solve(path, keypads[1:], memo); ok {
				if !found || subLength < shortest {
					shortest = subLength
					found = true
				}
			}
		}

		if !found {
			memo[key] = MemoValue{length: 0, ok: false}
			return 0, false
		}

		length += shortest
		pos = target
	}

	memo[key] = MemoValue{length: length, ok: true}
	return length, true
}

func solveForKeypads(codes []DoorCode, numDirKeypads int) int {
	numKeypad := NewKeypad(true)
	var keypads []*Keypad
	keypads = append(keypads, numKeypad)
	for i := 0; i < numDirKeypads; i++ {
		keypads = append(keypads, NewKeypad(false))
	}

	total := 0
	memo := make(map[MemoKey]MemoValue)
	for _, code := range codes {
		if length, ok := solve(code.numeric, keypads, memo); ok {
			total += length * code.numericPart
		}
	}
	return total
}

func main() {
	codes, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(codes)
	p2 := partTwo(codes)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(codes []DoorCode) int {
	return solveForKeypads(codes, 2)
}

func partTwo(codes []DoorCode) int {
	return solveForKeypads(codes, 25)
}

func parseInput(path string) ([]DoorCode, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}
	defer file.Close()

	var codes []DoorCode
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}

		value, err := strconv.Atoi(text[:len(text)-1])
		if err != nil {
			return nil, fmt.Errorf("parsing numeric value: %w", err)
		}

		codes = append(codes, DoorCode{
			numeric:     text,
			numericPart: value,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return codes, nil
}
