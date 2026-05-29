package model

import (
	"math/rand/v2"
	"time"
)

func (s *Simulation) Run() {
	for {
		// s.Tick()
		s.SpawnCars()
		s.UpdateTrafficLights()
		// s.MoveCars()
		// s.ResolveCollision()
		// s.Render()
		time.Sleep(100 * time.Millisecond)
	}
}
func (s *Simulation) SpawnCars() {
	for _, lane := range s.Lanes {

		if !lane.AllowSpawn {
			continue
		}

		random := rand.Float64()
		if random < lane.SpawnRate.InexactFloat64() {
			s.Cars = append(s.Cars, &Car{
				ID:           len(s.Cars) + 1,
				Position:     0,
				Speed:        0,
				DesiredSpeed: lane.MaxSpeed,
				Lane:         lane.ID,
			})
		}
	}
}

func (s *Simulation) UpdateTrafficLights() {
	for _, intersection := range s.Intersection {
		intersection.TrafficLight.Update()
	}
}
