package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Stone uint64

type StoneCache struct {
	cache map[cacheKey]uint64
}

type cacheKey struct {
	stone Stone
	depth uint8
}

func NewStoneCache() *StoneCache {
	return &StoneCache{
		cache: make(map[cacheKey]uint64),
	}
}

func (s *StoneCache) process(stone Stone, depth uint8) uint64 {
	if depth == 0 {
		return 1
	}

	key := cacheKey{stone, depth}
	if count, exists := s.cache[key]; exists {
		return count
	}

	var count uint64
	for _, next := range stone.getNextStones() {
		count += s.process(next, depth-1)
	}

	s.cache[key] = count
	return count
}

func (s Stone) getNextStones() []Stone {
	if s == 0 {
		return []Stone{1}
	}

	digits := s.countDigits()
	if digits%2 == 0 {
		exp := uint64(math.Pow10(int(digits) / 2))
		return []Stone{Stone(uint64(s) / exp), Stone(uint64(s) % exp)}
	}

	return []Stone{s * 2024}
}

func (s Stone) countDigits() uint8 {
	if s == 0 {
		return 1
	}
	return uint8(math.Floor(math.Log10(float64(s)))) + 1
}

func main() {
	stones, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(stones)
	p2 := partTwo(stones)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(stones []Stone) uint64 {
	return processStones(stones, 25)
}

func partTwo(stones []Stone) uint64 {
	return processStones(stones, 75)
}

func processStones(stones []Stone, depth uint8) uint64 {
	cache := NewStoneCache()
	var count uint64

	for _, stone := range stones {
		count += cache.process(stone, depth)
	}
	return count
}

func parseInput(inputPath string) ([]Stone, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var stones []Stone
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		for _, numStr := range strings.Fields(scanner.Text()) {
			num, err := strconv.ParseUint(numStr, 10, 64)
			if err != nil {
				return nil, err
			}
			stones = append(stones, Stone(num))
		}
	}

	return stones, scanner.Err()
}
