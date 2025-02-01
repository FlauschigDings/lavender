package lavender

type Event interface {
	// set the event name
	Name() Name
	// apply an event to the aggregate
	Apply(aggregate CustomAggregate[Event, Snapshot])
}
