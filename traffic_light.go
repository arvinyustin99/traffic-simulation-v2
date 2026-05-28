package main

func (tl *TrafficLight) Update() {
	tl.ElapsedTicks++

	switch tl.State {
	case GreenNS:
		if tl.ElapsedTicks >= tl.GreenDuration {
			tl.State = GreenEW
			tl.ElapsedTicks = 0
		}
	case GreenEW:
		if tl.ElapsedTicks >= tl.GreenDuration {
			tl.State = GreenNS
			tl.ElapsedTicks = 0
		}
	}
}
