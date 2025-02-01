package lavender

// EventIdentifier uniquely identifies an event and its associated aggregate.
type EventIdentifier struct {
	Event     Name // The name of the event
	Aggregate Name // The name of the aggregate
}

// EventId creates a new EventIdentifier for a given event and aggregate.
func EventId(event Name, aggregate Name) EventIdentifier {
	return EventIdentifier{Event: event, Aggregate: aggregate}
}
