package store

import (
	"sync"

	"github.com/FlauschigDings/lavender"
)

// InMemoryEventStore is a thread-safe, in-memory store for events and snapshots.
type InMemoryEventStore[E lavender.Event, S lavender.Snapshot] struct {
	Events    sync.Map // Store events as map[lavender.Name][]E
	Snapshots sync.Map // Store snapshots as map[lavender.Name]S
}

// Ensure InMemoryEventStore implements both EventStore and SnapshotStore interfaces.
var _ SnapshotStore[lavender.Event, lavender.Snapshot] = new(InMemoryEventStore[lavender.Event, lavender.Snapshot])
var _ EventStore[lavender.Event, lavender.Snapshot] = new(InMemoryEventStore[lavender.Event, lavender.Snapshot])

// NewInMemoryStore initalizes a new event and snapshot store with default types.
func NewInMemoryStore() *InMemoryEventStore[lavender.Event, lavender.Snapshot] {
	return NewInMemoryCustomStore[lavender.Event, lavender.Snapshot]()
}

// NewInMemoryCustomStore initializes a new custom generic event and snapshot store.
func NewInMemoryCustomStore[E lavender.Event, S lavender.Snapshot]() *InMemoryEventStore[E, S] {
	return &InMemoryEventStore[E, S]{}
}

// SaveEvents appends new events to the aggreagate's event store.
func (store *InMemoryEventStore[E, S]) SaveEvents(aggregate lavender.CustomAggregate[E, S], events []E) error {
	existing, _ := store.Events.Load(aggregate.Name())

	var eventList []E
	if existing != nil {
		eventList = existing.([]E) // Type assertion
	}
	eventList = append(eventList, events...)

	store.Events.Store(aggregate.Name(), eventList)
	return nil
}

// LoadEvents retrieves all stored events from a aggregate.
func (store *InMemoryEventStore[E, S]) LoadEvents(aggregate lavender.CustomAggregate[E, S]) ([]E, error) {
	existing, ok := store.Events.Load(aggregate.Name())
	if !ok {
		return nil, nil
	}
	return existing.([]E), nil
}

// SaveSnapshot stores a snapshot of the aggregate.
func (store *InMemoryEventStore[E, S]) SaveSnapshot(aggregate lavender.CustomAggregate[E, S], snapshot S) error {
	store.Snapshots.Store(aggregate.Name(), snapshot)
	return nil
}

// LoadSnapshot retrieves the last snapshot of the aggregate.
func (store *InMemoryEventStore[E, S]) LoadSnapshot(aggregate lavender.CustomAggregate[E, S]) (*S, error) {
	existing, ok := store.Snapshots.Load(aggregate.Name())
	if !ok {
		return nil, nil
	}
	snapshot := existing.(S) // Type assertion
	return &snapshot, nil
}

// ClearEvents removes all stored events from a aggregate.
func (store *InMemoryEventStore[E, S]) ClearEvents(aggregate lavender.CustomAggregate[E, S]) error {
	store.Events.Store(aggregate.Name(), []E{})
	return nil
}
