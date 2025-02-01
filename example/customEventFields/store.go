package customeventfields

import (
	"log"

	"github.com/FlauschigDings/lavender"
	"github.com/FlauschigDings/lavender/store"
)

// Wrapper of the InMemoryEventStore
type CustomMemoryStore[E Event, S lavender.Snapshot] struct {
	*store.InMemoryEventStore[E, S]
}

// Typesafety checks for SnapshotStore and EventStore.
var _ store.SnapshotStore[Event, lavender.Snapshot] = new(CustomMemoryStore[Event, lavender.Snapshot])
var _ store.EventStore[Event, lavender.Snapshot] = new(CustomMemoryStore[Event, lavender.Snapshot])

// Create a custom store instance.
func NewCustomMemoryStore() *CustomMemoryStore[Event, lavender.Snapshot] {
	return &CustomMemoryStore[Event, lavender.Snapshot]{
		store.NewInMemoryCustomStore[Event, lavender.Snapshot](),
	}
}

// Overwrite the old function and select the custom event field and write the operator in the chat.
func (c *CustomMemoryStore[E, S]) SaveEvents(aggregate lavender.CustomAggregate[E, S], events []E) error {
	for _, event := range events {
		log.Printf("event has been added from %#v", event.Operator())
	}
	return c.InMemoryEventStore.SaveEvents(aggregate, events)
}
