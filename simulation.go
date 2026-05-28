package main

import "math/rand/v2"

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
