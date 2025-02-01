package repo

import (
	"sync"

	"github.com/FlauschigDings/lavender"
	"github.com/FlauschigDings/lavender/store"
)

// Repository is an alias for CustomRepository with lavender.Event and lavender.Snapshot types.
type Repository = CustomRepository[lavender.Event, lavender.Snapshot]

// CustomRepository represents a repository that handles event and snapshot storage, including caching and auto-snapshot logic.
type CustomRepository[E lavender.Event, S lavender.Snapshot] struct {
	// aggregateCache is a thread-safe map to cache aggregates for faster access.
	aggregateCache sync.Map // map[lavender.Name]lavender.CustomAggregate[E, S]

	// aggregateCacheActive controls whether aggregate caching is enabled for faster access.
	aggregateCacheActive bool

	// EventStore is the store responsible for saving and loading events.
	EventStore store.EventStore[E, S]

	// SnapshotStore is the store responsible for saving and loading snapshots.
	SnapshotStore store.SnapshotStore[E, S]

	// AutoSnapshotHook defines the logic for when to automatically create snapshots.
	AutoSnapshotHook AutoSnapshotHook[E, S]
}

// NewRepository creates a new CustomRepository with event and snapshot stores, and caching enabled by default.
func NewRepository[E lavender.Event, S lavender.Snapshot](eventStore store.EventStore[E, S], snapshotStore store.SnapshotStore[E, S]) *CustomRepository[E, S] {
	return NewRepositoryConstructor[E, S](true, eventStore, snapshotStore)
}

// NewRepositoryConstructor initializes a CustomRepository with the specified settings.
// It accepts a boolean flag to activate or deactivate the aggregate cache.
func NewRepositoryConstructor[E lavender.Event, S lavender.Snapshot](aggregateCacheActive bool, eventStore store.EventStore[E, S], snapshotStore store.SnapshotStore[E, S]) *CustomRepository[E, S] {
	return &CustomRepository[E, S]{
		aggregateCacheActive: aggregateCacheActive,
		EventStore:           eventStore,
		SnapshotStore:        snapshotStore,
		AutoSnapshotHook: func(aggregate lavender.CustomAggregate[E, S], items []E) bool {
			// Hook to decide when to snapshot based on the number of events
			return len(items) > 100
		},
	}
}

// AutoSnapshot checks whether the aggregate should be snapshotted based on the number of events.
// If the condition is met, it creates a snapshot and clears the event log.
func (r *CustomRepository[E, S]) AutoSnapshot(aggregate lavender.CustomAggregate[E, S]) error {
	// Load the aggregate's events
	items, err := r.EventStore.LoadEvents(aggregate)
	if err != nil {
		return err
	}

	// Apply snapshot if the auto-snapshot hook condition is met
	if r.AutoSnapshotHook(aggregate, items) {
		if err := r.CreateSnapshot(aggregate); err != nil {
			return err
		}
		// Clear the event log after snapshotting
		err := r.EventStore.ClearEvents(aggregate)
		return err
	}
	return nil
}

// LoadAggregate loads the aggregate's state from either cache, snapshot, or events.
func (r *CustomRepository[E, S]) LoadAggregate(aggregate lavender.CustomAggregate[E, S]) error {
	// First, try to load from cache if caching is enabled
	if cache := r.LoadCache(aggregate); cache != nil {
		aggregate.ApplySnapshot((*cache).TakeSnapshot())
		return nil
	}

	// Load the snapshot from the snapshot store
	snapshot, err := r.SnapshotStore.LoadSnapshot(aggregate)
	if err != nil {
		return err
	}

	// Apply the snapshot if it exists
	if snapshot != nil {
		aggregate.ApplySnapshot(*snapshot)
	}

	// Load events and apply them to the aggregate
	events, err := r.EventStore.LoadEvents(aggregate)
	if err != nil {
		return err
	}

	// Apply all loaded events to the aggregate
	for _, event := range events {
		aggregate.ApplyEvent(event)
	}

	// Cache the aggregate for future access
	r.saveCache(aggregate)
	return nil
}

// CreateSnapshot creates a snapshot of the aggregate's current state and stores it in the snapshot store.
func (r *CustomRepository[E, S]) CreateSnapshot(aggregate lavender.CustomAggregate[E, S]) error {
	// Ensure the aggregate is fully loaded before snapshotting
	if err := r.LoadAggregate(aggregate); err != nil {
		return err
	}

	// Take a snapshot of the aggregate
	snapshot := aggregate.TakeSnapshot()

	// Save the snapshot in the snapshot store
	err := r.SnapshotStore.SaveSnapshot(aggregate, snapshot)
	return err
}

// ClearEventLog clears the event log for the given aggregate in the event store.
func (r *CustomRepository[E, S]) ClearEventLog(aggregate lavender.CustomAggregate[E, S]) error {
	// Clear the events stored for the aggregate
	err := r.EventStore.ClearEvents(aggregate)
	return err
}

// AddEvent appends events to the aggregate, potentially triggering a snapshot based on the auto-snapshot condition.
func (r *CustomRepository[E, S]) AddEvent(aggregate lavender.CustomAggregate[E, S], events ...E) error {
	// Automatically snapshot the aggregate if needed
	if err := r.AutoSnapshot(aggregate); err != nil {
		return err
	}

	// Load the aggregate to apply events
	if err := r.LoadAggregate(aggregate); err != nil {
		return err
	}

	// Apply each event to the aggregate
	for _, event := range events {
		aggregate.ApplyEvent(event)
	}

	// Cache the aggregate for future access
	r.saveCache(aggregate)

	// Save the new events to the event store
	return r.EventStore.SaveEvents(aggregate, events)
}
