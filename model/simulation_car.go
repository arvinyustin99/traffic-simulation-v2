package model

import "sort"

// BuildSpatialIndex builds a spatial index for cars (sorted by position ascending or descending based on direction)
func (s *Simulation) BuildSpatialIndex() {
	// Reset lane.Cars to rebuild from current s.Cars state (supports lane changes)
	for _, lane := range s.Lanes {
		lane.Cars = []*Car{}
	}
	for _, car := range s.Cars {
		if car == nil {
			continue
		}
		for _, lane := range s.Lanes {
			if lane.ID == car.Lane {
				lane.Cars = append(lane.Cars, car)
				break
			}
		}
	}
	for _, lane := range s.Lanes {
		if lane.Direction == Westbound || lane.Direction == Northbound {
			// Decreasing direction: sort descending (largest position at the back/start to smallest at the front/end)
			sort.Slice(lane.Cars, func(i, j int) bool {
				return lane.Cars[i].Position > lane.Cars[j].Position
			})
		} else {
			// Increasing direction: sort ascending (smallest position at the back/start to largest at the front/end)
			sort.Slice(lane.Cars, func(i, j int) bool {
				return lane.Cars[i].Position < lane.Cars[j].Position
			})
		}
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
				var gap int
				if lane.Direction == Westbound || lane.Direction == Northbound {
					// Decreasing position: A is behind B, so A.Position > B.Position
					gap = car.Position - frontCar.Position
				} else {
					// Increasing position: A is behind B, so B.Position > A.Position
					gap = frontCar.Position - car.Position
				}
				if gap < 0 {
					gap = 0
				}

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

// Find front car in the same lane (direction aware)
func (s *Simulation) FindFrontCar(car *Car) *Car {
	var frontCar *Car
	for _, lane := range s.Lanes {
		if lane.ID == car.Lane {
			isDecreasing := lane.Direction == Westbound || lane.Direction == Northbound
			for _, c := range lane.Cars {
				if c == car {
					continue
				}
				if isDecreasing {
					// Moving right to left, front car has smaller position
					if c.Position < car.Position {
						if frontCar == nil || c.Position > frontCar.Position {
							frontCar = c
						}
					}
				} else {
					// Moving left to right, front car has larger position
					if c.Position > car.Position {
						if frontCar == nil || c.Position < frontCar.Position {
							frontCar = c
						}
					}
				}
			}
		}
	}
	return frontCar
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
		if car == nil {
			continue
		}
		var carLane *Lane
		for _, l := range s.Lanes {
			if l.ID == car.Lane {
				carLane = l
				break
			}
		}

		keep := true
		if carLane != nil {
			switch carLane.Direction {
			case Eastbound:
				if car.Position >= width {
					keep = false
				}
			case Westbound:
				if car.Position <= 0 {
					keep = false
				}
			case Southbound:
				if car.Position >= height {
					keep = false
				}
			case Northbound:
				if car.Position <= 0 {
					keep = false
				}
			}
		}

		if keep {
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
		m := 0
		for _, car := range lane.Cars {
			keep := true
			switch lane.Direction {
			case Eastbound:
				if car.Position >= width {
					keep = false
				}
			case Westbound:
				if car.Position <= 0 {
					keep = false
				}
			case Southbound:
				if car.Position >= height {
					keep = false
				}
			case Northbound:
				if car.Position <= 0 {
					keep = false
				}
			}

			if keep {
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
