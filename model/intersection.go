package model

func (s *Simulation) ResolveIntersectionCrossing() {
	var width, height int = 60, 30 // default fallbacks
	if tr, ok := s.Renderer.(*TerminalRenderer); ok {
		width = tr.Width
		height = tr.Height
	}
	midX := width / 2
	midY := height / 2

	for _, ic := range s.Intersections {
		for _, lane := range ic.IncomingLanes {
			for _, car := range lane.Cars {
				if car == nil {
					continue
				}
				if ic.IsRed(lane.Direction) {
					switch lane.Direction {
					case Eastbound:
						// Moving left-to-right: stops before midX
						stopLine := midX - 2
						if car.Position >= stopLine && car.Position < midX {
							car.Intent.TargetSpeed = 0
						}
					case Westbound:
						// Moving right-to-left: stops before midX
						stopLine := midX + 2
						if car.Position <= stopLine && car.Position > midX {
							car.Intent.TargetSpeed = 0
						}
					case Southbound:
						// Moving top-to-bottom: stops before midY
						stopLine := midY - 2
						if car.Position >= stopLine && car.Position < midY {
							car.Intent.TargetSpeed = 0
						}
					case Northbound:
						// Moving bottom-to-top: stops before midY
						stopLine := midY + 2
						if car.Position <= stopLine && car.Position > midY {
							car.Intent.TargetSpeed = 0
						}
					}
				}
			}
		}
	}
}
