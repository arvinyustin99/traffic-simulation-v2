package model

import "github.com/shopspring/decimal"

type Lane struct {
	ID         int
	Cars       []*Car // Rotating array to avoid reallocating
	AllowSpawn bool
	SpawnRate  decimal.Decimal
	Direction  string
	MaxSpeed   int
}

type TrafficLight struct {
	ElapsedTicks  int
	State         TrafficLightState
	GreenDuration int
	RedDuration   int
	CurrentTick   int
}

type TrafficLightState string

const (
	GreenNS TrafficLightState = "green_ns"
	GreenEW TrafficLightState = "green_ew"
	// YellowNS TrafficLightState = "yellow"
	// RedNS    TrafficLightState = "red"
)

type Intersection struct {
	IncomingLanes []*Lane
	OutgoingLanes []*Lane
	TrafficLight  *TrafficLight
}

type Simulation struct {
	Tick         int
	Cars         []*Car // Rotating array to avoid reallocating
	Lanes        []*Lane
	Intersection []*Intersection
}
