package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Instruction struct {
	Opcode  int
	Operand int
}

type Computer struct {
	A, B, C int
	IP      int
	Prog    []Instruction
}

func (c *Computer) Run() []int {
	var out []int
	c.IP = 0
	for c.IP < len(c.Prog) {
		in := c.Prog[c.IP]
		v := c.OperandValue(in.Operand)
		switch in.Opcode {
		case 0:
			c.A >>= v
		case 1:
			c.B = c.B ^ in.Operand
		case 2:
			c.B = v % 8
		case 3:
			if c.A != 0 {
				c.IP = in.Operand
				continue
			}
		case 4:
			c.B = c.B ^ c.C
		case 5:
			out = append(out, v%8)
		case 6:
			c.B = c.A >> v
		case 7:
			c.C = c.A >> v
		}
		c.IP++
	}
	return out
}

func (c *Computer) OperandValue(op int) int {
	switch op {
	case 0, 1, 2, 3:
		return op
	case 4:
		return c.A
	case 5:
		return c.B
	case 6:
		return c.C
	}
	return 0
}

type ProgramInput struct {
	Comp         Computer
	Instructions []Instruction
}

func convertProgramToOutput(prog []Instruction) []int {
	res := make([]int, len(prog)*2)
	for i, x := range prog {
		res[i*2] = x.Opcode
		res[i*2+1] = x.Operand
	}
	return res
}

func main() {
	pr, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}
	p1 := partOne(pr)
	p2 := partTwo(pr)

	log.Printf("part1: %s", p1)
	log.Printf("part2: %d", p2)
}

func partOne(p *ProgramInput) string {
	c := Computer{p.Comp.A, p.Comp.B, p.Comp.C, p.Comp.IP, p.Instructions}
	o := c.Run()
	s := make([]string, len(o))
	for i, n := range o {
		s[i] = strconv.Itoa(n)
	}
	return strings.Join(s, ",")
}

func partTwo(p *ProgramInput) int {
	exp := convertProgramToOutput(p.Instructions)
	if len(p.Instructions) < 1 || p.Instructions[len(p.Instructions)-1].Opcode != 3 ||
		p.Instructions[len(p.Instructions)-1].Operand != 0 {
		return -1
	}

	base := p.Instructions[:len(p.Instructions)-1]
	getFirst := func(a int) int {
		comp := Computer{a, 0, 0, 0, base}
		out := comp.Run()
		if len(out) == 0 {
			return -1
		}
		return out[0]
	}

	a, ok := attemptReconstruct(exp, 0, getFirst)
	if !ok {
		return -1
	}

	if confirmReconstruction(p.Instructions, a) {
		return a
	}
	return -1
}

func attemptReconstruct(required []int, current int, testA func(int) int) (int, bool) {
	if len(required) == 0 {
		return current, true
	}
	lastIndex := len(required) - 1
	desired := required[lastIndex]

	base := current << 3
	for suffix := 0; suffix < 8; suffix++ {
		candidateA := base | suffix
		firstOut := testA(candidateA)
		if firstOut == desired {
			if found, ok := attemptReconstruct(required[:lastIndex], candidateA, testA); ok {
				return found, true
			}
		}
	}
	return 0, false
}

func confirmReconstruction(prog []Instruction, a int) bool {
	c := Computer{a, 0, 0, 0, prog}
	res := c.Run()
	expect := convertProgramToOutput(prog)
	return reflect.DeepEqual(res, expect)
}

func parseInput(path string) (*ProgramInput, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening: %w", err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	p := &ProgramInput{}
	for i := 0; i < 3 && s.Scan(); i++ {
		line := s.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			continue
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("reg parse: %w", err)
		}
		switch parts[0] {
		case "Register A":
			p.Comp.A = n
		case "Register B":
			p.Comp.B = n
		case "Register C":
			p.Comp.C = n
		}
	}
	s.Scan()
	for s.Scan() {
		line := s.Text()
		if !strings.HasPrefix(line, "Program: ") {
			continue
		}
		in := strings.Split(strings.TrimPrefix(line, "Program: "), ",")
		for i := 0; i < len(in)-1; i += 2 {
			op, err := strconv.Atoi(in[i])
			if err != nil {
				return nil, fmt.Errorf("op parse: %w", err)
			}
			arg, err := strconv.Atoi(in[i+1])
			if err != nil {
				return nil, fmt.Errorf("arg parse: %w", err)
			}
			p.Instructions = append(p.Instructions, Instruction{op, arg})
		}
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	p.Comp.Prog = p.Instructions
	return p, nil
}
