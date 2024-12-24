package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
		gateType:  gateType,
		input1:    input1,
		input2:    input2,
		output:    output,
		evaluated: false,
	})
}

func (c *Circuit) setWireValue(wire string, value int) {
	c.wireValues[wire] = value
}

func (c *Circuit) getInputNumber(prefix string) int {
	num := 0
	for wire, val := range c.wireValues {
		if strings.HasPrefix(wire, prefix) {
			pos := 0
			fmt.Sscanf(wire[1:], "%d", &pos)
			if val > 0 {
				num |= 1 << pos
			}
		}
	}
	return num
}

func (c *Circuit) copy() *Circuit {
	newCircuit := NewCircuit()
	for wire, value := range c.wireValues {
		newCircuit.wireValues[wire] = value
	}
	for _, gate := range c.gates {
		newGate := &Gate{
			gateType:  gate.gateType,
			input1:    gate.input1,
			input2:    gate.input2,
			output:    gate.output,
			evaluated: false,
		}
		newCircuit.gates = append(newCircuit.gates, newGate)
	}
	return newCircuit
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

func main() {
	circuit, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1 := partOne(circuit)
	log.Printf("part1: %d", p1)
}

func partOne(circuit *Circuit) int {
	circuit.evaluate()
	return circuit.getResult()
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
