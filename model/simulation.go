package model

import (
	"math/rand/v2"
	"time"

	"github.com/shopspring/decimal"
)

func NewSimulation() *Simulation {
	lane_WE_00 := &Lane{
		ID:         0,
		Cars:       make([]*Car, 0),
		AllowSpawn: true,
		SpawnRate:  decimal.NewFromFloat(0.1),
		MaxSpeed:   2,
		Direction:  Eastbound,
	}
	lane_WE_01 := &Lane{
		ID:         1,
		Cars:       make([]*Car, 0),
		AllowSpawn: true,
		SpawnRate:  decimal.NewFromFloat(0.1),
		MaxSpeed:   2,
		Direction:  Eastbound,
	}
	lane_EW_00 := &Lane{
		ID:         2,
		Cars:       make([]*Car, 0),
		AllowSpawn: true,
		SpawnRate:  decimal.NewFromFloat(0.2),
		MaxSpeed:   1,
		Direction:  Eastbound,
	}
	intersection_01 := NewIntersection()
	intersection_01.RegisterLane(lane_WE_00, true)
	intersection_01.RegisterLane(lane_WE_01, true)
	intersection_01.RegisterLane(lane_EW_00, false)

	return &Simulation{
		TickRate: 1000 * time.Millisecond,
		Running:  false,
		Cars:     make([]*Car, 0, 2), // Preallocate capacity to match MaxCars
		Lanes: []*Lane{
			lane_WE_00,
			lane_WE_01,
			lane_EW_00,
		},
		Intersections: []*Intersection{
			intersection_01,
		},
		MaxCars:      20,
		SafeDistance: 2,
		EventBus:     NewEventBus(),
		Metrics:      NewMetrics(),
		Renderer:     NewTerminalRenderer(),
	}

}
func (s *Simulation) Run() {
	s.Running = true

	ticker := time.NewTicker(s.TickRate)
	defer ticker.Stop()

	for s.Running {
		<-ticker.C

		s.Tick++
		s.Step()
	}
}
func (s *Simulation) Step() {
	s.SpawnCars()
	s.UpdateTrafficLights()
	s.ComputeCarIntentions()
	s.ApplySpeedAdjustments()
	s.MoveCars()
	s.CleanUpOffscreenCars()
	s.ResolveIntersectionCrossing()
	// s.DetectCongestion()
	// s.CollectMetrics()
	s.EmitEvents()
	s.Render()
}
func (s *Simulation) SpawnCars() {
	if len(s.Cars) >= s.MaxCars {
		return
	}
	for _, lane := range s.Lanes {

		if !lane.AllowSpawn {
			continue
		}

		random := rand.Float64()
		if random < lane.SpawnRate.InexactFloat64() {
			var (
				startPosition int = 0
			)
			switch lane.Direction {
			case Eastbound:
				startPosition = 0
			case Westbound:
				startPosition = s.Renderer.(*TerminalRenderer).Width
			case Northbound:
				startPosition = 0
			case Southbound:
				startPosition = s.Renderer.(*TerminalRenderer).Height
			}
			newCar := Car{
				ID:           len(s.Cars) + 1,
				Position:     startPosition,
				Speed:        0,
				SafeDistance: s.SafeDistance,
				DesiredSpeed: lane.MaxSpeed,
				Lane:         lane.ID,
			}
			s.Cars = append(s.Cars, &newCar)
			lane.Cars = append(lane.Cars, &newCar)
		}
	}
}

func (s *Simulation) EmitEvents() {
	// events := s.EventBus.Flush()
	// for _, e := range events {
	// 	fmt.Println(e)
	// }
}
