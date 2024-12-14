package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/fatih/color"
)

type Coordinate struct {
	x, y int64
}

func (c Coordinate) add(t int64, v Coordinate) Coordinate {
	return Coordinate{
		x: (c.x + v.x*t),
		y: (c.y + v.y*t),
	}
}

type Robot struct {
	Position Coordinate
	Velocity Coordinate
}

type RobotSwarm struct {
	Width   int64
	Height  int64
	Robots  []Robot
	Display *Display
}

func NewRobotSwarm(width, height int64) *RobotSwarm {
	return &RobotSwarm{
		Width:   width,
		Height:  height,
		Display: NewDisplay(),
	}
}

func (s *RobotSwarm) AddRobot(r Robot) {
	s.Robots = append(s.Robots, r)
}

func (s *RobotSwarm) calculateQuadrants(time int64) [4]int {
	midX, midY := s.Width/2, s.Height/2
	quadrants := [4]int{}

	for _, robot := range s.Robots {
		pos := robot.Position.add(time, robot.Velocity)
		x := (pos.x%s.Width + s.Width) % s.Width
		y := (pos.y%s.Height + s.Height) % s.Height

		if x == midX || y == midY {
			continue
		}

		if x > midX {
			if y > midY {
				quadrants[0]++
			} else {
				quadrants[1]++
			}
		} else {
			if y > midY {
				quadrants[2]++
			} else {
				quadrants[3]++
			}
		}
	}
	return quadrants
}

func (s *RobotSwarm) getRobotPositions(second int64) map[Coordinate][]Coordinate {
	positions := make(map[Coordinate][]Coordinate, len(s.Robots))
	for _, robot := range s.Robots {
		pos := robot.Position.add(second, robot.Velocity)
		x := (pos.x%s.Width + s.Width) % s.Width
		y := (pos.y%s.Height + s.Height) % s.Height
		newPos := Coordinate{x: x, y: y}
		positions[newPos] = append(positions[newPos], robot.Velocity)
	}
	return positions
}

type Display struct {
	tree     *color.Color
	box      *color.Color
	snow     *color.Color
	quarters []rune
}

func NewDisplay() *Display {
	return &Display{
		tree: color.New(color.FgGreen),
		box:  color.New(color.FgRed),
		snow: color.New(color.FgWhite),
		quarters: []rune{
			' ', // empty
			'▘', // tl
			'▝', // tr
			'▀', // top full
			'▖', // bl
			'▌', // left half
			'▞', // diag tr + bl
			'▛', // tl, tr, bl
			'▗', // br
			'▚', // tl, br
			'▐', // right half
			'▜', // tr, tl, br
			'▄', // bottom full
			'▙', // tl, bl, br
			'▟', // tr, bl, br
			'█', // full
		},
	}
}

func (d *Display) show(positions map[Coordinate][]Coordinate, width, height int64) {
	for y := int64(0); y < height; y += 2 {
		for x := int64(0); x < width; x += 2 {
			var bits, nBits int
			for dy := int64(0); dy < 2; dy++ {
				for dx := int64(0); dx < 2; dx++ {
					p := Coordinate{x: x + dx, y: y + dy}
					if len(positions[p]) > 0 {
						bits |= 1 << (dy*2 + dx)
						nBits++
					}
				}
			}

			o := d.snow
			switch nBits {
			case 2:
				o = d.box
			case 3, 4:
				o = d.tree
			}
			o.Printf("%c", d.quarters[bits])
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	swarm, err := parseInput("input.txt", 101, 103)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	p1 := partOne(swarm)
	p1Time := time.Since(start)

	start = time.Now()
	p2 := partTwo(swarm)
	p2Time := time.Since(start)

	log.Printf("part1: %d (took %v)", p1, p1Time)
	log.Printf("part2: %d (took %v)", p2, p2Time)
}

func partOne(s *RobotSwarm) int {
	quadrants := s.calculateQuadrants(100)
	result := 1
	for _, count := range quadrants {
		result *= count
	}
	return result
}

func partTwo(s *RobotSwarm) int64 {
	for second := int64(1); second <= 10000; second++ {
		positions := s.getRobotPositions(second)

		allUnique := true
		for _, velocities := range positions {
			if len(velocities) != 1 {
				allUnique = false
				break
			}
		}

		if allUnique {
			fmt.Printf("\nPattern found at second %d:\n", second)
			s.Display.show(positions, s.Width, s.Height)
			return second
		}

		if second%100 == 0 {
			fmt.Printf("\nSecond %d:\n", second)
			s.Display.show(positions, s.Width, s.Height)
		}
	}
	return 0
}

func parseInput(inputPath string, width, height int64) (*RobotSwarm, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}
	defer file.Close()

	swarm := NewRobotSwarm(width, height)
	pattern := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if matches := pattern.FindStringSubmatch(scanner.Text()); matches != nil {
			px, _ := strconv.ParseInt(matches[1], 10, 64)
			py, _ := strconv.ParseInt(matches[2], 10, 64)
			vx, _ := strconv.ParseInt(matches[3], 10, 64)
			vy, _ := strconv.ParseInt(matches[4], 10, 64)

			swarm.AddRobot(Robot{
				Position: Coordinate{x: px, y: py},
				Velocity: Coordinate{x: vx, y: vy},
			})
		}
	}

	return swarm, scanner.Err()
}
