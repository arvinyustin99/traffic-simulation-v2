package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Direction string

const (
	Northbound Direction = "northbound"
	Southbound Direction = "southbound"
	Eastbound  Direction = "eastbound"
	Westbound  Direction = "westbound"
)

type Lane struct {
	ID         int
	Cars       []*Car // Rotating array to avoid reallocating
	AllowSpawn bool
	SpawnRate  decimal.Decimal
	Direction  Direction
	MaxSpeed   int
	Congested  bool
}

func (l *Lane) AverageSpeed() int {
	if len(l.Cars) == 0 {
		return 0
	}
	sum := 0
	for _, car := range l.Cars {
		sum += car.Speed
	}
	return sum / len(l.Cars)
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
	StopLine      int
	TrafficLight  *TrafficLight
}

type Simulation struct {
	Tick          int
	Cars          []*Car // Rotating array to avoid reallocating
	Lanes         []*Lane
	Intersections []*Intersection

	TickRate     time.Duration
	SafeDistance int
	MaxCars      int

	Renderer   Renderer
	Metrics    *Metrics
	EventBus   *EventBus
	ChangeLane bool
	Running    bool
}

func NewIntersection() *Intersection {
	return &Intersection{
		IncomingLanes: make([]*Lane, 0),
		OutgoingLanes: make([]*Lane, 0),
		StopLine:      10,
		TrafficLight: &TrafficLight{
			State:         GreenEW,
			GreenDuration: 10,
			RedDuration:   5,
			ElapsedTicks:  0,
			CurrentTick:   0,
		},
	}
}
func (ic *Intersection) RegisterLane(l *Lane, isIncoming bool) {
	if isIncoming {
		ic.IncomingLanes = append(ic.IncomingLanes, l)
	} else {
		ic.OutgoingLanes = append(ic.OutgoingLanes, l)
	}
}
func (ic *Intersection) IsRed(d Direction) bool {
	switch d {
	case Northbound, Southbound:
		return ic.TrafficLight.State == GreenEW
	case Eastbound, Westbound:
		return ic.TrafficLight.State == GreenNS
	}
	return false
}
