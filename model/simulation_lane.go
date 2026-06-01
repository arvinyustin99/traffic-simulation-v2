package model

type LaneChangedEvent struct {
	OldLane int
	NewLane int
	CarID   int
}

func (s *Simulation) ResolveLaneChanges() {
	for _, car := range s.Cars {
		if !car.Intent.ChangeLane {
			continue
		}
		if s.ChangeLane {
			oldLane := car.Lane
			car.Lane = car.Intent.TargetLane

			s.EventBus.Emit(LaneChangedEvent{
				OldLane: oldLane,
				NewLane: car.Lane,
				CarID:   car.ID,
			})
		}
	}
}

func (s *Simulation) DetectCongestion() {
	for _, lane := range s.Lanes {
		avgSpeed := lane.AverageSpeed()

		if avgSpeed < lane.MaxSpeed/2 {
			lane.Congested = true
		} else {
			lane.Congested = false
		}
	}
}
