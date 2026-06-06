package model

type Car struct {
	ID           int
	Position     int
	Speed        int
	SafeDistance int
	Intent       CarIntent
	DesiredSpeed int
	Lane         int
}

type CarIntent struct {
	TargetSpeed int
	ChangeLane  bool
	TargetLane  int
}

type CarMovedEvent struct {
	CarID int
	From  int
	To    int
}

func (s *Simulation) ApplySpeedAdjustments() {
	for _, car := range s.Cars {
		target := car.Intent.TargetSpeed
		if target < car.Speed {
			car.Speed--
		} else if target > car.Speed {
			car.Speed++
		}
		if car.Speed < 0 {
			car.Speed = 0
		}
	}
}

func (s *Simulation) MoveCars() {
	for _, car := range s.Cars {
		if car == nil {
			continue
		}
		oldPos := car.Position

		// Find the lane for this car to determine direction
		var carLane *Lane
		for _, l := range s.Lanes {
			if l.ID == car.Lane {
				carLane = l
				break
			}
		}

		if carLane != nil && (carLane.Direction == Westbound || carLane.Direction == Northbound) {
			car.Position -= car.Speed
		} else {
			car.Position += car.Speed
		}

		s.EventBus.Emit(CarMovedEvent{
			CarID: car.ID,
			From:  oldPos,
			To:    car.Position,
		})
	}
}
