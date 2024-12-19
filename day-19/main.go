package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Onsen struct {
	towels  []string
	designs []string
}

func NewOnsen() *Onsen {
	return &Onsen{
		towels:  make([]string, 0),
		designs: make([]string, 0),
	}
}

func (o *Onsen) solve(design string, cache map[string]int) (bool, int) {
	if design == "" {
		return true, 1
	}

	if count, exists := cache[design]; exists {
		return true, count
	}

	count := 0
	for _, towel := range o.towels {
		if !strings.HasPrefix(design, towel) {
			continue
		}

		if possible, subCount := o.solve(design[len(towel):], cache); possible {
			count += subCount
		}
	}

	cache[design] = count
	return count > 0, count
}

func main() {
	onsen, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(onsen)
	p2 := partTwo(onsen)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(o *Onsen) int {
	count := 0

	for _, design := range o.designs {
		possible, _ := o.solve(design, make(map[string]int))
		if possible {
			count++
		}
	}
	return count
}

func partTwo(o *Onsen) int {
	total := 0

	for _, design := range o.designs {
		_, ways := o.solve(design, make(map[string]int))
		total += ways
	}
	return total
}

func parseInput(inputPath string) (*Onsen, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	onsen := NewOnsen()
	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		onsen.towels = strings.Split(scanner.Text(), ", ")
	}

	scanner.Scan()

	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			onsen.designs = append(onsen.designs, line)
		}
	}

	return onsen, scanner.Err()
}
