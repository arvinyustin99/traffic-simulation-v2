package model

import (
	"fmt"
	"strings"
)

type Renderer interface {
	Render(*Simulation)
}

func (s *Simulation) Render() {
	s.Renderer.Render(s)
}

type TerminalRenderer struct {
	Width  int
	Height int
}

func InitTerminal() {
	fmt.Print("\033[2J")
}

func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{
		Height: 30,
		Width:  60,
	}
}

func (tr *TerminalRenderer) Render(sim *Simulation) {
	fmt.Print("\033[2J")
	// Move cursor to top-left
	fmt.Print("\033[H")

	tr.RenderHeader(sim)
	tr.RenderIntersection(sim)
	tr.RenderLanes(sim)
	tr.RenderStats(sim)

}

func (tr *TerminalRenderer) RenderHeader(sim *Simulation) {
	fmt.Printf("Traffic Simulation v1.0\n")
	fmt.Printf("Tick: %d\n", sim.Tick)
	fmt.Printf("Cars: %d\n", sim.CountActiveCars())
}

func (tr *TerminalRenderer) RenderIntersection(sim *Simulation) {
	grid := make([][]rune, tr.Height)
	for y := range grid {
		grid[y] = make([]rune, tr.Width)
		for x := range grid[y] {
			grid[y][x] = ' '
		}
	}

	midY := tr.Height / 2
	midX := tr.Width / 2

	// Draw horizontal road
	for x := 0; x < tr.Width; x++ {
		grid[midY][x] = '_'
		grid[midY-1][x] = '1'
	}

	// Draw vertical road
	for y := 0; y < tr.Height; y++ {
		grid[y][midX] = '|'
		grid[y][midX-1] = '|'
	}

	// Draw intersection
	grid[midY][midX] = '+'
	grid[midY][midX+1] = '+'
	grid[midY-1][midX] = '+'
	grid[midY-1][midX+1] = '+'

}

func (tr *TerminalRenderer) RenderLanes(sim *Simulation) {
	var (
		isSameDirection bool = true
	)
	for i, lane := range sim.Lanes {
		var (
			isLeftMostLane  bool = i == 0
			leftMostFence   []string
			isRightMostLane bool = i == len(sim.Lanes)-1
			rightMostFence  []string
		)
		if i == 0 {
			isSameDirection = true
		} else if lane.Direction != sim.Lanes[i-1].Direction {
			isSameDirection = false
		}

		row := make([]string, tr.Width)

		// Fill empty road with "-"
		for j := range row {
			if isLeftMostLane {
				leftMostFence = append(leftMostFence, "+")
			}
			if isRightMostLane {
				rightMostFence = append(rightMostFence, "+")
			}
			row[j] = " "
		}
		if isLeftMostLane {
			fmt.Println(strings.Join(leftMostFence, ""))
		}
		// Draw cars
		for _, car := range lane.Cars {
			if car.Position >= 0 && car.Position < tr.Width {

				row[car.Position] = tr.SpeedSymbol(car.Speed)
			}
		}
		if !isSameDirection {
			for range tr.Width {
				fmt.Printf("-")
			}
			fmt.Println()
		}
		fmt.Println(strings.Join(row, ""))
		if isRightMostLane {
			fmt.Println(strings.Join(rightMostFence, ""))
		}
	}
	fmt.Printf("\n")
}
func (tr *TerminalRenderer) SpeedSymbol(speed int) string {
	switch {
	case speed < 1:
		return "S"
	case speed < 2:
		return "M"
	default:
		return "F"
	}
}

func (tr *TerminalRenderer) RenderStats(sim *Simulation) {

}
