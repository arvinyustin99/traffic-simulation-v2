package model

import "fmt"

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
	fmt.Printf("Cars: %d\n", len(sim.Cars))
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
		grid[midY][x] = '-'
		grid[midY-1][x] = '-'
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
	for _, lane := range sim.Lanes {
		row := make([]rune, tr.Width)

		// Fill empty road with "-"
		for i := range row {
			row[i] = '-'
		}

		// Draw cars
		for _, car := range lane.Cars {
			if car.Position >= 0 && car.Position < tr.Width {
				row[car.Position] = tr.SpeedSymbol(car.Speed)
			}
		}
		fmt.Println(string(row))
	}
	fmt.Printf("\n")
}
func (tr *TerminalRenderer) SpeedSymbol(speed int) rune {
	switch {
	case speed < 20:
		return 'S'
	case speed < 40:
		return 'M'
	default:
		return 'F'
	}
}

func (tr *TerminalRenderer) RenderStats(sim *Simulation) {

}
