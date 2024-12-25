package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Gate struct {
	gateType  string
	input1    string
	input2    string
	output    string
	evaluated bool
}

type Circuit struct {
	gates      []*Gate
	wireValues map[string]int
}

func NewCircuit() *Circuit {
	return &Circuit{
		gates:      make([]*Gate, 0),
		wireValues: make(map[string]int),
	}
}

func (c *Circuit) addGate(gateType, input1, input2, output string) {
	c.gates = append(c.gates, &Gate{
		gateType: gateType,
		input1:   input1,
		input2:   input2,
		output:   output,
	})
}

func (c *Circuit) setWireValue(wire string, value int) {
	c.wireValues[wire] = value
}

func (c *Circuit) evaluate() {
	for _, gate := range c.gates {
		gate.evaluated = false
	}

	for {
		progress := false
		allEvaluated := true

		for _, gate := range c.gates {
			if gate.evaluated {
				continue
			}

			in1, ok1 := c.wireValues[gate.input1]
			in2, ok2 := c.wireValues[gate.input2]

			if !ok1 || !ok2 {
				allEvaluated = false
				continue
			}

			var result int
			switch gate.gateType {
			case "AND":
				result = in1 & in2
			case "OR":
				result = in1 | in2
			case "XOR":
				result = in1 ^ in2
			}

			c.wireValues[gate.output] = result
			gate.evaluated = true
			progress = true
		}

		if !progress || allEvaluated {
			break
		}
	}
}

func (c *Circuit) getResult() int {
	result := 0
	maxPos := 0

	for wire := range c.wireValues {
		if strings.HasPrefix(wire, "z") {
			pos := 0
			fmt.Sscanf(wire[1:], "%d", &pos)
			if pos > maxPos {
				maxPos = pos
			}
		}
	}

	for i := 0; i <= maxPos; i++ {
		wire := fmt.Sprintf("z%02d", i)
		if val, ok := c.wireValues[wire]; ok {
			result |= val << i
		}
	}

	return result
}

func (c *Circuit) findBrokenConnections() string {
	inputMap := c.buildInputMap()
	broken := make(map[string]bool)

	for _, gate := range c.gates {
		switch gate.gateType {
		case "AND":
			if !c.isFirstAND(gate) && !inputMap[makeKey(gate.output, "OR")] {
				broken[gate.output] = true
			}
		case "OR":
			if c.isOutputWire(gate.output) && gate.output != "z45" {
				broken[gate.output] = true
			}
		case "XOR":
			if c.isFirstLevel(gate) {
				if !c.isFirstXOR(gate) && !inputMap[makeKey(gate.output, "XOR")] {
					broken[gate.output] = true
				}
			} else if !c.isOutputWire(gate.output) {
				broken[gate.output] = true
			}
		}
	}

	return c.formatResult(broken)
}

func (c *Circuit) buildInputMap() map[string]bool {
	m := make(map[string]bool)
	for _, g := range c.gates {
		m[makeKey(g.input1, g.gateType)] = true
		m[makeKey(g.input2, g.gateType)] = true
	}
	return m
}

func (c *Circuit) isFirstAND(g *Gate) bool {
	return g.input1 == "x00" || g.input2 == "x00"
}

func (c *Circuit) isFirstXOR(g *Gate) bool {
	return g.input1 == "x00" || g.input2 == "x00"
}

func (c *Circuit) isFirstLevel(g *Gate) bool {
	return strings.HasPrefix(g.input1, "x") || strings.HasPrefix(g.input2, "x")
}

func (c *Circuit) isOutputWire(wire string) bool {
	return strings.HasPrefix(wire, "z")
}

func (c *Circuit) formatResult(broken map[string]bool) string {
	var wires []string
	for wire := range broken {
		wires = append(wires, wire)
	}
	sort.Strings(wires)
	return strings.Join(wires, ",")
}

func makeKey(wire, gateType string) string {
	return fmt.Sprintf("%s,%s", wire, gateType)
}

func main() {
	circuit, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(circuit)
	p2 := partTwo(circuit)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %s", p2)
}

func partOne(c *Circuit) int {
	c.evaluate()
	return c.getResult()
}

func partTwo(c *Circuit) string {
	return c.findBrokenConnections()
}

func parseInput(path string) (*Circuit, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}
	defer file.Close()

	circuit := NewCircuit()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			continue
		}

		value := 0
		fmt.Sscanf(parts[1], "%d", &value)
		circuit.setWireValue(parts[0], value)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, " -> ")
		if len(parts) != 2 {
			continue
		}

		gateParts := strings.Split(parts[0], " ")
		if len(gateParts) != 3 {
			continue
		}

		circuit.addGate(gateParts[1], gateParts[0], gateParts[2], parts[1])
	}

	return circuit, scanner.Err()
}
