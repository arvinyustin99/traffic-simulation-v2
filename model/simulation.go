package model

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/shopspring/decimal"
)

func NewSimulation() *Simulation {
	lane_WE := &Lane{
		ID:         0,
		Cars:       make([]*Car, 0),
		AllowSpawn: true,
		SpawnRate:  decimal.NewFromFloat(0.1),
		MaxSpeed:   2,
		Direction:  Eastbound,
	}
	lane_EW := &Lane{
		ID:         1,
		Cars:       make([]*Car, 0),
		AllowSpawn: true,
		SpawnRate:  decimal.NewFromFloat(0.1),
		MaxSpeed:   1,
		Direction:  Westbound,
	}
	intersection_01 := NewIntersection()
	intersection_01.RegisterLane(lane_WE, true)
	intersection_01.RegisterLane(lane_EW, false)

	return &Simulation{
		TickRate: 2000 * time.Millisecond,
		Running:  false,
		Cars:     make([]*Car, 0),
		Lanes: []*Lane{
			lane_WE, lane_EW,
		},
		Intersections: []*Intersection{
			intersection_01,
		},
		MaxCars:      2,
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
	s.ResolveIntersectionCrossing()
	// s.DetectCongestion()
	// s.CollectMetrics()
	s.EmitEvents()
	s.Render()
}
func (s *Simulation) SpawnCars() {
	if len(s.Cars) == s.MaxCars {
		return
	}
	for _, lane := range s.Lanes {

		if !lane.AllowSpawn {
			continue
		}

		random := rand.Float64()
		if random < lane.SpawnRate.InexactFloat64() {
			newCar := Car{
				ID:           len(s.Cars) + 1,
				Position:     0,
				Speed:        0,
				DesiredSpeed: lane.MaxSpeed,
				Lane:         lane.ID,
			}
			s.Cars = append(s.Cars, &newCar)
			lane.Cars = append(lane.Cars, &newCar)
		}
	}
}

func (s *Simulation) EmitEvents() {
	events := s.EventBus.Flush()
	for _, e := range events {
		fmt.Println(e)
	}
}
