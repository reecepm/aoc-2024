package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Generator struct {
	secret int
}

type changeTracker struct {
	values    []int
	lastPrice int
}

func (ct *changeTracker) addChange(newPrice int) {
	copy(ct.values, ct.values[1:])
	ct.values[3] = newPrice - ct.lastPrice
	ct.lastPrice = newPrice
}

func newChangeTracker(initialPrice int) *changeTracker {
	return &changeTracker{
		values:    make([]int, 4),
		lastPrice: initialPrice,
	}
}

func (ct *changeTracker) getKey() string {
	return fmt.Sprintf("%d,%d,%d,%d", ct.values[0], ct.values[1], ct.values[2], ct.values[3])
}

type MonkeyMarket struct {
	initials []int
}

func NewMonkeyMarket() *MonkeyMarket {
	return &MonkeyMarket{
		initials: make([]int, 0),
	}
}

func (m *MonkeyMarket) findBestSequence() int {
	totalBananas := make(map[string]int)

	for _, initial := range m.initials {
		gen := NewGenerator(initial)
		for seq, price := range gen.findSequences() {
			totalBananas[seq] += price
		}
	}

	maxBananas := 0
	for _, total := range totalBananas {
		if total > maxBananas {
			maxBananas = total
		}
	}
	return maxBananas
}

func (m *MonkeyMarket) calculateDevicePrice() int {
	total := 0
	for _, initial := range m.initials {
		gen := NewGenerator(initial)
		for i := 0; i < 2000; i++ {
			gen.Next()
		}
		total += gen.secret
	}
	return total
}

func NewGenerator(initial int) *Generator {
	return &Generator{secret: initial}
}

func (g *Generator) Next() {
	g.secret = (g.secret ^ (g.secret * 64)) % 16777216
	g.secret = (g.secret ^ (g.secret / 32)) % 16777216
	g.secret = (g.secret ^ (g.secret * 2048)) % 16777216
}

func (g *Generator) findSequences() map[string]int {
	seenSequences := make(map[string]bool)
	sequencePrices := make(map[string]int)
	tracker := newChangeTracker(g.secret % 10)

	for i := 0; i < 3; i++ {
		g.Next()
		tracker.addChange(g.secret % 10)
	}

	for i := 3; i < 2000; i++ {
		g.Next()
		price := g.secret % 10
		tracker.addChange(price)

		key := tracker.getKey()
		if !seenSequences[key] {
			seenSequences[key] = true
			sequencePrices[key] = price
		}
	}

	return sequencePrices
}

func main() {
	market, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(market)
	p2 := partTwo(market)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(market *MonkeyMarket) int {
	return market.calculateDevicePrice()
}

func partTwo(market *MonkeyMarket) int {
	return market.findBestSequence()
}

func parseInput(path string) (*MonkeyMarket, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}
	defer file.Close()

	market := NewMonkeyMarket()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			num, err := strconv.Atoi(line)
			if err != nil {
				return nil, fmt.Errorf("parsing number: %w", err)
			}
			market.initials = append(market.initials, num)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return market, nil
}
