package model

type Car struct {
	ID           int
	Position     int
	Speed        int
	DesiredSpeed int
	Lane         int
}

type CarIntent struct {
	TargetSpeed int
	ChangeLane  bool
	TargetLane  int
}
