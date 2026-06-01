package model

func (s *Simulation) ResolveIntersectionCrossing() {
	for _, ic := range s.Intersections {
		for _, lane := range ic.IncomingLanes {
			for _, car := range lane.Cars {
				if ic.IsRed(lane.Direction) {
					if car.Position >= ic.StopLine {
						car.Intent.TargetSpeed = 0
					}
				}
			}
		}
	}
}
