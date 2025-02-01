package lavender

import "fmt"

// Aggregate is a convenient alias for CustomAggregate using standard Event and Snapshot types.
type Aggregate = CustomAggregate[Event, Snapshot]

// CustomAggregate defines a generic event-sourced aggregate.
type CustomAggregate[E Event, S Snapshot] interface {
	// Name returns the unique identifier of the aggregate.
	Name() Name

	// Version returns the current version of the aggregate.
	Version() Version

	// ApplyEvent processes and applies an event to the aggregate's state.
	ApplyEvent(event E)

	// Events retrieves all uncommitted events associated with the aggregate.
	Events() []E

	// TakeSnapshot creates a snapshot of the current aggregate state.
	TakeSnapshot() S

	// ApplySnapshot restores the aggregate's state from a given snapshot.
	ApplySnapshot(snapshot S)
}

// Parse attempts to convert a given aggregate into the specified type T.
// If the conversion fails, it panics with an error message.
func Parse[T CustomAggregate[E, S], E Event, S Snapshot](aggregate CustomAggregate[E, S]) T {
	accountAggregate, ok := aggregate.(T)
	if !ok {
		panic(fmt.Sprintf("can't load Aggregate: %#v", accountAggregate))
	}
	return accountAggregate
}
