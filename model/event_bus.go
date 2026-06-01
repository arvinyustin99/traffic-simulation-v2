package model

type EventBus struct {
	Events []Event
}

type Event interface {
}

func NewEventBus() *EventBus {
	return &EventBus{
		Events: make([]Event, 0),
	}
}

func (eb *EventBus) Flush() []Event {
	eventList := eb.Events
	eb.Events = []Event{}
	return eventList
}

func (eb *EventBus) Emit(e Event) {
	eb.Events = append(eb.Events, e)
}
