package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Coordinate struct {
	x, y int64
}

type Button Coordinate

type Prize Coordinate

type ClawMachine struct {
	ButtonA Button
	ButtonB Button
	Prize   Prize
}

type Arcade struct {
	machines []ClawMachine
}

func NewArcade() *Arcade {
	return &Arcade{
		machines: make([]ClawMachine, 0),
	}
}

func (a *Arcade) AddMachine(m ClawMachine) {
	a.machines = append(a.machines, m)
}

func (a *Arcade) SolvePuzzle(offset int64) int {
	machines := make([]ClawMachine, len(a.machines))
	copy(machines, a.machines)

	if offset > 0 {
		for i := range machines {
			machines[i].Prize.x += offset
			machines[i].Prize.y += offset
		}
	}

	total := 0
	for _, m := range machines {
		if !solutionExists(m.ButtonA, m.ButtonB, m.Prize) {
			continue
		}

		minTokens := findMinTokens(m)
		if minTokens != -1 {
			total += minTokens
		}
	}
	return total
}

func findMinTokens(m ClawMachine) int {
	det := float64(m.ButtonA.x*m.ButtonB.y - m.ButtonB.x*m.ButtonA.y)
	if det == 0 {
		return -1
	}

	a := float64(m.Prize.x*m.ButtonB.y-m.Prize.y*m.ButtonB.x) / det
	b := float64(m.ButtonA.x*m.Prize.y-m.ButtonA.y*m.Prize.x) / det

	if a != float64(int64(a)) || b != float64(int64(b)) || a < 0 || b < 0 {
		return -1
	}

	return int(3*int64(a) + int64(b))
}

func solutionExists(buttonA Button, buttonB Button, prize Prize) bool {
	gcdX := GCD(buttonA.x, buttonB.x)
	gcdY := GCD(buttonA.y, buttonB.y)

	if gcdX == 0 || gcdY == 0 {
		return false
	}

	return prize.x%gcdX == 0 && prize.y%gcdY == 0
}

func GCD(a, b int64) int64 {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}

	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func main() {
	arcade, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(arcade)
	p2 := partTwo(arcade)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(arcade *Arcade) int {
	return arcade.SolvePuzzle(0)
}

func partTwo(arcade *Arcade) int {
	return arcade.SolvePuzzle(10000000000000)
}

func parseInput(inputPath string) (*Arcade, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}
	defer file.Close()

	arcade := NewArcade()
	var currentMachine ClawMachine

	buttonPattern := regexp.MustCompile(`Button ([AB]): X\+(\d+), Y\+(\d+)`)
	prizePattern := regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if buttonMatch := buttonPattern.FindStringSubmatch(line); buttonMatch != nil {
			x, _ := strconv.ParseInt(buttonMatch[2], 10, 64)
			y, _ := strconv.ParseInt(buttonMatch[3], 10, 64)

			if buttonMatch[1] == "A" {
				currentMachine.ButtonA = Button{x: x, y: y}
			} else {
				currentMachine.ButtonB = Button{x: x, y: y}
			}
		} else if prizeMatch := prizePattern.FindStringSubmatch(line); prizeMatch != nil {
			x, _ := strconv.ParseInt(prizeMatch[1], 10, 64)
			y, _ := strconv.ParseInt(prizeMatch[2], 10, 64)
			currentMachine.Prize = Prize{x: x, y: y}

			arcade.AddMachine(currentMachine)
			currentMachine = ClawMachine{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return arcade, nil
}
