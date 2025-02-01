package store

import "github.com/FlauschigDings/lavender"

// EventStore defines a generic event persistence layer for event-sourced aggregates.
type EventStore[E lavender.Event, S lavender.Snapshot] interface {
	// SaveEvents stores new events for the given aggregate.
	SaveEvents(aggregate lavender.CustomAggregate[E, S], events []E) error

	// LoadEvents retrieves all stored events for the given aggregate.
	LoadEvents(aggregate lavender.CustomAggregate[E, S]) ([]E, error)

	// ClearEvents removes all stored events for the given aggregate (e.g., after snapshotting).
	ClearEvents(aggregate lavender.CustomAggregate[E, S]) error
}

// SnapshotStore provides an interface for managing aggregate snapshots.
type SnapshotStore[E lavender.Event, S lavender.Snapshot] interface {
	// SaveSnapshot stores a snapshot of the aggregate's state.
	SaveSnapshot(aggregate lavender.CustomAggregate[E, S], snapshot S) error

	// LoadSnapshot retrieves the most recent snapshot for the given aggregate.
	LoadSnapshot(aggregate lavender.CustomAggregate[E, S]) (*S, error)
}
