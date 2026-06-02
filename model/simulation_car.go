package model

import "sort"

// BuildSpatialIndex bu ilds a spatial index for cars (sorted by position ascending)
func (s *Simulation) BuildSpatialIndex() {
	// Reset lane.Cars to rebuild from current s.Cars state (supports lane changes)
	for _, lane := range s.Lanes {
		lane.Cars = []*Car{}
	}
	for _, car := range s.Cars {
		for _, lane := range s.Lanes {
			if lane.ID == car.Lane {
				lane.Cars = append(lane.Cars, car)
				break
			}
		}
	}
	for _, lane := range s.Lanes {
		sort.Slice(lane.Cars, func(i, j int) bool {
			return lane.Cars[i].Position < lane.Cars[j].Position
		})
	}
}

// Compute car intentions
func (s *Simulation) ComputeCarIntentions() {
	s.BuildSpatialIndex()

	for _, lane := range s.Lanes {
		// Iterate from front (highest index) to back (lowest index)
		for i := len(lane.Cars) - 1; i >= 0; i-- {
			car := lane.Cars[i]

			intent := CarIntent{
				TargetSpeed: car.Speed,
			}

			var frontCar *Car
			if i < len(lane.Cars)-1 {
				frontCar = lane.Cars[i+1]
			}

			if frontCar != nil {
				gap := frontCar.Position - car.Position
				safeDist := car.SafeDistance
				if safeDist <= 0 {
					safeDist = s.SafeDistance // Default fallback
				}

				speedA := car.Speed
				speedB := frontCar.Speed
				targetSpeedB := frontCar.Intent.TargetSpeed
				isBDecelerating := targetSpeedB < speedB

				if gap > safeDist {
					// 1. If A's speed > B's speed and longer than SafeDistance, A just needs to decelerate to B's speed
					if speedA > speedB {
						if isBDecelerating {
							// 3. If B is decelerating, A decelerates more aggressively to B's target speed
							intent.TargetSpeed = targetSpeedB
						} else {
							intent.TargetSpeed = speedB
						}
					} else {
						// A is safe and going slower/equal to B, can accelerate
						if intent.TargetSpeed < car.DesiredSpeed {
							intent.TargetSpeed++
						}
					}
				} else {
					// Within Safe Distance (gap <= safeDist)
					if speedA > speedB {
						// 2. If A's speed > B's speed and within Safe Distance, A needs to decelerate quickly (lower than B's speed)
						if isBDecelerating {
							// 3. If B is decelerating, A must decelerate more aggressively
							intent.TargetSpeed = targetSpeedB - 1
						} else {
							intent.TargetSpeed = speedB - 1
						}
					} else {
						// SpeedA <= SpeedB but within Safe Distance.
						// To restore safe distance, go slightly slower than B (unless B is stopped)
						if isBDecelerating {
							intent.TargetSpeed = targetSpeedB - 1
						} else {
							if speedB > 0 {
								intent.TargetSpeed = speedB - 1
							} else {
								intent.TargetSpeed = 0
							}
						}
					}
				}
			} else {
				// No car in front, accelerate/decelerate towards desired speed
				if intent.TargetSpeed < car.DesiredSpeed {
					intent.TargetSpeed++
				} else if intent.TargetSpeed > car.DesiredSpeed {
					intent.TargetSpeed--
				}
			}

			// Safety limits
			if intent.TargetSpeed < 0 {
				intent.TargetSpeed = 0
			}
			if intent.TargetSpeed > car.DesiredSpeed {
				intent.TargetSpeed = car.DesiredSpeed
			}

			car.Intent = intent
		}
	}
}

// Find front car in the same lane
func (s *Simulation) FindFrontCar(car *Car) *Car {
	for _, lane := range s.Lanes {
		if lane.ID == car.Lane {
			for _, c := range lane.Cars {
				if c.Position > car.Position {
					return c
				}
			}
		}
	}
	return nil
}

// CleanUpOffscreenCars removes cars that have moved beyond the viewport boundaries
func (s *Simulation) CleanUpOffscreenCars() {
	var width, height int = 60, 30 // default fallbacks
	if tr, ok := s.Renderer.(*TerminalRenderer); ok {
		width = tr.Width
		height = tr.Height
	}

	// Filter s.Cars in-place to avoid reallocating
	n := 0
	for _, car := range s.Cars {
		var laneLimit int
		var carLane *Lane
		for _, l := range s.Lanes {
			if l.ID == car.Lane {
				carLane = l
				break
			}
		}

		if carLane != nil && (carLane.Direction == Northbound || carLane.Direction == Southbound) {
			laneLimit = height
		} else {
			laneLimit = width
		}

		if car.Position < laneLimit {
			s.Cars[n] = car
			n++
		}
	}
	// Zero out remaining pointers to avoid memory leak
	for i := n; i < len(s.Cars); i++ {
		s.Cars[i] = nil
	}
	s.Cars = s.Cars[:n]

	// Also filter each lane's Cars in-place
	for _, lane := range s.Lanes {
		var laneLimit int
		if lane.Direction == Northbound || lane.Direction == Southbound {
			laneLimit = height
		} else {
			laneLimit = width
		}

		m := 0
		for _, car := range lane.Cars {
			if car.Position < laneLimit {
				lane.Cars[m] = car
				m++
			}
		}
		// Zero out remaining pointers to avoid memory leak
		for i := m; i < len(lane.Cars); i++ {
			lane.Cars[i] = nil
		}
		lane.Cars = lane.Cars[:m]
	}
}

func (s *Simulation) CountActiveCars() int {
	var (
		count int = 0
	)
	for _, car := range s.Cars {
		if car != nil {
			count++
		}
	}
	return count
}
