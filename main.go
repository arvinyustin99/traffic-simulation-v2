package main

import "time"

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
