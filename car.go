package main

import "sort"

func (s *Simulation) BuildSpatialIndex() {
	for _, lane := range s.Lanes {
		sort.Slice(lane.Cars, func(i, j int) bool {
			return lane.Cars[i].Position < lane.Cars[j].Position
		})
	}
}
