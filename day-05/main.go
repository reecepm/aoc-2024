package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Rule struct {
	beforePage, afterPage int
}

type Update struct {
	pages []int
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	rules, updates, err := parseInput(string(data))
	if err != nil {
		panic(err)
	}

	p1 := partOne(rules, updates)
	p2 := partTwo(rules, updates)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(rules []Rule, updates []Update) int {
	total := 0
	for _, update := range updates {
		if fixed, changed := validateAndFix(update.pages, rules); !changed {
			total += fixed[len(fixed)/2]
		}
	}
	return total
}

func partTwo(rules []Rule, updates []Update) int {
	total := 0
	for _, update := range updates {
		if fixed, changed := validateAndFix(update.pages, rules); changed {
			total += fixed[len(fixed)/2]
		}
	}
	return total
}

func validateAndFix(pages []int, rules []Rule) ([]int, bool) {
	result := make([]int, len(pages))
	copy(result, pages)
	changed := false

	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			for _, rule := range rules {
				if result[i] == rule.afterPage && result[j] == rule.beforePage {
					result[i], result[j] = result[j], result[i]
					changed = true
				}
			}
		}
	}

	return result, changed
}

func parseInput(input string) ([]Rule, []Update, error) {
	parts := strings.Split(strings.TrimSpace(input), "\n\n")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid input format")
	}

	ruleLines := strings.Split(strings.TrimSpace(parts[0]), "\n")
	rules := make([]Rule, 0, len(ruleLines))

	for _, line := range ruleLines {
		nums := strings.Split(line, "|")
		if len(nums) != 2 {
			return nil, nil, fmt.Errorf("invalid rule format: %s", line)
		}

		before, err := strconv.Atoi(nums[0])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid number in rule: %s", nums[0])
		}

		after, err := strconv.Atoi(nums[1])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid number in rule: %s", nums[1])
		}

		rules = append(rules, Rule{
			beforePage: before,
			afterPage:  after,
		})
	}

	updateLines := strings.Split(strings.TrimSpace(parts[1]), "\n")
	updates := make([]Update, 0, len(updateLines))

	for _, line := range updateLines {
		numStrs := strings.Split(line, ",")
		pages := make([]int, 0, len(numStrs))

		for _, numStr := range numStrs {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid number in update: %s", numStr)
			}
			pages = append(pages, num)
		}

		updates = append(updates, Update{pages: pages})
	}

	return rules, updates, nil
}
