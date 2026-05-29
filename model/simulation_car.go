package model

import "sort"

func (s *Simulation) BuildSpatialIndex() {
	for _, lane := range s.Lanes {
		sort.Slice(lane.Cars, func(i, j int) bool {
			return lane.Cars[i].Position < lane.Cars[j].Position
		})
	}
}

func (s *Simulation) ComputeCarIntentions() {
	for _, car := range s.Cars {
		frontCar := s.FindFrontCar(car)

		intent := CarIntent{
			TargetSpeed: car.Speed,
		}

		if frontCar != nil {

		}
	}
}
