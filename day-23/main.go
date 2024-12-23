package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Network struct {
	connections map[string][]string
}

func NewNetwork() *Network {
	return &Network{
		connections: make(map[string][]string),
	}
}

func (n *Network) addConnection(from, to string) {
	n.connections[from] = append(n.connections[from], to)
	n.connections[to] = append(n.connections[to], from)
}

func (n *Network) findTriples() [][]string {
	var triples [][]string
	seen := make(map[string]bool)

	for computer := range n.connections {
		for _, conn1 := range n.connections[computer] {
			for _, conn2 := range n.connections[conn1] {
				triple := []string{computer, conn1, conn2}
				key := makeTripleKey(triple)
				if seen[key] {
					continue
				}

				if n.areConnected(computer, conn2) {
					sort.Strings(triple)
					triples = append(triples, triple)
					seen[key] = true
				}
			}
		}
	}

	return triples
}

func (n *Network) findTriplesWithT() [][]string {
	var result [][]string
	for _, triple := range n.findTriples() {
		for _, computer := range triple {
			if strings.HasPrefix(computer, "t") {
				result = append(result, triple)
				break
			}
		}
	}
	return result
}

func (n *Network) findLargestConnectedGroup() []string {
	connections := make(map[string]map[string]bool)
	for comp, conns := range n.connections {
		connections[comp] = make(map[string]bool)
		for _, conn := range conns {
			connections[comp][conn] = true
		}
	}

	largestGroup := []string{}
	computers := make([]string, 0, len(connections))
	for comp := range connections {
		computers = append(computers, comp)
	}

	for _, start := range computers {
		group := []string{start}
		var candidates []string
		for neighbor := range connections[start] {
			candidates = append(candidates, neighbor)
		}

		for _, candidate := range candidates {
			isConnectedToAll := true
			for _, member := range group {
				if !connections[candidate][member] {
					isConnectedToAll = false
					break
				}
			}

			if isConnectedToAll {
				group = append(group, candidate)
			}
		}

		if len(group) > len(largestGroup) {
			largestGroup = make([]string, len(group))
			copy(largestGroup, group)
		}
	}

	sort.Strings(largestGroup)
	return largestGroup
}

func (n *Network) areConnected(comp1, comp2 string) bool {
	for _, conn := range n.connections[comp1] {
		if conn == comp2 {
			return true
		}
	}
	return false
}

func makeTripleKey(triple []string) string {
	sort.Strings(triple)
	return strings.Join(triple, ",")
}

func main() {
	network, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := len(network.findTriplesWithT())
	p2 := strings.Join(network.findLargestConnectedGroup(), ",")

	log.Printf("part1: %d", p1)
	log.Printf("part2: %s", p2)
}

func parseInput(path string) (*Network, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}
	defer file.Close()

	network := NewNetwork()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			parts := strings.Split(line, "-")
			if len(parts) == 2 {
				network.addConnection(parts[0], parts[1])
			}
		}
	}

	return network, scanner.Err()
}
